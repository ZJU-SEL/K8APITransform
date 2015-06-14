package controllers

import (
	"K8APITransform/ApiServer/Fti"
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/astaxie/beego"
	"io"
	"net/http"
	"os"

	"path"
	//"time"
)

// Operations about App
type AppController struct {
	beego.Controller
}

// IntstrKind represents the stored type of IntOrString.
func NewIntOrStringFromInt(val int) models.IntOrString {
	return models.IntOrString{Kind: models.IntstrInt, IntVal: val}
}

// @Title CreateEnv
// @Description createEnv

// @router /createEnv [post]
func (a *AppController) CreateEnv() {
	var env models.AppEnv
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &env)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = env.Validate()
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = models.AddAppEnv(&env)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	detail := &models.Detail{Name: env.Name, Status: 1, NodeType: 1, Context: []models.Detail{}, Children: []models.Detail{}}
	detail.Children = append(detail.Children, models.Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []models.Detail{
			models.Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	num, err := strconv.Atoi(env.NodeNum)
	tomcat := models.Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []models.Detail{}, Children: []models.Detail{}}
	for i := 1; i <= num; i++ {
		tomcat.Context = append(tomcat.Context, models.Detail{
			Name:     "Node" + strconv.Itoa(i),
			NodeType: 3,
		})
	}
	detail.Children = append(detail.Children, tomcat)
	a.Data["json"] = detail
	a.ServeJson()
}

// @Title GetUploadWars
// @Description GetUploadWars

