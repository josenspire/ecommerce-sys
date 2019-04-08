package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

type AdvertController struct {
	beego.Controller
}

// @Title Advert Insert
// @Description Insert a new advert
// @Param	advertUrl	body	string	true	"advert url"
// @Param	relativeId	body	float64	false	"advert relate to"
// @Param	remark		body	string	false	"advert remark"
// @Failure	200400	{object}	models.ResponseModel
// @Success	200000	{object}	models.ResponseModel
// @router	/insertAdvert	[post]
func (adv *AdvertController) InsertAdvert() {
	var response ResponseModel
	advert := new(Advert)
	err := json.Unmarshal(adv.Ctx.Input.RequestBody, &advert)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = advert.InsertAdvert()
		if err != nil {
			beego.Error(err.Error())
			response.HandleFail(REQUEST_FAIL, err.Error())
		} else {
			response.HandleSuccess(nil)
		}
	}
	adv.Data["json"] = response
	adv.ServeJSON()
}

// @Title Advert Update
// @Description Update a advert
// @Param	advertId	body	float64	 true	"advert id"
// @Param	advertUrl	body	string	true	"advert url"
// @Param	relativeId	body	float64	false	"advert relate to"
// @Param	remark		body	string	false	"advert remark"
// @Failure	200400	{object}	models.ResponseModel
// @Success	200000	{object}	models.ResponseModel
// @router	/updateAdvert		[put]
func (adv *AdvertController) UpdateAdvert() {
	var response ResponseModel
	advert := Advert{}
	err := json.Unmarshal(adv.Ctx.Input.RequestBody, &advert)
	if err != nil {
		beego.Warning(err.Error())
		response.HandleError(err, PARAMS_MISSING)
	} else {
		err = advert.UpdateAdvertByAdvertId()
		if err == gorm.ErrRecordNotFound {
			response.HandleFail(RECORD_NOT_FOUND, ErrRecordNotFound.Error())
		} else if err != nil {
			beego.Error(err.Error())
			response.HandleFail(REQUEST_FAIL, err.Error())
		} else {
			response.HandleSuccess(nil)
		}
	}
	adv.Data["json"] = response
	adv.ServeJSON()
}

// @Title Advert List
// @Description Get advert list
// @Failure	200400	{object}	models.ResponseModel
// @Success	200000	{object}	models.ResponseModel
// @router	/list		[get]
func (adv *AdvertController) GetAdvertList() {
	var response ResponseModel
	advert := Advert{}
	advertList, err := advert.QueryAdvertList()
	if err == gorm.ErrRecordNotFound {
		response.HandleSuccess([]Advert{})
	} else if err != nil {
		beego.Error(err.Error())
		response.HandleError(err)
	} else {
		response.HandleSuccess(&advertList)
	}
	adv.Data["json"] = response
	adv.ServeJSON()
}
