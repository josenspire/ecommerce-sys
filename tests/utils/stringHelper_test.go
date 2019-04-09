package utils

import (
	"ecommerce-sys/utils"
	"fmt"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateRandString(t *testing.T) {
	convey.Convey("Subject: String helper", t, func() {
		convey.Convey("Testing generate random string", func() {
			result := utils.GenerateRandString(6)
			fmt.Println("-------", result)
			convey.So(len(result), convey.ShouldEqual, 6)
		})
	})
}
