package main

import (
	_ "K8APITransform/ApiServer/docs"
	"K8APITransform/ApiServer/models"
	_ "K8APITransform/ApiServer/routers"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/go-etcd/etcd"
)

var kubenets = flag.String(
	"kubenetsip",
	"",
	"base URL of the cloud controller",
)

func main() {

	//models.KubernetesIp = beego.AppConfig.String("k8sip")
	machines := beego.AppConfig.Strings("etcdmachines")
	//fmt.Println("k8sip is ", models.KubernetesIp)
	Client, err := etcd.NewTLSClient(machines, "/home/zjw/etcdkey/devregistry.crt", "/home/zjw/etcdkey/devregistry.key", "/home/zjw/etcdkey/rootca.crt")
	if err != nil {
		fmt.Println(err.Error())
	}
	models.EtcdClient = Client
	//response := models.EtcdClient.CreateDir("/user")
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
