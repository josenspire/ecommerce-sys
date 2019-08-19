package controllers

import (
	. "ecommerce-sys/commons"
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type ClassifyController struct {
	beego.Controller
}

// @Title Create Classify
// @Description Create a new classify
// @Param	classifyName		body	string		true	"Classify name"
// @Param	classifyIcon		body	string		true	"Classify icon"
// @Param	classifyPriority	body	int			false	"Classify priority, default is 1"
// @Param	status				body	string		false	"Classify status, default is 'active'"
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object} models.ResponseModel
// @router /create [post]
func (cls *ClassifyController) CreateClassify() {
	var response ResponseModel
	var classify *Classify
	err := json.Unmarshal(cls.Ctx.Input.RequestBody, &classify)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = classify.CreateClassify()
		if err != nil {
			logs.Error(err.Error())
			response.HandleError(err)
		} else {
			response.HandleSuccess(nil, "Create new classify succeed")
		}
	}
	cls.Data["json"] = response
	cls.ServeJSON()
}

// @Title Create Category
// @Description Create a new category
// @Param	categoryName		body	string		true	"Category name"
// @Param	categoryIcon		body	string		true	"Category icon"
// @Param	categoryPriority	body	string		true	"Category priority, default is '1'"
// @Param	classifyId			body	string		true	"Classify's id"
// @Param	status				body	string		false	"Category status, default is 'active'"
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object} models.ResponseModel
// @router /category/create [post]
func (cls *ClassifyController) CreateCategory() {
	var response ResponseModel
	var category *Category
	err := json.Unmarshal(cls.Ctx.Input.RequestBody, &category)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleError(err, PARAMS_MISSING)
	} else {
		var classify *Classify
		err := classify.CreateCategory(category)
		if err != nil {
			logs.Error(err.Error())
			response.HandleError(err)
		} else {
			response.HandleSuccess(nil, "Create new category succeed")
		}
	}
	cls.Data["json"] = response
	cls.ServeJSON()
}

// @Title Classify
// @Description Query classifies information
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object} models.ResponseModel
// @router /list [get]
func (cls *ClassifyController) QueryClassifies() {
	var response ResponseModel
	var classify *Classify
	classifies, err := classify.QueryClassifies()
	if err != nil {
		logs.Error(err.Error())
		response.HandleError(err)
	} else {
		var message = ""
		if len(classifies) == 0 {
			message = WarnClassifiesMissing
		}
		response.HandleSuccess(classifies, message)
	}
	cls.Data["json"] = response
	cls.ServeJSON()
}
