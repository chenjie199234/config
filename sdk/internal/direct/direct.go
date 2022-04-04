package direct

import (
	"context"
	"time"

	"github.com/chenjie199234/config/model"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type updater func([]byte) error

func NewDirectSdk(selfgroup, selfname string, url string, updateapp, updatesource updater) error {
	client, e := newMongo(url)
	if e != nil {
		return e
	}
	watchfilter := mongo.Pipeline{bson.D{bson.E{Key: "$match", Value: bson.M{"fullDocument.index": 0}}}}
	watchop := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, e := client.Database(selfgroup).Collection(selfname).Watch(context.Background(), watchfilter, watchop)
	if e != nil {
		return e
	}
	col := client.Database(selfgroup, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection(selfname)
	//get first,then watch change stream
	s := &model.Summary{}
	c := &model.Config{}
	if e = col.FindOne(context.Background(), bson.M{"index": 0}).Decode(s); e != nil && e != mongo.ErrNoDocuments {
		return e
	}
	if s.CurVersion > 0 {
		if e = client.Database(selfgroup).Collection(selfname).FindOne(context.Background(), bson.M{"index": s.CurIndex}).Decode(c); e != nil {
			return e
		}
		if e = updateapp(common.Str2byte(c.AppConfig)); e != nil {
			return e
		}
		if e = updatesource(common.Str2byte(c.SourceConfig)); e != nil {
			return e
		}
	}
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 5)
				if stream, e = client.Database(selfgroup).Collection(selfname).Watch(context.Background(), watchfilter, watchop); e != nil {
					log.Error(nil, "[config.sdk.watch] reconnect error:", e)
					stream = nil
					continue
				}
			}
			for stream.Next(context.Background()) {
				tmps := &model.Summary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(tmps); e != nil {
					log.Error(nil, "[config.sdk.watch] summary info broken,error:", e)
					continue
				}
				if tmps.CurVersion <= s.CurVersion {
					continue
				}
				tmpc := &model.Config{}
				if e = col.FindOne(context.Background(), bson.M{"index": tmps.CurIndex}).Decode(tmpc); e != nil {
					log.Error(nil, "[config.sdk.watch] get config on index:", tmps.CurIndex, "error:", e)
					continue
				}
				s = tmps
				c = tmpc
				if e = updateapp(common.Str2byte(c.AppConfig)); e != nil {
					log.Error(nil, "[config.sdk.watch] update appconfig error:", e)
				}
				if e = updatesource(common.Str2byte(c.SourceConfig)); e != nil {
					log.Error(nil, "[config.sdk.watch] update sourceconfig error:", e)
				}
			}
			log.Error(nil, "[config.sdk.watch] error:", stream.Err())
			stream.Close(context.Background())
			stream = nil
		}
	}()
	return nil
}
func newMongo(url string) (*mongo.Client, error) {
	op := options.Client().ApplyURI(url)
	op = op.SetMaxPoolSize(2)
	op = op.SetHeartbeatInterval(time.Second * 5)
	op = op.SetReadPreference(readpref.SecondaryPreferred())
	op = op.SetReadConcern(readconcern.Local())
	db, e := mongo.Connect(context.Background(), op)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(context.Background(), nil); e != nil {
		return nil, e
	}
	return db, nil
}
