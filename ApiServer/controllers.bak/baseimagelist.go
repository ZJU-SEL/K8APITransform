//wz
//copy the base_image dir into the /controller/base_image
package controllers

import (
	//"bytes"

	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/context"
	"os"
	"os/exec"
	"path/filepath"
)

//type image struct{
//	imagename   string
//}
//	var imageslice []string

//search the base_image dir and add the info into the base_image
func Scandir(dirname string) {
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

	//fmt.Println(imageslice)
}

func Baseimagelist(ctx *context.Context) {
	fmt.Println("test base images")
	//output:="test base image list"

	//local path under APIServer
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fmt.Println(path)

	fmt.Println()
	Scandir("./controllers/base_image/")
	fmt.Println(imageslice)
	image_json, _ := json.Marshal(imageslice)
	ctx.Output.Body(image_json)

}
