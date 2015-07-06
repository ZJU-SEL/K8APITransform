package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"CreateEnv",
			`/createEnv`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Getuploadwars",
			`/getuploadwars`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Upload",
			`/upload`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Deploy",
			`/deploy`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"PartDetails",
			`/partDetails`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Details",
			`/details`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"RestartApp",
			`/restartApp`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"GetEnv",
			`/getEnv/:envname`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"DeleteEnv",
			`/deleteEnv/:envname`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Getpodsip",
			`/podsip/:sename`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Getseip",
			`/serviceip/:podip`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Scale",
			`/scaleApp`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"DeleteApp",
			`/delete`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppDeleteController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppDeleteController"],
		beego.ControllerComments{
			"Delete",
			`/`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppRollbackController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppRollbackController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppUpgradeController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppUpgradeController"],
		beego.ControllerComments{
			"Put",
			`/`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppViewController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppViewController"],
		beego.ControllerComments{
			"Get",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:BaseimageController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:BaseimageController"],
		beego.ControllerComments{
			"Baseimagelist",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:BaseimageController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:BaseimageController"],
		beego.ControllerComments{
			"Downloadimage",
			`/:imagename`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:LimitRangeController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NamespaceController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"],
		beego.ControllerComments{
			"Status",
			`/status`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"],
		beego.ControllerComments{
			"Addlabels",
			`/addlabels`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:NodesController"],
		beego.ControllerComments{
			"Getnodes",
			`/user/:username`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Get",
			`/:objectId`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Put",
			`/:objectId`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ObjectController"],
		beego.ControllerComments{
			"Delete",
			`/:objectId`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"],
		beego.ControllerComments{
			"Get",
			`/:name`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ResourceQuotaController"],
		beego.ControllerComments{
			"Delete",
			`/:name`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Post",
			`/`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"GetAll",
			`/`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Get",
			`/:uid`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Put",
			`/:uid`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Delete",
			`/:uid`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Login",
			`/login`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Auth",
			`/auth`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Logout",
			`/logout`,
			[]string{"get"},
			nil})

}
