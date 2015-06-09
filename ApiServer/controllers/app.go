package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/astaxie/beego"
	"net/http"
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
	a.Data["json"] = map[string]string{"msg": "SUCCESS"}
	a.ServeJson()
}

// @Title createApp
// @Description createEnv
// @Param	namespaces	path 	string	true		"The key for staticblock"
// @Param	body		body 	models.AppCreateRequest	 true		"body for user content"
// @Success 200 {string} "create success"
// @Failure 403 body is empty
// @router /createEnv [post]
func (a *AppController) Post() {
	namespace := a.Ctx.Input.Param(":namespaces")
	var app models.AppCreateRequest
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &app)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, err)
		return
	}
	err = app.Validate()
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, err)
		return
	}
	image := ""
	if app.Warpath == "" {
		////
		image = "reg-5000-gorouter"

	} else {
		image = "" //war to image
	}

	fmt.Printf("%s", image)
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
	containers := []models.Container{
		models.Container{
			Name:         app.Containername,
			Image:        app.Containerimage,
			Ports:        containerports,
			VolumeMounts: volumemount,
		},
	}
	var nodeSelector = map[string]string{}
	if app.Runlocal {
		nodeSelector["namespace"] = namespace
	} else {
		nodeSelector["ip"] = strings.Split(a.Ctx.Request.RemoteAddr, ":")[0]
	}
	var rc = &models.ReplicationController{
		TypeMeta: models.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name:   app.Name + "-" + "1",
			Labels: map[string]string{"name": app.Name},
		},
		Spec: models.ReplicationControllerSpec{
			Replicas: 1,
			Selector: map[string]string{"version": app.Name + "-" + "1"},
			Template: &models.PodTemplateSpec{
				ObjectMeta: models.ObjectMeta{
					Labels: map[string]string{"name": app.Name, "version": app.Name + "-" + "1"},
				},
				Spec: models.PodSpec{
					Containers:   containers,
					Volumes:      app.Volumes,
					NodeSelector: nodeSelector,
				},
			},
		},
	}
	body, _ := json.Marshal(rc)
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers"}, body)
	responsebody, _ := json.Marshal(result)
	if status != 201 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return
	}
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
			Name:   app.Name,
			Labels: map[string]string{"name": app.Name},
		},
		Spec: models.ServiceSpec{
			Selector: map[string]string{"name": app.Name},
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
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return
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
	a.Data["json"] = map[string]string{"messages": "create service successfully"}
	a.ServeJson()
}

// @Title get all apps
// @Description get all apps
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Success 200 {string} "get success"
// @router / [get]
func (a *AppController) GetAll() {
	namespaces := a.Ctx.Input.Param(":namespaces")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "services"}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var appListK8s models.ServiceList //service -> app

	var appList models.AppGetAllResponse
	var app models.AppGetAllResponseItem

	appList.Items = make([]models.AppGetAllResponseItem, 0, 60)

	json.Unmarshal([]byte(responsebodyK8s), &appListK8s)

	for index := 0; index < len(appListK8s.Items); index++ {
		app = models.AppGetAllResponseItem{
			Name: appListK8s.Items[index].ObjectMeta.Name,
		}
		appList.Items = append(appList.Items, app)
	}

	//appList.Kind = appListK8s.TypeMeta.Kind
	appList.Kind = "AppGetAllResponse"

	responsebody, _ := json.Marshal(appList)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))

	//a.Data["json"] = map[string]string{"status": "getall success"}
	//a.ServeJson()
}

// @Title Get App
// @Description get app by name and namespace
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Success 200 {string} "get success"
// @router /:service [get]
func (a *AppController) Get() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	name := a.Ctx.Input.Param(":service")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "services", name}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)

	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var appK8s models.Service //service -> app
	json.Unmarshal([]byte(responsebodyK8s), &appK8s)

	var app = models.AppGetResponse{
		Kind:              "AppGetResponse",
		Name:              appK8s.ObjectMeta.Name,
		Namespace:         appK8s.ObjectMeta.Namespace,
		CreationTimestamp: appK8s.ObjectMeta.CreationTimestamp,
		Labels:            appK8s.ObjectMeta.Labels,
		Spec:              appK8s.Spec,
		Status:            appK8s.Status,
	}
	responsebody, _ := json.Marshal(app)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))

	//a.Data["json"] = map[string]string{"status": "get success"}
	//a.ServeJson()
}

