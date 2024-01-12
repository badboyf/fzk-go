package main
//
//import (
//	"bytes"
//	"crypto/aes"
//	"crypto/cipher"
//	"encoding/base64"
//	"fmt"
//)
//
//type AesCryptor struct {
//	key []byte
//	iv  []byte
//}
//
////加密数据
//func (a *AesCryptor) Encrypt(data []byte) ([]byte, error) {
//	aesBlockEncrypter, err := aes.NewCipher(a.key)
//	content := PKCS5Padding(data, aesBlockEncrypter.BlockSize())
//	encrypted := make([]byte, len(content))
//	if err != nil {
//		println(err.Error())
//		return nil, err
//	}
//	aesEncrypter := cipher.NewCBCEncrypter(aesBlockEncrypter, a.iv)
//	aesEncrypter.CryptBlocks(encrypted, content)
//	return encrypted, nil
//}
//
////解密数据
//func (a *AesCryptor) Decrypt(src []byte) (data []byte, err error) {
//	decrypted := make([]byte, len(src))
//	var aesBlockDecrypter cipher.Block
//	aesBlockDecrypter, err = aes.NewCipher(a.key)
//	if err != nil {
//		println(err.Error())
//		return nil, err
//	}
//	aesDecrypter := cipher.NewCBCDecrypter(aesBlockDecrypter, a.iv)
//	aesDecrypter.CryptBlocks(decrypted, src)
//	return PKCS5Trimming(decrypted), nil
//}
//
///**
//PKCS5包装
//*/
//func PKCS5Padding(cipherText []byte, blockSize int) []byte {
//	padding := blockSize - len(cipherText)%blockSize
//	padText := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(cipherText, padText...)
//}
//
///*
//解包装
//*/
//func PKCS5Trimming(encrypt []byte) []byte {
//	padding := encrypt[len(encrypt)-1]
//	return encrypt[:len(encrypt)-int(padding)]
//}
//
//func main() {
//	var a =AesCryptor{
//		key: []byte("1234567890123456"),
//		iv: []byte("1234567890123456"),
//	}
//
//	var value = "hello world"
//	var ret, err = a.Encrypt([]byte(value))
//	var b, err2 = base64.StdEncoding.DecodeString(string(ret))
//	fmt.Printf("%s %s %s", b, err, err2)
//}