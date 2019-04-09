package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
)

type SMSController struct {
	beego.Controller
}

func (sms *SMSController) ObtainSecurityCode() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(sms.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		beego.Error(err)
		response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
	} else {
		telephone := reqArgs["telephone"].(string)
		userId := reqArgs["userId"].(float64)
		operationMode := reqArgs["operationMode"].(string)
		var sms *SMS
		smsProfile, err := sms.ObtainSecurityCode(telephone, uint64(userId), operationMode)
		if err != nil {
			beego.Error(err)
			response.HandleError(err, REQUEST_FAIL)
		} else {
			response.HandleSuccess(smsProfile, "security code send success, will expire at 15 min after")
		}
	}
	sms.Data["json"] = response
	sms.ServeJSON()
}

func (sms *SMSController) VerifySecurityCode() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(sms.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		beego.Error(err)
		response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
	} else {
		userId := reqArgs["userId"].(float64)
		telephone := reqArgs["telephone"].(string)
		securityCode := reqArgs["securityCode"].(string)
		operationMode := reqArgs["operationMode"].(string)
		var sms *SMS
		verifyResult, err := sms.VerifySecurityCode(telephone, uint64(userId), securityCode, operationMode)
		if err != nil {
			beego.Error(err)
			response.HandleError(err, REQUEST_FAIL)
		} else {
			if verifyResult {
				response.HandleSuccess(nil, "security code verify success")
			} else {
				response.HandleFail(SECURITY_CODE_INVALID, ErrSecurityCodeInvalid)
			}
		}
	}
	sms.Data["json"] = response
	sms.ServeJSON()
}
