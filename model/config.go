package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
| config_groupname1(database)
|      appname1(collection)
|      appname2
|      appname3
| config_groupname2(database)
|      appnameN(collection)
*/
//every collection has two kinds of data,they are seprated by index
type Summary struct {
	ID         primitive.ObjectID `bson:"_id"`
	CurIndex   uint32             `bson:"cur_index"`
	MaxIndex   uint32             `bson:"max_index"`
	CurVersion uint32             `bson:"cur_version"`
	Index      uint32             `bson:"index"` //this is always 0 for summary
}
type Config struct {
	Index        uint32 `bson:"index"` //this is always >0  for Config
	AppConfig    string `bson:"app_config"`
	SourceConfig string `bson:"source_config"`
}
type Current struct {
	ID           string //summary's ID
	GroupName    string
	AppName      string
	CurVersion   uint32
	AppConfig    string
	SourceConfig string
}
