package config

import (
	"context"
	"database/sql"
	"encoding/json"
	"os"
	"time"

	"github.com/chenjie199234/Corelib/log"
	"github.com/chenjie199234/Corelib/redis"
	ctime "github.com/chenjie199234/Corelib/util/time"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//sourceConfig can't hot update
type sourceConfig struct {
	CGrpcServer *CGrpcServerConfig      `json:"cgrpc_server"`
	CGrpcClient *CGrpcClientConfig      `json:"cgrpc_client"`
	CrpcServer  *CrpcServerConfig       `json:"crpc_server"`
	CrpcClient  *CrpcClientConfig       `json:"crpc_client"`
	WebServer   *WebServerConfig        `json:"web_server"`
	WebClient   *WebClientConfig        `json:"web_client"`
	Mongo       map[string]*MongoConfig `json:"mongo"`     //key example:xxx_mongo
	Sql         map[string]*SqlConfig   `json:"sql"`       //key example:xx_sql
	Redis       map[string]*RedisConfig `json:"redis"`     //key example:xx_redis
	KafkaPub    []*KafkaPubConfig       `json:"kafka_pub"`
	KafkaSub    []*KafkaSubConfig       `json:"kafka_sub"`
}

//CGrpcServerConfig
type CGrpcServerConfig struct {
	ConnectTimeout ctime.Duration `json:"connect_timeout"` //default 500ms,max time to finish the handshake
	GlobalTimeout  ctime.Duration `json:"global_timeout"`  //default 500ms,max time to handle the request,unless the specific handle timeout is used in HandlerTimeout in AppConfig,handler's timeout will also be effected by caller's deadline
	HeartProbe     ctime.Duration `json:"heart_probe"`     //default 1.5s
}

//CGrpcClientConfig
type CGrpcClientConfig struct {
	ConnectTimeout ctime.Duration `json:"conn_timeout"`   //default 500ms,max time to finish the handshake
	GlobalTimeout  ctime.Duration `json:"global_timeout"` //max time to handle the request,0 means no default timeout
	HeartProbe     ctime.Duration `json:"heart_probe"`    //default 1.5s
}

//CrpcServerConfig -
type CrpcServerConfig struct {
	ConnectTimeout ctime.Duration `json:"connect_timeout"` //default 500ms,max time to finish the handshake
	GlobalTimeout  ctime.Duration `json:"global_timeout"`  //default 500ms,max time to handle the request,unless the specific handle timeout is used in HandlerTimeout in AppConfig,handler's timeout will also be effected by caller's deadline
	HeartProbe     ctime.Duration `json:"heart_probe"`     //default 1.5s
}

//CrpcClientConfig -
type CrpcClientConfig struct {
	ConnectTimeout ctime.Duration `json:"conn_timeout"`   //default 500ms,max time to finish the handshake
	GlobalTimeout  ctime.Duration `json:"global_timeout"` //max time to handle the request,0 means no default timeout
	HeartProbe     ctime.Duration `json:"heart_probe"`    //default 1.5s
}

//WebServerConfig -
type WebServerConfig struct {
	ConnectTimeout ctime.Duration `json:"connect_timeout"` //default 500ms,max time to finish the handshake and read each whole request
	GlobalTimeout  ctime.Duration `json:"global_timeout"`  //default 500ms,max time to handle the request,unless the specific handle timeout is used in HandlerTimeout in AppConfig,handler's timeout will also be effected by caller's deadline
	IdleTimeout    ctime.Duration `json:"idle_timeout"`    //default 5s
	HeartProbe     ctime.Duration `json:"heart_probe"`     //default 1.5s
	StaticFilePath string         `json:"static_file_path"`
	//cors
	Cors *WebCorsConfig `json:"cors"`
}

//WebCorsConfig -
type WebCorsConfig struct {
	CorsOrigin []string `json:"cors_origin"`
	CorsHeader []string `json:"cors_header"`
	CorsExpose []string `json:"cors_expose"`
}

//WebClientConfig -
type WebClientConfig struct {
	ConnectTimeout ctime.Duration `json:"conn_timeout"`   //default 500ms,max time to finish the handshake
	GlobalTimeout  ctime.Duration `json:"global_timeout"` //max time to handle the request,0 means no default timeout
	IdleTimeout    ctime.Duration `json:"idle_timeout"`   //default 5s
	HeartProbe     ctime.Duration `json:"heart_probe"`    //default 1.5s
}

//RedisConfig -
type RedisConfig struct {
	URL         string         `json:"url"`           //[redis/rediss]://[[username:]password@]host/[dbindex]
	MaxOpen     int            `json:"max_open"`      //default 100   //this will overwrite the param in url
	MaxIdletime ctime.Duration `json:"max_idletime"`  //default 10min //this will overwrite the param in url
	IOTimeout   ctime.Duration `json:"io_timeout"`    //default 500ms //this will overwrite the param in url
	ConnTimeout ctime.Duration `json:"conn_timeout"`  //default 250ms //this will overwrite the param in url
}

//SqlConfig -
type SqlConfig struct {
	URL         string         `json:"url"`           //[username:password@][protocol(address)]/[dbname][?param1=value1&...&paramN=valueN]
	MaxOpen     int            `json:"max_open"`      //default 100   //this will overwrite the param in url
	MaxIdletime ctime.Duration `json:"max_idletime"`  //default 10min //this will overwrite the param in url
	IOTimeout   ctime.Duration `json:"io_timeout"`    //default 500ms //this will overwrite the param in url
	ConnTimeout ctime.Duration `json:"conn_timeout"`  //default 250ms //this will overwrite the param in url
}

//MongoConfig -
type MongoConfig struct {
	URL         string         `json:"url"`           //[mongodb/mongodb+srv]://[username:password@]host1,...,hostN/[dbname][?param1=value1&...&paramN=valueN]
	MaxOpen     uint64         `json:"max_open"`      //default 100   //this will overwrite the param in url
	MaxIdletime ctime.Duration `json:"max_idletime"`  //default 10min //this will overwrite the param in url
	IOTimeout   ctime.Duration `json:"io_timeout"`    //default 500ms //this will overwrite the param in url
	ConnTimeout ctime.Duration `json:"conn_timeout"`  //default 250ms //this will overwrite the param in url
}

//KafkaPubConfig -
type KafkaPubConfig struct {
	Addrs          []string       `json:"addrs"`
	Username       string         `json:"username"`
	Passwd         string         `json:"password"`
	AuthMethod     int            `json:"auth_method"`     //1-plain,2-scram sha256,3-scram sha512
	CompressMethod int            `json:"compress_method"` //0-none,1-gzip,2-snappy,3-lz4,4-zstd
	TopicName      string         `json:"topic_name"`
	IOTimeout      ctime.Duration `json:"io_timeout"`   //default 500ms
	ConnTimeout    ctime.Duration `json:"conn_timeout"` //default 250ms
}

//KafkaSubConfig -
type KafkaSubConfig struct {
	Addrs       []string       `json:"addrs"`
	Username    string         `json:"username"`
	Passwd      string         `json:"password"`
	AuthMethod  int            `json:"auth_method"` //1-plain,2-scram sha256,3-scram sha512
	TopicName   string         `json:"topic_name"`
	GroupName   string         `json:"group_name"`
	ConnTimeout ctime.Duration `json:"conn_timeout"` //default 250ms
	//when there is no offset in a partition(add partition or first time to use the topic)
	//-1 will sub from the newest
	//-2 will sub from the firt
	//if this is 0,default -2 will be used
	StartOffset int64 `json:"start_offset"`
	//if this is 0,commit is synced,and effective is slow.
	//if this is not 0,commit is asynced,effective is high,but will cause duplicate sub when the program crash
	CommitInterval ctime.Duration `json:"commit_interval"`
}

//SC total source config instance
var sc *sourceConfig

var mongos map[string]*mongo.Client

var sqls map[string]*sql.DB

var rediss map[string]*redis.Pool

var kafkaSubers map[string]*kafka.Reader

var kafkaPubers map[string]*kafka.Writer

func initsource() {
	data, e := os.ReadFile("./SourceConfig.json")
	if e != nil {
		log.Error(nil, "[config.initsource] read config file error:", e)
		Close()
		os.Exit(1)
	}
	sc = &sourceConfig{}
	if e = json.Unmarshal(data, sc); e != nil {
		log.Error(nil, "[config.initsource] config file format error:", e)
		Close()
		os.Exit(1)
	}

	initgrpcserver()
	initgrpcclient()
	initcrpcserver()
	initcrpcclient()
	initwebserver()
	initwebclient()
	initredis()
	initmongo()
	initsql()
	initkafkapub()
	initkafkasub()
}
func initgrpcserver() {
	if sc.CGrpcServer == nil {
		sc.CGrpcServer = &CGrpcServerConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			HeartProbe:     ctime.Duration(1500 * time.Millisecond),
		}
	} else {
		if sc.CGrpcServer.ConnectTimeout <= 0 {
			sc.CGrpcServer.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CGrpcServer.GlobalTimeout <= 0 {
			sc.CGrpcServer.GlobalTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CGrpcServer.HeartProbe <= 0 {
			sc.CGrpcServer.HeartProbe = ctime.Duration(1500 * time.Millisecond)
		}
	}
}
func initgrpcclient() {
	if sc.CGrpcClient == nil {
		sc.CGrpcClient = &CGrpcClientConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			HeartProbe:     ctime.Duration(time.Millisecond * 1500),
		}
	} else {
		if sc.CGrpcClient.ConnectTimeout <= 0 {
			sc.CGrpcClient.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CGrpcClient.GlobalTimeout < 0 {
			sc.CGrpcClient.GlobalTimeout = 0
		}
		if sc.CGrpcClient.HeartProbe <= 0 {
			sc.CGrpcClient.HeartProbe = ctime.Duration(time.Millisecond * 1500)
		}
	}
}
func initcrpcserver() {
	if sc.CrpcServer == nil {
		sc.CrpcServer = &CrpcServerConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			HeartProbe:     ctime.Duration(1500 * time.Millisecond),
		}
	} else {
		if sc.CrpcServer.ConnectTimeout <= 0 {
			sc.CrpcServer.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CrpcServer.GlobalTimeout <= 0 {
			sc.CrpcServer.GlobalTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CrpcServer.HeartProbe <= 0 {
			sc.CrpcServer.HeartProbe = ctime.Duration(1500 * time.Millisecond)
		}
	}
}
func initcrpcclient() {
	if sc.CrpcClient == nil {
		sc.CrpcClient = &CrpcClientConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			HeartProbe:     ctime.Duration(time.Millisecond * 1500),
		}
	} else {
		if sc.CrpcClient.ConnectTimeout <= 0 {
			sc.CrpcClient.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.CrpcClient.GlobalTimeout < 0 {
			sc.CrpcClient.GlobalTimeout = 0
		}
		if sc.CrpcClient.HeartProbe <= 0 {
			sc.CrpcClient.HeartProbe = ctime.Duration(time.Millisecond * 1500)
		}
	}

}
func initwebserver() {
	if sc.WebServer == nil {
		sc.WebServer = &WebServerConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			IdleTimeout:    ctime.Duration(time.Second * 5),
			HeartProbe:     ctime.Duration(time.Millisecond * 1500),
			StaticFilePath: "./src",
			Cors: &WebCorsConfig{
				CorsOrigin: []string{"*"},
				CorsHeader: []string{"*"},
				CorsExpose: nil,
			},
		}
	} else {
		if sc.WebServer.ConnectTimeout <= 0 {
			sc.WebServer.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.WebServer.GlobalTimeout <= 0 {
			sc.WebServer.GlobalTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.WebServer.IdleTimeout <= 0 {
			sc.WebServer.IdleTimeout = ctime.Duration(time.Second * 5)
		}
		if sc.WebServer.HeartProbe <= 0 {
			sc.WebServer.HeartProbe = ctime.Duration(time.Millisecond * 1500)
		}
		if sc.WebServer.Cors == nil {
			sc.WebServer.Cors = &WebCorsConfig{
				CorsOrigin: []string{"*"},
				CorsHeader: []string{"*"},
				CorsExpose: nil,
			}
		}
	}
}
func initwebclient() {
	if sc.WebClient == nil {
		sc.WebClient = &WebClientConfig{
			ConnectTimeout: ctime.Duration(time.Millisecond * 500),
			GlobalTimeout:  ctime.Duration(time.Millisecond * 500),
			IdleTimeout:    ctime.Duration(time.Second * 5),
			HeartProbe:     ctime.Duration(time.Millisecond * 1500),
		}
	} else {
		if sc.WebClient.ConnectTimeout <= 0 {
			sc.WebClient.ConnectTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sc.WebClient.GlobalTimeout < 0 {
			sc.WebClient.GlobalTimeout = 0
		}
		if sc.WebClient.IdleTimeout <= 0 {
			sc.WebClient.IdleTimeout = ctime.Duration(time.Second * 5)
		}
		if sc.WebClient.HeartProbe <= 0 {
			sc.WebClient.HeartProbe = ctime.Duration(time.Millisecond * 1500)
		}
	}
}
func initredis(){
	for k, redisc := range sc.Redis {
		if k == "example_redis" {
			continue
		}
		if redisc.MaxOpen == 0 {
			redisc.MaxOpen = 100
		}
		if redisc.MaxIdletime == 0 {
			redisc.MaxIdletime = ctime.Duration(time.Minute * 10)
		}
		if redisc.IOTimeout == 0 {
			redisc.IOTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if redisc.ConnTimeout == 0 {
			redisc.ConnTimeout = ctime.Duration(time.Millisecond * 250)
		}
	}
	rediss = make(map[string]*redis.Pool, len(sc.Redis))
	for k, redisc := range sc.Redis {
		if k == "example_redis" {
			continue
		}
		tempredis := redis.NewRedis(&redis.Config{
			RedisName:   k,
			URL:         redisc.URL,
			MaxOpen:     redisc.MaxOpen,
			MaxIdletime: redisc.MaxIdletime.StdDuration(),
			ConnTimeout: redisc.ConnTimeout.StdDuration(),
			IOTimeout:   redisc.IOTimeout.StdDuration(),
		})
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if e := tempredis.Ping(ctx); e != nil {
			cancel()
			log.Error(nil, "[config.initsource] ping redis:", k, "error:", e)
			Close()
			os.Exit(1)
		}
		cancel()
		rediss[k] = tempredis
	}
}
func initmongo(){
	for k, mongoc := range sc.Mongo {
		if k == "example_mongo" {
			continue
		}
		if mongoc.MaxOpen == 0 {
			mongoc.MaxOpen = 100
		}
		if mongoc.MaxIdletime == 0 {
			mongoc.MaxIdletime = ctime.Duration(time.Minute * 10)
		}
		if mongoc.IOTimeout == 0 {
			mongoc.IOTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if mongoc.ConnTimeout == 0 {
			mongoc.ConnTimeout = ctime.Duration(time.Millisecond * 250)
		}
	}
	mongos = make(map[string]*mongo.Client, len(sc.Mongo))
	for k, mongoc := range sc.Mongo {
		if k == "example_mongo" {
			continue
		}
		op := options.Client().ApplyURI(mongoc.URL)
		op = op.SetConnectTimeout(mongoc.ConnTimeout.StdDuration())
		op = op.SetMaxConnIdleTime(mongoc.MaxIdletime.StdDuration())
		op = op.SetMaxPoolSize(mongoc.MaxOpen)
		op = op.SetSocketTimeout(mongoc.IOTimeout.StdDuration())
		tempdb, e := mongo.Connect(nil, op)
		if e != nil {
			log.Error(nil, "[config.initsource] open mongodb:", k, "error:", e)
			Close()
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		e = tempdb.Ping(ctx, readpref.Primary())
		if e != nil {
			cancel()
			log.Error(nil, "[config.initsource] ping mongodb:", k, "error:", e)
			Close()
			os.Exit(1)
		}
		cancel()
		mongos[k] = tempdb
	}
}
func initsql(){
	for _, sqlc := range sc.Sql {
		if sqlc.MaxOpen == 0 {
			sqlc.MaxOpen = 100
		}
		if sqlc.MaxIdletime == 0 {
			sqlc.MaxIdletime = ctime.Duration(time.Minute * 10)
		}
		if sqlc.IOTimeout == 0 {
			sqlc.IOTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if sqlc.ConnTimeout == 0 {
			sqlc.ConnTimeout = ctime.Duration(time.Millisecond * 250)
		}
	}
	sqls = make(map[string]*sql.DB, len(sc.Sql))
	for k, sqlc := range sc.Sql {
		if k == "example_sql" {
			continue
		}
		tempdb, e := sql.Open("mysql", sqlc.URL)
		if e != nil {
			log.Error(nil, "[config.initsource] open mysql:", k, "error:", e)
			Close()
			os.Exit(1)
		}
		tempdb.SetMaxOpenConns(sqlc.MaxOpen)
		tempdb.SetMaxIdleConns(sqlc.MaxOpen)
		tempdb.SetConnMaxIdleTime(sqlc.MaxIdletime.StdDuration())
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		e = tempdb.PingContext(ctx)
		if e != nil {
			cancel()
			log.Error(nil, "[config.initsource] ping mysql:", k, "error:", e)
			Close()
			os.Exit(1)
		}
		cancel()
		sqls[k] = tempdb
	}
}
func initkafkapub(){
	for _, pubc := range sc.KafkaPub {
		if pubc.TopicName == "example_topic" || pubc.TopicName == "" {
			continue
		}
		if len(pubc.Addrs) == 0 {
			pubc.Addrs = []string{"127.0.0.1:9092"}
		}
		if (pubc.AuthMethod == 1 || pubc.AuthMethod == 2 || pubc.AuthMethod == 3) && (pubc.Username == "" || pubc.Passwd == "") {
			log.Error(nil, "[config.initsource] pub topic:", pubc.TopicName, "username or password missing")
			Close()
			os.Exit(1)
		}
		if pubc.IOTimeout == 0 {
			pubc.IOTimeout = ctime.Duration(time.Millisecond * 500)
		}
		if pubc.ConnTimeout == 0 {
			pubc.IOTimeout = ctime.Duration(time.Millisecond * 250)
		}
	}
	kafkaPubers = make(map[string]*kafka.Writer, len(sc.KafkaPub))
	for _, pubc := range sc.KafkaPub {
		if pubc.TopicName == "example_topic" || pubc.TopicName == "" {
			continue
		}
		dialer := &kafka.Dialer{
			Timeout:   pubc.ConnTimeout.StdDuration(),
			DualStack: true,
		}
		var e error
		switch pubc.AuthMethod {
		case 1:
			dialer.SASLMechanism = plain.Mechanism{Username: pubc.Username, Password: pubc.Passwd}
		case 2:
			dialer.SASLMechanism, e = scram.Mechanism(scram.SHA256, pubc.Username, pubc.Passwd)
		case 3:
			dialer.SASLMechanism, e = scram.Mechanism(scram.SHA512, pubc.Username, pubc.Passwd)
		}
		if e != nil {
			log.Error(nil, "[config.initsource] kafka topic:", pubc.TopicName, "pub username and password parse error:", e)
			Close()
			os.Exit(1)
		}
		var compressor kafka.CompressionCodec
		switch pubc.CompressMethod {
		case 1:
			compressor = kafka.Gzip.Codec()
		case 2:
			compressor = kafka.Snappy.Codec()
		case 3:
			compressor = kafka.Lz4.Codec()
		case 4:
			compressor = kafka.Zstd.Codec()
		}
		writer := kafka.NewWriter(kafka.WriterConfig{
			Brokers:          pubc.Addrs,
			Topic:            pubc.TopicName,
			Dialer:           dialer,
			ReadTimeout:      pubc.IOTimeout.StdDuration(),
			WriteTimeout:     pubc.IOTimeout.StdDuration(),
			Balancer:         &kafka.Hash{},
			MaxAttempts:      3,
			RequiredAcks:     int(kafka.RequireAll),
			Async:            false,
			CompressionCodec: compressor,
		})
		kafkaPubers[pubc.TopicName] = writer
	}
}
func initkafkasub(){
	for _, subc := range sc.KafkaSub {
		if subc.TopicName == "example_topic" || subc.TopicName == "" {
			continue
		}
		if len(subc.Addrs) == 0 {
			subc.Addrs = []string{"127.0.0.1:9092"}
		}
		if (subc.AuthMethod == 1 || subc.AuthMethod == 2 || subc.AuthMethod == 3) && (subc.Username == "" || subc.Passwd == "") {
			log.Error(nil, "[config.initsource] sub topic:", subc.TopicName, "username or password missing")
			Close()
			os.Exit(1)
		}
		if subc.GroupName == "" {
			log.Error(nil, "[config.initsource] sub topic:", subc.TopicName, "groupname missing")
			Close()
			os.Exit(1)
		}
		if subc.ConnTimeout == 0 {
			subc.ConnTimeout = ctime.Duration(time.Millisecond * 250)
		}
	}
	kafkaSubers = make(map[string]*kafka.Reader, len(sc.KafkaSub))
	for _, subc := range sc.KafkaSub {
		if subc.TopicName == "example_topic" || subc.TopicName == "" {
			continue
		}
		dialer := &kafka.Dialer{
			Timeout:   subc.ConnTimeout.StdDuration(),
			DualStack: true,
		}
		var e error
		switch subc.AuthMethod {
		case 1:
			dialer.SASLMechanism = plain.Mechanism{Username: subc.Username, Password: subc.Passwd}
		case 2:
			dialer.SASLMechanism, e = scram.Mechanism(scram.SHA256, subc.Username, subc.Passwd)
		case 3:
			dialer.SASLMechanism, e = scram.Mechanism(scram.SHA512, subc.Username, subc.Passwd)
		}
		if e != nil {
			log.Error(nil, "[config.initsource] kafka topic:", subc.TopicName, "sub username and password parse error:", e)
			Close()
			os.Exit(1)
		}
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:                subc.Addrs,
			Dialer:                 dialer,
			Topic:                  subc.TopicName,
			GroupID:                subc.GroupName,
			StartOffset:            subc.StartOffset,
			MinBytes:               1,
			MaxBytes:               1024 * 1024 * 10,
			CommitInterval:         time.Duration(subc.CommitInterval),
			IsolationLevel:         kafka.ReadCommitted,
			PartitionWatchInterval: time.Second,
			WatchPartitionChanges:  true,
			MaxAttempts:            3,
		})
		kafkaSubers[subc.TopicName+subc.GroupName] = reader
	}
}

//GetCGrpcServerConfig get the grpc net config
func GetCGrpcServerConfig() *CGrpcServerConfig {
	return sc.CGrpcServer
}

//GetCGrpcClientConfig get the grpc net config
func GetCGrpcClientConfig() *CGrpcClientConfig {
	return sc.CGrpcClient
}

//GetCrpcServerConfig get the crpc net config
func GetCrpcServerConfig() *CrpcServerConfig {
	return sc.CrpcServer
}

//GetCrpcClientConfig get the crpc net config
func GetCrpcClientConfig() *CrpcClientConfig {
	return sc.CrpcClient
}

//GetWebServerConfig get the web net config
func GetWebServerConfig() *WebServerConfig {
	return sc.WebServer
}

//GetWebClientConfig get the web net config
func GetWebClientConfig() *WebClientConfig {
	return sc.WebClient
}

//GetMongo get a mongodb client by db's instance name
//return nil means not exist
func GetMongo(mongoname string) *mongo.Client {
	return mongos[mongoname]
}

//GetSql get a mysql db client by db's instance name
//return nil means not exist
func GetSql(mysqlname string) *sql.DB {
	return sqls[mysqlname]
}

//GetRedis get a redis client by redis's instance name
//return nil means not exist
func GetRedis(redisname string) *redis.Pool {
	return rediss[redisname]
}

//GetKafkaSuber get a kafka sub client by topic and groupid
func GetKafkaSuber(topic string, groupid string) *kafka.Reader {
	return kafkaSubers[topic+groupid]
}

//GetKafkaPuber get a kafka pub client by topic name
func GetKafkaPuber(topic string) *kafka.Writer {
	return kafkaPubers[topic]
}