package controllers

import (
	//"K8APITransform/K8APITransform/ApiServer/models"
	//"encoding/json"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/astaxie/beego"
)

// Operations about Users
type AppViewController struct {
	beego.Controller
}

// @Title createApp
// @Description create app
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [get]
func (v *AppViewController) Get() {

	//json.Unmarshal(u.Ctx.Input.RequestBody, &app)

	v.Data["json"] = map[string]string{"status": "success"}
	v.ServeJson()
}
