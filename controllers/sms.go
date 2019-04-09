package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"strings"
)

type SMSController struct {
	beego.Controller
}

func (s *SMSController) ObtainSecurityCode() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		beego.Error(err)
		response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
	} else {
		telephone := reqArgs["telephone"].(string)
		userId := reqArgs["userId"].(float64)
		operationMode := strings.ToUpper(reqArgs["operationMode"].(string))

		var sms *SMS
		smsProfile, err := sms.ObtainSecurityCode(telephone, uint64(userId), operationMode)
		if err == WarnTelephoneNotRegistered {
			response.HandleFail(TELEPHONE_HAS_NOT_REGISTERED, WarnTelephoneNotRegistered.Error())
		} else if err == WarnTelephoneAlreadyRegistered {
			response.HandleFail(TELEPHONE_HAS_BEEN_USED, WarnTelephoneAlreadyRegistered.Error())
		} else if err != nil {
			beego.Error(err)
			response.HandleError(err, REQUEST_FAIL)
		} else {
			response.HandleSuccess(smsProfile, "security code send success, will expire at 15 min after")
		}
	}
	s.Data["json"] = response
	s.ServeJSON()
}

func (s *SMSController) VerifySecurityCode() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(s.Ctx.Input.RequestBody, &reqArgs)
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
	s.Data["json"] = response
	s.ServeJSON()
}
