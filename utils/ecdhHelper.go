package utils

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/md5"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/astaxie/beego"
	"hash"
	"io"
	"math/big"
)

type ECDH interface {
	GenerateECKeyPair() (crypto.PrivateKey, crypto.PublicKey, error)
	Marshal(pub crypto.PublicKey) []byte
	Unmarshal(data []byte) (crypto.PublicKey, bool)
	ParsePKCS8ECPrivateKey(privateKeyDerBytes []byte) (*EllipticPrivateKey, *EllipticPublicKey, error)
	ParsePKIXECPublicKeyFrom(publicKeyDerBytes []byte) (*EllipticPublicKey, error)
	DecodePEMToDERBytes(pemBytes []byte) []byte
	ComputeSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error)
}

type EllipticECDH struct {
	PrivateKey EllipticPrivateKey
	PublicKey  EllipticPublicKey
}

type EllipticPrivateKey struct {
	D []byte
}

type EllipticPublicKey struct {
	X *big.Int
	Y *big.Int
}

var curve = elliptic.P256()

func (e *EllipticECDH) GenerateECKeyPair() (crypto.PrivateKey, crypto.PublicKey, error) {
	priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, err
	}
	privateKey := &EllipticPrivateKey{
		D: priv,
	}
	publicKey := &EllipticPublicKey{
		X: x,
		Y: y,
	}
	return privateKey, publicKey, nil
}

// Marshal ec public key
func (e *EllipticECDH) Marshal(pub crypto.PublicKey) []byte {
	publicKey := pub.(*EllipticPublicKey)
	return elliptic.Marshal(curve, publicKey.X, publicKey.Y)
}

// Unmarshal ec public key
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

// private pem file, parse PKCS#8 private key
func (e *EllipticECDH) ParsePKCS8ECPrivateKey(privateKeyDerBytes []byte) (*EllipticPrivateKey, *EllipticPublicKey, error) {
	var privateKey *ecdsa.PrivateKey
	// privateKey, err := x509.ParseECPrivateKey(privateKeyDer)
	key, err := x509.ParsePKCS8PrivateKey(privateKeyDerBytes)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, err
	}
	if key == nil {
		return nil, nil, ErrPrivKeyParseFailedFromPEM
	}
	privateKey = key.(*ecdsa.PrivateKey)

	privKey := EllipticPrivateKey{
		D: privateKey.D.Bytes(),
	}
	pubKey := EllipticPublicKey{
		X: privateKey.X,
		Y: privateKey.Y,
	}
	return &privKey, &pubKey, nil
}

func (e *EllipticECDH) ParsePKIXECPublicKeyFrom(publicKeyDerBytes []byte) (*EllipticPublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKeyDerBytes)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	fmt.Println(pub)
	publicKey := pub.(*ecdsa.PublicKey)

	pubKey := EllipticPublicKey{
		X: publicKey.X,
		Y: publicKey.Y,
	}
	return &pubKey, nil
}

func (e *EllipticECDH) DecodePEMToDERBytes(pemBytes []byte) []byte {
	block, _ := pem.Decode(pemBytes)
	return block.Bytes
}

// ECDH Compute Secret
func (e *EllipticECDH) ComputeSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error) {
	priv := privKey.(*EllipticPrivateKey)
	pub := pubKey.(*EllipticPublicKey)
	x, _ := curve.ScalarMult(pub.X, pub.Y, priv.D)
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
