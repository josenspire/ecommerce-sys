package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "E-Commerce system"
	c.Data["Email"] = "josenspire@gmail.com"
	c.TplName = "index.tpl"
}
