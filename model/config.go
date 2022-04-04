package model

/*
| groupname1(database)
|      appname1(collection)
|      appname2
|      appname3
| groupname2(database)
|      appnameN(collection)
*/
//every collection has two kinds of data,they are seprated by index
type Summary struct {
	CurIndex   uint32 `bson:"cur_index"`
	MaxIndex   uint32 `bson:"max_index"`
	CurVersion uint32 `bson:"cur_version"`
	Index      uint32 `bson:"index"` //this is always 0 for summary
}
type Config struct {
	Index        uint32 `bson:"index"` //this is always >0  for Config
	AppConfig    string `bson:"app_config"`
	SourceConfig string `bson:"source_config"`
}
