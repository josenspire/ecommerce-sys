package controllers

import (
	. "ecommerce-sys/models"
	"encoding/json"
	"github.com/astaxie/beego"
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
