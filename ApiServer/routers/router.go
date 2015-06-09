// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	//"K8APITransform/ApiServer/controllers"
	"K8APITransform/ApiServer/controllers"
	"github.com/astaxie/beego"
)

func init() {

	ns := beego.NewNamespace("v1",
		beego.NSNamespace("/application",
			beego.NSInclude(
				&controllers.AppController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		beego.NSNamespace("/s",
			beego.NSInclude(
				&controllers.BaseimageController{},
			),
		),
		beego.NSPost("/namespaces/:namespace/upload/:appname", controllers.Appuploadandtoimage),
		//beego.NSGet("/baseimage/search", controllers.Baseimagelist),
		//beego.NSGet("/baseimage/pull/:imagename", controllers.Baseimagepull),
	)
	//beego.Router("/application/v1", &controllers.AppController{})
	//beego.Router("/user", &controllers.UserController{})
	//beego.NSPost("/v1/namespaces/:namespace/upload/:appname", controllers.Appuploadandtoimage)
	beego.AddNamespace(ns)
}