// @router /getuploadwars [get]
func (a *AppController) Getuploadwars() {
	username := "cxy"
	dirhandle, err := os.Open("applications/" + username)
	//fmt.Println(dirname)
	//fmt.Println(reflect.TypeOf(dirhandle))
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer dirhandle.Close()

	//fis, err := ioutil.ReadDir(dir)
	fis, err := dirhandle.Readdir(0)
	//fis的类型为 []os.FileInfo
	//fmt.Println(reflect.TypeOf(fis))
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
		fmt.Println(filename)
		fileinfo := strings.Split(filename, "_")
		if fileinfo[len(fileinfo)-1] == "deploy" {
			filename = strings.TrimRight(filename, "_deploy")
			filename = strings.TrimRight(filename, ".war")
			fileinfo = strings.Split(filename, "-")
			version := fileinfo[len(fileinfo)-1]
			warname := strings.TrimRight(filename, "-"+version) + ".war"
			data := `{"id": 1,"name": "` + warname + `","nodeType": 0,"resource": [{"name": "app_version","value": "` + version + `"}]}`
			mapdata := map[string]interface{}{}
			json.Unmarshal([]byte(data), &mapdata)
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

// @Title upload war
// @Description upload

// @router /upload [post]
func (a *AppController) Upload() {
	//a.ParseForm()
	file, _, err := a.GetFile("filePath")
	version := a.GetString("version")
	appName := a.GetString("appName")
	fmt.Println(version)
	date := []byte(appName)
	date = date[0 : len(date)-4]
	//todo :use regx
	app_part := string(date)
	appName = app_part + "-" + version + ".war"

	username := "cxy"
	//uploaddir := "applications/" + username + "/" + appName + "-" + version + "_deploy/"
	uploaddir := "applications/" + username + "/" + appName + "_deploy/"
	Fti.Createdir(uploaddir)
	//version := a.GetString("version")

	fmt.Println(version)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	//f, err := os.OpenFile("applications/"+handle.Filename+version, os.O_WRONLY|os.O_CREATE, 0666)
	fmt.Println(uploaddir)
	//f, err := os.OpenFile(uploaddir+appName+"-"+version, os.O_WRONLY|os.O_CREATE, 0666)
	f, err := os.OpenFile(uploaddir+appName, os.O_WRONLY|os.O_CREATE, 0666)
	io.Copy(f, file)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer f.Close()
	defer file.Close()
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title deploy
// @Description deploy

// @router /deploy [post]
func (a *AppController) Deploy() {
	namespace := "default"
	deployReq := models.DeployRequest{}
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &deployReq)
	fmt.Println(deployReq)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	err = deployReq.Validate()
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	env, err := models.GetAppEnv(deployReq.EnvName)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}

	uploadfilename := deployReq.WarName
	//var uploadfilename string
	//if deployReq.AppVersion != "" {
	//	uploadfilename = deployReq.WarName + "-" + deployReq.AppVersion
	//} else {
	//	uploadfilename = deployReq.WarName
	//}

	username := "cxy"
	//newimage := uploadfilename

	//newimage_part := strings.Split(uploadfilename, "-")[0]
	if deployReq.IsGreyUpdating == "0" {
		//namespace := "default"
		url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/services" + "?labelSelector=env%3D" + deployReq.EnvName
		//fmt.Println(url)
		rsp, _ := http.Get(url)
		var serviceslist models.ServiceList
		body, _ := ioutil.ReadAll(rsp.Body)
		json.Unmarshal(body, &serviceslist)
		for _, v := range serviceslist.Items {
			a.deleteapp(v.ObjectMeta.Labels["name"])
		}

	}
	newimage_name_temp := []byte(uploadfilename)
	newimage_name := string(newimage_name_temp[0 : len(newimage_name_temp)-4])
	//newimage_name := strings.Split(newimage_part, ".")[0]

	newimage_version := deployReq.AppVersion
	newimage := newimage_name + "-" + newimage_version + ".war"

	fmt.Println("newimagename:", newimage)
	//deployReq imagename string, uploaddir string) error
	dockerdeamon := "unix:///var/run/docker.sock"
	//dockerdeamon := "http://10.211.55.10:2376"

	imageprefix := username + "reg:5000"

	//deployReq imagename string, uploaddir string) error
	//dockerdeamon := "unix:///var/run/docker.sock"
	baseimage := "jre7" + "-" + "tomcat7"
	//baseimage = env.JdkV + "-" + env.TomcatV
	//baseimage := "jre" + strconv(env.JdkV) + "-" + "tomcat" + strconv(env.TomcatV)
	newimage, err = Fti.Wartoimage(dockerdeamon, imageprefix, username, baseimage, newimage)

	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	imagename := imageprefix + "/" + newimage
	fmt.Println("newimage:", imagename)

	//imagename := "7-jre-customize"
	//wartoimage
	//uploadimage
	//createapplication imagename = ""
	replicas, err := strconv.Atoi(env.NodeNum)
	app := &models.AppCreateRequest{
		Name:    env.Name,
		Version: deployReq.WarName + "-" + deployReq.AppVersion,
		Ports: []models.Port{
			models.Port{
				Port:       8080,
				TargetPort: 8080,
				Protocol:   "tcp",
			},
		},
		Replicas: replicas,
		ContainerPort: []models.Containerport{
			models.Containerport{
				Port:     8080,
				Protocol: "tcp",
			},
		},
		Containername:  env.Name,
		Containerimage: imagename,
	}
	service, err := a.CreateApp(app)
	if err != nil {
		a.deleteapp(app.Name + "-" + app.Version)
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	env.Used++
	err = models.UpdateAppEnv(env.Name, env)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/pods" + "?labelSelector=name%3D" + app.Name + "-" + app.Version
	//fmt.Println(url)
	rsp, _ := http.Get(url)
	var podslist models.PodList
	body, _ := ioutil.ReadAll(rsp.Body)
	json.Unmarshal(body, &podslist)
	fmt.Println(string(body))
	detail := &models.Detail{Name: env.Name, Status: 1, NodeType: 1, Context: []models.Detail{}, Children: []models.Detail{}}
	detail.Children = append(detail.Children, models.Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []models.Detail{
			models.Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := models.Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []models.Detail{}, Children: []models.Detail{}}
	for k, v := range podslist.Items {
		status := 0
		if v.Status.Phase == models.PodRunning {
			status = 1
		}

		tomcat.Context = append(tomcat.Context, models.Detail{
			Name: "Node" + strconv.Itoa(k+1),
			//AppVersion: v.ObjectMeta.Labels["name"],
			AppVersion: deployReq.AppVersion,
			Status:     status,
			NodeType:   3,
		})
	}
	tomcat.Children = append(tomcat.Children, models.Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []models.Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, models.Detail{
		Name:         service.ObjectMeta.Labels["name"],
		NodeType:     4,
		Status:       1,
		Resource:     []models.Detail{models.Detail{Name: "IP", Value: service.Spec.PortalIP + ":8080"}},
		OriginalName: deployReq.WarName,
	})
	detail.Children = append(detail.Children, tomcat)
	a.Data["json"] = detail
	a.ServeJson()
}

// @Title get partDetails
// @Description get partDetails

// @router /partDetails [get]
func (a *AppController) PartDetails() {
	//a.ParseForm()
	envName := a.GetString("envName")
	fmt.Println(envName)
	env, err := models.GetAppEnv(envName)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	detail := a.getdetails(env)
	a.Data["json"] = detail
	a.ServeJson()

}

func (a *AppController) getdetails(env *models.AppEnv) *models.Detail {
	namespace := "default"
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/services" + "?labelSelector=env%3D" + env.Name
	//fmt.Println(url)
	rsp, _ := http.Get(url)
	var serviceslist models.ServiceList
	body, _ := ioutil.ReadAll(rsp.Body)
	json.Unmarshal(body, &serviceslist)

	url = "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/pods" + "?labelSelector=env%3D" + env.Name
	//fmt.Println(url)
	rsp, _ = http.Get(url)
	var podslist models.PodList
	body, _ = ioutil.ReadAll(rsp.Body)
	json.Unmarshal(body, &podslist)

	detail := &models.Detail{Name: env.Name, Status: 1, NodeType: 1, Context: []models.Detail{}, Children: []models.Detail{}}
	detail.Children = append(detail.Children, models.Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []models.Detail{
			models.Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := models.Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []models.Detail{}, Children: []models.Detail{}}
	if len(podslist.Items) == 0 {
		num, _ := strconv.Atoi(env.NodeNum)
		for k := 0; k < num; k++ {
			//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
			tomcat.Context = append(tomcat.Context, models.Detail{
				Name:     "Node" + strconv.Itoa(k+1),
				NodeType: 3,
			})
		}
	} else {
		for k, v := range podslist.Items {
			status := 0
			if v.Status.Phase == models.PodRunning {
				status = 1
			}
			//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
			tomcat.Context = append(tomcat.Context, models.Detail{
				Name:       "Node" + strconv.Itoa(k+1),
				AppVersion: v.ObjectMeta.Labels["name"],
				Status:     status,
				NodeType:   3,
			})
		}
	}
	apps := []models.Detail{}
	for _, v := range serviceslist.Items {
		//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
		apps = append(apps, models.Detail{
			Name:     v.ObjectMeta.Labels["name"],
			NodeType: 4,
			Status:   1,
			Resource: []models.Detail{models.Detail{Name: "IP", Value: v.Spec.PortalIP + ":8080"}},
		})
	}
	tomcat.Children = append(tomcat.Children, models.Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []models.Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, apps...)
	detail.Children = append(detail.Children, tomcat)
	return detail
}

// @Title get Details
// @Description get Details

// @router /details [get]
func (a *AppController) Details() {
	envs, err := models.GetAllAppEnv()
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	response := []*models.Detail{}
	for _, v := range envs {
		response = append(response, a.getdetails(v))
	}
	a.Data["json"] = response
	a.ServeJson()
}

// @Title restartApp
// @Description restartApp

// @router /restartApp [post]
func (a *AppController) RestartApp() {
	namespace := "default"
	req := map[string]string{}
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &req)
	//fmt.Println(deployReq)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	fmt.Println(req["appName"])
	if _, exist := req["appName"]; exist == false {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+"request not has appName "+`"}`, 406)
		return
	}
	appName := req["appName"]
	//appName := app.Name + "-" + app.Version
	appName = strings.ToLower(appName)
	appName = strings.Replace(appName, ".", "", -1)
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers/" + appName
	rsp, _ := http.Get(url)
	var rc models.ReplicationController
	body, _ := ioutil.ReadAll(rsp.Body)
	err = json.Unmarshal(body, &rc)
	fmt.Println(string(body))
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	replicas := rc.Spec.Replicas
	rc.Spec.Replicas = 0
	body, _ = json.Marshal(rc)
	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", rc.ObjectMeta.Name}, body)
	fmt.Println(status)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, result)
		return
	}
	body, _ = json.Marshal(result)
	json.Unmarshal(body, &rc)
	rc.Spec.Replicas = replicas
	body, _ = json.Marshal(rc)
	status, result = lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", rc.ObjectMeta.Name}, body)
	fmt.Println(status)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, result)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title getEnv
