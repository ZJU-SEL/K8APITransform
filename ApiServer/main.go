package main

import (
	"K8APITransform/ApiServer/controllers"
	_ "K8APITransform/ApiServer/docs"
	"K8APITransform/ApiServer/models"
	_ "K8APITransform/ApiServer/routers"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/go-etcd/etcd"
)

func main() {
	beego.SessionOn = true
	models.KubernetesIp = beego.AppConfig.String("k8sip")
	machines := beego.AppConfig.Strings("etcdmachines")
	serverCrt := beego.AppConfig.String("serverCrt")
	serverKey := beego.AppConfig.String("serverKey")
	rootCrt := beego.AppConfig.String("rootCrt")
	//fmt.Println("k8sip is ", models.KubernetesIp)
	controllers.K8sBackend, _ = models.NewBackend(models.KubernetesIp, "v1beta3")
	Client, err := etcd.NewTLSClient(machines, serverCrt, serverKey, rootCrt)
	if err != nil {
		fmt.Println(err.Error())
	}
	models.EtcdClient = Client
	models.IdPools = models.NewIdPools()
	//response := models.EtcdClient.CreateDir("/user")
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
