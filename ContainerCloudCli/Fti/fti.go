package Fti

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	//	"io/ioutil"
	"github.com/fsouza/go-dockerclient"
	"os"
	//"reflect"
	"bytes"
	"log"
	"os/exec"
)

func Filecompress(tw *tar.Writer, dir string, fi os.FileInfo) {

	//打开文件 open当中是 目录名称/文件名称 构成的组合
	fmt.Println(dir + fi.Name())
	fr, err := os.Open(dir + fi.Name())
	fmt.Println(fr.Name())
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	hdr, err := tar.FileInfoHeader(fi, "")

	hdr.Name = fr.Name()
	if err = tw.WriteHeader(hdr); err != nil {
		panic(err)
	}

	_, err = io.Copy(tw, fr)
	if err != nil {
		panic(err)
	}
	//打印文件名称
	fmt.Println("add the file: " + fi.Name())

}

func Dircompress(tw *tar.Writer, dir string) {

	//打开文件夹
	dirhandle, err := os.Open(dir)
	//fmt.Println(dir.Name())
	//fmt.Println(reflect.TypeOf(dir))
	if err != nil {
		panic(nil)
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//fmt.Println(reflect.TypeOf(fis))
	if err != nil {
		panic(err)
	}

	//遍历文件列表 每一个文件到要写入一个新的*tar.Header
	//var fi os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {

			//			//如果再加上这段的内容 就会多生成一层目录
			//			hdr, err := tar.FileInfoHeader(fi, "")
			//			if err != nil {
			//				panic(err)
			//			}
			//			hdr.Name = fi.Name()
			//			err = tw.WriteHeader(hdr)
			//			if err != nil {
			//				panic(err)
			//			}

			newname := dir + fi.Name()
			fmt.Println("using dir")
			fmt.Println(newname)
			//这个样直接continue就将所有文件写入到了一起 没有层级结构了
			//Filecompress(tw, dir, fi)
			Dircompress(tw, newname+"/")

		} else {
			//如果是普通文件 直接写入 dir 后面已经有了 /
			Filecompress(tw, dir, fi)
		}

	}

}

func Dirtotar(sourcedir string, tardir string, newimage string) {
	//file write 在tardir目录下创建
	_, err := os.Stat(sourcedir)
	if err != nil {
		fmt.Println("please create the deploy dir")
		return
	}
	fw, err := os.Create(tardir + "/" + newimage + ".tar.gz")
	//type of fw is *os.File
	//	fmt.Println(reflect.TypeOf(fw))
	if err != nil {
		panic(err)

	}
	defer fw.Close()

	//gzip writer
	gw := gzip.NewWriter(fw)
	defer gw.Close()

	//tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()
	//	fmt.Println(reflect.TypeOf(tw))
	//add the deployments contens
	//Dircompress(tw, "deployments/")
	Dircompress(tw, sourcedir+"/")
	//	// add the dockerfile
	//	fr, err := os.Open("Dockerfile")

	//do not package the dockerfile individual
	//write into the dockerfile
	//fileinfo, err := os.Stat(tardir + "/" + newimage + "_Dockerfile")
	//fileinfo, err := os.Stat("./Dockerfile")
	//if err != nil {
	//	panic(err)

	//}
	//fmt.Println(reflect.TypeOf(os.FileInfo(fileinfo)))
	//dockerfile要单独放在根目录下 和其他archivefile并列
	//Filecompress(tw, "", fileinfo)

	fmt.Println("tar.gz packaging OK")

}

//return a tar stream
func SourceTar(filename string) *os.File {
	//"tardir/deployments.tar.gz"
	fw, _ := os.Open(filename)
	//fmt.Println(reflect.TypeOf(fw))
	return fw

}

func Systemexec(s string) {
	cmd := exec.Command("/bin/sh", "-c", s)
	fmt.Println(s)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", out.String())
}

func Createdockerfile(baseimage string, newimage string) {
	targetDocker := newimage + "_deploy/" + newimage + "_Dockerfile"
	_, err := os.Stat(targetDocker)
	if err == nil {
		os.Remove(targetDocker)
	}
	//recreate the file
	dst, err := os.OpenFile(targetDocker, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("please create the <newimagename>_deploy dir and add war file into it")
		return
	}
	src, err := os.Open("./Dockerfile")
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(dst, src)

	if err != nil {
		panic(err)
	}

	defer src.Close()
	defer dst.Close()

	modifybase := `sed -i "s/baseimage/` + baseimage + `/g"` + " ./" + targetDocker

	Systemexec(modifybase)
	modifynew := `sed -i "s/newimage/` + newimage + `/g"` + " ./" + targetDocker

	Systemexec(modifynew)
	//modify the docker file
}

