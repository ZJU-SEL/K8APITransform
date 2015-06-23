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
	//"bytes"
	"bufio"
	"log"
	"os/exec"
	"strings"
)

const (
	applications = "applications"
)

func Filecompress(tw *tar.Writer, dir string, fi os.FileInfo) error {

	//打开文件 open当中是 目录名称/文件名称 构成的组合
	fmt.Println(dir + fi.Name())
	fr, err := os.Open(dir + fi.Name())
	fmt.Println(fr.Name())
	if err != nil {
		return err
	}
	defer fr.Close()

	hdr, err := tar.FileInfoHeader(fi, "")

	hdr.Name = fr.Name()
	if err = tw.WriteHeader(hdr); err != nil {
		return err
	}

	_, err = io.Copy(tw, fr)
	if err != nil {
		return err
	}
	//打印文件名称
	fmt.Println("add the file: " + fi.Name())
	return nil

}

func Dircompress(tw *tar.Writer, dir string) error {

	//打开文件夹
	dirhandle, err := os.Open(dir)
	//fmt.Println(dir.Name())
	//fmt.Println(reflect.TypeOf(dir))
	if err != nil {
		return err
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//fmt.Println(reflect.TypeOf(fis))
	if err != nil {
		return err
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
			err = Dircompress(tw, newname+"/")
			if err != nil {
				return err
			}

		} else {
			//如果是普通文件 直接写入 dir 后面已经有了 /
			err = Filecompress(tw, dir, fi)
			if err != nil {
				return err
			}
		}

	}
	return nil

}

func Dirtotar(sourcedir string, tardir string, newimage string) error {
	//file write 在tardir目录下创建
	_, err := os.Stat(sourcedir)
	if err != nil {
		fmt.Println("please create the deploy dir")
		return err
	}
	fw, err := os.Create(tardir + "/" + newimage + ".tar.gz")
	//type of fw is *os.File
	//	fmt.Println(reflect.TypeOf(fw))
	if err != nil {
		return err

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
	err = Dircompress(tw, sourcedir+"/")
	if err != nil {
		return err
	}
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
	return nil

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
	out, err := cmd.StdoutPipe()
	go func() {
		o := bufio.NewReader(out)
		for {
			line, _, err := o.ReadLine()
			if err == io.EOF {
				break
			} else {
				fmt.Println(string(line))
			}
		}
	}()
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func Createdockerfile(username string, baseimage string, newimage string, warName string) error {
	targetDocker := applications + "/" + username + "/" + newimage + "_deploy" + "/" + "Dockerfile"
	fmt.Println("tardocker:", targetDocker)
	_, err := os.Stat(targetDocker)
	if err == nil {
		os.Remove(targetDocker)
	}
	//recreate the file
	dst, err := os.OpenFile(targetDocker, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	src, err := os.Open(applications + "/Dockerfile")
	if err != nil {
		return err
	}

	_, err = io.Copy(dst, src)

	if err != nil {
		return err
	}

	defer src.Close()
	defer dst.Close()

	modifybase := `sed -i "s/baseimage/` + baseimage + `/g" ` + targetDocker

	Systemexec(modifybase)
	//modifynew := `sed -i "s/newimage/` + username + "/" + newimage + "_deploy /" + `/g" ` + targetDocker
	modifynew := `sed -i "s/newimage/` + applications + `\/cxy\/` + newimage + `_deploy\/` + warName + ` /g" ` + targetDocker

	//newimage + "_deploy" + "/"

	Systemexec(modifynew)
	//modify the docker file
	return nil
}

//the image will be covered if the image already exist
//dockerdeamon, username,imageprefix, baseimage, newimage
func Wartoimage(dockerdeamon string, imageprefix string, username string, baseimage string, newimage string, warName string) (string, error) {
	// put the war file into the _deploy dir

	sourcedir := applications + "/" + username + "/" + newimage + "_deploy"
	//put the baseimage_Dockerfile and the deployments.tar.gz into the baseimage_tar
	tardir := applications + "/" + username + "/" + newimage + "_tar"

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
	//the template is located in the dir username
	err := Createdockerfile(username, baseimage, newimage, warName)
	if err != nil {
		return "", err
	}

	err = Dirtotar(sourcedir, tardir, newimage)
	if err != nil {
		return "", err
	}

	//using go-docker client
	endpoint := dockerdeamon
	client, _ := docker.NewClient(endpoint)
	//fmt.Println(client)
	filename := tardir + "/" + newimage + ".tar.gz"
	//filename := "tardir/Dockerfile"
	tarStream := SourceTar(filename)
	defer tarStream.Close()
	fmt.Println(tarStream)

	//dockerhub的认证信息
	auth := docker.AuthConfiguration{
	//	Username:      "",
	//	Password:      "",
	//	Email:         "w_hessen@126.com",
	//	ServerAddress: "https://10.211.55.5",
	}
	opts := docker.BuildImageOptions{

		Name:         imageprefix + "/" + strings.ToLower(newimage),
		InputStream:  tarStream,
		OutputStream: os.Stdout,
		Auth:         auth,
		Dockerfile:   applications + "/" + username + "/" + newimage + "_deploy" + "/" + "Dockerfile",
	}

	//error
	err = client.BuildImage(opts)
	if err != nil {
		return "", err

	}

	pushopts := docker.PushImageOptions{
		Name:         newimage,
		Tag:          "latest",
		Registry:     imageprefix,
		OutputStream: os.Stdout,
	}

	client.PushImage(pushopts, auth)

	//send the image to the private registry
	//pushcommand := `docker push ` + imageprefix + "/" + strings.ToLower(newimage)
	//Systemexec(pushcommand)
	return strings.ToLower(newimage), nil
}

// 检查文件或目录是否存在
// 如果由 filename 指定的文件或目录存在则返回 true，否则返回 false
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

//create the temp dir and return this dir
//input image name
func Createdir(imagename string) (string, error) {

	//if the file already exist , delete and recreate it
	exist := Exist(imagename)

	if exist {
		fmt.Println("the folder exist , remove it")
		Cleandir(imagename)
	}
	dirname := imagename
	err := os.MkdirAll(dirname, 0777)
	if err != nil {
		return "", err
	}
	fmt.Println("create succesful: " + dirname)
	return dirname, nil

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
