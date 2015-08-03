package controllers

import (
	"K8APITransform/ApiServer/models"
	//"fmt"
	"encoding/json"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
)

type ClusterController struct {
	beego.Controller
}

// @Title get nodes
// @Description get nodes

// @router /nodes [get]
func (c *ClusterController) Nodes() {
	ip := c.Ctx.Request.Header.Get("Authorization")
	username, cluster, err := models.Ip2UC(ip)
	client := models.NewMonitClient(username, cluster)
	request, err := http.NewRequest("GET", "https://"+ip+":50000/api/cluster/nodes", nil)
	request.Header.Set("token", "qwertyuiopasdfghjklzxcvbnm1234567890")
	response, err := client.Do(request)
	if err != nil {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(c.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	var w = []string{}
	json.Unmarshal(body, &w)
	c.Data["json"] = w
	c.ServeJson()
}

// @Title get clusterInfo
// @Description get clusterInfo

// @router /nodeInfo [get]
func (c *ClusterController) NodeInfo() {
	ip := c.Ctx.Request.Header.Get("Authorization")
	node := c.GetString("node")
	flag := c.GetString("isFirst")
	username, cluster, err := models.Ip2UC(ip)
	client := models.NewMonitClient(username, cluster)
	request, err := http.NewRequest("GET", "https://"+ip+":50000/api/node/status", nil)
	request.Header.Set("token", "qwertyuiopasdfghjklzxcvbnm1234567890")
	request.Header.Set("node", node)
	request.Header.Set("flag", flag)
	response, err := client.Do(request)
	if err != nil {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(c.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	var w = map[string]interface{}{}
	json.Unmarshal(body, &w)
	c.Data["json"] = w
	c.ServeJson()
}

// @Title get clusterInfo
// @Description get clusterInfo

// @router /clusterInfo [get]
func (c *ClusterController) ClusterInfo() {
	ip := c.Ctx.Request.Header.Get("Authorization")
	flag := c.GetString("isFirst")
	username, cluster, err := models.Ip2UC(ip)
	client := models.NewMonitClient(username, cluster)
	request, err := http.NewRequest("GET", "https://"+ip+":50000/api/cluster/status", nil)
	request.Header.Set("token", "qwertyuiopasdfghjklzxcvbnm1234567890")
	request.Header.Set("flag", flag)
	response, err := client.Do(request)
	if err != nil {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
		http.Error(c.Ctx.ResponseWriter, `{"errorMessage":"`+err.Error()+`"}`, 406)
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))
	var w = map[string]interface{}{}
	json.Unmarshal(body, &w)
	//w, err := simplejson.NewJson(body)
	//result := map[string]string{}
	//result["name"] = cluster
	//clusters, err := w.Get("Cluster").Map()
	//w = w.Get("Cluster")
	//var memoryusage uint64
	//var diskusage uint64
	//var memorycapacity uint64
	//var diskcapacity uint64
	//for k, _ := range clusters {
	//	disk, err := w.Get(k).Get("Diskcapacity").Uint64()
	//	diskcapacity += disk
	//	memory, err := w.Get(k).Get("Memorycapacity").Uint64()
	//	memorycapacity += memory
	//	memory, err = w.Get(k).Get("Spec").Get("MemoryAvg").Uint64()
	//	memoryusage += memory
	//	disk, err = w.Get(k).Get("Spec").Get("DiskAvg").Uint64()
	//	diskusage += disk
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}
	//result["memoryUsage"] = fmt.Sprintf("%v", memoryusage)
	//result["diskUsage"] = fmt.Sprintf("%v", diskusage)
	//result["memoryLimit"] = fmt.Sprintf("%v", memorycapacity)
	//result["diskLimit"] = fmt.Sprintf("%v", diskcapacity)

	c.Data["json"] = w
	c.ServeJson()
}
