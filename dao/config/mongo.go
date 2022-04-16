package config

import (
	"context"
	"time"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/config/ecode"
	"github.com/chenjie199234/config/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func (d *Dao) MongoGetAllGroups(ctx context.Context, searchfilter string) ([]string, error) {
	regex := "^config_"
	if searchfilter != "" {
		regex += ".*" + searchfilter + ".*"
	}
	r, e := d.mongo.ListDatabaseNames(ctx, bson.M{"name": bson.M{"$regex": regex}})
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

func (d *Dao) MongoCreate(ctx context.Context, groupname, appname, cipher string, encrypt datahandler) (e error) {
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
			if mongo.IsDuplicateKeyError(e) {
				e = ecode.ErrAppAlreadyExist
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	index := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "index", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, e = col.Indexes().CreateOne(sctx, index); e != nil {
		return
	}
	appconfig := "{}"
	sourceconfig := "{}"
	if cipher != "" {
		appconfig = encrypt(cipher, appconfig)
		sourceconfig = encrypt(cipher, sourceconfig)
	}
	if _, e = col.InsertOne(sctx, bson.M{
		"index":             0,
		"cipher":            cipher,
		"cur_index":         0,
		"max_index":         0,
		"cur_version":       0,
		"cur_app_config":    appconfig,
		"cur_source_config": sourceconfig,
	}); e != nil {
		return
	}
	return
}

type datahandler func(cipher string, origindata string) (newdata string)

func (d *Dao) MongoUpdateCipher(ctx context.Context, groupname, appname, oldcipher, newcipher string, decrypt, encrypt datahandler) (e error) {
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
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	summary := &model.Summary{}
	if e = col.FindOne(sctx, bson.M{"index": 0}).Decode(summary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if summary.Cipher != oldcipher {
		e = ecode.ErrWrongCipher
		return
	}
	if oldcipher != "" {
		summary.CurAppConfig = decrypt(oldcipher, summary.CurAppConfig)
		summary.CurSourceConfig = decrypt(oldcipher, summary.CurSourceConfig)
	}
	if newcipher != "" {
		summary.CurAppConfig = encrypt(newcipher, summary.CurAppConfig)
		summary.CurSourceConfig = encrypt(newcipher, summary.CurSourceConfig)
	}
	if _, e = col.UpdateOne(sctx, bson.M{"index": 0}, bson.M{"$set": bson.M{"cipher": newcipher, "cur_app_config": summary.CurAppConfig, "cur_source_config": summary.CurSourceConfig}}); e != nil {
		return
	}
	cursor, e := col.Find(sctx, bson.M{"index": bson.M{"$gt": 0}}, options.Find().SetSort(bson.M{"index": -1}))
	if e != nil {
		return
	}
	for cursor.Next(sctx) {
		tmp := &model.Config{}
		if e = cursor.Decode(tmp); e != nil {
			return
		}
		if oldcipher != "" {
			tmp.AppConfig = decrypt(oldcipher, tmp.AppConfig)
			tmp.SourceConfig = decrypt(oldcipher, tmp.SourceConfig)
		}
		if newcipher != "" {
			tmp.AppConfig = encrypt(newcipher, tmp.AppConfig)
			tmp.SourceConfig = encrypt(newcipher, tmp.SourceConfig)
		}
		if _, e = col.UpdateOne(sctx, bson.M{"index": tmp.Index}, bson.M{"$set": bson.M{"app_config": tmp.AppConfig, "source_config": tmp.SourceConfig}}); e != nil {
			return
		}
	}
	e = cursor.Err()
	return
}

//index == 0 get the current index's config
func (d *Dao) MongoGetConfig(ctx context.Context, groupname, appname string, index uint32, decrypt datahandler) (*model.Summary, *model.Config, error) {
	col := d.mongo.Database("config_"+groupname, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(appname)
	var summary *model.Summary
	var config *model.Config
	if index != 0 {
		filter := bson.M{"$or": bson.A{bson.M{"index": 0}, bson.M{"index": index}}}
		cursor, e := col.Find(ctx, filter, options.Find().SetSort(bson.M{"index": 1}))
		if e != nil {
			return nil, nil, e
		}
		for cursor.Next(ctx) {
			if summary == nil {
				tmps := &model.Summary{}
				if e = cursor.Decode(tmps); e != nil {
					return nil, nil, e
				}
				summary = tmps
			} else {
				tmpc := &model.Config{}
				if e = cursor.Decode(tmpc); e != nil {
					return nil, nil, e
				}
				config = tmpc
			}
		}
		if e := cursor.Err(); e != nil {
			return nil, nil, e
		}
	} else {
		summary := &model.Summary{}
		if e := col.FindOne(ctx, bson.M{"index": 0}).Decode(summary); e != nil {
			if e == mongo.ErrNoDocuments {
				e = nil
			}
			return nil, nil, e
		}
		config = &model.Config{
			Index:        summary.CurIndex,
			AppConfig:    summary.CurAppConfig,
			SourceConfig: summary.CurSourceConfig,
		}
	}
	if summary != nil && summary.Cipher != "" {
		summary.CurAppConfig = decrypt(summary.Cipher, summary.CurAppConfig)
		summary.CurSourceConfig = decrypt(summary.Cipher, summary.CurSourceConfig)
		if config != nil {
			config.AppConfig = decrypt(summary.Cipher, config.AppConfig)
			config.SourceConfig = decrypt(summary.Cipher, config.SourceConfig)
		}
	}
	return summary, config, nil
}
func (d *Dao) MongoSetConfig(ctx context.Context, groupname, appname, appconfig, sourceconfig string, encrypt datahandler) (e error) {
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
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	summary := &model.Summary{}
	if e = col.FindOne(sctx, bson.M{"index": 0}).Decode(summary); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		}
		return
	}
	if summary.Cipher != "" {
		appconfig = encrypt(summary.Cipher, appconfig)
		sourceconfig = encrypt(summary.Cipher, sourceconfig)
	}
	updateSummary := bson.M{
		"cur_version":       summary.CurVersion + 1,
		"max_index":         summary.MaxIndex + 1,
		"cur_index":         summary.MaxIndex + 1,
		"cur_app_config":    appconfig,
		"cur_source_config": sourceconfig,
	}
	if _, e = col.UpdateOne(sctx, bson.M{"index": 0}, bson.M{"$set": updateSummary}); e != nil {
		return
	}
	updateConfig := bson.M{
		"app_config":    appconfig,
		"source_config": sourceconfig,
	}
	if _, e = col.UpdateOne(sctx, bson.M{"index": summary.MaxIndex + 1}, bson.M{"$set": updateConfig}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	return
}
func (d *Dao) MongoRollbackConfig(ctx context.Context, groupname, appname string, index uint32) (e error) {
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
	col := d.mongo.Database("config_" + groupname).Collection(appname)
	config := &model.Config{}
	if e = col.FindOne(sctx, bson.M{"index": index}).Decode(config); e != nil {
		if e == mongo.ErrNoDocuments {
			e = ecode.ErrIndexNotExist
		}
		return
	}
	updateSummary := bson.M{
		"$set": bson.M{
			"cur_index":         index,
			"cur_app_config":    config.AppConfig,
			"cur_source_config": config.SourceConfig,
		},
		"$inc": bson.M{
			"cur_version": 1,
		},
	}
	if r := col.FindOneAndUpdate(sctx, bson.M{"index": 0}, updateSummary); r.Err() != nil {
		if r.Err() == mongo.ErrNoDocuments {
			e = ecode.ErrAppNotExist
		} else {
			e = r.Err()
		}
	}
	return
}

//first key groupname,second key appname,value curconfig
type WatchRefreshHandler func(map[string]map[string]*model.Summary)
type WatchUpdateHandler func(string, string, *model.Summary)
type WatchDeleteGroupHandler func(groupname string)
type WatchDeleteAppHandler func(groupname, appname string)
type WatchDeleteConfigHandler func(groupname, appname string, id string)

func (d *Dao) getall(decrypt datahandler) (map[string]map[string]*model.Summary, error) {
	groups, e := d.MongoGetAllGroups(context.Background(), "")
	if e != nil {
		return nil, e
	}
	result := make(map[string]map[string]*model.Summary, len(groups))
	for _, group := range groups {
		tmpgroup := make(map[string]*model.Summary)
		apps, e := d.MongoGetAllApps(context.Background(), group, "")
		if e != nil {
			return nil, e
		}
		for _, app := range apps {
			summary, _, e := d.MongoGetConfig(context.Background(), group, app, 0, decrypt)
			if e != nil {
				return nil, e
			}
			if summary == nil || summary.CurVersion == 0 {
				continue
			}
			tmpgroup[app] = summary
		}
		if len(tmpgroup) != 0 {
			result[group] = tmpgroup
		}
	}
	return result, nil
}
func (d *Dao) MongoWatchConfig(refresh WatchRefreshHandler, update WatchUpdateHandler, delG WatchDeleteGroupHandler, delA WatchDeleteAppHandler, delC WatchDeleteConfigHandler, decrypt datahandler) error {
	watchfilter := mongo.Pipeline{bson.D{primitive.E{Key: "$match", Value: bson.M{"ns.db": bson.M{"$regex": "^config_"}}}}}
	stream, e := d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup))
	if e != nil {
		return e
	}
	datas, e := d.getall(decrypt)
	if e != nil {
		return e
	}
	refresh(datas)
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 5)
				if stream, e = d.mongo.Watch(context.Background(), watchfilter, options.ChangeStream().SetFullDocument(options.UpdateLookup)); e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] reconnect stream error:", e)
					stream = nil
					continue
				}
				datas, e = d.getall(decrypt)
				if e != nil {
					log.Error(nil, "[dao.MongoWatchConfig] refresh after reconnect stream error:", e)
					stream.Close(context.Background())
					stream = nil
					continue
				}
				refresh(datas)
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
					fallthrough
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
					if s.Cipher != "" {
						s.CurAppConfig = decrypt(s.Cipher, s.CurAppConfig)
						s.CurSourceConfig = decrypt(s.Cipher, s.CurSourceConfig)
					}
					update(groupname, appname, s)
				case "delete":
					//delete document
					groupname := stream.Current.Lookup("ns").Document().Lookup("db").StringValue()[7:]
					appname := stream.Current.Lookup("ns").Document().Lookup("coll").StringValue()
					id := stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex()
					delC(groupname, appname, id)
				}
			}
			if stream.Err() != nil {
				log.Error(nil, "[dao.MongoWatchConfig]", stream.Err())
			}
			stream.Close(context.Background())
			stream = nil
		}
	}()
	return nil
}
