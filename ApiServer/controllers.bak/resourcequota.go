package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

// Operations about ResourceQuota
type ResourceQuotaController struct {
	beego.Controller
}

// @Title createResourceQuota
// @Description create resourcequota
// @Param	namespaces	path 	string	true		"The key for staticblock"
// @Param	body		body 	models.ResourceQuota	 true		"body for user content"
// @Success 201 {string} "create success"
// @router / [post]
func (a *ResourceQuotaController) Post() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	var resourcequota models.ResourceQuotaCreateRequest

	err := json.Unmarshal(a.Ctx.Input.RequestBody, &resourcequota)
	fmt.Println("err is ", err)
	//if err != nil {
	//	//a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	//a.Ctx.ResponseWriter.WriteHeader(500)
	//	//fmt.Fprintln(a.Ctx.ResponseWriter, `{"err":"`+err.Error()+`"}`)
	//	fmt.Fprintln(a.Ctx.ResponseWriter, err)
	//	return
	//}
	fmt.Println("name is " + resourcequota.Name)

	var resourcequotaK8s = models.ResourceQuota{
		TypeMeta: models.TypeMeta{
			Kind:       "ResourceQuota",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name: resourcequota.Name,
		},
		Spec: models.ResourceQuotaSpec{
			Hard: models.ResourceList{"memory": "1Gi", "cpu": "20",
				"pods": "10", "services": "5",
				"replicationcontrollers": "20", "resourcequotas": "5",
				"secrets": "10", "persistentvolumeclaims": "10"},
		},
	}

	body, _ := json.Marshal(resourcequotaK8s)

	fmt.Println(string(body))
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "resourcequotas"}, body)
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

// @Title get all resourcequotas
// @Description get all resourcequotas
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Success 200 {string} "get success"
// @router / [get]
func (a *ResourceQuotaController) GetAll() {
	namespaces := a.Ctx.Input.Param(":namespaces")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "resourcequotas"}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var resourcequotaListK8s models.ResourceQuotaList
	var resourcequotaList models.ResourceQuotaGetAllResponse
	var resourcequota models.ResourceQuotaGetAllResponseItem

	resourcequotaList.Items = make([]models.ResourceQuotaGetAllResponseItem, 0, 60)

	json.Unmarshal([]byte(responsebodyK8s), &resourcequotaListK8s)

	for index := 0; index < len(resourcequotaListK8s.Items); index++ {
		resourcequota = models.ResourceQuotaGetAllResponseItem{
			Name: resourcequotaListK8s.Items[index].ObjectMeta.Name,
		}
		resourcequotaList.Items = append(resourcequotaList.Items, resourcequota)
	}

	//resourcequotaList.Kind = resourcequotaListK8s.TypeMeta.Kind
	resourcequotaList.Kind = "ResourceQuotaGetAllResponse"

	responsebody, _ := json.Marshal(resourcequotaList)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title Get
// @Description get resourcequota by name and namespace
// @Param	namespaces	path 	string	true		"The key for namespaces"
// @Param	name		path 	string	true		"The key for name"
// @Success 200 {string} "get success"
// @router /:name [get]
func (a *ResourceQuotaController) Get() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	name := a.Ctx.Input.Param(":name")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "resourcequotas", name}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)

	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	//var resourcequota_get models.ResourceQuotaGetResponse
	var resourcequotaK8s models.ResourceQuota
	json.Unmarshal([]byte(responsebodyK8s), &resourcequotaK8s)

	var resourcequota = models.ResourceQuotaGetResponse{
		Kind:              "ResourceQuotaGetResponse",
		Name:              resourcequotaK8s.ObjectMeta.Name,
		Namespace:         resourcequotaK8s.ObjectMeta.Namespace,
		CreationTimestamp: resourcequotaK8s.ObjectMeta.CreationTimestamp,
		Spec:              resourcequotaK8s.Spec,
		Status:            resourcequotaK8s.Status,
	}
	responsebody, _ := json.Marshal(resourcequota)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title delete
// @Description delete the resourcequota
// @Param	namespaces	path 	string	true		"The namespaces you want to delete"
// @Param	name		path 	string	true		"The name you want to delete"
// @Success 200 {string} delete success!
// @router /:name [delete]
func (a *ResourceQuotaController) Delete() {
	namespaces := a.Ctx.Input.Param(":namespaces")
	name := a.Ctx.Input.Param(":name")

	fmt.Println("delete resourcequota" + name)

	status, result := lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", namespaces, "resourcequotas", name}, []byte{})
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
