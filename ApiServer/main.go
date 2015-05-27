package main

import (
	_ "K8APITransform/ApiServer/docs"
	"K8APITransform/ApiServer/models"
	_ "K8APITransform/ApiServer/routers"
	"flag"
	"fmt"
	"github.com/astaxie/beego"
)

var kubenets = flag.String(
	"kubenetsip",
	"",
	"base URL of the cloud controller",
)

func main() {
	flag.Parse()
	fmt.Println(*kubenets)
	models.KubenetesIp = *kubenets
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
