package utils

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"reflect"
)

// can not handle key as Hump named
func TransformStructToMap(st interface{}) map[string]interface{} {
	t := reflect.TypeOf(st)
	v := reflect.ValueOf(st)

	var params = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		params[t.Field(i).Name] = v.Field(i).String()
	}
	return params
}

func TransformInterfaceToMap(origin interface{}) map[string]interface{} {
	var obj map[string]interface{}
	err := json.Unmarshal(origin.([]byte), &obj)
	if err != nil {
		return nil
	}
	return obj
}

func TransformStructToJSONMap(model interface{}) (map[string]interface{}, error) {
	if params, err := json.Marshal(model); err != nil {
		return nil, err
	} else {
		var paramsMap map[string]interface{}
		if err = json.Unmarshal([]byte(params), &paramsMap); err != nil {
			return nil, err
		}
		return paramsMap, nil
	}
}

func TransformByteToJSON(str []byte) interface{} {
	var tsJson interface{}
	if err := json.Unmarshal(str, &tsJson); err != nil {
		beego.Error(err.Error())
		return nil
	} else {
		beego.Info(tsJson)
		return tsJson
	}
}

func MergeMaps(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	for k, v := range m2 {
		m1[k] = v
	}
	return m1
}
