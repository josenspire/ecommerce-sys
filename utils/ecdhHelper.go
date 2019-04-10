package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"github.com/astaxie/beego"
	"hash"
	"io"
	"math/big"
)

type EllipticECDH struct {
	PrivateKey []byte
	PublicKey  EllipticPublicKey
}

type EllipticPublicKey struct {
	X *big.Int
	Y *big.Int
}

var curve = elliptic.P256()

func (e *EllipticECDH) GenerateKeyPair() (crypto.PrivateKey, crypto.PublicKey, error) {
	privateKey, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, err
	}
	publicKey := &EllipticPublicKey{
		X: x,
		Y: y,
	}
	return privateKey, publicKey, nil
}

// Marshal用于公钥的序列化
func (e *EllipticECDH) Marshal(pub crypto.PublicKey) []byte {
	publicKey := pub.(*EllipticPublicKey)
	return elliptic.Marshal(curve, publicKey.X, publicKey.Y)
}

// Unmarshal用于公钥的反序列化
func (e *EllipticECDH) Unmarshal(data []byte) (crypto.PublicKey, bool) {
	x, y := elliptic.Unmarshal(curve, data)
	if x == nil || y == nil {
		return nil, false
	}
	publicKey := &EllipticPublicKey{
		X: x,
		Y: y,
	}
	return publicKey, true
}

// ComputeSecret
func (e *EllipticECDH) ComputeSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	priv := privKey.([]byte)
	pub := pubKey.(*EllipticPublicKey)
	x, _ := curve.ScalarMult(pub.X, pub.Y, priv)
	return x.Bytes(), nil
}

// https://blog.yumaojun.net/2017/02/19/go-crypto/

// SignData 用于保存签名的数据
type SignData struct {
	r         *big.Int
	s         *big.Int
	signhash  *[]byte
	signature *[]byte
}

func (e *EllipticECDH) Signature(message string, privateKey *ecdsa.PrivateKey) (signData *SignData, err error) {
	// 签名数据
	var h hash.Hash
	h = md5.New()
	r := big.NewInt(0)
	s := big.NewInt(0)
	io.WriteString(h, message)
	signhash := h.Sum(nil)
	r, s, serr := ecdsa.Sign(rand.Reader, privateKey, signhash)
	if serr != nil {
		return nil, serr
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	signData = &SignData{
		r:         r,
		s:         s,
		signhash:  &signhash,
		signature: &signature,
	}
	return
}

// 校验数字签名
func verifySign(signData *SignData, publicKey *ecdsa.PublicKey) (status bool) {
	status = ecdsa.Verify(publicKey, *signData.signhash, signData.r, signData.s)
	return
}
