package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/jinzhu/gorm"
	"strconv"
)

type OrderController struct {
	beego.Controller
}

// @Title Order List
// @Description Query user all order records
// @Param	userId			query		float64		true	"User Id"
// @Param	orderType		query		string		true	"Order type, default is 'all'"
// @Param	pageIndex		query		uint		false	"Page index, default is: 1"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/list		[post]
func (or *OrderController) QueryOrders() {
	// TODO waiting for testing
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(or.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		userId := reqArgs["userId"].(float64)
		orderType := reqArgs["orderType"].(string)
		pageIndex := int(reqArgs["pageIndex"].(float64))
		if IsEmptyString(strconv.Itoa(int(userId))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		} else {
			if IsEmptyString(orderType) {
				orderType = "ALL"
			}
			if IsEmptyString(strconv.Itoa(pageIndex)) {
				pageIndex = 1
			}
			var order *OrderForm
			orders, err := order.QueryOrders(uint64(userId), orderType, pageIndex)
			if err != nil && err != gorm.ErrRecordNotFound {
				response.HandleError(err)
			} else {
				response.HandleSuccess(orders)
			}
		}
	}
	or.Data["json"] = response
	or.ServeJSON()
}

// @Title PlaceOrder
// @Description Place a new order
// @Param	userId			query		float64			true	"User Id"
// @Param	addressId		query		float64			true	"Address Id"
// @Param	orders			query		interface		true	"Order array details"
// @Param	discount		query		string			false	"Total discount"
// @Param	remark			query		string			false	"Order remark"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/place		[post]
func (or *OrderController) PlaceOrder() {
	var response ResponseModel
	var dto *PlaceOrderDTO
	err := json.Unmarshal(or.Ctx.Input.RequestBody, &dto)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		var order *OrderForm
		err := order.PlaceOrder(dto)
		if err != nil {
			logs.Error(err)
			response.HandleError(err)
		} else {
			response.HandleSuccess(nil, "Placing order succeed")
		}
	}
	or.Data["json"] = response
	or.ServeJSON()
}

// @Title OrderCompleted
// @Description Order is already completed
// @Param	userId			query		float64			true	"User Id"
// @Param	orderId			query		float64			true	"Order Id"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/completed		[post]
func (or *OrderController) OrderCompleted() {
	var response ResponseModel
	reqArgs := make(map[string]float64)
	err := json.Unmarshal(or.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		userId := int(reqArgs["userId"])
		orderId := int(reqArgs["orderId"])
		if IsEmptyString(strconv.Itoa(userId), strconv.Itoa(orderId)) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		} else {
			var order *OrderForm
			err := order.OrderCompleted(uint64(userId), uint64(orderId))
			if err == gorm.ErrRecordNotFound {
				response.HandleFail(ORDER_NOT_FOUND, ErrOrderNotFound)
			} else if err != nil {
				logs.Error(err)
				response.HandleError(err)
			} else {
				response.HandleSuccess(nil, "Order is already completed")
			}
		}
	}
	or.Data["json"] = response
	or.ServeJSON()
}
