package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego/logs"
	"math/big"
)

const (
	iv = "0102030405060708"
)

type Crypto struct {
	SecretKey  string
	OriginData interface{}
}

func (ct *Crypto) CreateSecretKey(size int) string {
	rs := GenerateRandString(16)
	ct.SecretKey = rs
	return rs
}

func (ct *Crypto) RSAEncrypt(secKey string, pubKey string, modulus string) string {
	encSecKey := rsaEncrypt(secKey, pubKey, modulus)
	return encSecKey
}

// AES加密的具体算法为: AES-128-CBC，输出格式为 base64
// AES加密时需要指定 iv：0102030405060708
// AES加密时需要 padding
// either 16, 24, or 32 bytes to select
// AES-128, AES-192, or AES-256.
// https://github.com/darknessomi/musicbox/wiki/%E7%BD%91%E6%98%93%E4%BA%91%E9%9F%B3%E4%B9%90%E6%96%B0%E7%99%BB%E5%BD%95API%E5%88%86%E6%9E%90
func AESEncrypt(encodeStr string, secretKeyStr string) (string, error) {
	secretKey := []byte(secretKeyStr)
	encodeBytes := []byte(encodeStr)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		logs.Error(err)
		return "", err
	}

	blockSize := block.BlockSize()
	encodeBytes = pKCS5Padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(iv))
	cipherText := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(cipherText, encodeBytes)

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func AESDecrypt(decodeStr string, secretKeyStr string) (string, error) {
	// decode base64
	decodeBytes, _ := base64.StdEncoding.DecodeString(decodeStr)

	secretKey := []byte(secretKeyStr)
	block, _ := aes.NewCipher(secretKey)

	blockMode := cipher.NewCBCDecrypter(block, []byte(iv))
	originData := make([]byte, len(decodeBytes))

	blockMode.CryptBlocks(originData, decodeBytes)
	originData = pKCS5UnPadding(originData)
	return string(originData[:]), nil
}

func rsaEncrypt(secKey string, pubKey string, modulus string) string {
	// 倒序 key
	rKey := ""
	for i := len(secKey) - 1; i >= 0; i-- {
		rKey += secKey[i : i+1]
	}
	// 将 key 转 ascii 编码 然后转成 16 进制字符串
	hexRKey := ""
	for _, char := range []rune(rKey) {
		hexRKey += fmt.Sprintf("%x", int(char))
	}
	// 将 16进制 的 三个参数 转为10进制的 bigint
	bigRKey, _ := big.NewInt(0).SetString(hexRKey, 16)
	bigPubKey, _ := big.NewInt(0).SetString(pubKey, 16)
	bigModulus, _ := big.NewInt(0).SetString(modulus, 16)
	// 执行幂乘取模运算得到最终的bigint结果
	bigRs := bigRKey.Exp(bigRKey, bigPubKey, bigModulus)
	// 将结果转为 16进制字符串
	hexRs := fmt.Sprintf("%x", bigRs)
	// 可能存在不满256位的情况，要在前面补0补满256位
	return addRSAPadding(hexRs, modulus)
}

// 补0步骤
func addRSAPadding(encText string, modulus string) string {
	ml := len(modulus)
	for i := 0; ml > 0 && modulus[i:i+1] == "0"; i++ {
		ml--
	}
	num := ml - len(encText)
	prefix := ""
	for i := 0; i < num; i++ {
		prefix += "0"
	}
	return prefix + encText
}

func pKCS5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize // 16, 32, 48 etc..
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

func pKCS5UnPadding(originData []byte) []byte {
	length := len(originData)
	unPadding := int(originData[length-1])
	return originData[:(length - unPadding)]
}