// @Description getEnv

// @router /getEnv/:envname [get]
func (a *AppController) GetEnv() {
	name := a.Ctx.Input.Param(":envname")
	env, err := models.GetAppEnv(name)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = env
	a.ServeJson()
}

// @Title deleteEnv
// @Description deleteEnv

// @router /deleteEnv/:envname [delete]
func (a *AppController) DeleteEnv() {
	name := a.Ctx.Input.Param(":envname")
	err := models.DeleteAppEnv(name)
	//env, err := models.GetAppEnv(name)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}
func (a *AppController) CreateApp(app *models.AppCreateRequest) (*models.Service, error) {
	//var result = map[string]interface{}{}
	namespace := "default"
	err := app.Validate()
	if err != nil {
		return nil, err
	}
	containerports := []models.ContainerPort{}
	for _, v := range app.ContainerPort {
		containerports = append(containerports, models.ContainerPort{
			ContainerPort: v.Port, //
			Protocol:      models.Protocol(v.Protocol),
		})
	}
	volumemount := []models.VolumeMount{}
	for _, v := range app.Volumes {
		volumemount = append(volumemount, models.VolumeMount{
			Name:      v.Name,
			MountPath: "/usr/local/tomcat/webapps/" + path.Base(v.VolumeSource.HostPath.Path),
		})
	}
	name := app.Name + "-" + app.Version
	name = strings.ToLower(name)
	name = strings.Replace(name, ".", "", -1)
	containers := []models.Container{
		models.Container{
			Name:         name,
			Image:        app.Containerimage,
			Ports:        containerports,
			VolumeMounts: volumemount,
		},
	}
	//var nodeSelector = map[string]string{}
	//if app.Runlocal {
	//	nodeSelector["namespace"] = namespace
	//} else {
	//	nodeSelector["ip"] = strings.Split(a.Ctx.Request.RemoteAddr, ":")[0]
	//}

	var rc = &models.ReplicationController{
		TypeMeta: models.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"env": app.Name, "name": app.Name + "-" + app.Version},
		},
		Spec: models.ReplicationControllerSpec{
			Replicas: app.Replicas,
			Selector: map[string]string{"env": app.Name, "name": app.Name + "-" + app.Version},
			Template: &models.PodTemplateSpec{
				ObjectMeta: models.ObjectMeta{
					Labels: map[string]string{"env": app.Name, "name": app.Name + "-" + app.Version},
				},
				Spec: models.PodSpec{
					Containers: containers,
					Volumes:    app.Volumes,
				},
			},
		},
	}
	body, _ := json.Marshal(rc)
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers"}, body)
	responsebody, _ := json.Marshal(result)
	if status != 201 {
		//fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return nil, models.ErrResponse{Response: string(responsebody)}
	}
	//result["RC"]=
	var Ports = []models.ServicePort{}
	for k, v := range app.Ports {
		Ports = append(Ports, models.ServicePort{
			Name:     "default" + strconv.Itoa(k),
			Port:     v.Port,
			Protocol: models.Protocol(v.Protocol),
		})
	}
	var service = models.Service{
		TypeMeta: models.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name:   name,
			Labels: map[string]string{"env": app.Name, "name": app.Name + "-" + app.Version},
		},
		Spec: models.ServiceSpec{
			Selector: map[string]string{"env": app.Name, "name": app.Name + "-" + app.Version},
			Ports:    Ports,
		},
	}
	for k, v := range app.Ports {
		targetPort := v.TargetPort
		if targetPort == 0 {
			if k < len(app.ContainerPort) {
				targetPort = app.ContainerPort[k].Port
			}
		}
		if targetPort != 0 {
			service.Spec.Ports[k].TargetPort = targetPort
		} else {
			service.Spec.Ports[k].TargetPort = v.Port
		}
	}
	if len(app.PublicIPs) != 0 {
		service.Spec.PublicIPs = app.PublicIPs
	}
	body, _ = json.Marshal(service)
	fmt.Println(string(body))
	status, result = lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "services"}, body)
	responsebody, _ = json.Marshal(result)
	if status != 201 {
		return nil, models.ErrResponse{Response: string(responsebody)}
	}

	_, exist := models.Appinfo[namespace]
	if !exist {
		models.Appinfo[namespace] = models.NamespaceInfo{}
	}
	_, exist = models.Appinfo[namespace][app.Name]
	if !exist {
		models.Appinfo[namespace][app.Name] = &models.AppMetaInfo{
			Name:     app.Name,
			Replicas: app.Replicas,
			Status:   1,
		}
	}
	json.Unmarshal(responsebody, &service)
	return &service, nil
}

