package config

import (
	"context"
	"time"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/config/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoGetAllGroups(ctx context.Context, searchfilter string) ([]string, error) {
	r, e := d.mongo.ListDatabaseNames(ctx, bson.M{"name": bson.M{"$regex": "^config_.*" + searchfilter + ".*"}})
	if e != nil {
		return nil, e
	}
	for i := range r {
		r[i] = r[i][7:]
	}
	return r, nil
}
func (d *Dao) MongoGetAllApps(ctx context.Context, groupname, searchfilter string) ([]string, error) {
	return d.mongo.Database("config_"+groupname).ListCollectionNames(ctx, bson.M{"name": bson.M{"$regex": searchfilter}})
}

//index == 0 get the current index's config
func (d *Dao) MongoGetCOnfig(ctx context.Context, groupname, appname string, index uint32) (summary *model.Summary, config *model.Config, e error) {
	if index != 0 {
		col := d.mongo.Database("config_"+groupname, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(appname)
		filter := bson.M{"$or": bson.A{bson.M{"index": 0}, bson.M{"index": index}}}
		var cursor *mongo.Cursor
		cursor, e = col.Find(ctx, filter, options.Find().SetSort(bson.M{"index": 1}))
		if e != nil {
			if e == mongo.ErrNoDocuments {
				e = nil
			}
			return
		}
		for cursor.Next(ctx) {
			if summary == nil {
				tmps := &model.Summary{}
				if e = cursor.Decode(tmps); e != nil {
					return
				}
				summary = tmps
			} else {
				tmpc := &model.Config{}
				if e = cursor.Decode(tmpc); e != nil {
					return
				}
				config = tmpc
			}
		}
		e = cursor.Err()
	} else {
		tmps := &model.Summary{}
		if e = d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"index": 0}).Decode(tmps); e != nil {
			if e == mongo.ErrNoDocuments {
				e = nil
			}
			return
		}
		summary = tmps
		if summary.CurVersion > 0 {
			tmpc := &model.Config{}
			if e = d.mongo.Database("config_"+groupname).Collection(appname).FindOne(ctx, bson.M{"index": summary.CurIndex}).Decode(tmpc); e != nil {
				return
			}
			config = tmpc
		}
	}
	return
}
func (d *Dao) MongoSetConfig(ctx context.Context, groupname, appname, appconfig, sourceconfig string) (e error) {
	cursor, e := d.mongo.Database("config_" + groupname).Collection(appname).Indexes().List(ctx)
	if e != nil {
		return
	}
	if cursor.RemainingBatchLength() == 0 {
		index := mongo.IndexModel{
			Keys:    bson.D{primitive.E{Key: "index", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		if _, e = d.mongo.Database("config_"+groupname).Collection(appname).Indexes().CreateOne(ctx, index); e != nil {
			return
		}
	}
	var s mongo.Session
	s, e = d.mongo.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local()))
	if e != nil {
		return
	}
	defer s.EndSession(ctx)
	sctx := mongo.NewSessionContext(ctx, s)
	if e = s.StartTransaction(); e != nil {
		return
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	filter := bson.M{"index": 0}
	update1 := bson.A{
		bson.M{
			"$set": bson.M{
				"cur_version": bson.M{
					"$ifNull": bson.A{
						bson.M{"$add": bson.A{"$cur_version", 1}},
						1,
					},
				},
				"max_index": bson.M{
					"$ifNull": bson.A{
						bson.M{"$add": bson.A{"$max_index", 1}},
						1,
					},
				},
			},
		},
		bson.M{
			"$set": bson.M{
				"cur_index": bson.M{
					"$toInt": "$max_index",
				},
			},
		},
	}
	summary := &model.Summary{}
	r := d.mongo.Database("config_"+groupname).Collection(appname).FindOneAndUpdate(sctx, filter, update1, options.FindOneAndUpdate().SetUpsert(true))
	if r.Err() != nil && r.Err() != mongo.ErrNoDocuments {
		e = r.Err()
		return
	} else if r.Err() == nil {
		if e = r.Decode(summary); e != nil {
			return
		}
	}
	filter["index"] = summary.MaxIndex + 1
	update2 := bson.M{"$set": bson.M{"app_config": appconfig, "source_config": sourceconfig}}
	_, e = d.mongo.Database("config_"+groupname).Collection(appname).UpdateOne(sctx, filter, update2, options.Update().SetUpsert(true))
	return
}
func (d *Dao) MongoRollbackConfig(ctx context.Context, groupname, appname string, index uint32) (bool, error) {
	filter := bson.M{"index": 0, "max_index": bson.M{"$gte": index}}
	update := bson.M{
		"$set": bson.M{"cur_index": index},
		"$inc": bson.M{"cur_version": 1},
	}
	r, e := d.mongo.Database("config_"+groupname).Collection(appname).UpdateOne(ctx, filter, update)
	if e != nil {
		return false, e
	}
	if r.MatchedCount == 0 {
		return false, nil
	}
	return true, nil
}

type WatchInitHandler func([]*model.Current)
type WatchUpdateHandler func(*model.Current)
type WatchDeleteGroupHandler func(groupname string)
type WatchDeleteAppHandler func(groupname, appname string)
type WatchDeleteConfigHandler func(groupname, appname string, id string)

func (d *Dao) MongoWatchConfig(init WatchInitHandler, update WatchUpdateHandler, delG WatchDeleteGroupHandler, delA WatchDeleteAppHandler, delC WatchDeleteConfigHandler) error {
	watchfilter := mongo.Pipeline{bson.D{primitive.E{Key: "$match", Value: bson.M{"ns.db": bson.M{"$regex": "^config_"}}}}}
	stream, e := d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if e != nil {
		return e
	}
	//TODO
	// init()
	for {
		for stream == nil {
			//reconnect
			time.Sleep(time.Millisecond * 5)
			if stream, e = d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup)); e != nil {
				log.Error(nil, "[dao.MongoWatchConfig]", e)
				stream = nil
				continue
			}
			//TODO
			//init()
		}
		for stream.Next(context.Background()) {
			switch stream.Current.Lookup("operationType").StringValue() {
			case "dropDatabase":
				//drop database
				groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
				delG(groupname)
			case "drop":
				//drop collection
				groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
				appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
				delA(groupname, appname)
			case "insert":
				//insert document
				groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
				appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
				index, ok := stream.Current.Lookup("fullDocument").Document().Lookup("index").Int32OK()
				if !ok {
					//unknown doc
					continue
				}
				if index != 0 {
					//this is not the summary
					continue
				}
				//this is the summary
				s := &model.Summary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(s); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] group:", groupname, "app:", appname, "summary data broken:", e)
					continue
				}
				c := &model.Config{}
				if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(context.Background(), bson.M{"index": s.CurIndex}).Decode(c); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] group:", groupname, "app:", appname, "index:", s.CurIndex, "config data broken:", e)
					continue
				}
				update(&model.Current{
					ID:           s.ID.Hex(),
					GroupName:    groupname,
					AppName:      appname,
					CurVersion:   s.CurVersion,
					AppConfig:    c.AppConfig,
					SourceConfig: c.SourceConfig,
				})
			case "update":
				//update document
				groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
				appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
				index, ok := stream.Current.Lookup("fullDocument").Document().Lookup("index").Int32OK()
				if !ok {
					//unknown doc
					continue
				}
				if index != 0 {
					//this is not the summary
					continue
				}
				//this is the summary
				s := &model.Summary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(s); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] group:", groupname, "app:", appname, "summary data broken:", e)
					continue
				}
				c := &model.Current{}
				if e := d.mongo.Database("config_"+groupname).Collection(appname).FindOne(context.Background(), bson.M{"index": s.CurIndex}).Decode(c); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] group:", groupname, "app:", appname, "index:", s.CurIndex, "config data broken:", e)
					continue
				}
				update(&model.Current{
					ID:           s.ID.Hex(),
					GroupName:    groupname,
					AppName:      appname,
					CurVersion:   s.CurVersion,
					AppConfig:    c.AppConfig,
					SourceConfig: c.SourceConfig,
				})
			case "delete":
				//delete document
				groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
				appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
				id := stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex()
				delC(groupname, appname, id)
			}
		}
	}
}
