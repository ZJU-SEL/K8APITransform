package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

// Operations about LimitRange
type LimitRangeController struct {
	beego.Controller
}

// @Title createLimitRange
// @Description create limitrange
// @Param	namespaces	path 	string	true		"The key for staticblock"
// @Param	body		body 	models.LimitRangeCreateRequest	 true		"body for user content"
// @Success 201 {string} "create success"
// @router / [post]
func (a *LimitRangeController) Post() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	var limitrange models.LimitRangeCreateRequest

	err := json.Unmarshal(a.Ctx.Input.RequestBody, &limitrange)
	fmt.Println("err is ", err)
	//if err != nil {
	//	//a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	//a.Ctx.ResponseWriter.WriteHeader(500)
	//	//fmt.Fprintln(a.Ctx.ResponseWriter, `{"err":"`+err.Error()+`"}`)
	//	fmt.Fprintln(a.Ctx.ResponseWriter, err)
	//	return
	//}
	fmt.Println("name is " + limitrange.Name)

	var limitrangeK8s = models.LimitRange{
		TypeMeta: models.TypeMeta{
			Kind:       "LimitRange",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name: limitrange.Name,
		},
		Spec: models.LimitRangeSpec{
			Limits: []models.LimitRangeItem{
				{
					Type:    "Pod",
					Max:     models.ResourceList{"memory": "2Gi", "cpu": "3"},
					Min:     models.ResourceList{"memory": "5Mi", "cpu": "250m"},
					Default: models.ResourceList{"memory": "5Mi", "cpu": "250m"},
				},
				{
					Type:    "Container",
					Max:     models.ResourceList{"memory": "1Gi", "cpu": "2"},
					Min:     models.ResourceList{"memory": "5Mi", "cpu": "250m"},
					Default: models.ResourceList{"memory": "5Mi", "cpu": "250m"},
				},
			},
		},
	}

	body, _ := json.Marshal(limitrangeK8s)

	fmt.Println(string(body))
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "limitranges"}, body)
	fmt.Println("status is ", status)
	responsebody, _ := json.Marshal(result)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)

	if status != 201 {
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return
	}
	a.Data["json"] = map[string]string{"messages": "create namespace successfully"}
	a.ServeJson()
}

// @Title get all limitranges
// @Description get all limitranges
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Success 200 {string} "get success"
// @router / [get]
func (a *LimitRangeController) GetAll() {
	namespaces := a.Ctx.Input.Param(":namespaces")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "limitranges"}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var limitrangeListK8s models.LimitRangeList

	var limitrangeList models.LimitRangeGetAllResponse
	var limitrange models.LimitRangeGetAllResponseItem

	limitrangeList.Items = make([]models.LimitRangeGetAllResponseItem, 0, 60)

	json.Unmarshal([]byte(responsebodyK8s), &limitrangeListK8s)

	for index := 0; index < len(limitrangeListK8s.Items); index++ {

		limitrange = models.LimitRangeGetAllResponseItem{
			Name: limitrangeListK8s.Items[index].ObjectMeta.Name,
		}
		limitrangeList.Items = append(limitrangeList.Items, limitrange)
	}

	//limitrangeList.Kind = limitrangeListK8s.TypeMeta.Kind
	limitrangeList.Kind = "LimitRangeGetAllResponse"

	responsebody, _ := json.Marshal(limitrangeList)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title Get
// @Description get limitrange by name
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	name		path 	string	true		"The key for name"
// @Success 200 {string} "get success"
// @router /:name [get]
func (a *LimitRangeController) Get() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	name := a.Ctx.Input.Param(":name")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "limitranges", name}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)

	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var limitrangeK8s models.LimitRange
	json.Unmarshal([]byte(responsebodyK8s), &limitrangeK8s)

	var limitrange = models.LimitRangeGetResponse{
		Kind:              "LimitRangeGetResponse",
		Name:              limitrangeK8s.ObjectMeta.Name,
		Namespace:         limitrangeK8s.ObjectMeta.Namespace,
		CreationTimestamp: limitrangeK8s.ObjectMeta.CreationTimestamp,
		Spec:              limitrangeK8s.Spec,
	}
	responsebody, _ := json.Marshal(limitrange)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title delete
// @Description delete the limitrange
// @Param	namespaces	path 	string	true		"The namespaces you want to delete"
// @Param	name		path 	string	true		"The name you want to delete"
// @Success 200 {string} delete success!
// @router /:name [delete]
func (a *LimitRangeController) Delete() {

	namespaces := a.Ctx.Input.Param(":namespaces")
	name := a.Ctx.Input.Param(":name")

	fmt.Println("delete limitrange" + name)

	status, result := lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "limitranges", name}, []byte{})
	responsebody, _ := json.Marshal(result)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)

	if status != 200 {
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
		return
	}

	a.Data["json"] = map[string]string{"status": "delete success"}
	a.ServeJson()
}