// @Title ScaleApp
// @Description Scale app
// @Param       namespaces      path    string  true            "The key for namespaces"
// @Param       service         path    string  true            "The key for name"
// @Param       body            body    models.AppUpgradeRequest         true           "body for user content"
// @Success 200 {string} "scale success"
// @Failure 403 body is empty
// @router /scaleApp [put]
func (a *AppController) Scale() {
	namespace := "default"

	var appScale models.AppScale
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &appScale)

	fmt.Println(appScale.Num)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, err)
		return
	}
	appName := appScale.Name
	//appName := app.Name + "-" + app.Version
	appName = strings.ToLower(appName)
	appName = strings.Replace(appName, ".", "", -1)
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers/" + appName

	rsp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(rsp.Body)
	fmt.Println(string(body))
	var rc models.ReplicationController

	json.Unmarshal(body, &rc)
	rc.Spec.Replicas, _ = strconv.Atoi(appScale.Num)

	body, _ = json.Marshal(rc)
	fmt.Println(string(body))
	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", appName}, body)
	fmt.Println(status)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, result)
		return
	}
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title createApp
// @Description create app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Success 200 {string} "create success"
// @Failure 403 body is empty
// @router /delete [delete]
func (a *AppController) DeleteApp() {
	//namespace := "default"
	//service := a.Ctx.Input.Param(":service")
	var app = map[string]string{}
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &app)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	//fmt.Println(appScale.Num)

	if _, exist := app["appName"]; exist == false {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(a.Ctx.ResponseWriter, `{"errorMessage":"`+"not send appName"+`"}`, 406)
		return
	}
	appName := app["appName"]
	a.deleteapp(appName)
	//re := map[string]interface{}{}
	//re["delete rc"] = result
	//delete(models.Appinfo[namespace], service)
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()

}
func (a *AppController) deleteapp(appName string) {
	namespace := "default"
	appName = strings.ToLower(appName)
	appName = strings.Replace(appName, ".", "", -1)
	lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "services", appName}, []byte{})
	//re["delete service"] = result
	lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", appName}, []byte{})
}

