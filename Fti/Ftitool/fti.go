package fti

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	//	"io/ioutil"
	"github.com/fsouza/go-dockerclient"
	"os"
	//"reflect"
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
	//bad way
	//	//信息头部 生成tar文件的时候要先写入tar结构体
	//	h := new(tar.Header)
	//	//fmt.Println(reflect.TypeOf(h))

	//	h.Name = fi.Name()
	//	h.Size = fi.Size()
	//	h.Mode = int64(fi.Mode())
	//	h.ModTime = fi.ModTime()

	//	//将信息头部的内容写入
	//	err = tw.WriteHeader(h)
	//	if err != nil {
	//		panic(err)
	//	}

	//copy(dst Writer,src Reader)
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

func Dirtotar() {
	//file write 在tardir目录下创建
	fw, err := os.Create("tardir/deployments.tar.gz")
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
	Dircompress(tw, "deployments/")
	//	// add the dockerfile
	//	fr, err := os.Open("Dockerfile")

	//write into the dockerfile
	fileinfo, err := os.Stat("Dockerfile")
	if err != nil {
		panic(err)

	}
	//fmt.Println(reflect.TypeOf(os.FileInfo(fileinfo)))
	//dockerfile要单独放在根目录下 和其他archivefile并列
	Filecompress(tw, "", fileinfo)

	fmt.Println("tar.gz packaging OK")

}

//return a tar stream
func SourceTar(filename string) *os.File {
	//"tardir/deployments.tar.gz"
	fw, _ := os.Open(filename)
	//fmt.Println(reflect.TypeOf(fw))
	return fw

}

//the image will be covered if the image already exist
func Wartoimage(imagename string) error {
	Dirtotar()

	//using go-docker client
	endpoint := "http://10.211.55.5:2375"
	client, _ := docker.NewClient(endpoint)
	//fmt.Println(client)
	filename := "tardir/deployments.tar.gz"
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

		Name:         imagename,
		InputStream:  tarStream,
		OutputStream: os.Stdout,
		Auth:         auth,
		Dockerfile:   "Dockerfile",
	}

	//error
	error := client.BuildImage(opts)
	if error != nil {
		panic(error)

	}
	return error

}
