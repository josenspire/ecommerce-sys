package controllers

import (
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

// @Title Login
// @Description User login api
// @Param	phone	query	string	true	"Login by cellphone"
// @Param	password	query	string	true	"User password, length need to more then 6"
// @Success	200	{json}	login success
// @Failure	403	{json}	user not exist
// @router	/login	[post]
func (u *UserController) Login() {
	// resParams := models.User{}
	// u.Ctx.Input.Bind(&resParams.Cellphone, "cellphone")
	// u.Ctx.Input.Bind(&resParams.Password, "password")

	u.Data["json"] = "success"
	u.ServeJSON()
}