//// @Title get all apps
//// @Description get all apps
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Success 200 {string} "get success"
//// @router / [get]
//func (a *AppController) GetAll() {
//	namespaces := a.Ctx.Input.Param(":namespaces")

//	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "services"}, []byte{})
//	responsebodyK8s, _ := json.Marshal(result)
//	if status != 200 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(status)
//		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
//		return
//	}

//	var appListK8s models.ServiceList //service -> app

//	var appList models.AppGetAllResponse
//	var app models.AppGetAllResponseItem

//	appList.Items = make([]models.AppGetAllResponseItem, 0, 60)

//	json.Unmarshal([]byte(responsebodyK8s), &appListK8s)

//	for index := 0; index < len(appListK8s.Items); index++ {
//		app = models.AppGetAllResponseItem{
//			Name: appListK8s.Items[index].ObjectMeta.Name,
//		}
//		appList.Items = append(appList.Items, app)
//	}

//	//appList.Kind = appListK8s.TypeMeta.Kind
//	appList.Kind = "AppGetAllResponse"

//	responsebody, _ := json.Marshal(appList)

//	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
//	a.Ctx.ResponseWriter.WriteHeader(status)
//	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))

//	//a.Data["json"] = map[string]string{"status": "getall success"}
//	//a.ServeJson()
//}

