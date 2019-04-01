package controllers

import (
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
// @Param	models.Classify		query	object		true	"classify request model"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/create		[post]
func (cls *ClassifyController) CreateClassify() {
	var response ResponseModel
	var classify *Classify
	err := json.Unmarshal(cls.Ctx.Input.RequestBody, &classify)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = classify.CreateClassify()
		if err != nil {
			logs.Error(err)
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
// @Param	models.Category		query	object		true	"Create a new category
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/category/create		[post]
func (cls *ClassifyController) CreateCategory() {
	var response ResponseModel
	var category *Category
	err := json.Unmarshal(cls.Ctx.Input.RequestBody, &category)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		var classify *Classify
		err := classify.CreateCategory(category)
		if err != nil {
			logs.Error(err)
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
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/list		[get]
func (cls *ClassifyController) QueryClassifies() {
	var response ResponseModel
	var classify *Classify
	classifies, err := classify.QueryClassifies()
	if err != nil {
		logs.Error(err)
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
