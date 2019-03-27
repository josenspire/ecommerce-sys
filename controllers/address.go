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

type AddressController struct {
	beego.Controller
}

// @Title Create Address
// @Description Create user address
// @Param	models.AddressDTO		query	object		true	"Create a new address"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/create [post]
func (addr *AddressController) CreateAddress() {
	var response ResponseModel
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

// @Title Query Address Details
// @Description Query user address details
// @Param	userId		query	 float64		true	"user id"
// @Param	addressId		query	 float64		true	"address id"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/details [post]
func (addr *AddressController) QueryDetails() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleError(err)
	} else {
		addressId := reqArgs["addressId"].(float64)
		userId := reqArgs["userId"].(float64)
		if IsEmptyString(strconv.Itoa(int(addressId)), strconv.Itoa(int(userId))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		}
		var address *Address
		addressDetails, err := address.QueryAddressByAddressId(uint64(userId), uint64(addressId))
		if err == gorm.ErrRecordNotFound {
			response.HandleFail(RECORD_NOT_FOUND, ErrRecordNotFound.Error())
		} else if err != nil {
			logs.Error(err)
			response.HandleError(err)
		} else {
			response.HandleSuccess(&addressDetails)
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}

// @Title Update Address Details
// @Description Update user address details
// @Param	models.AddressDTO		query	object		true	"Create a new address"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/update [put]
func (addr *AddressController) UpdateAddress() {
	var response ResponseModel
	var dto *AddressDTO
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &dto)
	if err != nil {
		response.HandleError(err)
	} else {
		var address *Address
		err = address.UpdateAddress(dto)
		if err == gorm.ErrRecordNotFound {
			response.HandleFail(ADDRESS_NOT_FOUND, ErrAddressNotFound.Error())
		} else if err != nil {
			response.HandleError(err)
		} else {
			response.HandleSuccess(address)
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}

// @Title Delete Address Details
// @Description Delete user address details
// @Param	models.AddressDTO		query	object		true	"Delete a address"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/delete [delete]
func (addr *AddressController) DeleteAddress() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleError(err)
	} else {
		addressId := reqArgs["addressId"].(float64)
		userId := reqArgs["userId"].(float64)
		if IsEmptyString(strconv.Itoa(int(addressId)), strconv.Itoa(int(userId))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		}
		var address *Address
		err := address.DeleteAddressByAddressId(uint64(userId), uint64(addressId))
		if err == gorm.ErrRecordNotFound {
			response.HandleFail(RECORD_NOT_FOUND, ErrRecordNotFound.Error())
		} else if err != nil {
			logs.Error(err)
			response.HandleError(err)
		} else {
			response.HandleSuccess(nil, "address remove success")
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}

// @Title Set As Default Address
// @Description Set user default address
// @Param	userId		query		float64		true	"userId"
// @Param	addressId		query		float64		true	"addressId"
// @Success	200000	{object}	models.ResponseModel
// @Failure	200400
// @router	/delete [put]
func (addr *AddressController) SetAsDefaultAddress() {
	var response ResponseModel
	reqArgs := make(map[string]interface{})
	err := json.Unmarshal(addr.Ctx.Input.RequestBody, &reqArgs)
	if err != nil {
		response.HandleError(err)
	} else {
		userId := reqArgs["userId"].(float64)
		addressId := reqArgs["addressId"].(float64)
		if IsEmptyString(strconv.Itoa(int(addressId)), strconv.Itoa(int(userId))) {
			response.HandleFail(PARAMS_MISSING, ErrParamsMissing.Error())
		}
		var address *Address
		err := address.SetDefaultAddress(uint64(userId), uint64(addressId))
		if err == gorm.ErrRecordNotFound {
			response.HandleFail(RECORD_NOT_FOUND, ErrRecordNotFound.Error())
		} else if err != nil {
			logs.Error(err)
			response.HandleError(err)
		} else {
			response.HandleSuccess(nil, "set default address succeed")
		}
	}
	addr.Data["json"] = response
	addr.ServeJSON()
}
