package utils

import (
	"crypto/ecdsa"
	. "ecommerce-sys/utils"
	"encoding/base64"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	convey.Convey("Subject: EllipticECDH", t, func() {
		convey.Convey("Testing parse PKCS#8 private key and compute secret", func() {
			var ellipticECDH *EllipticECDH
			var privKey1 *EllipticPrivateKey
			var pubKey1 *EllipticPublicKey

			privBytes := OsFileReader("./../pem/ecdh_priv.pem")
			privDerBytes := ellipticECDH.DecodePEMToDERBytes(privBytes)

			ellipticECDH, _ = ellipticECDH.ParsePKCS8ECPrivateKey(privDerBytes)
			privKey1 = ellipticECDH.PrivateKey
			pubKey1 = ellipticECDH.PublicKey
			privKey2, pubKey2, _ := ellipticECDH.GenerateECKeyPair()

			secret1, _ := ellipticECDH.ComputeSecret(privKey1, pubKey2)
			secret2, _ := ellipticECDH.ComputeSecret(privKey2, pubKey1)

			secretStr1 := base64.StdEncoding.EncodeToString(secret1)
			secretStr2 := base64.StdEncoding.EncodeToString(secret2)

			fmt.Println(secretStr1, secretStr2)
			convey.So(secretStr1, convey.ShouldEqual, secretStr2)
		})

		convey.Convey("Testing parse PKCS#8 public key and compute secret", func() {
			var ellipticECDH *EllipticECDH
			var privKey1 *EllipticPrivateKey
			var pubKey1 *EllipticPublicKey
			var pubKeyTemp *EllipticPublicKey

			privBytes := OsFileReader("./../pem/ecdh_priv.pem")
			privDerBytes := ellipticECDH.DecodePEMToDERBytes(privBytes)

			pubBytes := OsFileReader("./../pem/ecdh_pub.pem")
			pubDerBytes := ellipticECDH.DecodePEMToDERBytes(pubBytes)

			ellipticECDH, _ = ellipticECDH.ParsePKCS8ECPrivateKey(privDerBytes)
			privKey1, pubKeyTemp = ellipticECDH.PrivateKey, ellipticECDH.PublicKey
			pubKey1, _, _ = ellipticECDH.ParsePKIXECPublicKey(pubDerBytes)

			privKey2, pubKey2, _ := ellipticECDH.GenerateECKeyPair()

			secret1, _ := ellipticECDH.ComputeSecret(privKey1, pubKey2)
			secret2, _ := ellipticECDH.ComputeSecret(privKey2, pubKey1)

			secretStr1 := base64.StdEncoding.EncodeToString(secret1)
			secretStr2 := base64.StdEncoding.EncodeToString(secret2)

			fmt.Println(secretStr1, secretStr2)
			convey.So(pubKey1.X.String(), convey.ShouldEqual, pubKeyTemp.X.String())
			convey.So(secretStr1, convey.ShouldEqual, secretStr2)
		})

		convey.Convey("Testing generate PKIX stander public key", func() {
			var ellipticECDH *EllipticECDH
			var inputPublicKey = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE60BkU5fcacDtqV6Co2rPgxzfXdmLcnVNau6JE84AVPRz3x/cZFlJK6aSrSgzqxUPAU8NBNj1J4Z2oHdsjzZpMg=="
			var publicKeyPEM = ellipticECDH.GeneratePKIXPublicKey(inputPublicKey)
			convey.So(publicKeyPEM, convey.ShouldContainSubstring, `MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE60BkU5fcacDtqV6Co2rPgxzfXdmL`)
		})

		convey.Convey("Testing get public key form pem and format to base64", func() {
			var ellipticECDH *EllipticECDH

			pubBytes := OsFileReader("./../pem/ecdh_pub.pem")
			publicKeyStr := ellipticECDH.GetPKIXPublicKeyBlockFromPEM(pubBytes)

			var expectation = "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE60BkU5fcacDtqV6Co2rPgxzfXdmLcnVNau6JE84AVPRz3x/cZFlJK6aSrSgzqxUPAU8NBNj1J4Z2oHdsjzZpMg=="
			convey.So(publicKeyStr, convey.ShouldStartWith, expectation)
		})

		convey.Convey("Testing private key signature data and public verify signature", func() {
			var ellipticECDH *EllipticECDH
			var ecdsaPrivateKey *ecdsa.PrivateKey

			privBytes := OsFileReader("./../pem/ecdh_priv.pem")
			privDerBytes := ellipticECDH.DecodePEMToDERBytes(privBytes)
			ellipticECDH, _ = ellipticECDH.ParsePKCS8ECPrivateKey(privDerBytes)
			ecdsaPrivateKey = ellipticECDH.ECDSAPrivateKey
			signatureData, _ := ellipticECDH.Signature("德玛西亚", ecdsaPrivateKey)

			verifyResult := ellipticECDH.VerifySignature(signatureData, &ecdsaPrivateKey.PublicKey)
			convey.So(verifyResult, convey.ShouldBeTrue)
		})
	})
}
