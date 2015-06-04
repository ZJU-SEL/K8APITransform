package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	//"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
)

// Operations about Namespace
type NodesController struct {
	beego.Controller
}

// @Title get node status
// @Description get node status
// @Param	body		body 	nil	 true		"body for user content"
// @Success 201 {string} ""
// @router /status [get]
func (n *NodesController) Status() {
	ip := strings.Split(n.Ctx.Request.RemoteAddr, ":")[0] //n.Ctx.Input.Param(":ip")
	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"nodes", ip}, nil)
	//fmt.Println(status)
	//fmt.Println(result)
	if status != 200 {
		n.Data["json"] = map[string]string{"status": "Failure"}
		n.ServeJson()

	} else {
		date, _ := json.Marshal(result)
		node := models.Node{}
		json.Unmarshal(date, &node)
		response := map[string]string{}
		response["type"] = string(node.Status.Conditions[0].Type)
		response["status"] = string(node.Status.Conditions[0].Status)
		n.Data["json"] = response
		n.ServeJson()
	}
}

// @Title add node label
// @Description add node label
// @Param	body		body 	nil	 "body for user content"
// @Success 201 {string} ""
// @router /addlabels [get]
func (n *NodesController) Addlabels() {

}

// @Title get node status
// @Description get node status
// @Param	body		body 	nil	 true		"body for user content"
// @Success 201 {string} ""
// @router /user/:username [get]
func (n *NodesController) Getnodes() {
	username := n.Ctx.Input.Param(":username")
	url := "http://" + models.KubernetesIp + ":8080/api/v1beta3/nodes" + "?labelSelector=namespace%3D" + username
	//fmt.Println(url)
	rsp, _ := http.Get(url)

	body, _ := ioutil.ReadAll(rsp.Body)
	var response = []string{}
	if rsp.StatusCode != 200 {
		n.Data["json"] = response
		n.ServeJson()
	} else {
		nodelist := models.NodeList{}
		json.Unmarshal(body, &nodelist)
		for _, v := range nodelist.Items {
			response = append(response, v.ObjectMeta.Name)
		}
		n.Data["json"] = response
		n.ServeJson()
	}
}
