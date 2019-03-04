// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"ecommerce-sys/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	ns := beego.NewNamespace("/v1/api",
		//	api cache checking
		//beego.NSBefore(models.ReadApiCache),

		beego.NSNamespace("/user",
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
		),
	)

	beego.AddNamespace(ns)
}
