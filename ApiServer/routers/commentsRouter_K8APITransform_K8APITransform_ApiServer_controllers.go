package routers

import (
	"github.com/astaxie/beego"
)

func init() {
	
	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppRollbackController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppRollbackController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppUpgradeController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppUpgradeController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppViewController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppViewController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Get",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Put",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Delete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Get",
			`/:uid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Put",
			`/:uid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Delete",
			`/:uid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Logout",
			`/logout`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Get",
			`/:service`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Delete",
			`/:service`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Upgrade",
			`/:service/upgrade`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Rollback",
			`/:service/rollback`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppDeleteController"] = append(beego.GlobalControllerRouter["K8APITransform/K8APITransform/ApiServer/controllers:AppDeleteController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

}