//// @Title Get App
//// @Description get app by name and namespace
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Param	service		path 	string	true		"The key for name"
//// @Success 200 {string} "get success"
//// @router /:service [get]
//func (a *AppController) Get() {
//	namespaces := a.Ctx.Input.Param(":namespaces")
//	name := a.Ctx.Input.Param(":service")

//	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "services", name}, []byte{})
//	responsebodyK8s, _ := json.Marshal(result)

//	if status != 200 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(status)
//		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
//		return
//	}

//	var appK8s models.Service //service -> app
//	json.Unmarshal([]byte(responsebodyK8s), &appK8s)

//	var app = models.AppGetResponse{
//		Kind:              "AppGetResponse",
//		Name:              appK8s.ObjectMeta.Name,
//		Namespace:         appK8s.ObjectMeta.Namespace,
//		CreationTimestamp: appK8s.ObjectMeta.CreationTimestamp,
//		Labels:            appK8s.ObjectMeta.Labels,
//		Spec:              appK8s.Spec,
//		Status:            appK8s.Status,
//	}
//	responsebody, _ := json.Marshal(app)

//	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
//	a.Ctx.ResponseWriter.WriteHeader(status)
//	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))

//	//a.Data["json"] = map[string]string{"status": "get success"}
//	//a.ServeJson()
//}

//// @Title createApp
//// @Description create app
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Param	service		path 	string	true		"The key for name"
//// @Success 200 {string} "create success"
//// @Failure 403 body is empty
//// @router /:service [delete]
//func (a *AppController) Deleteapp() {
//	namespace := a.Ctx.Input.Param(":namespaces")
//	service := a.Ctx.Input.Param(":service")
//	re := map[string]interface{}{}
//	_, result := lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "services", service}, []byte{})
//	re["delete service"] = result
//	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
//	//fmt.Println(url)
//	rsp, _ := http.Get(url)
//	var rclist models.ReplicationControllerList
//	//var oldrc models.ReplicationController
//	body, _ := ioutil.ReadAll(rsp.Body)
//	//fmt.Println(string(body))
//	json.Unmarshal(body, &rclist)
//	//fmt.Println(rclist.Items[0].Spec)
//	if len(rclist.Items) == 0 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, string("service with no rc"))
//		return
//	}
//	oldrc := rclist.Items[0]
//	oldrc.TypeMeta.Kind = "ReplicationController"
//	oldrc.TypeMeta.APIVersion = "v1beta3"
//	oldrc.Spec.Replicas = 0
//	body, _ = json.Marshal(oldrc)
//	fmt.Println(string(body))
//	_, result = lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
//	re["delete pod"] = result
//	//time.Sleep(5 * time.Second)

//	_, result = lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, []byte{})
//	re["delete rc"] = result
//	delete(models.Appinfo[namespace], service)
//	a.Data["json"] = re
//	a.ServeJson()

//}

//// @Title get App state
//// @Description get App state
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Success 200 {string} "get App state success"
//// @Failure 403 body is empty
//// @router /:service/state [get]
//func (a *AppController) Getstate() {
//	namespace := a.Ctx.Input.Param(":namespaces")
//	service := a.Ctx.Input.Param(":service")
//	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/pods" + "?labelSelector=name%3D" + service
//	//fmt.Println(url)
//	rsp, _ := http.Get(url)

