package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"github.com/chenjie199234/Corelib/util/common"
)

//this is the server's secret
//must be aes.BlockSize length
const servercipher = "chenjie_1992_3_4"

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
	cipher.NewCBCDecrypter(block, common.Str2byte(servercipher)).CryptBlocks(data, data)
	data = pkcs7UnPadding(data)
	return common.Byte2str(data)
}
func Encrypt(cipherkey, origin string) string {
	data := pkcs7Padding([]byte(origin), aes.BlockSize)
	block, _ := aes.NewCipher(common.Str2byte(cipherkey))
	cipher.NewCBCEncrypter(block, common.Str2byte(servercipher)).CryptBlocks(data, data)
	return hex.EncodeToString(data)
}