//the image will be covered if the image already exist
func Wartoimage(dockerdeamon string, baseimage string, newimage string) {
	// put the war file into the _deploy dir
	sourcedir := newimage + "_deploy"
	//put the baseimage_Dockerfile and the deployments.tar.gz into the baseimage_tar
	tardir := newimage + "_tar"

	//upload the war file from remote server to the deploy dir and add some scripts
	//todo: add a rest api which could receive the tar file and put the war file into the _deploy dir
	//a war->tar->war add scripts（such as dockerfile） -> tar -> image
	//Createdir(deploydir)
	fmt.Println(tardir)
	Createdir(tardir)
	defer os.RemoveAll(tardir)
	//delete the temp dir at last
	//defer Cleandir(imagename)

	//create the dockerfile according to the baseimage in the baseimage_tar
	Createdockerfile(baseimage, newimage)

	Dirtotar(sourcedir, tardir, newimage)

	//using go-docker client
	endpoint := dockerdeamon
	client, _ := docker.NewClient(endpoint)
	//fmt.Println(client)
	filename := tardir + "/" + newimage + ".tar.gz"
	//filename := "tardir/Dockerfile"
	tarStream := SourceTar(filename)
	defer tarStream.Close()
	fmt.Println(tarStream)
	//  test the basic using
	//	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	//	for _, img := range imgs {
	//		fmt.Println("ID: ", img.ID)
	//		fmt.Println("RepoTags: ", img.RepoTags)
	//		fmt.Println("Created: ", img.Created)
	//		fmt.Println("Size: ", img.Size)
	//		fmt.Println("VirtualSize: ", img.VirtualSize)
	//		fmt.Println("ParentId: ", img.ParentID)
	//	}

	//dockerhub的认证信息
	auth := docker.AuthConfiguration{
	//	Username:      "wangzhe",
	//	Password:      "3.1415",
	//	Email:         "w_hessen@126.com",
	//	ServerAddress: "https://10.211.55.5",
	}

	opts := docker.BuildImageOptions{

		Name:         newimage,
		InputStream:  tarStream,
		OutputStream: os.Stdout,
		Auth:         auth,
		Dockerfile:   newimage + "_deploy/" + newimage + "_Dockerfile",
	}

	//error
	error := client.BuildImage(opts)
	if error != nil {
		panic(error)

	}
	//return error
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//create the temp dir and return this dir
//input image name
func Createdir(imagename string) string {

	//if the file already exist , delete and recreate it
	exist := Exist(imagename)

	if exist {
		fmt.Println("the folder exist , remove it")
		Cleandir(imagename)
	}
	dirname := imagename
	err := os.Mkdir(dirname, 0777)
	if err != nil {
		panic(err)
	}
	fmt.Println("create succesful: " + dirname)
	return dirname

}

//递归删除文件夹
func Cleandir(dirname string) {

	//打开文件夹
	dirhandle, err := os.Open(dirname)
	//fmt.Println(dirname)
	//fmt.Println(reflect.TypeOf(dir))
	if err != nil {
		panic(nil)
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//fmt.Println(reflect.TypeOf(fis))
	if err != nil {
		panic(err)
	}

	//遍历文件列表 每一个文件到要写入一个新的*tar.Header
	//var fi os.FileInfo
	for _, fi := range fis {
		if fi.IsDir() {
			newname := dirname + "/" + fi.Name()
			//fmt.Println("using dir")
			//fmt.Println(newname)
			//这个样直接continue就将所有文件写入到了一起 没有层级结构了
			//Filecompress(tw, dir, fi)
			Cleandir(newname)

		} else {
			//如果是普通文件 直接写入 dir 后面已经有了 /
			filename := dirname + "/" + fi.Name()
			fmt.Println(filename)
			err := os.Remove(filename)
			if err != nil {
				panic(err)
			}
			fmt.Println("delete " + filename)
		}

	}
	//递归结束 删除当前文件夹
	err = os.Remove(dirname)
	fmt.Println("delete " + dirname)
	if err != nil {
		panic(err)
	}

}
