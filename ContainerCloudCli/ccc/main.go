package main

import (
	"ContainerCloudCli/lib"
	"ContainerCloudCli/models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
)

const (
	serverip = "10.10.105.112"
)

func sendGet(host string, port string, version string, getcommands []string) ([]byte, int) {
	url := "http://" + host + ":" + port + "/" + version

	for _, str := range getcommands {
		url = url + "/" + str

	}

	fmt.Println("send request:" + url)

	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	response, _ := client.Do(reqest)

	//body为 []byte类型
	body, _ := ioutil.ReadAll(response.Body)
	status := response.StatusCode

	return body, status
}

func sendDelete(host string, port string, version string, getcommands []string) ([]byte, int) {
	url := "http://" + host + ":" + port + "/" + version

	for _, str := range getcommands {
		url = url + "/" + str

	}

	fmt.Println("send request:" + url)

	client := &http.Client{}
	reqest, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		panic(err.Error())
	}

	response, _ := client.Do(reqest)

	//body为 []byte类型
	body, _ := ioutil.ReadAll(response.Body)
	status := response.StatusCode

	return body, status
}

func sendPost(host string, port string, version string, getcommands []string, filename string) ([]byte, int) {
	url := "http://" + host + ":" + port + "/" + version
	for _, str := range getcommands {
		url = url + "/" + str
	}
	fmt.Println("send post request:" + url)
	client := &http.Client{}
	byt, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err.Error())
	}

	reqest, _ := http.NewRequest("POST", url, bytes.NewBuffer(byt))
	reqest.Header.Set("Content-Type", "application/json")
	response, _ := client.Do(reqest)

	//body为 []byte类型
	body, _ := ioutil.ReadAll(response.Body)
	status := response.StatusCode

	return body, status
}

func Scandir(dirname string, imageslice []string) {
	//imageslice=[]string{}
	//fmt.Println("add image list")

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
		//filename := dirname + fi.Name()
		//fmt.Println(filename)
		//err := os.Remove(filename)
		//temp_image:=&image{imagename:fi.Name(),}
		if strings.Contains(fi.Name(), ".tar") {
			imageslice = append(imageslice, string(fi.Name()))
			if err != nil {
				panic(err)
			}
		}

	}

	fmt.Println(imageslice)
}

func newCmdList() *cobra.Command {
	var listtype string
	Listcmd := &cobra.Command{
		Use:   "list",
		Short: "list the image name which could be used on the server",
		Long:  "list the image name which could be used on the server ... detail info",
		Run: func(cmd *cobra.Command, args []string) {

			if status := Auth(); status == 200 {
				//fmt.Println(strings.EqualFold(listtype, "server"))
				//fmt.Println("listtype:", listtype)
				if strings.EqualFold(listtype, "server") {
					//	namespace := "localnamespace"
					getcommands := []string{"baseimages"}
					responsebody, status := sendGet(serverip, "8080", "v1", getcommands)

					if status == 200 {
						fmt.Println("the avaliable image in server")
						fmt.Println(string(responsebody))
					} else {
						fmt.Println("error")
					}
				} else if strings.EqualFold(listtype, "local") {
					fmt.Println("serch the base image in local repo")

					//file, _ := exec.LookPath(os.Args[0])
					//path, _ := filepath.Abs(file)
					//fmt.Println(path)

					var imageslice []string
					Scandir("./base_image_repo", imageslice)
					//fmt.Println(imageslice)

				} else {
					fmt.Println("error in location")
				}
			} else {
				fmt.Println("auth err")
			}
		},
	}

	Listcmd.Flags().StringVarP(&listtype, "location", "l", "server", "list the aviliable base image in \"local\" or \"server\"")
	return Listcmd
}

