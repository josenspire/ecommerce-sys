package controllers

import (
	. "ecommerce-sys/utils"
)

type ResponseModel struct {
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func (res *ResponseModel) HandleSuccess(data interface{}, message ...string) {
	// responseModel := ResponseModel{}
	res.Code = REQUEST_SUCCESS
	res.Data = data
	res.Message = message[0]
	// return &responseModel
}

func (res *ResponseModel) HandleFail(errorCode uint, message ...string) {
	// responseModel := ResponseModel{}
	res.Code = errorCode
	res.Message = message[0]
	// return &responseModel
}

func (res *ResponseModel) HandleError(err error, errorCode ...uint) {
	// responseModel := ResponseModel{}
	if len(errorCode) > 0 {
		res.Code = errorCode[0]
	} else {
		res.Code = SERVER_UNKNOW_ERROR
	}
	res.Message = err.Error()
	// return &responseModel
}
