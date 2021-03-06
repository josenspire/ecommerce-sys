// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	. "ecommerce-sys/commons"
	. "ecommerce-sys/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &MainController{})

	var aspectControl IAspectControl = &AspectControlWithoutEcdh{}

	beego.InsertFilter("/*", beego.BeforeExec, aspectControl.HandleRequestWithoutEcdh)
	beego.InsertFilter("/*", beego.FinishRouter, aspectControl.HandleResponseWithoutEcdh)

	ns := beego.NewNamespace("/v1/api",
		// 	api cache checking
		// beego.NSBefore(models.ReadApiCache),

		/**
		* Note: swagger api doc, just support to `NSNamespace + NSInclude`, or will not work in others way
		 */
		beego.NSNamespace("/advert",
			beego.NSRouter("/insert", &AdvertController{}, "post:InsertAdvert"),
			beego.NSRouter("/update", &AdvertController{}, "put:UpdateAdvert"),
			beego.NSRouter("/list", &AdvertController{}, "post:GetAdvertList"),
		),

		beego.NSNamespace("/auth",
			beego.NSRouter("/obtainSecurityCode", &SMSController{}, "post:ObtainSecurityCode"),
			beego.NSRouter("/verifySecurityCode", &SMSController{}, "post:VerifySecurityCode"),
		),

		beego.NSNamespace("/user",
			beego.NSRouter("/register", &UserController{}, "post:Register"),
			beego.NSRouter("/loginByTelephone", &UserController{}, "post:LoginByTelephone"),
			beego.NSRouter("/loginByWechat", &UserController{}, "post:LoginByWechat"),
			beego.NSRouter("/teams", &UserController{}, "post:QueryUserTeams"),
		),

		beego.NSNamespace("/address",
			beego.NSRouter("/list", &AddressController{}, "post:QueryAddresses"),
			beego.NSRouter("/create", &AddressController{}, "post:CreateAddress"),
			beego.NSRouter("/details", &AddressController{}, "post:QueryDetails"),
			beego.NSRouter("/update", &AddressController{}, "put:UpdateAddress"),
			beego.NSRouter("/delete", &AddressController{}, "delete:DeleteAddress"),
			beego.NSRouter("/setDefault", &AddressController{}, "put:SetAsDefaultAddress"),
		),

		beego.NSNamespace("/product",
			beego.NSRouter("/insert", &ProductController{}, "post:InsertProduct"),
			beego.NSRouter("/insertMultiple", &ProductController{}, "post:InsertMultipleProducts"),
			beego.NSRouter("/list", &ProductController{}, "post:QueryProducts"),
			beego.NSRouter("/details", &ProductController{}, "post:QueryProductDetails"),
			beego.NSRouter("/details/specification", &ProductController{}, "post:QuerySpecificationDetails"),
		),

		beego.NSNamespace("/classify",
			beego.NSRouter("/create", &ClassifyController{}, "post:CreateClassify"),
			beego.NSRouter("/category/create", &ClassifyController{}, "post:CreateCategory"),
			beego.NSRouter("/list", &ClassifyController{}, "post:QueryClassifies"),
		),

		beego.NSNamespace("/order",
			beego.NSRouter("/list", &OrderController{}, "post:QueryOrders"),
			beego.NSRouter("/place", &OrderController{}, "post:PlaceOrder"),
			beego.NSRouter("/completed", &OrderController{}, "put:OrderCompleted"),
			beego.NSRouter("/cancel", &OrderController{}, "put:OrderCancel"),
			beego.NSRouter("/details", &OrderController{}, "post:QueryProductDetails"),
		),
	)
	beego.AddNamespace(ns)
}
