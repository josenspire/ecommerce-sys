package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type WechatMP struct {
	appid      string
	secret     string
	js_code    string
	grant_type string
}

var (
	WechatAppId     string
	WechatSecret    string
	WechatGrantType string
)

const (
	ConnectTimeOut   = time.Second * 30
	ReadWriteTimeOut = time.Second * 30

	UrlJsCode2Session string = "https://api.weixin.qq.com/sns/jscode2session"
)

func init() {
	WechatAppId = beego.AppConfig.String("WechatAppId")
	WechatSecret = beego.AppConfig.String("WechatSecret")
	WechatGrantType = beego.AppConfig.String("WechatGrantType")
}

func JsCode2Session(jsCode string) (string, string, error) {
	req := httplib.NewBeegoRequest(UrlJsCode2Session, http.MethodGet).SetTimeout(ConnectTimeOut, ReadWriteTimeOut).SetEnableCookie(true)
	paramsBody := fmt.Sprintf("appid=%s&secret=%s&js_code=%s&grant_type=%s", WechatAppId, WechatSecret, jsCode, WechatGrantType)
	req.Body(paramsBody)

	jsonObj := make(map[string]interface{})
	err := req.ToJSON(&jsonObj)
	if err != nil {
		logs.Error(err.Error())
		return "", "", err
	}
	if strconv.FormatFloat(jsonObj["errcode"].(float64), 'E', -1, 64) == "" {
		return jsonObj["openid"].(string), jsonObj["session_key"].(string), nil
	}
	logs.Informational(jsonObj)
	return "", "", errors.New(jsonObj["errmsg"].(string))
}
