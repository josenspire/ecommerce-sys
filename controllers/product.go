package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"strconv"
)

type ProductController struct {
	beego.Controller
}

func (pd *ProductController) InsertProduct() {
	var response ResponseModel
	reqArgs := make(map[string]ProductDTO)
	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleError(err)
	} else {
		dto := reqArgs["product"]
		if dto == (ProductDTO{}) {
			response.HandleFail(PARAMS_MISSING, ErrCreateRecordsIsEmpty.Error())
		} else {
			var product *Product
			err := product.InsertProduct(&dto)
			if err != nil {
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
		response.HandleError(err)
	} else {
		dtos := reqArgs["products"]
		if len(dtos) == 0 {
			response.HandleFail(PARAMS_MISSING, ErrCreateRecordsIsEmpty.Error())
		} else {
			var product *Product
			err := product.InsertMultipleProducts(&dtos)
			if err != nil {
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
// @Params	productType		Query	true	"product type, default is 'recommend'"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/recommend 	[get]
func (pd *ProductController) QueryProducts() {
	var response ResponseModel
	var productType = pd.Input().Get("type")
	var pageIndexStr = pd.Input().Get("page")
	if IsEmptyString(productType) {
		productType = "normal"
	}
	pageIndexInt, err := strconv.Atoi(pageIndexStr)
	if err != nil || pageIndexInt <= 0 {
		pageIndexInt = 1
	}
	var products *[]Product
	var product *Product
	products, err = product.QueryProductsByProductType(productType, pageIndexInt)
	if err != nil && err != gorm.ErrRecordNotFound {
		response.HandleError(err)
	} else {
		response.HandleSuccess(products)
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}

// @Title Product - details
// @Description Query product details
// @Params	productId		Query		float64  	true		"product id"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/recommend 	[get]
func (pd *ProductController) QueryProductDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})

	fmt.Println(reqArgs)

	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
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
// @Params	inventoryId		 Query		float64 	true		"inventoryId"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/recommend 	[get]
func (pd *ProductController) QuerySpecificationDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})

	fmt.Println(reqArgs)

	err := json.Unmarshal(pd.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
	} else {
		inventoryId := reqArgs["inventoryId"]
		if inventoryId == nil {
			response.HandleFail(PARAMS_MISSING, ErrParamsInValid.Error())
		} else {
			var product *Product
			productDetails, err := product.QueryInventoryDetails(uint64(inventoryId.(float64)))
			if err != nil && err != gorm.ErrRecordNotFound {
				response.HandleError(err)
			} else {
				response.HandleSuccess(productDetails)
			}
		}
	}
	pd.Data["json"] = response
	pd.ServeJSON()
}
