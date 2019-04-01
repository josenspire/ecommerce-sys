package utils

import (
	"ecommerce-sys/utils"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestAESEncrypt(t *testing.T) {
	encodeStr, secretKeyStr := "德玛西亚", "德玛西亚永不言弃"
	convey.Convey("Subject: Date utils", t, func() {
		convey.Convey("Testing encrypt plainText by aes", func() {
			result, _ := utils.AESEncrypt(encodeStr, secretKeyStr)
			fmt.Println(result)
			convey.So(result, convey.ShouldStartWith, "bJTJMCaSC0NTdVcLpat/5Q==")
		})
	})
}

func TestAESDecrypt(t *testing.T) {
	deCodeStr, secretKeyStr := `bJTJMCaSC0NTdVcLpat/5Q==`, "德玛西亚永不言弃"
	convey.Convey("Subject: Date utils", t, func() {
		convey.Convey("Testing decrypt text data by aes", func() {
			result, _ := utils.AESDecrypt(deCodeStr, secretKeyStr)
			fmt.Println(result)
			convey.So(result, convey.ShouldStartWith, "德玛西亚")
		})
	})
}
