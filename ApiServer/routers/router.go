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
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		//beego.NSNamespace("/namespaces",
		//	beego.NSNamespace("/:namespaces",
		//		beego.NSNamespace("/resourcequotas",
		//			beego.NSInclude(
		//				&controllers.ResourceQuotaController{},
		//			),
		//		),
		//		beego.NSNamespace("/limitranges",
		//			beego.NSInclude(
		//				&controllers.LimitRangeController{},
		//			),
		//		),
		//		beego.NSNamespace("/services",
		//			beego.NSInclude(
		//				&controllers.AppController{},
		//			),
		//		),
		//	),
		//	beego.NSInclude(
		//		&controllers.NamespaceController{},
		//	),
		//),
		//beego.NSNamespace("/namespaces",
		//	beego.NSInclude(
		//		&controllers.NamespaceController{},
		//	),
		//),
		//beego.NSNamespace("/namespaces/:namespaces/limitranges",
		//	beego.NSInclude(
		//		&controllers.LimitRangeController{},
		//	),
		//),
		//beego.NSNamespace("/namespaces/:namespaces/resourcequotas",
		//	beego.NSInclude(
		//		&controllers.ResourceQuotaController{},
		//	),
		//),
		beego.NSNamespace("/namespaces/:namespaces/services",
			beego.NSInclude(
				&controllers.AppController{},
			),
		),
		beego.NSNamespace("/baseimages",
			beego.NSInclude(
				&controllers.BaseimageController{},
			),
		),
		beego.NSPost("/namespaces/:namespace/upload/:appname", controllers.Appuploadandtoimage),
		//beego.NSGet("/baseimage/search", controllers.Baseimagelist),
		//beego.NSGet("/baseimage/pull/:imagename", controllers.Baseimagepull),
	)
	beego.AddNamespace(ns)
}
