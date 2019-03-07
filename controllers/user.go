package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

// @Title Register
// @Description User register api by telephone
// @Param	telephone	query	string	true	"Register by cellphone"
// @Param	username	query	string	false	"User's username"
// @Param	password	query	string	true	"User password, length need to more then 6"
// @Param	nickname	query	string	true	"User nickname"
// @Param	signature	query	string	false	"User signature"
// @Param	male		query	bool	false	"Male/Female"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/register	[post]
func (u *UserController) Register() {
	var response ResponseModel
	var user = new(User)

	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user.UserProfile)
	if err != nil {
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = user.Register()
		if err != nil {
			response.HandleFail(REQUEST_FAIL, err.Error())
		} else {
			response.HandleSuccess(user, "Registration Successful")
		}
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @Title Login
// @Description User login api
// @Param	telephone	query	string	true	"Login by telephone"
// @Param	password	query	string	true	"User password, length need to more then 6"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/loginByTelephone	[post]
func (u *UserController) LoginByTelephone() {
	var response ResponseModel
	reqParams := make(map[string]interface{})
	json.Unmarshal(u.Ctx.Input.RequestBody, &reqParams)
	telephone := reqParams["telephone"].(string)
	password := reqParams["password"].(string)

	user := new(User)
	err := user.LoginByTelephone(telephone, password)
	if err != nil {
		response.HandleError(err)
	} else if user == nil {
		response.HandleError(ErrTelOrPswInvalid, USER_TELEPHONE_PSW_INVALID)
	} else {
		response.HandleSuccess(user)
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// // @Title LoginByWechat
// // @Description User use wechat login api
// // @Param		query	string	true	"Login by cellphone"
// // @Param	password	query	string	true	"User password, length need to more then 6"
// // @Success	200000	{object}	models.ResponseModel
// // @Failure	200400
// // @router	/loginByWechat	[post]
// func (u *UserController) LoginByWechat() {
// 	var response *ResponseModel
// 	var user *User
// 	jsCode := u.GetString("jsCode")
// 	userInfo := u.GetString("userInfo")
// 	invitationCode := u.GetString("invitationCode")
//
// 	if jsCode == "" {
// 		response = HandleFail(REQUEST_FAIL, "Sorry, params missing")
// 	} else {
// 		user, err := user.LoginByWechat(jsCode, userInfo, invitationCode)
// 		if err != nil {
// 			// TODO handle error
// 		} else {
// 			response = HandleSuccess(user, "")
// 		}
// 	}
// 	u.Data["json"] = response
// 	u.ServeJSON()
// }
