package commons

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
)

type RequestModel struct {
	Data      string `json:"data"`
	SecretKey string `json:"secretKey"`
	Signature string `json:"signature"`
}

type AspectControl struct{}

func (asp *AspectControl) HandleRequest(ct *context.Context) {
	var requestModel *RequestModel

	inputContent := ct.Input.RequestBody
	if len(inputContent) != 0 {
		err := json.Unmarshal(ct.Input.RequestBody, requestModel)
		if err != nil {
			beego.Error(err.Error())
			asp.HandleResponse(ct)
		}
		// TODO: decrypt and calculate length
		var requestContent = make([]byte, len(inputContent))
		ct.Input.RequestBody = requestContent
		log.Println("request body: ", ct.Input.RequestBody)
	}
}

func (asp *AspectControl) HandleResponse(ct *context.Context) {
	// 	TODO: response
	// resArgs := make(map[string]interface{})
	// log.Println("response body: ", ct.Output.JSON(&resArgs, true, true))
}