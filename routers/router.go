// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	. "ecommerce-sys/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &MainController{})

	ns := beego.NewNamespace("/v1/api",
		// 	api cache checking
		// beego.NSBefore(models.ReadApiCache),

		beego.NSNamespace("/advert",
			beego.NSRouter("/insert", &AdvertController{}, "post:InsertAdvert"),
			beego.NSRouter("/update", &AdvertController{}, "put:UpdateAdvert"),
			beego.NSRouter("/list", &AdvertController{}, "get:GetAdvertList"),
		),

		beego.NSNamespace("/user",
			beego.NSRouter("/register", &UserController{}, "post:Register"),
			beego.NSRouter("/loginByTelephone", &UserController{}, "post:LoginByTelephone"),
			beego.NSRouter("/loginByWechat", &UserController{}, "post:LoginByWechat"),
		),

		beego.NSNamespace("/address",
			beego.NSRouter("/create", &AddressController{}, "post:CreateAddress"),
			beego.NSRouter("/details", &AddressController{}, "post:QueryDetails"),
			beego.NSRouter("/update", &AddressController{}, "put:UpdateAddress"),
			beego.NSRouter("/delete", &AddressController{}, "delete:DeleteAddress"),
			beego.NSRouter("/setDefault", &AddressController{}, "put:SetAsDefaultAddress"),
		),

		beego.NSNamespace("/product",
			beego.NSRouter("/insert", &ProductController{}, "post:InsertProduct"),
			beego.NSRouter("/insertMultiple", &ProductController{}, "post:InsertMultipleProducts"),
			beego.NSRouter("/list", &ProductController{}, "get:QueryProducts"),
			beego.NSRouter("/details", &ProductController{}, "post:QueryProductDetails"),
		),
	)

	beego.AddNamespace(ns)
}
