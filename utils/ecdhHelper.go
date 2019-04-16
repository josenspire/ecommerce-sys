package utils

import (
	"bytes"
	"compress/gzip"
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/astaxie/beego"
	"hash"
	"io"
	"math/big"
)

type ECDH interface {
	GenerateECKeyPair() (crypto.PrivateKey, crypto.PublicKey, error)
	GenerateECKeyPairToPEM(curve elliptic.Curve) ([]byte, []byte, error)
	GeneratePKIXPublicKey(publicKeyBlock string) string
	Marshal(pub crypto.PublicKey) []byte
	Unmarshal(data []byte) (crypto.PublicKey, bool)
	ParsePKCS8ECPrivateKey(privateKeyDerBytes []byte) (*EllipticECDH, error)
	ParsePKIXECPublicKey(publicKeyDerBytes []byte) (*EllipticPublicKey, error)
	GetPKIXPublicKeyBlockFromPEM(pemBytes []byte) string
	DecodePEMToDERBytes(pemBytes []byte) []byte
	ComputeSecret(privKey crypto.PrivateKey, pubKey crypto.PublicKey) ([]byte, error)

	Signature(message string, privateKey *ecdsa.PrivateKey) (signatureData *SignatureData, err error)
	VerifySignature(signatureData *SignatureData, publicKey *ecdsa.PublicKey) (status bool)
}

type EllipticECDH struct {
	PrivateKey      *EllipticPrivateKey
	PublicKey       *EllipticPublicKey
	ECDSAPrivateKey *ecdsa.PrivateKey
}

type EllipticPrivateKey struct {
	D []byte
}

type EllipticPublicKey struct {
	X *big.Int
	Y *big.Int
}

const (
	EC_PUBLIC_KEY_BLOCK_BEGIN = "-----BEGIN PUBLIC KEY-----"
	EC_PUBLIC_KEY_BLOCK_END   = "-----END PUBLIC KEY-----"
)

// elliptic.P224(), elliptic.P384(), elliptic.P521()
var curve = elliptic.P256()

// generate ec key pair and return
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

// generate ec key pair and write them to PEM file, then return them
func (e *EllipticECDH) GenerateECKeyPairToPEM(curve elliptic.Curve) ([]byte, []byte, error) {
	// priv, x, y, err := elliptic.GenerateKey(curve, rand.Reader)
	// if err != nil {
	// 	beego.Error(err.Error())
	// 	return nil, nil, err
	// }
	// 	TODO: GEN key
	return nil, nil, nil
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
func (e *EllipticECDH) ParsePKCS8ECPrivateKey(privateKeyDerBytes []byte) (*EllipticECDH, error) {
	var privateKey *ecdsa.PrivateKey
	// privateKey, err := x509.ParseECPrivateKey(privateKeyDer)
	key, err := x509.ParsePKCS8PrivateKey(privateKeyDerBytes)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	if key == nil {
		return nil, ErrPrivKeyParseFailedFromPEM
	}
	privateKey = key.(*ecdsa.PrivateKey)

	privKey := EllipticPrivateKey{
		D: privateKey.D.Bytes(),
	}
	pubKey := EllipticPublicKey{
		X: privateKey.X,
		Y: privateKey.Y,
	}
	ecdhKey := &EllipticECDH{
		PrivateKey:      &privKey,
		PublicKey:       &pubKey,
		ECDSAPrivateKey: privateKey,
	}
	return ecdhKey, nil
}

func (e *EllipticECDH) ParsePKIXECPublicKey(publicKeyDerBytes []byte) (*EllipticPublicKey, *ecdsa.PublicKey, error) {
	pub, err := x509.ParsePKIXPublicKey(publicKeyDerBytes)
	if err != nil {
		beego.Error(err.Error())
		return nil, nil, err
	}
	publicKey := pub.(*ecdsa.PublicKey)

	pubKey := EllipticPublicKey{
		X: publicKey.X,
		Y: publicKey.Y,
	}
	return &pubKey, publicKey, nil
}

// decode PEM, support private/public key, return block data
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

// handle client public key which is less of begin and end
func (e *EllipticECDH) GeneratePKIXPublicKey(publicKeyBlock string) string {
	return fmt.Sprintf("%s\n%s\n%s", EC_PUBLIC_KEY_BLOCK_BEGIN, publicKeyBlock, EC_PUBLIC_KEY_BLOCK_END)
}

// get public key form pem and format to base64
func (e *EllipticECDH) GetPKIXPublicKeyBlockFromPEM(pemBytes []byte) string {
	blockBytes := e.DecodePEMToDERBytes(pemBytes)
	if len(blockBytes) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString(blockBytes)
}

// SignatureData 用于保存签名的数据
type SignatureData struct {
	r         *big.Int
	s         *big.Int
	signHash  *[]byte
	signature *[]byte
}

func (e *EllipticECDH) Signature(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	var h hash.Hash
	h = sha1.New()
	r := big.NewInt(0)
	s := big.NewInt(0)
	_, err := io.WriteString(h, message)
	if err != nil {
		beego.Error(err.Error())
		return "", err
	}
	signHash := h.Sum(nil)
	r, s, err = ecdsa.Sign(rand.Reader, privateKey, signHash)
	if err != nil {
		return "", err
	}
	signature := r.Bytes()
	signature = append(signature, s.Bytes()...)
	// signatureData = &SignatureData{
	// 	r:         r,
	// 	s:         s,
	// 	signHash:  &signHash,
	// 	signature: &signature,
	// }
	rt, rErr := r.MarshalText()
	if rErr != nil {
		beego.Error(rErr.Error())
		return "", err
	}
	st, sErr := r.MarshalText()
	if sErr != nil {
		beego.Error(rErr.Error())
		return "", err
	}

	var certBytes bytes.Buffer
	writer := gzip.NewWriter(&certBytes)
	defer writer.Close()
	_, err = writer.Write(rt)
	_, err = writer.Write(st)
	if err != nil {
		return "", err
	}
	writer.Flush()
	return base64.StdEncoding.EncodeToString(certBytes.Bytes()), nil
}

func (e *EllipticECDH) VerifySignature(signatureData *SignatureData, publicKey *ecdsa.PublicKey) (status bool) {
	status = ecdsa.Verify(publicKey, *signatureData.signHash, signatureData.r, signatureData.s)
	return
}
