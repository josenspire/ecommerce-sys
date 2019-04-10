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
		convey.Convey("Testing decrypt text data by aes", func() {
			var ellipticECDH *EllipticECDH
			privKey1, pubKey1, _ := ellipticECDH.GenerateKeyPair()
			privKey2, pubKey2, _ := ellipticECDH.GenerateKeyPair()

			secret1, _ := ellipticECDH.ComputeSecret(privKey1, pubKey2)
			secret2, _ := ellipticECDH.ComputeSecret(privKey2, pubKey1)

			// pubStr := ellipticECDH.Marshal(pubKey1)
			pubStr2 := ellipticECDH.Marshal(pubKey2)

			publicKey := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEOs56iDUHeQ5tcdFlRKzHKvSdR8Y/pFJMbYlZWF90dqGXVFfCf0/3ZiKYrKAeOvR3HqXcxMvQudLn+Y99X0FMQw=="

			result, _ := ellipticECDH.Unmarshal([]byte(publicKey))
			fmt.Println("secret: ", result, "=======", base64.StdEncoding.EncodeToString(pubStr2))

			secretStr1 := base64.StdEncoding.EncodeToString(secret1)
			secretStr2 := base64.StdEncoding.EncodeToString(secret2)
			fmt.Println(secretStr1, secretStr2)
			convey.So(secretStr1, convey.ShouldEqual, secretStr2)
		})
	})
}

// func TestPEM(t *testing.T) {
// 	convey.Convey("Subject: PEM", t, func() {
// 		convey.Convey("Testing pem decode", func() {
// 			publicKey := "MFYwEAYHKoZIzj0CAQYFK4EEAAoDQgAEOs56iDUHeQ5tcdFlRKzHKvSdR8Y/pFJMbYlZWF90dqGXVFfCf0/3ZiKYrKAeOvR3HqXcxMvQudLn+Y99X0FMQw=="
//
// 			pubBlock, _ := pem.Decode([]byte(publicKey))
// 			pubKeyValue, err := x509.ParsePKIXPublicKey(pubBlock.Bytes)
// 			if err != nil {
// 				panic(err)
// 			}
// 			pub := pubKeyValue.(*rsa.PublicKey)
// 			fmt.Println("------", pubKeyValue)
// 			fmt.Println("======", pub)
// 			convey.So(pub, convey.ShouldStartWith, "MF")
// 		})
// 	})
// }
