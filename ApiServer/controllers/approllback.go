package controllers

import (
	//"K8APITransform/K8APITransform/ApiServer/models"
	//"encoding/json"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/astaxie/beego"
)

// Operations about Users
type AppRollbackController struct {
	beego.Controller
}

// @Title createApp
// @Description create app
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [put]
func (r *AppRollbackController) Put() {

	//json.Unmarshal(u.Ctx.Input.RequestBody, &app)

	r.Data["json"] = map[string]string{"status": "rollback success"}
	r.ServeJson()
}
