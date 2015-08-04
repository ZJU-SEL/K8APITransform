package models

import (
	//"encoding/json"
	"errors"
	"fmt"
	api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api/resource"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	//"net/url"
	"sort"
	"time"
	//"sort"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"time"
)

// PodInterface has methods to work with Pod resources.
type AppsInterface interface {
	//List() (*Detail, error)
	Get(name string) (*Detail, error)
	Delete(name string) error
	DeleteAll() error
	Stop(name string) error
	Start(name string) error
	Create(app *AppConfig) error
	Scale(name string, replicas int) error
	Restart(name string) error
	Monit(name string, id string, flag string) (map[string]interface{}, error)
	Rmi(imagename string) error
	Debug(name string) error
	CloseDebug(name string) error
}
type Services struct {
	api.ServiceList
}

func (list *Services) Len() int {
	return len(list.Items)
}

func (list *Services) Less(i, j int) bool {
	if list.Items[i].ObjectMeta.CreationTimestamp.Unix() < list.Items[j].ObjectMeta.CreationTimestamp.Unix() {
		return true
	} else {
		return false
	}
}

func (list *Services) Swap(i, j int) {
	var temp = list.Items[i]
	list.Items[i] = list.Items[j]
	list.Items[j] = temp
}

type Pods struct {
	api.PodList
}

func (list *Pods) Len() int {
	return len(list.Items)
}

func (list *Pods) Less(i, j int) bool {
	if list.Items[i].ObjectMeta.CreationTimestamp.Unix() < list.Items[j].ObjectMeta.CreationTimestamp.Unix() {
		return true
	} else {
		return false
	}
}

func (list *Pods) Swap(i, j int) {
	var temp = list.Items[i]
	list.Items[i] = list.Items[j]
	list.Items[j] = temp
}

// pods implements PodsNamespacer interface
type apps struct {
	e   *Envs
	b   *client.Client
	env string
}

// newPods returns a pods
func newApps(e *Envs, b *client.Client, env string) *apps {
	return &apps{
		e:   e,
		b:   b,
		env: env,
	}
}
func (a *apps) Create(app *AppConfig) (err error) {
	label := map[string]string{"env": a.env, "name": a.env + "-" + app.Name + "-" + app.Version}
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(label))
	if err != nil {
		return err
	}
	errmsg := ""
	if len(sevicelist.Items) != 0 {
		errmsg += "app service is exist"
	}
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(label))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 0 {
		errmsg += "app rc is exist"
	}
	if errmsg != "" {
		return fmt.Errorf(errmsg)
	}
	containerports := []api.ContainerPort{}
	for _, v := range app.ContainerPort {
		containerports = append(containerports, api.ContainerPort{
			ContainerPort: v.Port, //
			Protocol:      api.Protocol(v.Protocol),
		})
	}
	env, err := a.e.Get(a.env)
	if err != nil {
		return err
	}
	id, err := a.e.C.IdPools(a.e.Cluster).GetId(a.env)
	if err != nil {
		return err
	}

	containers := []api.Container{
		api.Container{
			Name:  a.env + "-" + id,
			Image: app.Containerimage,
			Ports: containerports,
			Resources: api.ResourceRequirements{
				Limits: api.ResourceList{},
			},
			Env: []api.EnvVar{},
		},
		api.Container{
			Name:  "packbeat",
			Image: "packetbeat/v2",
			Resources: api.ResourceRequirements{
				Limits: api.ResourceList{},
			},
			Env: []api.EnvVar{},
		},
	}
	//var cores, memorySize, diskSize *resource.Quantity

	cores, err := resource.ParseQuantity(env.Cpu)
	if err != nil {
		return err
	}
	containers[0].Resources.Limits[api.ResourceCPU] = *cores

	memorySize, err := resource.ParseQuantity(env.Memory + "Gi")
	if err != nil {
		return err
	}
	containers[0].Resources.Limits[api.ResourceMemory] = *memorySize

	diskSize, err := resource.ParseQuantity(env.Disk + "Gi")
	if err != nil {
		return err
	}
	containers[0].Resources.Limits[api.ResourceStorage] = *diskSize

	log.Println(containers[0].Resources.Limits)

	replicas, err := strconv.Atoi(env.Instance)
	if err != nil {
		return err
	}
	var rc = &api.ReplicationController{
		TypeMeta: api.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: ApiVersion,
		},
		ObjectMeta: api.ObjectMeta{
			Name:   a.env + "-" + id,
			Labels: label,
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: replicas,
			Selector: label,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: label,
				},
				Spec: api.PodSpec{
					Containers: containers,
					Volumes:    app.Volumes,
				},
			},
		},
	}
	data, err := json.Marshal(rc)
	log.Println(string(data))

	data, err = json.Marshal(rc)
	log.Println(string(data))

	var Ports = []api.ServicePort{}
	for k, v := range app.Ports {
		Ports = append(Ports, api.ServicePort{
			Name:     "default" + strconv.Itoa(k),
			Port:     v.Port,
			Protocol: api.Protocol(v.Protocol),
		})
	}
	var service = &api.Service{
		TypeMeta: api.TypeMeta{
			Kind:       "Service",
			APIVersion: ApiVersion,
		},
		ObjectMeta: api.ObjectMeta{
			Name:   a.env + "-" + id,
			Labels: label,
		},
		Spec: api.ServiceSpec{
			Selector: label,
			Ports:    Ports,
		},
	}
	defer func() {
		if err != nil {
			a.b.ReplicationControllers("default").Delete(rc.ObjectMeta.Name)
			a.b.Services("default").Delete(service.ObjectMeta.Name)
		}
	}()
	service, err = a.b.Services("default").Create(service)
	if err != nil {
		return err
	}
	rc.Spec.Template.Spec.Containers[0].Env = append(rc.Spec.Template.Spec.Containers[0].Env, api.EnvVar{"serviceIp", service.Spec.ClusterIP, nil})
	rc.Spec.Template.Spec.Containers[1].Env = append(rc.Spec.Template.Spec.Containers[1].Env, api.EnvVar{"serviceIp", service.Spec.ClusterIP, nil})

	rc, err = a.b.ReplicationControllers("default").Create(rc)
	if err != nil {
		return err
	}

	fmt.Println(rc)
	fmt.Println(service)
	t := time.After(time.Minute)
	ip, err := Host2Ip(a.e.C.UserName, a.e.Cluster)
	if err != nil {
		return err
	}