// @Title createApp
// @Description create app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Success 200 {string} "create success"
// @Failure 403 body is empty
// @router /:service [delete]
func (a *AppController) Deleteapp() {
	namespace := a.Ctx.Input.Param(":namespaces")
	service := a.Ctx.Input.Param(":service")
	re := map[string]interface{}{}
	_, result := lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "services", service}, []byte{})
	re["delete service"] = result
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
	//fmt.Println(url)
	rsp, _ := http.Get(url)
	var rclist models.ReplicationControllerList
	//var oldrc models.ReplicationController
	body, _ := ioutil.ReadAll(rsp.Body)
	//fmt.Println(string(body))
	json.Unmarshal(body, &rclist)
	//fmt.Println(rclist.Items[0].Spec)
	if len(rclist.Items) == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, string("service with no rc"))
		return
	}
	oldrc := rclist.Items[0]
	oldrc.TypeMeta.Kind = "ReplicationController"
	oldrc.TypeMeta.APIVersion = "v1beta3"
	oldrc.Spec.Replicas = 0
	body, _ = json.Marshal(oldrc)
	fmt.Println(string(body))
	_, result = lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
	re["delete pod"] = result
	//time.Sleep(5 * time.Second)

	_, result = lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, []byte{})
	re["delete rc"] = result
	delete(models.Appinfo[namespace], service)
	a.Data["json"] = re
	a.ServeJson()

}

// @Title get App state
// @Description get App state
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Success 200 {string} "get App state success"
// @Failure 403 body is empty
// @router /:service/state [get]
func (a *AppController) Getstate() {
	namespace := a.Ctx.Input.Param(":namespaces")
	service := a.Ctx.Input.Param(":service")
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/pods" + "?labelSelector=name%3D" + service
	//fmt.Println(url)
	rsp, _ := http.Get(url)

	var rclist models.PodList
	body, _ := ioutil.ReadAll(rsp.Body)
	json.Unmarshal(body, &rclist)
	fmt.Println(rclist.Items)
	var res = map[models.PodPhase]int{}
	for _, v := range rclist.Items {
		res[v.Status.Phase]++
	}
	a.Data["json"] = res
	a.ServeJson()
}

// @Title stop app
// @Description stop app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Success 200 {string} "stop success"
// @Failure 403 body is empty
// @router /:service/stop [get]
func (a *AppController) Stop() {
	namespace := a.Ctx.Input.Param(":namespaces")
	service := a.Ctx.Input.Param(":service")

	_, exist := models.Appinfo[namespace]
	if !exist {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"no namespace`+namespace+`"}`)
		return
	}
	_, exist = models.Appinfo[namespace][service]
	if !exist {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"no service`+service+`"}`)
		return
	}
	if models.Appinfo[namespace][service].Status == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, `{"error":"service `+service+` has already been stopped"}`)
		return
	}
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
	//fmt.Println(url)
	rsp, _ := http.Get(url)
	var rclist models.ReplicationControllerList
	//var oldrc models.ReplicationController
	body, _ := ioutil.ReadAll(rsp.Body)
	//fmt.Println(string(body))
	json.Unmarshal(body, &rclist)
	//fmt.Println(rclist.Items[0].Spec)
	if len(rclist.Items) == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, string("service with no rc"))
		return
	}
	oldrc := rclist.Items[0]
	oldrc.TypeMeta.Kind = "ReplicationController"
	oldrc.TypeMeta.APIVersion = "v1beta3"
	oldrc.Spec.Replicas = 0
	body, _ = json.Marshal(oldrc)
	fmt.Println(string(body))
	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
	fmt.Println(status)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, result)
		return
	} else {
		models.Appinfo[namespace][service].Status = 0
	}
	a.Data["json"] = map[string]string{"messages": "start service successfully"}
	a.ServeJson()
}

// @Title start app
// @Description start app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Success 200 {string} "start success"
// @Failure 403 body is empty
// @router /:service/start [get]
func (a *AppController) Start() {
	namespace := a.Ctx.Input.Param(":namespaces")
	service := a.Ctx.Input.Param(":service")
	_, exist := models.Appinfo[namespace]
	if !exist {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, "no namespace "+namespace)
		return
	}
	_, exist = models.Appinfo[namespace][service]
	if !exist {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, "no service"+service)
		return
	}
	if models.Appinfo[namespace][service].Status == 1 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+" has already been started")
		return
	}
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
	//fmt.Println(url)
	rsp, _ := http.Get(url)
	var rclist models.ReplicationControllerList
	//var oldrc models.ReplicationController
	body, _ := ioutil.ReadAll(rsp.Body)
	//fmt.Println(string(body))
	json.Unmarshal(body, &rclist)
	//fmt.Println(rclist.Items[0].Spec)
	if len(rclist.Items) == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+"with no rc")
		return
	}
	oldrc := rclist.Items[0]
	oldrc.TypeMeta.Kind = "ReplicationController"
	oldrc.TypeMeta.APIVersion = "v1beta3"
	oldrc.Spec.Replicas = models.Appinfo[namespace][service].Replicas
	body, _ = json.Marshal(oldrc)
	fmt.Println(string(body))
	status, result := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)

	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, result)
		return
	} else {
		models.Appinfo[namespace][service].Status = 1
	}
	a.Data["json"] = map[string]string{"messages": "start service successfully"}
	a.ServeJson()
}

