package selfsdk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/chenjie199234/config/model"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/util/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var status int32

//url format [mongodb/mongodb+srv]://[username:password@]host1,...,hostN[/dbname][?param1=value1&...&paramN=valueN]
func NewDirectSdk(selfgroup string, url string) error {
	if !atomic.CompareAndSwapInt32(&status, 0, 1) {
		return nil
	}
	client, e := newMongo(url, selfgroup)
	if e != nil {
		return e
	}
	watchfilter := mongo.Pipeline{bson.D{bson.E{Key: "$match", Value: bson.M{"$or": bson.A{bson.M{"operationType": "delete"}, bson.M{"fullDocument.index": 0}}}}}}
	watchop := options.ChangeStream().SetFullDocument(options.UpdateLookup)
	stream, e := client.Database("config_"+selfgroup).Collection("config").Watch(context.Background(), watchfilter, watchop)
	if e != nil {
		return e
	}
	col := client.Database("config_"+selfgroup, options.Database().SetReadPreference(readpref.Primary()).SetReadConcern(readconcern.Local())).Collection("config")
	//get first,then watch change stream
	s := &model.Summary{}
	if e = col.FindOne(context.Background(), bson.M{"index": 0}).Decode(s); e != nil && e != mongo.ErrNoDocuments {
		return e
	} else if e != nil {
		//not exist
		s.CurAppConfig = "{}"
		s.CurSourceConfig = "{}"
	}
	if s.Cipher != "" {
		s.CurAppConfig = model.Decrypt(s.Cipher, s.CurAppConfig)
		s.CurSourceConfig = model.Decrypt(s.Cipher, s.CurSourceConfig)
	}
	if e = updateApp(common.Str2byte(s.CurAppConfig)); e != nil {
		return e
	}
	if e = updateSource(common.Str2byte(s.CurSourceConfig)); e != nil {
		return e
	}
	go func() {
		for {
			for stream == nil {
				//reconnect
				time.Sleep(time.Millisecond * 100)
				if stream, e = client.Database("config_"+selfgroup).Collection("config").Watch(context.Background(), watchfilter, watchop); e != nil {
					log.Error(nil, "[config.selfsdk.watch] reconnect error:", e)
					stream = nil
					continue
				}
				//refresh after reconnect
				tmps := &model.Summary{}
				if e = col.FindOne(context.Background(), bson.M{"index": 0}).Decode(tmps); e != nil && e != mongo.ErrNoDocuments {
					log.Error(nil, "[config.selfsdk.watch] refresh after reconnect error:", e)
					stream.Close(context.Background())
					stream = nil
					continue
				} else if e != nil {
					//not exist
					tmps.CurAppConfig = "{}"
					tmps.CurSourceConfig = "{}"
				}
				s = tmps
				if s.Cipher != "" {
					s.CurAppConfig = model.Decrypt(s.Cipher, s.CurAppConfig)
					s.CurSourceConfig = model.Decrypt(s.Cipher, s.CurSourceConfig)
				}
				if e = updateApp(common.Str2byte(s.CurAppConfig)); e != nil {
					log.Error(nil, "[config.selfsdk.watch] refresh after reconnect error:", e)
					stream.Close(context.Background())
					stream = nil
					continue
				}
				if e = updateSource(common.Str2byte(s.CurSourceConfig)); e != nil {
					log.Error(nil, "[config.selfsdk.watch] refresh after reconnect error:", e)
					stream.Close(context.Background())
					stream = nil
					continue
				}
			}
			for stream.Next(context.Background()) {
				if stream.Current.Lookup("operationType").StringValue() == "delete" {
					if s.ID.Hex() == stream.Current.Lookup("documentKey").Document().Lookup("_id").ObjectID().Hex() {
						//delete the summary,need to refresh
						s.ID = primitive.ObjectID{}
						s.Cipher = ""
						s.CurIndex = 0
						s.MaxIndex = 0
						s.CurVersion = 0
						s.CurAppConfig = "{}"
						s.CurSourceConfig = "{}"
						if e = updateApp(common.Str2byte(s.CurAppConfig)); e != nil {
							log.Error(nil, "[config.selfsdk.watch] update appconfig error:", e)
							break
						}
						if e = updateSource(common.Str2byte(s.CurSourceConfig)); e != nil {
							log.Error(nil, "[config.selfsdk.watch] update sourceconfig error:", e)
							break
						}
					}
					continue
				}
				tmps := &model.Summary{}
				if e := stream.Current.Lookup("fullDocument").Unmarshal(tmps); e != nil {
					log.Error(nil, "[config.selfsdk.watch] summary info broken,error:", e)
					continue
				}
				if tmps.CurVersion <= s.CurVersion {
					continue
				}
				s = tmps
				if s.Cipher != "" {
					s.CurAppConfig = model.Decrypt(s.Cipher, s.CurAppConfig)
					s.CurSourceConfig = model.Decrypt(s.Cipher, s.CurSourceConfig)
				}
				if e = updateApp(common.Str2byte(s.CurAppConfig)); e != nil {
					log.Error(nil, "[config.selfsdk.watch] update appconfig error:", e)
					break
				}
				if e = updateSource(common.Str2byte(s.CurSourceConfig)); e != nil {
					log.Error(nil, "[config.selfsdk.watch] update sourceconfig error:", e)
					break
				}
			}
			if stream.Err() != nil {
				log.Error(nil, "[config.selfsdk.watch] error:", stream.Err())
			}
			stream.Close(context.Background())
			stream = nil
		}
	}()
	return nil
}