A:
	for {
		select {
		case <-t:
			log.Println("time out to allocate ip")
			//delete the se which deploy failed
			return errors.New(`{"errorMessage":"` + "deploy error : time out" + `"}`)
			break A
		default:
			//log.Println("logout:", <-timeout)
			sename := service.ObjectMeta.Labels["name"]
			podslist, err := a.Podip(ip, sename)
			if err == nil {
				if len(podslist) == 0 {
					continue
				} else {
					log.Println("allocation ok ......")
					break A
				}
			} else {
				log.Println(err.Error())
				return errors.New(`{"errorMessage":"` + err.Error() + `"}`)
				//delayok <- 0
				break A
			}
		}
	}

	log.Println("waing pods ip allocation....")
	//detail, err := a.Get(labels["name"])
	return nil
}
func (a *apps) Podip(clusterip, sename string) ([]string, error) {
	namespace := "default"
	port := "8080"
	label := map[string]string{}
	label["name"] = sename
	podlist, err := a.b.Pods(namespace).List(labels.SelectorFromSet(label), nil)
	if err != nil {
		return nil, err
	}
	//json.Unmarshal(body, &podlist)
	//log.Println(string(body))
	var iplist []string
	//var tmppodip = "null"
	if len(podlist.Items) == 0 {
		return iplist, nil
	}
	tmppodip := podlist.Items[0].Status.PodIP
	//log.Println("tmppodip:", tmppodip)
	if tmppodip == "" {
		return iplist, nil
	}
	for _, pod := range podlist.Items {
		podip := pod.Status.PodIP
		iplist = append(iplist, podip+":"+port)
	}
	servicelist, err := a.b.Services(namespace).List(labels.SelectorFromSet(label))
	if err != nil {
		return nil, err
	}

	service := servicelist.Items[0]
	serviceip := service.Spec.ClusterIP + ":" + port
	log.Println("podlist:", iplist)
	for _, podip := range iplist {
		err := AddPodtoSe(clusterip, podip, serviceip)
		if err != nil {
			return nil, err
		}
	}
	return iplist, nil

}

