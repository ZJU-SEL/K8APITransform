package controllers

import (
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	Hosts []string
)

type UserController struct {
	beego.Controller
}

// @Title checkUser
// @Description checkuser and store the ca.crt
// @router /checkuser [post]
func (a *UserController) Checkuser() {
	var clusterinfo = map[string]string{}
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &clusterinfo)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	username := clusterinfo["userName"]
	cluster := clusterinfo["cloudName"]
	masterip := clusterinfo["masterIp"]
	cafile := clusterinfo["cacrt"]
	data := []byte(cafile)
	filename := path.Join(models.CertRoot, username, cluster, "ca.crt")
	os.MkdirAll(path.Join(models.CertRoot, username, cluster), 0777)
	err = ioutil.WriteFile(filename, data, 0666)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	//for _, v := range Hosts {
	//	fmt.Println(v)
	//}
	ret := []string{}
	for _, v := range Hosts {
		if v != "" && !strings.Contains(v, cluster+"."+username) && !strings.Contains(v, masterip) {
			ret = append(ret, v)
		}
	}
	//fmt.Println(strings.Join(ret, ",\n"))
	ret = append(ret, masterip+" "+cluster+"."+username)
	back := strings.Join(ret, "\n")
	err = ioutil.WriteFile("/etc/hosts", []byte(back), 0777)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	models.EtcdClient.Set(path.Join(models.HostRoot, masterip), cluster+"."+username, 0)
	models.EtcdClient.Set(path.Join("hosttoip", username, cluster), masterip, 0)
	models.EtcdClient.SetDir(path.Join(models.EnvRoot, username, cluster), 0)
	Hosts = ret
	a.Ctx.ResponseWriter.WriteHeader(200)
}
