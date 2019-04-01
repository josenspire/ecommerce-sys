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
// @router	/create		[post]
func (or *OrderController) QueryOrders() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(or.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		logs.Error(err)
		response.HandleError(err, PARAMS_MISSING)
	} else {
		userId, orderType, pageIndex := reqArgs["userId"], reqArgs["orderType"], reqArgs["pageIndex"]
		if IsEmptyString(strconv.Itoa(int(userId.(float64)))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		} else {
			if IsEmptyString(orderType.(string)) {
				orderType = "ALL"
			}
			if IsEmptyString(strconv.Itoa(pageIndex.(int))) {
				pageIndex = 1
			}
			var order *OrderForm
			orders, err = order.QueryOrders(userId, orderType, pageIndex)
			if err == gorm.ErrRecordNotFound {

			} else if err != nil {
				response.HandleError(err)
			} else {
				response.HandleSuccess(orders)
			}
		}
	}
}
