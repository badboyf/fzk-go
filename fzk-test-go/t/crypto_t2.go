package main
//
//import (
//	"bytes"
//	"crypto/aes"
//	"crypto/cipher"
//	"encoding/base64"
//	"fmt"
//	"strings"
//)
//
//func main() {
//	/*
//	 *src 要加密的字符串
//	 *key 用来加密的密钥 密钥长度可以是128bit、192bit、256bit中的任意一个
//	 *16位key对应128bit
//	 */
//	src := "hello world"
//	key := "0123456789abcdef"
//
//	crypted := AesEncrypt(src, key)
//	var dec = AesDecrypt(crypted, []byte(key))
//
//	fmt.Println("加密:", base64.StdEncoding.EncodeToString(crypted))
//	fmt.Println("解密:", string(dec))
//	Base64URLDecode("39W7dWTd_SBOCM8UbnG6qA")
//}
//
//func Base64URLDecode(data string) ([]byte, error) {
//	var missing = (4 - len(data)%4) % 4
//	data += strings.Repeat("=", missing)
//	res, err := base64.URLEncoding.DecodeString(data)
//	fmt.Println("  decodebase64urlsafe is :", string(res), err)
//	return base64.URLEncoding.DecodeString(data)
//}
//
//func AesDecrypt(crypted, key []byte) []byte {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		fmt.Println("err is:", err)
//	}
//	blockMode := NewECBDecrypter(block)
//	origData := make([]byte, len(crypted))
//	blockMode.CryptBlocks(origData, crypted)
//	origData = PKCS5UnPadding(origData)
//	return origData
//}
//
//func AesEncrypt(src, key string) []byte {
//	block, err := aes.NewCipher([]byte(key))
//	if err != nil {
//		fmt.Println("key error1", err)
//	}
//	if src == "" {
//		fmt.Println("plain content empty")
//	}
//	ecb := NewECBEncrypter(block)
//	content := []byte(src)
//	content = PKCS5Padding(content, block.BlockSize())
//	crypted := make([]byte, len(content))
//	ecb.CryptBlocks(crypted, content)
//	return crypted
//}
//
//func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
//	padding := blockSize - len(ciphertext)%blockSize
//	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
//	return append(ciphertext, padtext...)
//}
//
//func PKCS5UnPadding(origData []byte) []byte {
//	length := len(origData)
//	// 去掉最后一个字节 unpadding 次
//	unpadding := int(origData[length-1])
//	return origData[:(length - unpadding)]
//}
//
//type ecb struct {
//	b         cipher.Block
//	blockSize int
//}
//
//func newECB(b cipher.Block) *ecb {
//	return &ecb{
//		b:         b,
//		blockSize: b.BlockSize(),
//	}
//}
//
//type ecbEncrypter ecb
//
//// NewECBEncrypter returns a BlockMode which encrypts in electronic code book
//// mode, using the given Block.
//func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
//	return (*ecbEncrypter)(newECB(b))
//}
//func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
//func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
//	if len(src)%x.blockSize != 0 {
//		panic("crypto/cipher: input not full blocks")
//	}
//	if len(dst) < len(src) {
//		panic("crypto/cipher: output smaller than input")
//	}
//	for len(src) > 0 {
//		x.b.Encrypt(dst, src[:x.blockSize])
//		src = src[x.blockSize:]
//		dst = dst[x.blockSize:]
//	}
//}
//
//type ecbDecrypter ecb
//
//// NewECBDecrypter returns a BlockMode which decrypts in electronic code book
//// mode, using the given Block.
//func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
//	return (*ecbDecrypter)(newECB(b))
//}
//func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
//func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
//	if len(src)%x.blockSize != 0 {
//		panic("crypto/cipher: input not full blocks")
//	}
//	if len(dst) < len(src) {
//		panic("crypto/cipher: output smaller than input")
//	}
//	for len(src) > 0 {
//		x.b.Decrypt(dst, src[:x.blockSize])
//		src = src[x.blockSize:]
//		dst = dst[x.blockSize:]
//	}
//}