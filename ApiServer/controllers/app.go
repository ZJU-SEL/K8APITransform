package controllers

import (
	"K8APITransform/ApiServer/Fti"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/coreos/go-etcd/etcd"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

type AppController struct {
	beego.Controller
}

//// @Title list envinfo
//// @Description list envinfo

//// @router /list [get]
//func (a *AppController) ListInfo() {
//	username := a.Ctx.Request.Header.Get("Authorization")
//	cluster := a.GetString("target")
//	if cluster == "" {
//		response, err := models.EtcdClient.Get(path.Join(models.IpRoot, username), false, true)
//		if err != nil {
//			log.Println("Get User's Clusters Error :", err.Error())
//			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
//			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
//			return
//		}
//		log.Println("User :", username)
//		result := []map[string]string{}
//		for _, v := range response.Node.Nodes {
//			cluster = path.Base(v.Key)
//			log.Println("Cluster :", cluster)
//			envlist, err := models.NewUserClient(username).Envs(cluster).ListInfo()
//			if err != nil {
//				log.Printf("List User %v Cluster %v Error :%v\n", username, cluster, err.Error())
//				a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
//				http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
//				return
//			}
//			result = append(result, envlist...)
//			//models.NewUserClient(username).Envs(cluster).ListInfo()
//		}
//		a.Data["json"] = result
//		a.ServeJson()
//	} else {
//		envlist, err := models.NewUserClient(username).Envs(cluster).ListInfo()
//		if err != nil {
//			log.Printf("List User %v Cluster %v Error :%v\n", username, cluster, err.Error())
//			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
//			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
//			return
//		}
//		result := []map[string]string{}
//		result = append(result, envlist...)
//		//models.NewUserClient(username).Envs(cluster).ListInfo()

//		a.Data["json"] = result
//		a.ServeJson()
//	}
//}

// @Title list envinfo
// @Description list envinfo

// @router /list [get]
func (a *AppController) ListInfo2() {
	username := a.Ctx.Request.Header.Get("Authorization")
	cluster := a.GetString("target")
	if cluster == "" {
		response, err := models.EtcdClient.Get(path.Join(models.IpRoot, username), false, true)
		if err != nil {
			log.Println("Get User's Clusters Error :", err.Error())
			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
			return
		}
		log.Println("User :", username)
		result := []map[string]string{}
		var channel = []chan int{}
		mutex := &sync.Mutex{}
		for _, v := range response.Node.Nodes {
			c := make(chan int)
			channel = append(channel, c)
			go func(v *etcd.Node, c chan int) {
				cluster = path.Base(v.Key)
				log.Println("Cluster :", cluster)
				envlist, err := models.NewUserClient(username).Envs(cluster).ListInfo2()
				if err != nil {
					log.Printf("List User %v Cluster %v Error :%v\n", username, cluster, err.Error())
					a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
					http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
					return
				}
				mutex.Lock()
				result = append(result, envlist...)
				mutex.Unlock()
				close(c)
			}(v, c)
			//models.NewUserClient(username).Envs(cluster).ListInfo()
		}
		for _, v := range channel {
			<-v
		}
		a.Data["json"] = result
		a.ServeJson()
	} else {
		envlist, err := models.NewUserClient(username).Envs(cluster).ListInfo2()
		if err != nil {
			log.Printf("List User %v Cluster %v Error :%v\n", username, cluster, err.Error())
			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
			return
		}
		result := []map[string]string{}
		result = append(result, envlist...)
		//models.NewUserClient(username).Envs(cluster).ListInfo()

		a.Data["json"] = result
		a.ServeJson()
	}
}

// @Title details
// @Description details

// @router /details [get]
func (a *AppController) Details() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	log.Println("User :", username)
	log.Println("Cluster :", cluster)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envName := a.GetString("envName")
	appName := a.GetString("appName")
	fmt.Println(envName, appName)
	result, err := models.NewUserClient(username).Envs(cluster).Apps(envName).Get(appName)
	if err != nil {
		log.Printf("Get User %v Cluster %v Env %v App %v Details Error :%v\n", username, cluster, envName, appName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = result
	a.ServeJson()
}

// @Title monitor
// @Description monitor

// @router /monitor [get]
func (a *AppController) Monit() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)

	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envName := a.GetString("envName")
	appName := a.GetString("appName")
	id := a.GetString("id")
	flag := a.GetString("isFirst")
	result, err := models.NewUserClient(username).Envs(cluster).Apps(envName).Monit(appName, id, flag)
	if err != nil {
		log.Printf("Get User %v Cluster %v Env %v App %v Monit Error :%v\n", username, cluster, envName, appName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = result
	a.ServeJson()
}

// @Title status
// @Description status

// @router /status [get]
func (a *AppController) Status() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	envName := a.GetString("envName")
	appName := a.GetString("appName")
	//id := a.GetString("id")
	fmt.Println(envName, appName)
	detail, err := models.NewUserClient(username).Envs(cluster).Apps(envName).Get(appName)
	if err != nil {
		log.Printf("Get User %v Cluster %v Env %v App %v Status Error :%v\n", username, cluster, envName, appName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	result := map[string]string{"status": detail.Status}
	a.Data["json"] = result
	a.ServeJson()
}

// @Title details
// @Description details

// @router /deploy [post]
func (a *AppController) Deploy() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	deployReq := models.DeployRequest{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &deployReq)
	//log.Println(deployReq)
	if err != nil {
		log.Printf("Get User %v Cluster %v Deploy Request Error :%v\n", username, cluster, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = deployReq.Validate()
	if err != nil {
		log.Printf("Validate User %v Cluster %v Deploy Request Error :%v\n", username, cluster, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	appls := models.NewUserClient(username).Envs(cluster).Apps(deployReq.EnvName)
	if deployReq.IsGreyUpdating == "0" {
		label := map[string]string{}
		label["env"] = deployReq.EnvName
		err := appls.DeleteAll()
		if err != nil {
			log.Printf("Get User %v Cluster %v Env %v GreyUpdating Error :%v\n", username, cluster, deployReq.EnvName, err.Error())
			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
			return
		}
	}
	warName := deployReq.WarName
	//newimage_name_temp := []byte(uploadfilename)
	newimage_name := strings.TrimSuffix(warName, ".war") //(newimage_name_temp[0 : len(newimage_name_temp)-4])
	newimage_version := deployReq.Version
	newimage := newimage_name + "-" + newimage_version + ".war"
	log.Println("newimagename:", newimage)

	dockerdeamon := "http://" + ip + ":2376"
	imageprefix := username + "reg:5000"

	baseimage := imageprefix + `\/apm-jre7-tomcat7:v4`
	newimage, err = Fti.Wartoimage(dockerdeamon, deployReq.EnvName, imageprefix, username, baseimage, newimage, warName)

	if err != nil {
		log.Printf("User %v Cluster %v Env %v WarToimage Error :%v\n", username, cluster, deployReq.EnvName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	imagename := newimage
	log.Println("newimage:", imagename)
	app := &models.AppConfig{
		Name:    deployReq.WarName,
		Version: deployReq.Version,
		Ports: []models.Port{
			models.Port{
				Port:       8080,
				TargetPort: 8080,
				Protocol:   "TCP",
			},
		},
		ContainerPort: []models.Containerport{
			models.Containerport{
				Port:     8080,
				Protocol: "TCP",
			},
		},
		Containerimage: imagename,
	}
	err = appls.Create(app)
	if err != nil {
		log.Printf("User %v Cluster %v Env %v Create App Error :%v\n", username, cluster, deployReq.EnvName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	result, err := models.NewUserClient(username).Envs(cluster).GetInfo(deployReq.EnvName)
	if err != nil {
		log.Printf("Get User %v Cluster %v Env %v Info Error :%v\n", username, cluster, deployReq.EnvName, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = result
	a.ServeJson()
}

// @Title getservice
// @Description getservice
// @router /serviceip1/:podip1 [get]
func (a *AppController) Getseip1() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	podip := a.Ctx.Input.Param(":podip1")
	fmt.Println(podip)
	podip = strings.TrimSuffix(podip, ":8080")
	//todo:watch the etcd
	//seip := serviceipmap[podip]
	seip, err := models.GetPodtoSeFromEtcd(ip, podip)
	if err != nil {
		log.Println(err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = seip + ":8080"
	a.ServeJson()

	//return nil
}

// @Title getservice
// @Description getservice
// @router /serviceip/:podip [get]
func (a *AppController) Getseip() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	podip := a.Ctx.Input.Param(":podip")
	fmt.Println(podip)
	//todo:watch the etcd
	//seip := serviceipmap[podip]
	seip, err := models.GetPodtoSe(ip, podip)
	if err != nil {
		log.Println(err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = seip
	a.ServeJson()

	//return nil
}

// @Title stop app
// @Description stop app
// @router /stop [post]
func (a *AppController) Stop() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Stop(input["appName"])
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title start app
// @Description start app
// @router /start [post]
func (a *AppController) Start() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Start(input["appName"])
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title delete app
// @Description delete app
// @router /delete [delete]
func (a *AppController) DeleteApp() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Delete(input["appName"])
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title restart app
// @Description restart app
// @router /restart [post]
func (a *AppController) Restart() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Restart(input["appName"])
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title debug app
// @Description debug app
// @router /debug [post]
func (a *AppController) Debug() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Debug(input["appName"])
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title Scale app
// @Description Scale app
// @router /scale [put]
func (a *AppController) Scale() {
	ip := a.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	if err != nil {
		log.Printf("Get User name and cluster by ip  %v Error :%v\n", ip, err.Error())
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	var input = map[string]string{}
	err = json.Unmarshal(a.Ctx.Input.RequestBody, &input)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	replicas, err := strconv.Atoi(input["num"])
	if replicas == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+"num can not be zero"+`"}`, 406)
		return
	}
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.NewUserClient(username).Envs(cluster).Apps(input["envName"]).Scale(input["name"], replicas)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title upload war
// @Description upload
// @router /upload [post]
func (a *AppController) Upload() {
	username := a.Ctx.Request.Header.Get("Authorization")
	file, _, err := a.GetFile("file")
	version := a.GetString("version")
	appName := a.GetString("name")
	log.Println(version)
	app_part := strings.TrimSuffix(appName, ".war")
	appName_tmp := app_part + "-" + version + ".war"

	uploaddir := "applications/" + username + "/" + appName_tmp + "_deploy/"
	if Fti.Exist(uploaddir) {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+"version :"+version+" of war "+appName+" is exist"+`"}`, 406)
		return
	}
	Fti.Createdir(uploaddir)

	log.Println(version)
	if err != nil {
		os.RemoveAll(uploaddir)
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}

	log.Println(uploaddir)

	f, err := os.OpenFile(uploaddir+appName, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		os.RemoveAll(uploaddir)
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer f.Close()
	defer file.Close()
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title GetUploadWars
// @Description GetUploadWars

// @router /listWar [get]
func (a *AppController) Getuploadwars() {
	username := a.Ctx.Request.Header.Get("Authorization")
	dirhandle, err := os.Open("applications/" + username)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//log.Println(reflect.TypeOf(fis))
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	result := []interface{}{}
	//遍历文件列表 (no dir inside) 每一个文件到要写入一个新的*tar.Header
	//var fi os.FileInfo
	for _, fi := range fis {

		//如果是普通文件 直接写入 dir 后面已经有了/
		filename := fi.Name()

		log.Println(filename)
		fileinfo := strings.Split(filename, "_")
		if fileinfo[len(fileinfo)-1] == "deploy" {

			filename = strings.TrimSuffix(filename, "_deploy")
			filename = strings.TrimSuffix(filename, ".war")
			fileinfo = strings.Split(filename, "-")
			mapdata := map[string]interface{}{}
			mapdata["version"] = fileinfo[len(fileinfo)-1]
			//mapdata["uploadTime"] = fi.ModTime().Unix()
			name := strings.TrimSuffix(filename, "-"+mapdata["version"].(string)) + ".war"
			mapdata["name"] = name
			warfile, _ := os.Open("applications/" + username + "/" + fi.Name() + "/" + name)
			warfi, _ := warfile.Stat()
			mapdata["uploadTime"] = warfi.ModTime().Unix()
			//data := `{"id": 1,"name": "` + warname + `","nodeType": 0,"resource": [{"name": "app_version","value": "` + version + `"}]}`
			//json.Unmarshal([]byte(data), &mapdata)
			result = append(result, mapdata)
		}
		if err != nil {
			a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
			http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
			return
		}
	}
	a.Data["json"] = result
	a.ServeJson()
}

// @Title DeleteUploadWars
// @Description DeleteUploadWars

// @router /deleteWar [delete]
func (a *AppController) DeleteUploadwars() {
	username := a.Ctx.Request.Header.Get("Authorization")
	var warinfo = map[string]string{}

	err := json.Unmarshal(a.Ctx.Input.RequestBody, &warinfo)
	log.Println(warinfo)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	prefix := strings.TrimSuffix(warinfo["name"], ".war")
	err = os.RemoveAll("applications/" + username + "/" + prefix + "-" + warinfo["version"] + ".war" + "_deploy")
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = os.RemoveAll("applications/" + username + "/" + prefix + "-" + warinfo["version"] + ".war" + "_tar")
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}
