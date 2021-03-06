package utils

import (
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"os"
)

func OsFileReader(fileName string) []byte {
	if fileObj, err := os.Open(fileName); err == nil {
		defer fileObj.Close()
		if contents, err := ioutil.ReadAll(fileObj); err == nil {
			return contents
		}
	} else {
		logs.Error(err.Error())
	}
	return nil
}
