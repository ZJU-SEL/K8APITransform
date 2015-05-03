package main

import (
	"Fti/Ftitool"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io"
	"os"
	//"reflect"
)

type BuildImageOptions struct {
	Name                string             `qs:"t"`
	Dockerfile          string             `qs:"dockerfile"`
	NoCache             bool               `qs:"nocache"`
	SuppressOutput      bool               `qs:"q"`
	RmTmpContainer      bool               `qs:"rm"`
	ForceRmTmpContainer bool               `qs:"forcerm"`
	InputStream         io.Reader          `qs:"-"`
	OutputStream        io.Writer          `qs:"-"`
	RawJSONStream       bool               `qs:"-"`
	Remote              string             `qs:"remote"`
	Auth                AuthConfiguration  `qs:"-"` // for older docker X-Registry-Auth header
	AuthConfigs         AuthConfigurations `qs:"-"` // for newer docker X-Registry-Config header
	ContextDir          string             `qs:"-"`
}

type AuthConfiguration struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
}

type AuthConfigurations struct {
	Configs map[string]AuthConfiguration `json:"configs"`
}

// changing the parameter into the file name
// SourceTar produces a tar archive containing application source and stream it
// return type *os.File
func SourceTar(filename string) *os.File {
	//"tardir/deployments.tar.gz"
	fw, _ := os.Open(filename)
	//fmt.Println(reflect.TypeOf(fw))
	return fw

}

//把文件打包 之后 生成新的镜像
func main() {
	fti.Dirtotar()

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
		Username:      "wangzhe",
		Password:      "3.1415",
		Email:         "w_hessen@126.com",
		ServerAddress: "https://10.211.55.5",
	}

	opts := docker.BuildImageOptions{

		Name:         "mytest-tomcat",
		InputStream:  tarStream,
		OutputStream: os.Stdout,
		Auth:         auth,
		Dockerfile:   "Dockerfile",
	}

	fmt.Println(client.BuildImage(opts))
}
