//wz
package controllers

import (
	//"bytes"
	"fmt"
	"github.com/astaxie/beego/context"
	"io"
	"os"
)

//http://127.0.0.1:8080/v1/namespaces/test/upload/testwar
func Appupload(ctx *context.Context) {

	r := ctx.Request
	w := ctx.ResponseWriter
	para := ctx.Input.Params
	fmt.Println("namespace:", para[":namespace"])
	fmt.Println("warname:", para[":warname"])

	output := "upload test \n" + "namespaces:" + string(para[":namespace"]) + " warname:" + string(para[":warname"])

	ctx.Output.Body([]byte(output))

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile("uploadtest/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

}
