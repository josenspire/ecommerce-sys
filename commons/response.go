package commons

import (
	. "ecommerce-sys/utils"
)

type ResponseModel struct {
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (res *ResponseModel) HandleSuccess(data interface{}, message ...string) {
	res.Code = REQUEST_SUCCESS
	res.Data = data
	if len(message) > 0 {
		res.Message = message[0]
	} else {
		res.Message = ""
	}
}

func (res *ResponseModel) HandleFail(errorCode uint, message ...string) {
	res.Code = errorCode
	if len(message) > 0 {
		res.Message = message[0]
	} else {
		res.Message = ""
	}
}

func (res *ResponseModel) HandleError(err error, errorCode ...uint) {
	if len(errorCode) > 0 {
		res.Code = errorCode[0]
	} else {
		res.Code = SERVER_UNKNOW_ERROR
	}
	res.Message = err.Error()
}