//	var rclist models.PodList
//	body, _ := ioutil.ReadAll(rsp.Body)
//	json.Unmarshal(body, &rclist)
//	fmt.Println(rclist.Items)
//	var res = map[models.PodPhase]int{}
//	for _, v := range rclist.Items {
//		res[v.Status.Phase]++
//	}
//	a.Data["json"] = res
//	a.ServeJson()
//}

//// @Title stop app
//// @Description stop app
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Param	service		path 	string	true		"The key for name"
//// @Success 200 {string} "stop success"
//// @Failure 403 body is empty
//// @router /:service/stop [get]
//func (a *AppController) Stop() {
//	namespace := a.Ctx.Input.Param(":namespaces")
//	service := a.Ctx.Input.Param(":service")

//	_, exist := models.Appinfo[namespace]
//	if !exist {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"no namespace`+namespace+`"}`)
//		return
//	}
//	_, exist = models.Appinfo[namespace][service]
//	if !exist {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"no service`+service+`"}`)
//		return
//	}
//	if models.Appinfo[namespace][service].Status == 0 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"service `+service+` has already been stopped"}`)
//		return
//	}
//	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
//	//fmt.Println(url)
//	rsp, _ := http.Get(url)
//	var rclist models.ReplicationControllerList
//	//var oldrc models.ReplicationController
//	body, _ := ioutil.ReadAll(rsp.Body)
//	//fmt.Println(string(body))
//	json.Unmarshal(body, &rclist)
//	//fmt.Println(rclist.Items[0].Spec)
//	if len(rclist.Items) == 0 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, string("service with no rc"))
//		return
//	}
//	oldrc := rclist.Items[0]
//	oldrc.TypeMeta.Kind = "ReplicationController"
//	oldrc.TypeMeta.APIVersion = "v1beta3"
//	oldrc.Spec.Replicas = 0
//	body, _ = json.Marshal(oldrc)
//	fmt.Println(string(body))
//	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
//	fmt.Println(status)
//	if status != 200 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, result)
//		return
//	} else {
//		models.Appinfo[namespace][service].Status = 0
//	}
//	a.Data["json"] = map[string]string{"messages": "start service successfully"}
//	a.ServeJson()
//}

//// @Title start app
//// @Description start app
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Param	service		path 	string	true		"The key for name"
//// @Success 200 {string} "start success"
//// @Failure 403 body is empty
//// @router /:service/start [get]
//func (a *AppController) Start() {
//	namespace := a.Ctx.Input.Param(":namespaces")
//	service := a.Ctx.Input.Param(":service")
//	_, exist := models.Appinfo[namespace]
//	if !exist {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, "no namespace "+namespace)
//		return
//	}
//	_, exist = models.Appinfo[namespace][service]
//	if !exist {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, "no service"+service)
//		return
//	}
//	if models.Appinfo[namespace][service].Status == 1 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+" has already been started")
//		return
//	}
//	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
//	//fmt.Println(url)
//	rsp, _ := http.Get(url)
//	var rclist models.ReplicationControllerList
//	//var oldrc models.ReplicationController
//	body, _ := ioutil.ReadAll(rsp.Body)
//	//fmt.Println(string(body))
//	json.Unmarshal(body, &rclist)
//	//fmt.Println(rclist.Items[0].Spec)
//	if len(rclist.Items) == 0 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+"with no rc")
//		return
//	}
//	oldrc := rclist.Items[0]
//	oldrc.TypeMeta.Kind = "ReplicationController"
//	oldrc.TypeMeta.APIVersion = "v1beta3"
//	oldrc.Spec.Replicas = models.Appinfo[namespace][service].Replicas
//	body, _ = json.Marshal(oldrc)
//	fmt.Println(string(body))
//	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)

//	if status != 200 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, result)
//		return
//	} else {
//		models.Appinfo[namespace][service].Status = 1
//	}
//	a.Data["json"] = map[string]string{"messages": "start service successfully"}
//	a.ServeJson()
//}

