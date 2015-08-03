package controllers

import (
	"K8APITransform/ApiServer/models"
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego"
	//"github.com/bitly/go-simplejson"
	"io/ioutil"
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
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var env models.Env
	err = json.Unmarshal(e.Ctx.Input.RequestBody, &env)
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
	env.Target = cluster
	if env.Cpu == "" {
		env.Cpu = "1"
	}
	if env.Memory == "" {
		env.Memory = "0.5"
	}
	if env.Disk == "" {
		env.Disk = "0.5"
	}
	err = models.NewUserClient(username).Envs(cluster).Create(&env)
	//err = models.models.AddEnv(ip, &env)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}

	e.Data["json"] = map[string]string{"msg": "success"}
	e.ServeJson()
}

// @Title List Env
// @Description List Env

// @router /list [get]
func (e *EnvController) ListEnv() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envlist, err := models.NewUserClient(username).Envs(cluster).List()

	e.Data["json"] = envlist
	e.ServeJson()
}

// @Title List Info
// @Description List Info

// @router /listinfo [get]
func (e *EnvController) ListInfo() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envlist, err := models.NewUserClient(username).Envs(cluster).ListInfo()

	e.Data["json"] = envlist
	e.ServeJson()
}

// @Title List Info
// @Description List Info

// @router /getinfo [get]
func (e *EnvController) GetInfo() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envName := e.GetString("envName")
	envlist, err := models.NewUserClient(username).Envs(cluster).GetInfo(envName)
	e.Data["json"] = envlist
	e.ServeJson()
}

// @Title Delete Env
// @Description Delete Env

// @router /delete [delete]
func (e *EnvController) DeleteEnv() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	input := map[string]string{}
	err = json.Unmarshal(e.Ctx.Input.RequestBody, &input)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Delete(input["envName"])
	//err = models.DeleteEnv(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	e.Data["json"] = map[string]string{"msg": "success"}
	e.ServeJson()
}

// @Title Delete c
// @Description Delete All

// @router /deleteall [delete]
func (e *EnvController) DeleteAllEnv() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).DeleteAll()
	//err = models.DeleteEnv(ip)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	e.Data["json"] = map[string]string{"msg": "success"}
	e.ServeJson()
}

// @Title get clusterInfo
// @Description get clusterInfo

// @router /clusterInfo [get]
func (e *EnvController) ClusterInfo() {
	ip := e.Ctx.Request.Header.Get("Authorization")
	//flag := c.GetString("isFirst")
	username, cluster, err := models.Ip2UC(ip)
	client := models.NewMonitClient(username, cluster)
	request, err := http.NewRequest("GET", "https://"+ip+":50000/api/cluster/status", nil)
	request.Header.Set("token", "qwertyuiopasdfghjklzxcvbnm1234567890")
	//request.Header.Set("flag", flag)
	response, err := client.Do(request)
	if err != nil {
		e.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(e.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	var w = map[string]interface{}{}
	json.Unmarshal(body, &w)
	//w, err := simplejson.NewJson(body)
	//result := map[string]string{}
	//result["name"] = cluster
	//clusters, err := w.Get("Cluster").Map()
	//w = w.Get("Cluster")
	//var memoryusage uint64
	//var diskusage uint64
	//var memorycapacity uint64
	//var diskcapacity uint64
	//for k, _ := range clusters {
	//	disk, err := w.Get(k).Get("Diskcapacity").Uint64()
	//	diskcapacity += disk
	//	memory, err := w.Get(k).Get("Memorycapacity").Uint64()
	//	memorycapacity += memory
	//	memory, err = w.Get(k).Get("Spec").Get("MemoryAvg").Uint64()
	//	memoryusage += memory
	//	disk, err = w.Get(k).Get("Spec").Get("DiskAvg").Uint64()
	//	diskusage += disk
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	//result["memoryUsage"] = fmt.Sprintf("%v", memoryusage)
	//result["diskUsage"] = fmt.Sprintf("%v", diskusage)
	//result["memoryLimit"] = fmt.Sprintf("%v", memorycapacity)
	//result["diskLimit"] = fmt.Sprintf("%v", diskcapacity)

	e.Data["json"] = w
	e.ServeJson()
}
