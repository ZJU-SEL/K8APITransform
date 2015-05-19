package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"os"
	"os/exec"
	//"path"
	"path/filepath"
)

// Operations about App
type BaseimageController struct {
	beego.Controller
}

var imageslice []string

// @Title get image list
// @Description create app
// @Param	body		body 	null	 true		"body for user content"
// @Success 200 {json} "image list"
// @Failure 403 body is empty
// @router / [get]
func (images *BaseimageController) Baseimagelist() {
	fmt.Println("test base images")
	//output:="test base image list"

	//local path under APIServer
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fmt.Println(path)

	//fmt.Println()
	images.Scandir("/home/zjw/base_image/")
	fmt.Println(imageslice)
	image_json, _ := json.Marshal(imageslice)
	images.Ctx.Output.Body(image_json)
}
func (images *BaseimageController) Scandir(dirname string) {
	imageslice = []string{}
	fmt.Println("add image list")

	//打开文件夹
	dirhandle, err := os.Open(dirname)
	//fmt.Println(dirname)
	//fmt.Println(reflect.TypeOf(dirhandle))
	if err != nil {
		panic(err)
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//fmt.Println(reflect.TypeOf(fis))
	if err != nil {
		panic(err)
	}

	//遍历文件列表 (no dir inside) 每一个文件到要写入一个新的*tar.Header
	//var fi os.FileInfo
	for _, fi := range fis {

		//如果是普通文件 直接写入 dir 后面已经有了/
		filename := dirname + fi.Name()
		fmt.Println(filename)
		//err := os.Remove(filename)
		//temp_image:=&image{imagename:fi.Name(),}
		imageslice = append(imageslice, string(fi.Name()))
		if err != nil {
			panic(err)
		}
	}

}

// @Title get image
// @Description create app
// @Param	body		body 	null	 true		"body for user content"
// @Success 200 {json} "image matedata"
// @Failure 403 body is empty
// @router /:imagename [get]
func (images *BaseimageController) Downloadimage() {
	imagename := images.Ctx.Input.Param(":imagename")
	imagepath := "/home/zjw/base_image/" + imagename
	images.Ctx.ResponseWriter.Header().Set("Content-Type", "file")
	//images.
	http.ServeFile(images.Ctx.ResponseWriter, images.Ctx.Request, imagepath)
}
