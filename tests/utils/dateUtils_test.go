package utils

import (
	"ecommerce-sys/utils"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestGenerateNowDateString(t *testing.T) {
	convey.Convey("Subject: Date utils", t, func() {
		convey.Convey("Testing generate current date string", func() {
			result := utils.GenerateNowDateString()
			convey.So(result, convey.ShouldStartWith, "2019-03-28")
		})
	})
}