func systemexec(s string) {
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

func newCmdPull() *cobra.Command {

	Pullcmd := &cobra.Command{
		Use:   "pull [image name]",
		Short: "pull the base image on the server",
		Long:  `pull the base image on the server and store it detail info...`,
		Run: func(cmd *cobra.Command, args []string) {
			if status := Auth(); status == 200 {
				fmt.Println("test pull")
				// 获取命令行输入参数 命令行参数被自动存在args中 第一个输入的值 是镜像的名字
				//	namespace := "localnamespace"

				for _, value := range args {
					fmt.Println("input args:" + value)
				}
				base_dir := "./base_image_repo/"
				//检查这个镜像是否已经存在在本地
				_, err := os.Stat(base_dir + args[0])
				var status int
				var responsebody []byte
				//  var responsebody http*re
				if err == nil {
					fmt.Println("image file already exist ")

				} else {
					getcommands := []string{"baseimages", args[0]}
					responsebody, status = sendGet(serverip, "8080", "v1", getcommands)
					//download镜像

					if status == 200 {
						fmt.Println(base_dir + args[0])
						ioutil.WriteFile(base_dir+args[0], responsebody, 0666)

					} else {
						fmt.Println("err")
					}
				}

				//把镜像解压开 在本地（minion节点上 生成image）
				//systemexec("cd base_image_repo/")
				//systemexec("pwd")
				//sudo docker load < base_image_repo/
				systemexec("sudo docker load  < " + "./base_image_repo/" + args[0])
			} else {
				fmt.Println("auth err")
			}
		},
	}

	return Pullcmd

}

func newCmdInfo() *cobra.Command {
	Infocmd := &cobra.Command{
		Use:   "info",
		Short: "show the info running in server",
		Long:  `show the info running in server details...`,
		Run: func(cmd *cobra.Command, args []string) {
			if status := Auth(); status == 200 {
				fmt.Println("test info")
				//send get api
				getcommands := []string{"namespaces", "default", "services", args[0], "state"}
				responsebody, status := sendGet(serverip, "8080", "v1", getcommands)
				fmt.Println(string(responsebody), status)
			} else {
				fmt.Println("auth err")
			}
		},
	}
	return Infocmd

}

func newCmdDelete() *cobra.Command {
	Deletecmd := &cobra.Command{
		Use:   "delete",
		Short: "delete the services running in server",
		Long:  `delete the services running in server details...`,
		Run: func(cmd *cobra.Command, args []string) {
			if status := Auth(); status == 200 {
				fmt.Println("test delete")
				//send get api
				getcommands := []string{"namespaces", "default", "services", args[0]}
				responsebody, status := sendDelete(serverip, "8080", "v1", getcommands)
				fmt.Println(string(responsebody), status)
			} else {
				fmt.Println("auth err")
			}
		},
	}
	return Deletecmd

}

func newCmdLogin() *cobra.Command {

	var (
		name     string
		password string
	)
	logincmd := &cobra.Command{
		Use:   "login",
		Short: "login to the server",
		Long:  `login to the server detail info...`,
		Run: func(cmd *cobra.Command, args []string) {
			userinfo := &models.UserInfo{
				Username: name,
				Password: password,
			}
			body, _ := json.Marshal(userinfo)
			status, result := lib.Sendapi("POST", serverip, "8080", []string{"v1", "user", "login"}, body)
			//fmt.Println(status, string(result))
			if status == 200 {
				user, _ := user.Current()
				Dir := user.HomeDir + "/.blackPaaS/"
				file := Dir + "/config.json"
				os.MkdirAll(Dir, 0777)
				os.Create(file)
				ioutil.WriteFile(file, []byte(strings.Split(string(result), "@")[0]), 0666)
				fmt.Println(strings.Split(string(result), "@")[1])
			} else {
				fmt.Println("login error")
			}
		},
	}
	logincmd.Flags().StringVarP(&name, "name", "n", "", "name")
	logincmd.Flags().StringVarP(&password, "password", "p", "", "password")
	return logincmd
}

func Auth() int {
	user, _ := user.Current()
	Dir := user.HomeDir + "/.blackPaaS/"
	file := Dir + "/config.json"
	body, _ := ioutil.ReadFile(file)

	status, _ := lib.Sendapi("POST", serverip, "8080", []string{"v1", "user", "auth"}, body)
	return status
}

func newCmdStart() *cobra.Command {
	Startcmd := &cobra.Command{
		Use:   "start [image name]",
		Short: "start the base image on the server",
		Long:  `start the base image with the local volum files`,
		Run: func(cmd *cobra.Command, args []string) {
			//the first arg is the name of image
			if status := Auth(); status == 200 {
				fmt.Println("test start")
				if len(args) == 0 {
					fmt.Println("please input the image name")

				} else if len(args) == 1 {
					fmt.Println("please input the service name")
				} else {
					//modify the .startconfig change the containerimage tobe the arg[0] attention to user ` `

					modifyimage := `sed -i "s/\"containerimage\":.*/\"containerimage\": \"` + args[0] + `\",/g" .startconfig`
					//modify:=`sed -i "s/\"containerimage\":.*/\"containerimage\": \"test3\"/g" .startconfig`
					systemexec(modifyimage)

					modifyservice := `sed -i "s/\"name\":.*/\"name\": \"` + args[1] + `\",/g" .startconfig`
					systemexec(modifyservice)

					getcommands := []string{"namespaces", "default", "services"}
					responsebody, status := sendPost(serverip, "8080", "v1", getcommands, "./.startconfig")
					fmt.Println(string(responsebody), status)
				}
			} else {
				fmt.Println("auth err")
			}
		},
	}
	return Startcmd
}

//short 会显示比较简单的内容
//Long  在具体help的时候会显示对应的内容
func main() {

	CCCCmd := &cobra.Command{
		Use:  "ccc",
		Long: "Container-Cloud-Cli (STI) is a tool for runnig your app based on k8s.\n\n",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	//add persistance source
	//stiCmd.PersistentFlags().StringVarP(&(req.DockerSocket), "url", "U", dockerSocket(), "Set the url of the docker socket to use")
	//CCCCmd.AddCommand(newCmdVersion())

	var echoTimes int

	var cmdTimes = &cobra.Command{
		Use:   "times [# times] [string to echo]",
		Short: "Echo anything to the screen more times",
		Long: `echo things multiple times back to the user by providing
	        a count and a string.`,
		Run: func(cmd *cobra.Command, args []string) {
			for i := 0; i < echoTimes; i++ {
				fmt.Println("Echo: " + strings.Join(args, " "))
			}
		},
	}

	//use to show how to use tags intvarp 表示 参数的类型
	//后面的参数分别是 给哪个命令赋予flag flag全称 flag缩写 flag默认值 以及输出的具体帮助信息
	cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

	//	var CCCCmd = &cobra.Command{Use: "ccc"}
	//CCCCmd.AddCommand(cmdTimes)
	CCCCmd.AddCommand(newCmdLogin(), newCmdList(), newCmdPull(), newCmdStart(), newCmdInfo(), newCmdDelete())
	//cmdEcho.AddCommand(cmdTimes)
	CCCCmd.Execute()
}
