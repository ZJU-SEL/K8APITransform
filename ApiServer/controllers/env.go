package controllers

import (
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"net/http"
)

type EnvController struct {
	beego.Controller
}

// @Title CreateEnv
// @Description createEnv

// @router /createEnv [post]
func (e *EnvController) CreateEnv() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	var env models.AppEnv
	err := json.Unmarshal(e.Ctx.Input.RequestBody, &env)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = env.Validate()
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.AddAppEnv(ip, &env)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}

	e.Data["json"] = map[string]string{"msg": "success"}
	e.ServeJson()
}