//// @Title UpgradeApp
//// @Description Upgrade app
//// @Param	namespaces	path 	string	true		"The key for namespaces"
//// @Param	service		path 	string	true		"The key for name"
//// @Param	body		body 	models.AppUpgradeRequest	 true		"body for user content"
//// @Success 200 {string} "upgrade success"
//// @Failure 403 body is empty
//// @router /:service/upgrade [put]
//func (a *AppController) Upgrade() {
//	namespace := a.Ctx.Input.Param(":namespaces")
//	service := a.Ctx.Input.Param(":service")
//	var upgradeRequest models.AppUpgradeRequest
//	err := json.Unmarshal(a.Ctx.Input.RequestBody, &upgradeRequest)
//	fmt.Println(upgradeRequest.Containerimage)
//	if err != nil {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, err)
//		return
//	}
//	image := ""
//	//fmt.Println("%v", []byte(upgradeRequest.Warpath))
//	if upgradeRequest.Warpath == "" {
//		////
//		image = upgradeRequest.Containerimage
//	} else {
//		image = "" //war to image
//	}
//	//fmt.Println(image)
//	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
//	//fmt.Println(url)
//	rsp, err := http.Get(url)
//	var rclist models.ReplicationControllerList
//	//var oldrc models.ReplicationController
//	body, _ := ioutil.ReadAll(rsp.Body)
//	//fmt.Println(string(body))
//	json.Unmarshal(body, &rclist)
//	//fmt.Println(rclist.Items[0].Spec)
//	if len(rclist.Items) == 0 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(406)
//		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+"with no rc")
//		return
//	}
//	oldrc := rclist.Items[0]
//	oldrc.TypeMeta.Kind = "ReplicationController"
//	oldrc.TypeMeta.APIVersion = "v1beta3"
//	//fmt.Println(rclist.Items[0])
//	//fmt.Println(oldrc.Spec.Template)
//	//var newrc ReplicationController
//	//fmt.Println(strings.Split(oldrc.ObjectMeta.Name, "-"))
//	oldversion, _ := strconv.Atoi(strings.Split(oldrc.ObjectMeta.Name, "-")[1])
//	newversion := service + "-" + strconv.Itoa(oldversion+1)

//	containers := []models.Container{
//		models.Container{
//			Name:  upgradeRequest.Containername,
//			Image: image,
//			Ports: oldrc.Spec.Template.Spec.Containers[0].Ports,
//		},
//	}

//	var newrc = &models.ReplicationController{
//		TypeMeta: models.TypeMeta{
//			Kind:       "ReplicationController",
//			APIVersion: "v1beta3",
//		},
//		ObjectMeta: models.ObjectMeta{
//			Name:   newversion,
//			Labels: map[string]string{"name": service},
//		},
//		Spec: models.ReplicationControllerSpec{
//			Replicas: oldrc.Spec.Replicas,
//			Selector: map[string]string{"version": newversion},
//			Template: &models.PodTemplateSpec{
//				ObjectMeta: models.ObjectMeta{
//					Labels: map[string]string{"name": service, "version": newversion},
//				},
//				Spec: models.PodSpec{
//					Containers: containers,
//				},
//			},
//		},
//	}

//	body, _ = json.Marshal(newrc)
//	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers"}, body)
//	responsebody, _ := json.Marshal(result)
//	if status != 201 {
//		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
//		a.Ctx.ResponseWriter.WriteHeader(status)
//		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
//		return
//	}
//	//
//	var re = map[string]interface{}{}
//	re["create new rc"] = result
//	oldrc.Spec.Replicas = 0
//	body, _ = json.Marshal(oldrc)
//	fmt.Println(string(body))
//	_, result = lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
//	re["close old pod"] = result
//	//time.Sleep(5 * time.Second)

//	_, result = lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, []byte{})
//	re["delete old rc"] = result

//	_, exist := models.Appinfo[namespace]
//	if !exist {
//		models.Appinfo[namespace] = models.NamespaceInfo{}
//	}
//	_, exist = models.Appinfo[namespace][service]
//	if !exist {
//		models.Appinfo[namespace][service] = &models.AppMetaInfo{
//			Name:     service,
//			Replicas: newrc.Spec.Replicas,
//			Status:   1,
//		}
//	}
//	a.Data["json"] = re
//	a.ServeJson()
//}
