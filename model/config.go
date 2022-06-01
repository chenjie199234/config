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
//every collection has two kinds of data
type Summary struct {
	ID              primitive.ObjectID `bson:"_id"`
	Cipher          string             `bson:"cipher"`
	CurIndex        uint32             `bson:"cur_index"`
	MaxIndex        uint32             `bson:"max_index"`
	CurVersion      uint32             `bson:"cur_version"`
	Index           uint32             `bson:"index"`             //this is always 0 for summary
	CurAppConfig    string             `bson:"cur_app_config"`    //if Cipher is not empty,this field is encrypt
	CurSourceConfig string             `bson:"cur_source_config"` //if Cipher is not empty,this field is encrypt
}
type Config struct {
	Index        uint32 `bson:"index"`         //this is always > 0  for Config
	AppConfig    string `bson:"app_config"`    //if Cipher is not empty,this field is encrypt
	SourceConfig string `bson:"source_config"` //if Cipher is not empty,this field is encrypt
}
