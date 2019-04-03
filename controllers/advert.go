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
// @Param	advertUrl	query	string	true	"advert url"
// @Param	relativeId	query	float64	false	"advert relate to"
// @Param	remark		query	string	false	"advert remark"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
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
// @Param	advertId		query	float64	 true	"advert id"
// @Param	advertUrl	query	string	true	"advert url"
// @Param	relativeId	query	float64	false	"advert relate to"
// @Param	remark		query	string	false	"advert remark"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
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
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
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
