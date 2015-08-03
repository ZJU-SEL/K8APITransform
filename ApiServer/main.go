package main

import (
	"K8APITransform/ApiServer/controllers"
	_ "K8APITransform/ApiServer/docs"
	"K8APITransform/ApiServer/models"
	_ "K8APITransform/ApiServer/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/coreos/go-etcd/etcd"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func init() {
	filedata, _ := ioutil.ReadFile("/etc/hosts")

	controllers.Hosts = strings.Split(string(filedata), string(10))
	log.Println(controllers.Hosts)
}
func main() {
	//beego.SessionOn = true
	//models.KubernetesIp = beego.AppConfig.String("k8sip")
	machines := beego.AppConfig.Strings("etcdmachines")
	serverCrt := beego.AppConfig.String("serverCrt")
	serverKey := beego.AppConfig.String("serverKey")
	rootCrt := beego.AppConfig.String("rootCrt")

	//controllers.DockerBuilddeamon = beego.AppConfig.String("DOCKER_BUILD_DEAMON")
	//models.ApiVersion = beego.AppConfig.String("APIVERSION")

	//fmt.Println("k8sip is ", models.KubernetesIp)
	//controllers.K8sBackend, _ = models.NewBackend(models.KubernetesIp+":8080", models.ApiVersion)
	//controllers.K8sBackend, _ = models.NewBackendTLS("https://k8master:8081", models.ApiVersion, "certs")

	Client, err := etcd.NewTLSClient(machines, serverCrt, serverKey, rootCrt)
	if err != nil {
		fmt.Println(err.Error())
	}
	models.EtcdClient = Client
	//models.IdPools = models.NewIdPools()
	//response := models.EtcdClient.CreateDir("/user")
	var UrlManager = func(ctx *context.Context) {
		//read urlMapping data from database
		//urlMapping := model.GetUrlMapping()
		fmt.Println(ctx.Request.RequestURI)
		flag := 0
		for baseurl, _ := range map[string]string{"/v1/user/checkuser": "POST"} {
			if baseurl == ctx.Request.RequestURI {
				flag = 1
				break

			}
		}
		if flag == 0 {
			ip := ctx.Request.Header.Get("Authorization")
			fmt.Println(ip)
			//_, ok := ctx.Input.Session("user").(string)
			if ip == "" {
				http.Error(ctx.ResponseWriter, "Authorization not set", 500)
			}
		}
	}
	beego.InsertFilter("/*", beego.BeforeRouter, UrlManager)
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
