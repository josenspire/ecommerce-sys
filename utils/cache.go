package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"net/http"
	"reflect"
	"time"
)

var (
	host   = ""
	db     = ""
	key    = ""
	config = ""
)

func init() {
	host = beego.AppConfig.String("redis_host")
	db = beego.AppConfig.String("redis_db")
	key = beego.AppConfig.String("redis_key")

	configMap := make(map[string]string)
	configMap["conn"] = host
	configMap["key"] = key
	configMap["dbNum"] = db

	if configByte, err := json.Marshal(configMap); err != nil {
		fmt.Println(err.Error())
	} else {
		config = string(configByte)
	}
}

func GetRedis() (cache.Cache, error) {
	return cache.NewCache("redis", config)
}

func ReadApiCache(ct *context.Context) {
	input := ct.Input

	var reqJson interface{}
	json.Unmarshal(input.RequestBody, &reqJson)

	fmt.Println("[RequestBody]:", reqJson)

	if redis, err := GetRedis(); err != nil {
		logs.Error(err.Error())
		ct.Abort(http.StatusInternalServerError, err.Error())
	} else {
		if cacheJsonObj := redis.Get(input.URI()); cacheJsonObj != nil {
			cacheJson := TransformInterfaceToMap(cacheJsonObj)
			cacheRequestBody := cacheJson["requestBody"]
			requestBody := TransformInterfaceToMap(input.RequestBody)

			if reflect.DeepEqual(cacheRequestBody, requestBody) {
				beego.Info("[API Cache]")
				ct.Output.SetStatus(http.StatusOK)
				ct.Output.JSON(cacheJson["responseBody"], true, true)
			}
		}
	}
}

func WriteApiCache(ct *context.Context, response interface{}) {
	input := ct.Input

	if redis, err := GetRedis(); err != nil {
		logs.Error(err.Error())
	} else {
		var requestBody interface{}
		json.Unmarshal(input.RequestBody, &requestBody)

		cacheData := make(map[string]interface{})
		cacheData["requestBody"] = requestBody
		cacheData["responseBody"] = response
		cacheByte, _ := json.Marshal(cacheData)

		err := redis.Put(input.URI(), cacheByte, time.Second*60*2)
		if err != nil {
			logs.Error(err.Error())
		}
	}
}

func ReadCacheDataByKey(cacheKey string) interface{} {
	if redis, err := GetRedis(); err != nil {
		logs.Error(err.Error())
		return err.Error()
	} else {
		if cacheJsonObj := redis.Get(cacheKey); cacheJsonObj != nil {
			cacheJson := TransformInterfaceToMap(cacheJsonObj)
			cacheResponseBody := cacheJson["responseBody"]
			return cacheResponseBody
		}
		return nil
	}
}
