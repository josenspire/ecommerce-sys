package controllers

import (
	. "ecommerce-sys/commons"
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
)

type ProductController struct {
	beego.Controller
}

func (pd *ProductController) InsertProduct() {
	var response ResponseModel
	reqArgs := make(map[string]ProductDTO)
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleError(err)
	} else {
		dto := reqArgs["product"]
		if dto == (ProductDTO{}) {
			response.HandleFail(PARAMS_MISSING, ErrCreateRecordsIsEmpty.Error())
		} else {
			var product *Product
			err := product.InsertProduct(&dto)
			if err != nil {
				logs.Error(err.Error())
				response.HandleError(err)
			} else {
				response.HandleSuccess(nil, "record inserted success")
			}
		}
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}

func (pd *ProductController) InsertMultipleProducts() {
	var response ResponseModel
	reqArgs := make(map[string][]ProductDTO)
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleError(err)
	} else {
		dtos := reqArgs["products"]
		if len(dtos) == 0 {
			response.HandleFail(PARAMS_MISSING, ErrCreateRecordsIsEmpty.Error())
		} else {
			var product *Product
			err := product.InsertMultipleProducts(&dtos)
			if err != nil {
				logs.Error(err.Error())
				response.HandleError(err)
			} else {
				successMessage := fmt.Sprintf("%d records are inserted success", len(dtos))
				response.HandleSuccess(nil, successMessage)
			}
		}
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}

// @Title Product list
// @Description Query product list by product type.
// @Params	productType		body	true	"Product type, default is 'recommend'"
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /recommend [get]
func (pd *ProductController) QueryProducts() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleError(err)
	} else {
		productType := reqArgs["productType"].(string)
		pageIndex := int(reqArgs["pageIndex"].(float64))

		if IsEmptyString(productType) {
			productType = "normal"
		}
		if pageIndex <= 0 {
			pageIndex = 1
		}
		var products *[]Product
		var product *Product
		products, err = product.QueryProductsByProductType(productType, pageIndex)
		if err != nil && err != gorm.ErrRecordNotFound {
			logs.Error(err.Error())
			response.HandleError(err)
		} else {
			response.HandleSuccess(products)
		}
	}

	pd.Data["json"] = response
	pd.ServeJSON()
}

// @Title Product - details
// @Description Query product details
// @Params	productId		body		float64  	true	"Product id"
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /recommend [post]
func (pd *ProductController) QueryProductDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
	} else {
		productId := reqArgs["productId"]
		if productId == nil {
			response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
		} else {
			var product *Product
			productDetails, err := product.QueryProductDetails(uint64(productId.(float64)))
			if err == gorm.ErrRecordNotFound {
				response.HandleFail(RECORD_NOT_FOUND, ErrProductNotFound)
			} else if err != nil {
				logs.Error(err.Error())
				response.HandleError(err, RECORD_NOT_FOUND)
			} else {
				response.HandleSuccess(productDetails)
			}
		}
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}

// @Title Product - Specification details
// @Description Query product specification details
// @Params	inventoryId		 body		float64 	true	"Inventory's id"
// @Success	200000 {object} models.ResponseModel
// @Failure	200400 {object}	models.ResponseModel
// @router /recommend [post]
func (pd *ProductController) QuerySpecificationDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Warning(err.Error())
		response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
	} else {
		inventoryId := reqArgs["inventoryId"]
		if inventoryId == nil {
			response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
		} else {
			var product *Product
			productDetails, err := product.QueryInventoryDetails(uint64(inventoryId.(float64)))
			if err != nil && err != gorm.ErrRecordNotFound {
				logs.Error(err.Error())
				response.HandleError(err)
			} else {
				response.HandleSuccess(productDetails)
			}
		}
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}
