package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"strconv"
)

type UserController struct {
	beego.Controller
}

// @Title Register
// @Description User register api by telephone
// @Param	telephone	body	string	true	"Register by cellphone"
// @Param	username	body	string	false	"User's username"
// @Param	password	body	string	true	"User password, length need to more then 6"
// @Param	nickname	body	string	true	"User nickname"
// @Param	signature	body	string	false	"User signature"
// @Param	male		body	bool	false	"Male/Female"
// @Param	invitationCode	query	string	true	"User's agent invitation code"
// @Success	200000 {object}	models.ResponseModel
// @Failure	200400 {object} models.ResponseModel
// @router /register [post]
func (u *UserController) Register() {
	var response ResponseModel
	user := new(User)
	dto := UserRegisterDTO{}
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &dto)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = user.Register(dto)
		if err != nil {
			beego.Error(err.Error())
			response.HandleFail(REQUEST_FAIL, err.Error())
		} else {
			response.HandleSuccess(nil, "Registration Successful")
		}
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @Title Login
// @Description User login api
// @Param	telephone	body	string	true	"Login by telephone"
// @Param	password	body	string	true	"User password, length need to more then 6"
// @Success	200000 {object}	models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /loginByTelephone [post]
func (u *UserController) LoginByTelephone() {
	var response ResponseModel
	reqParams := make(map[string]string)
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &reqParams)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleError(err)
	} else {
		telephone := reqParams["telephone"]
		password := reqParams["password"]
		user := new(User)
		err = user.LoginByTelephone(telephone, password)
		if err == gorm.ErrRecordNotFound {
			response.HandleError(ErrTelOrPswInvalid, USER_TELEPHONE_PSW_INVALID)
		} else if err == WarnAccountNeedVerify {
			// TODO: need to return the notification
			response.HandleFail(TELEPHONE_VERIFY, WarnAccountNeedVerify.Error())
		} else if err != nil {
			beego.Error(err.Error())
			response.HandleError(err)
		} else {
			response.HandleSuccess(user)
		}
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @Title LoginByWechat
// @Description User use wechat login api
// @Param	jsCode			body	string	true	"User's wechat 'jsCode'"
// @Param	userInfo		body	string	true	"User's wechat base profile"
// @Param	invitationCode	body	string	false	"Invitation code"
// @Success	200000 {object}	models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /loginByWechat [post]
func (u *UserController) LoginByWechat() {
	var response ResponseModel
	user := new(User)
	reqParams := make(map[string]string)
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &reqParams)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleError(err)
	} else {
		jsCode := reqParams["jsCode"]
		userInfo := reqParams["userInfo"]
		invitationCode := reqParams["invitationCode"]
		if jsCode == "" {
			response.HandleError(ErrParamsMissing)
		} else {
			user, err := user.LoginByWechat(jsCode, userInfo, invitationCode)
			if err != nil {
				beego.Error(err.Error())
				response.HandleError(err, REQUEST_FAIL)
			} else {
				response.HandleSuccess(&user, "")
			}
		}
	}
	u.Data["json"] = response
	u.ServeJSON()
}

// @Title Query User Teams
// @Description Query user's team information
// @Param	userId		body	float64		true	"User's Id"
// @Success	200000 {object}	models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /teams [post]
func (u *UserController) QueryUserTeams() {
	var response ResponseModel
	reqParams := make(map[string]float64)
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &reqParams)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleFail(REQUEST_FAIL, err.Error())
	} else {
		userId := int(reqParams["userId"])
		if IsEmptyString(strconv.Itoa(userId)) {
			response.HandleError(ErrParamsMissing)
		} else {
			var team *Team
			teamProfile, err := team.QueryUserTeams(uint64(userId))
			if err == gorm.ErrRecordNotFound {
				response.HandleSuccess(nil, WarnUserTeamMissing)
			} else if err != nil {
				beego.Error(err)
				response.HandleError(err)
			} else {
				response.HandleSuccess(teamProfile)
			}
		}
	}
	u.Data["json"] = response
	u.ServeJSON()
}
