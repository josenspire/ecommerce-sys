package utils

import (
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

			privBytes := OsFileReader("ecdh_priv.pem")
			privDerBytes := ellipticECDH.DecodePEMToDERBytes(privBytes)

			privKey1, pubKey1, _ = ellipticECDH.ParsePKCS8ECPrivateKey(privDerBytes)
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

			privBytes := OsFileReader("ecdh_priv.pem")
			privDerBytes := ellipticECDH.DecodePEMToDERBytes(privBytes)

			pubBytes := OsFileReader("ecdh_pub.pem")
			pubDerBytes := ellipticECDH.DecodePEMToDERBytes(pubBytes)

			privKey1, pubKeyTemp, _ = ellipticECDH.ParsePKCS8ECPrivateKey(privDerBytes)
			pubKey1, _ = ellipticECDH.ParsePKIXECPublicKeyFrom(pubDerBytes)

			privKey2, pubKey2, _ := ellipticECDH.GenerateECKeyPair()

			secret1, _ := ellipticECDH.ComputeSecret(privKey1, pubKey2)
			secret2, _ := ellipticECDH.ComputeSecret(privKey2, pubKey1)

			secretStr1 := base64.StdEncoding.EncodeToString(secret1)
			secretStr2 := base64.StdEncoding.EncodeToString(secret2)

			fmt.Println(secretStr1, secretStr2)
			convey.So(pubKey1.X.String(), convey.ShouldEqual, pubKeyTemp.X.String())
			convey.So(secretStr1, convey.ShouldEqual, secretStr2)
		})
	})
}
