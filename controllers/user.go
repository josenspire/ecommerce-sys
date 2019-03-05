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
// @Param	password	query	string	true	"User password, length need to more then 6"
// @Param	nickname	query	string	true	"User nickname"
// @Param	signature	query	string	false	"User signature"
// @Param	male		query	bool	true	"Male/Female"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/login	[post]
func (u *UserController) Register() {
	var user User
	err := json.Unmarshal(u.Ctx.Input.RequestBody, &user.UserProfile)
	if err != nil {
		u.Data["json"] = HandleError(err, PARAMS_MISSING)
	} else {
		err = user.Register()
		if err != nil {
			u.Data["json"] = HandleFail(REQUEST_FAIL, err.Error())
		} else {
			u.Data["json"] = HandleSuccess(user, "Registration Successful")
		}
	}
	u.ServeJSON()
}

// @Title Login
// @Description User login api
// @Param	phone	query	string	true	"Login by cellphone"
// @Param	password	query	string	true	"User password, length need to more then 6"
// @Success	200	{object}	models.ResponseModel
// @Failure	403
// @router	/login	[post]
func (u *UserController) Login() {
	// resParams := models.User{}
	// u.Ctx.Input.Bind(&resParams.Cellphone, "cellphone")
	// u.Ctx.Input.Bind(&resParams.Password, "password")

	u.Data["json"] = "success"
	u.ServeJSON()
}
