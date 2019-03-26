package controllers

import (
	. "ecommerce-sys/models"
	. "ecommerce-sys/utils"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"strconv"
)

type AddressController struct {
	beego.Controller
}

// @Title Create Address
// @Description Create user address
// @Param	telephone	query	string	true	"Create a new address"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/create [post]
func (addr *AddressController) CreateAddress() {
	var response ResponseModel
	// reqParams := make(map[string]string)
	var dto *AddressDTO
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &dto)
	if err != nil {
		response.HandleError(err)
	} else {
		var address *Address
		err = address.CreateAddress(dto)
		if err != nil {
			response.HandleError(err)
		} else {
			response.HandleSuccess(address)
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}

func (addr *AddressController) QueryDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleError(err)
	} else {
		addressId := reqArgs["addressId"].(uint64)
		userId := reqArgs["userId"].(uint64)
		if IsEmptyString(strconv.Itoa(int(addressId)), strconv.Itoa(int(userId))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		}
		var address *Address
		addressDetails, err := address.QueryAddressByAddressId(userId, addressId)
		if err != nil {
			logs.Error(err)
			response.HandleError(err)
		} else {
			response.HandleSuccess(&addressDetails)
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}