var defaultAppConfig = `{
	"handler_timeout":{
		"/config.status/ping":{
			"GET":"200ms",
			"CRPC":"200ms",
			"GRPC":"200ms"
		},
		"/config.config/watch":{
			"POST":"0s"
		},
		"/config.config/updatechiper":{
			"POST":"0s"
		}
	},
	"handler_rate":[{
		"Path":"/config.status/ping",
		"Method":["GET","GRPC","CRPC"],
		"MaxPerSec":10
	}],
	"access_keys":{
		"default":"default_sec_key",
		"/config.status/ping":"specific_sec_key"
	},
	"service":{

	}
}`
var defaultSourceConfig = `{
	"cgrpc_server":{
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"heart_probe":"1.5s"
	},
	"cgrpc_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"heart_probe":"1.5s"
	},
	"crpc_server":{
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"heart_probe":"1.5s"
	},
	"crpc_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"heart_probe":"1.5s"
	},
	"web_server":{
		"close_mode":1,
		"connect_timeout":"200ms",
		"global_timeout":"200ms",
		"idle_timeout":"5s",
		"heart_probe":"1.5s",
		"static_file":"./src",
		"web_cors":{
			"cors_origin":["*"],
			"cors_header":["*"],
			"cors_expose":[]
		}
	},
	"web_client":{
		"connect_timeout":"200ms",
		"global_timeout":"0",
		"idle_timeout":"5s",
		"heart_probe":"1.5s"
	},
	"mongo":{
		"config_mongo":{
			"url":"%s",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"500ms",
			"conn_timeout":"500ms"
		}
	},
	"sql":{
		"example_sql":{
			"url":"[username:password@][protocol(address)][/dbname][?param1=value1&...&paramN=valueN]",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"200ms",
			"conn_timeout":"200ms"
		}
	},
	"redis":{
		"example_redis":{
			"url":"[redis/rediss]://[[username:]password@]host[/dbindex]",
			"max_open":100,
			"max_idletime":"10m",
			"io_timeout":"200ms",
			"conn_timeout":"200ms"
		}
	},
	"kafka_pub":[
		{
			"addrs":["127.0.0.1:12345"],
			"username":"example",
			"password":"example",
			"auth_method":3,
			"compress_method":2,
			"topic_name":"example_topic",
			"io_timeout":"500ms",
			"conn_timeout":"200ms"
		}
	],
	"kafka_sub":[
		{
			"addrs":["127.0.0.1:12345"],
			"username":"example",
			"password":"example",
			"auth_method":3,
			"topic_name":"example_topic",
			"group_name":"example_group",
			"conn_timeout":"200ms",
			"start_offset":-2,
			"commit_interval":"0s"
		}
	]
}`

