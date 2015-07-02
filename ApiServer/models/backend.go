package models

import (
	//"K8APITransform/ApiServer/models"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	//"encoding/json"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	//"io/ioutil"
	//"errors"
	"log"
	//"net/http"
	//"net/url"
	//"path"
)

type Backend struct {
	*client.Client
}

func NewBackend(host string, apiVersion string) (*Backend, error) {
	Client, err := client.New(&client.Config{Host: host, Version: apiVersion})
	if err != nil {
		return nil, err
	}
	return &Backend{Client}, nil
}
func (c *Backend) Applications(env string) ApplicationInterface {
	return newApplications(c, env)
}

func (c *Backend) Podip(sename string) ([]string, error) {
	namespace := "default"
	//todo : get info from the sys dynamically
	port := "8080"
	//url := "http://" + KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/pods" + "?labelSelector=name%3D" + sename
	//log.Println(url)
	//rsp, _ := http.Get(url)
	//body, _ := ioutil.ReadAll(rsp.Body)
	label := map[string]string{}
	label["name"] = sename
	podlist, err := c.Pods(namespace).List(labels.SelectorFromSet(label), nil)
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

	//----------------------------------------------------
	//service ip
	//url = "http://" + KubernetesIp + ":8080/api/v1beta3/namespaces/" + namespace + "/services" + "?labelSelector=name%3D" + sename
	//rsp, _ = http.Get(url)
	//body, _ = ioutil.ReadAll(rsp.Body)
	//var servicelist ServiceList
	//json.Unmarshal(body, &servicelist)
	servicelist, err := c.Services(namespace).List(labels.SelectorFromSet(label))
	if err != nil {
		return nil, err
	}

	service := servicelist.Items[0]
	serviceip := service.Spec.ClusterIP + ":" + port

	//log.Println("servicePortalIP:", serviceip)
	//servicename := service.ObjectMeta.Labels["name"]
	//range the podip and store it
	// pod could not be allocated the ip immediately watch etcd???
	//question here could not get the pod ip?????
	//time.Sleep(time.Second * 5)
	log.Println("podlist:", iplist)
	for _, podip := range iplist {
		//serviceipmap[podip] = serviceip
		//store the info into etcd
		err := AddPodtoSe(podip, serviceip)
		//return nil, err
		if err != nil {
			return nil, err
		}
	}

	//	log.Println("servicemap:", serviceipmap)
	//----------------------------------------------------

	return iplist, nil

}
