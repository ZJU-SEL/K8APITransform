package controllers

import (
	//"bytes"
	
    "fmt"
	"github.com/astaxie/beego/context"


)


func Baseimagelist(ctx *context.Context){
	fmt.Println("test base images")
	//output:="test base image list"
	ctx.Output.Body([]byte("test base image list"))
}