func (a *apps) Scale(name string, replicas int) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	if _, exist := rc.ObjectMeta.Labels["stopped"]; exist {
		return fmt.Errorf("app %s is stopped", name)
	}
	rc.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		return err
	}
	return nil
}
func (a *apps) Get(name string) (*Detail, error) {
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return nil, err
	}
	if len(sevicelist.Items) != 1 {
		return nil, ErrResponse{fmt.Sprintf("a app with %d services", len(sevicelist.Items))}
	}
	service := sevicelist.Items[0]
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return nil, err
	}
	if len(rclist.Items) != 1 {
		return nil, ErrResponse{fmt.Sprintf("a app with %d rcs", len(sevicelist.Items))}
	}
	rc := rclist.Items[0]
	podslist, err := a.b.Pods("default").List(labels.SelectorFromSet(service.ObjectMeta.Labels), nil)
	if err != nil {
		return nil, err
	}
	sort.Sort(&Pods{*podslist})
	e, err := a.e.Get(a.env)
	if err != nil {
		return nil, err
	}
	detail := &Detail{
		Name:         a.env,
		Cpu:          e.Cpu,
		Memory:       e.Memory,
		Disk:         e.Disk,
		InstanceInfo: []*InstanceInfo{},
	}

	count := 0
	for k, v := range podslist.Items {
		instanceInfo := &InstanceInfo{}
		instanceInfo.Status = string(v.Status.Phase)
		if instanceInfo.Status == "Running" {
			count++
		}
		fmt.Println("abc", k, v.ObjectMeta.CreationTimestamp.Unix())
		instanceInfo.NodeIp = v.Status.PodIP
		instanceInfo.Id = v.ObjectMeta.Name
		detail.InstanceInfo = append(detail.InstanceInfo, instanceInfo)
	}
	if rc.ObjectMeta.Labels["stopped"] != "" {
		detail.Status = "Stopped"
	} else if count == 0 {
		detail.Status = "Unavailable"
	} else if count == rc.Spec.Replicas {
		detail.Status = "Running"
	} else {
		detail.Status = "Pending"
	}
	detail.Instance = rc.Spec.Replicas
	return detail, nil
}
func (a *apps) Stop(name string) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	//replicas := rc.Spec.Replicas
	if _, exist := rc.ObjectMeta.Labels["stopped"]; exist {
		return fmt.Errorf("app %s is stopped", name)
	}
	rc.ObjectMeta.Labels["stopped"] = fmt.Sprintf("%d", rc.Spec.Replicas)
	rc.Spec.Replicas = 0
	_, err = a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		return err
	}
	return nil
}
func (a *apps) Start(name string) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	if _, exist := rc.ObjectMeta.Labels["stopped"]; !exist {
		return fmt.Errorf("app %s is started", name)
	}
	replicas, err := strconv.Atoi(rc.ObjectMeta.Labels["stopped"])
	if err != nil {
		return err
	}
	delete(rc.ObjectMeta.Labels, "stopped")
	rc.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		return err
	}
	return nil
}
func (a *apps) Restart(name string) error {
	err := a.Stop(name)
	if err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	err = a.Start(name)
	if err != nil {
		return err
	}
	return nil
}
func (a *apps) Delete(name string) (err error) {
	log.Println("Delete App :", name)
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return err
	}
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": a.env, "name": name}))
	if err != nil {
		return err
	}
	replicas := map[string]int{}
	defer func() {
		if err != nil {
			//a.b.ReplicationControllers("default").Delete(rc.ObjectMeta.Name)
			for _, v := range sevicelist.Items {
				a.b.Services("default").Create(&v)
			}
			for _, v := range rclist.Items {
				a.b.ReplicationControllers("default").Create(&v)
			}
			//a.b.Services("default").Delete(service.ObjectMeta.Name)
		} else {
			//tag := strings.TrimPrefix(name, a.env)
			tag := name
			data := strings.Split(tag, "-")
			tag = strings.TrimSuffix(tag, ".war-"+data[len(data)-1]) + "-" + data[len(data)-1] + ".war"
			err = a.Rmi(a.e.C.UserName + "reg:5000/apm-jre7-tomcat7:" + tag)
		}
	}()
	for _, v := range sevicelist.Items {
		log.Println("Delete Sevice :", v.ObjectMeta.Name)
		err = a.b.Services("default").Delete(v.ObjectMeta.Name)
		if err != nil {
			return err
		}
	}

	for _, v := range rclist.Items {
		replicas[v.ObjectMeta.Name] = v.Spec.Replicas
		v.Spec.Replicas = 0
		rc, err := a.b.ReplicationControllers("default").Update(&v)
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond)
		v = *rc
		err = a.b.ReplicationControllers("default").Delete(v.ObjectMeta.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *apps) DeleteAll() error {
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(map[string]string{"env": a.env}))
	if err != nil {
		return err
	}

	for _, v := range sevicelist.Items {
		log.Println("Delete App :", v.ObjectMeta.Labels["name"])
		name := v.ObjectMeta.Labels["name"]
		err = a.Delete(name)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *apps) Monit(name string, id string, flag string) (map[string]interface{}, error) {
	podslist, err := a.b.Pods("default").List(labels.SelectorFromSet(map[string]string{"env": a.env}), nil)
	if err != nil {
		return nil, err
	}
	containerid := ""
A:
	for _, v := range podslist.Items {
		if v.ObjectMeta.Name == id {
			for _, container := range v.Status.ContainerStatuses {
				log.Println("Container Name:", container.Name)
				if container.Name != "packbeat" {
					containerid = container.ContainerID
					break A
				}
			}
		}
	}
	if containerid == "" {
		return nil, fmt.Errorf("pod downs")
	}
	containerid = strings.TrimPrefix(containerid, "docker://")
	fmt.Println(containerid)
	client := NewMonitClient(a.e.C.UserName, a.e.Cluster)

	request, err := http.NewRequest("GET", "https://"+a.e.Cluster+"."+a.e.C.UserName+":50000/api/container/status", nil)
	request.Header.Set("token", "qwertyuiopasdfghjklzxcvbnm1234567890")
	request.Header.Set("flag", flag)
	request.Header.Set("container", containerid)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(body))
	//fmt.Println(string(body))
	var w = map[string]interface{}{}
	json.Unmarshal(body, &w)
	return w, nil
}
func (a *apps) Rmi(imagename string) error {
	fmt.Println(imagename)
	client := NewMonitClient(a.e.C.UserName, a.e.Cluster)

	request, err := http.NewRequest("POST", "https://"+a.e.Cluster+"."+a.e.C.UserName+":50000/rmi?imagesname="+imagename, nil)
	request.Header.Set("Authorization", "qwertyuiopasdfghjklzxcvbnm1234567890")
	//request.Form = url.Values{}
	//request.Form.Add("imagesname", imagename)
	//request.Header.Set("container", containerid)
	response, err := client.Do(request)
	if err != nil {
		return err
	} else {
		response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
		return nil
	}

}
func (a *apps) Debug(name string) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	replicas := rc.Spec.Replicas
	rc.Spec.Replicas = 0
	rcnew, err := a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		fmt.Println(1)
		return err
	}
	time.Sleep(time.Second)
	//rcnew := rclist.Items[0]
	rcnew.Spec.Template.Spec.Containers[0].Env = append(rc.Spec.Template.Spec.Containers[0].Env, api.EnvVar{"apm_debug", "true", nil})
	rcnew.Spec.Template.Spec.Containers[1].Env = append(rc.Spec.Template.Spec.Containers[1].Env, api.EnvVar{"apm_debug", "true", nil})
	rcnew.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(rcnew)
	if err != nil {
		fmt.Println("afasdfasdfasdf3")
		return err
	}
	return nil
}
func (a *apps) CloseDebug(name string) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}

	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	replicas := rc.Spec.Replicas
	rc.Spec.Replicas = 0
	rcnew, err := a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		fmt.Println(1)
		return err
	}

	//rcnew := rclist.Items[0]
	envs := []api.EnvVar{}
	for _, v := range rcnew.Spec.Template.Spec.Containers[0].Env {
		if v.Name != "apm_debug" {
			envs = append(envs, v)
		}
	}
	rcnew.Spec.Template.Spec.Containers[0].Env = envs
	envs = []api.EnvVar{}
	for _, v := range rcnew.Spec.Template.Spec.Containers[1].Env {
		if v.Name != "apm_debug" {
			envs = append(envs, v)
		}
	}
	rcnew.Spec.Template.Spec.Containers[1].Env = envs
	rcnew.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(rcnew)
	if err != nil {
		fmt.Println("afasdfasdfasdf3")
		return err
	}
	return nil
}