func newMongo(url string, groupname string) (db *mongo.Client, e error) {
	op := options.Client().ApplyURI(url)
	op = op.SetMaxPoolSize(2)
	op = op.SetHeartbeatInterval(time.Second * 5)
	op = op.SetReadPreference(readpref.SecondaryPreferred())
	op = op.SetReadConcern(readconcern.Majority())
	if db, e = mongo.Connect(context.Background(), op); e != nil {
		return
	}
	if e = db.Ping(context.Background(), nil); e != nil {
		return nil, e
	}
	//init self mongo
	var s mongo.Session
	if s, e = db.StartSession(options.Session().SetDefaultReadPreference(readpref.Primary()).SetDefaultReadConcern(readconcern.Local())); e != nil {
		return
	}
	defer s.EndSession(context.Background())
	sctx := mongo.NewSessionContext(context.Background(), s)
	if e = s.StartTransaction(); e != nil {
		return nil, e
	}
	defer func() {
		if e != nil {
			s.AbortTransaction(sctx)
			if mongo.IsDuplicateKeyError(e) {
				e = nil
			}
		} else if e = s.CommitTransaction(sctx); e != nil {
			s.AbortTransaction(sctx)
		}
	}()
	col := db.Database("config_" + groupname).Collection("config")
	appconfig := defaultAppConfig
	sourceconfig := fmt.Sprintf(defaultSourceConfig, url)
	if _, e = col.InsertOne(sctx, bson.M{
		"index":             0,
		"cipher":            "",
		"cur_index":         1,
		"max_index":         1,
		"cur_version":       1,
		"cur_app_config":    appconfig,
		"cur_source_config": sourceconfig,
	}); e != nil {
		return
	}
	index := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "index", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, e = col.Indexes().CreateOne(sctx, index); e != nil {
		return
	}
	if _, e = col.UpdateOne(sctx, bson.M{"index": 1}, bson.M{"$set": bson.M{"app_config": appconfig, "source_config": sourceconfig}}, options.Update().SetUpsert(true)); e != nil {
		return
	}
	return db, nil
}

func updateApp(appconfig []byte) error {
	if len(appconfig) == 0 {
		appconfig = []byte{'{', '}'}
	}
	if len(appconfig) < 2 || appconfig[0] != '{' || appconfig[len(appconfig)-1] != '}' || !json.Valid(appconfig) {
		return errors.New("[config.selfsdk.updateApp] data format error")
	}
	appfile, e := os.OpenFile("./AppConfig_tmp.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		log.Error(nil, "[config.selfsdk.updateApp] open tmp file error:", e)
		return e
	}
	_, e = appfile.Write(appconfig)
	if e != nil {
		log.Error(nil, "[config.selfsdk.updateApp] write temp file error:", e)
		return e
	}
	if e = appfile.Sync(); e != nil {
		log.Error(nil, "[config.selfsdk.updateApp] sync tmp file to disk error:", e)
		return e
	}
	if e = appfile.Close(); e != nil {
		log.Error(nil, "[config.selfsdk.updateApp] close tmp file error:", e)
		return e
	}
	if e = os.Rename("./AppConfig_tmp.json", "./AppConfig.json"); e != nil {
		log.Error(nil, "[config.selfsdk.updateApp] rename error:", e)
		return e
	}
	return nil
}
func updateSource(sourceconfig []byte) error {
	if len(sourceconfig) == 0 {
		sourceconfig = []byte{'{', '}'}
	}
	if len(sourceconfig) < 2 || sourceconfig[0] != '{' || sourceconfig[len(sourceconfig)-1] != '}' || !json.Valid(sourceconfig) {
		return errors.New("[config.selfsdk.updateSource] data format error")
	}
	sourcefile, e := os.OpenFile("./SourceConfig_tmp.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		log.Error(nil, "[config.selfsdk.updateSource] open tmp file error:", e)
		return e
	}
	_, e = sourcefile.Write(sourceconfig)
	if e != nil {
		log.Error(nil, "[config.selfsdk.updateSource] write tmp file error:", e)
		return e
	}
	if e = sourcefile.Sync(); e != nil {
		log.Error(nil, "[config.selfsdk.updateSource] sync tmp file to disk error:", e)
		return e
	}
	if e = sourcefile.Close(); e != nil {
		log.Error(nil, "[config.selfsdk.updateSource] close tmp file error:", e)
		return e
	}
	if e = os.Rename("./SourceConfig_tmp.json", "./SourceConfig.json"); e != nil {
		log.Error(nil, "[config.selfsdk.updateSource] rename error:", e)
		return e
	}
	return nil
}