// @Title UpgradeApp
// @Description Upgrade app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Param	body		body 	models.AppUpgradeRequest	 true		"body for user content"
// @Success 200 {string} "upgrade success"
// @Failure 403 body is empty
// @router /:service/upgrade [put]
func (a *AppController) Upgrade() {
	namespace := a.Ctx.Input.Param(":namespaces")
	service := a.Ctx.Input.Param(":service")
	var upgradeRequest models.AppUpgradeRequest
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &upgradeRequest)
	fmt.Println(upgradeRequest.Containerimage)
	if err != nil {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, err)
		return
	}
	image := ""
	//fmt.Println("%v", []byte(upgradeRequest.Warpath))
	if upgradeRequest.Warpath == "" {
		////
		image = upgradeRequest.Containerimage
	} else {
		image = "" //war to image
	}
	//fmt.Println(image)
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/replicationcontrollers" + "?labelSelector=name%3D" + service
	//fmt.Println(url)
	rsp, err := http.Get(url)
	var rclist models.ReplicationControllerList
	//var oldrc models.ReplicationController
	body, _ := ioutil.ReadAll(rsp.Body)
	//fmt.Println(string(body))
	json.Unmarshal(body, &rclist)
	//fmt.Println(rclist.Items[0].Spec)
	if len(rclist.Items) == 0 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(a.Ctx.ResponseWriter, "service "+service+"with no rc")
		return
	}
	oldrc := rclist.Items[0]
	oldrc.TypeMeta.Kind = "ReplicationController"
	oldrc.TypeMeta.APIVersion = "v1beta3"
	//fmt.Println(rclist.Items[0])
	//fmt.Println(oldrc.Spec.Template)
	//var newrc ReplicationController
	//fmt.Println(strings.Split(oldrc.ObjectMeta.Name, "-"))
	oldversion, _ := strconv.Atoi(strings.Split(oldrc.ObjectMeta.Name, "-")[1])
	newversion := service + "-" + strconv.Itoa(oldversion+1)

	containers := []models.Container{
		models.Container{
			Name:  upgradeRequest.Containername,
			Image: image,
			Ports: oldrc.Spec.Template.Spec.Containers[0].Ports,
		},
	}

	var newrc = &models.ReplicationController{
		TypeMeta: models.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name:   newversion,
			Labels: map[string]string{"name": service},
		},
		Spec: models.ReplicationControllerSpec{
			Replicas: oldrc.Spec.Replicas,
			Selector: map[string]string{"version": newversion},
			Template: &models.PodTemplateSpec{
				ObjectMeta: models.ObjectMeta{
					Labels: map[string]string{"name": service, "version": newversion},
				},
				Spec: models.PodSpec{
					Containers: containers,
				},
			},
		},
	}

	body, _ = json.Marshal(newrc)
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers"}, body)
	responsebody, _ := json.Marshal(result)
	if status != 201 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return
	}
	//
	var re = map[string]interface{}{}
	re["create new rc"] = result
	oldrc.Spec.Replicas = 0
	body, _ = json.Marshal(oldrc)
	fmt.Println(string(body))
	_, result = lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, body)
	re["close old pod"] = result
	//time.Sleep(5 * time.Second)

	_, result = lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespace, "replicationcontrollers", oldrc.ObjectMeta.Name}, []byte{})
	re["delete old rc"] = result

	_, exist := models.Appinfo[namespace]
	if !exist {
		models.Appinfo[namespace] = models.NamespaceInfo{}
	}
	_, exist = models.Appinfo[namespace][service]
	if !exist {
		models.Appinfo[namespace][service] = &models.AppMetaInfo{
			Name:     service,
			Replicas: newrc.Spec.Replicas,
			Status:   1,
		}
	}
	a.Data["json"] = re
	a.ServeJson()
}

// @Title Roll back App
// @Description roll back app
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	service		path 	string	true		"The key for name"
// @Param	body		body 	models.AppCreateRequest	 true		"body for user content"
// @Success 200 {string} "create success"
// @Failure 403 body is empty
// @router /:service/rollback [put]
func (a *AppController) Rollback() {

	a.Data["json"] = map[string]string{"status": "rollback success"}
	a.ServeJson()
}
