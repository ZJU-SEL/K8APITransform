package main

import (
	_ "K8APITransform/ApiServer/docs"
	_ "K8APITransform/ApiServer/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.RunMode == "dev" {
		beego.DirectoryIndex = true
		beego.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
