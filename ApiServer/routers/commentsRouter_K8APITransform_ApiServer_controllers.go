package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"ListInfo2",
			`/list`,
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
			"Monit",
			`/monitor`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Status",
			`/status`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Deploy",
			`/deploy`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Getseip1",
			`/serviceip1/:podip1`,
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
			"Stop",
			`/stop`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Start",
			`/start`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"DeleteApp",
			`/delete`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Restart",
			`/restart`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"CloseDebugApp",
			`/closedebugApp`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Debug",
			`/debug`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Scale",
			`/scale`,
			[]string{"put"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Upload",
			`/upload`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"Getuploadwars",
			`/listWar`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:AppController"],
		beego.ControllerComments{
			"DeleteUploadwars",
			`/deleteWar`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"],
		beego.ControllerComments{
			"Nodes",
			`/nodes`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"],
		beego.ControllerComments{
			"NodeInfo",
			`/nodeInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:ClusterController"],
		beego.ControllerComments{
			"ClusterInfo",
			`/clusterInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"CreateEnv",
			`/createEnv`,
			[]string{"post"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"ListEnv",
			`/list`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"ListInfo",
			`/listinfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"GetInfo",
			`/getinfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"DeleteEnv",
			`/delete`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"DeleteAllEnv",
			`/deleteall`,
			[]string{"delete"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:EnvController"],
		beego.ControllerComments{
			"ClusterInfo",
			`/clusterInfo`,
			[]string{"get"},
			nil})

	beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"] = append(beego.GlobalControllerRouter["K8APITransform/ApiServer/controllers:UserController"],
		beego.ControllerComments{
			"Checkuser",
			`/checkuser`,
			[]string{"post"},
			nil})

}
