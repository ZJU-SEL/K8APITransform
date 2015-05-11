//wz
package controllers

import (
	//"bytes"
	"K8APITransform/K8APITransform/Fti/Ftitool"
	"fmt"
	"github.com/astaxie/beego/context"
	"io"
	"os"
)

//http://127.0.0.1:8080/v1/namespaces/test/upload/testwar
func Appupload(dirname string, ctx *context.Context) bool {

	r := ctx.Request
	w := ctx.ResponseWriter
	para := ctx.Input.Params
	namespace := para[":namespace"]
	appname := para[":appname"]

	fmt.Println("namespace:", namespace)
	fmt.Println("appname:", appname)

	output := "upload test \n" + "namespaces:" + string(para[":namespace"]) + " appname:" + string(para[":appname"])

	ctx.Output.Body([]byte(output))

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	//create the dir(without version info)
	//dirname := namespace + "_" + appname
	//add suffix automatically
	//the systempdir already exist
	//Ftitool.Createdir("systempdir")
	uploaddir := Ftitool.Createdir("systempdir/" + dirname + "_deploy")

	//默认的路径是从API server 这个目录下开始的
	f, err := os.OpenFile(uploaddir+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	fmt.Println("ok")
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer f.Close()
	io.Copy(f, file)
	return true

}

func Appuploadandtoimage(ctx *context.Context) {
	para := ctx.Input.Params
	namespace := para[":namespace"]
	appname := para[":appname"]
	imagename := namespace + "_" + appname
	uploaddirname := imagename
	uploadok := Appupload(uploaddirname, ctx)
	fmt.Println(uploadok)

	if uploadok {
		//read the war from the {namespace}_{appname}_deploy , add the dockerfile and tar them
		Ftitool.Wartoimage(imagename, uploaddirname)

	} else {
		fmt.Println("upload fail")
	}
}
