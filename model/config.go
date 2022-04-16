package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/chenjie199234/Corelib/util/common"
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
	Index        uint32 `bson:"index"`         //this is always >0  for Config
	AppConfig    string `bson:"app_config"`    //if Cipher is not empty,this field is encrypt
	SourceConfig string `bson:"source_config"` //if Cipher is not empty,this field is encrypt
}

//this is the server's secret
//must be aes.BlockSize length
const iv = "chenjie_1992_3_4"

func pkcs7Padding(origin []byte, blockSize int) []byte {
	padding := blockSize - len(origin)%blockSize
	if padding == 0 {
		return origin
	}
	return append(origin, bytes.Repeat([]byte{byte(padding)}, padding)...)
}
func pkcs7UnPadding(origin []byte) []byte {
	length := len(origin)
	unpadding := int(origin[length-1])
	if unpadding >= aes.BlockSize {
		return nil
	}
	return origin[:(length - unpadding)]
}
func Decrypt(cipherkey, origin string) string {
	data, e := hex.DecodeString(origin)
	if e != nil {
		return ""
	}
	if len(data)%aes.BlockSize != 0 {
		return ""
	}
	block, _ := aes.NewCipher(common.Str2byte(cipherkey))
	cipher.NewCBCDecrypter(block, common.Str2byte(iv)).CryptBlocks(data, data)
	data = pkcs7UnPadding(data)
	return common.Byte2str(data)
}
func Encrypt(cipherkey, origin string) string {
	data := pkcs7Padding([]byte(origin), aes.BlockSize)
	block, _ := aes.NewCipher(common.Str2byte(cipherkey))
	cipher.NewCBCEncrypter(block, common.Str2byte(iv)).CryptBlocks(data, data)
	return hex.EncodeToString(data)
}
