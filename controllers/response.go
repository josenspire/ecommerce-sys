package controllers

import (
	. "ecommerce-sys/utils"
)

type ResponseModel struct {
	Code    uint        `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func HandleSuccess(data interface{}, message string) *ResponseModel {
	responseModel := ResponseModel{}
	responseModel.Code = REQUEST_SUCCESS
	responseModel.Data = data
	responseModel.Message = message
	return &responseModel
}

func HandleFail(errorCode uint, message string) *ResponseModel {
	responseModel := ResponseModel{}
	responseModel.Code = errorCode
	responseModel.Message = message
	return &responseModel
}

func HandleError(err error, errorCode ...uint) *ResponseModel {
	responseModel := ResponseModel{}
	if len(errorCode) > 0 {
		responseModel.Code = errorCode[0]
	} else {
		responseModel.Code = SERVER_UNKNOW_ERROR
	}
	responseModel.Message = err.Error()
	return &responseModel
}
