package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

// Operations about Namespace
type NamespaceController struct {
	beego.Controller
}

// @Title createNamespace
// @Description create namespace
// @Param	body		body 	models.NamespaceCreateRequest	 true		"body for user content"
// @Success 201 {string} "create success"
// @router / [post]
func (a *NamespaceController) Post() {
	var namespace models.NamespaceCreateRequest
	err := json.Unmarshal(a.Ctx.Input.RequestBody, &namespace)
	fmt.Println(err)
	//if err != nil {
	//	//a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	//	//a.Ctx.ResponseWriter.WriteHeader(500)
	//	//fmt.Fprintln(a.Ctx.ResponseWriter, `{"err":"`+err.Error()+`"}`)
	//	fmt.Fprintln(a.Ctx.ResponseWriter, err)
	//	return
	//}

	fmt.Println("name is " + namespace.Name)

	var namespaceK8s = models.Namespace{
		TypeMeta: models.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1beta3",
		},
		ObjectMeta: models.ObjectMeta{
			Name:   namespace.Name,
			Labels: map[string]string{"name": namespace.Name},
		},
	}

	body, _ := json.Marshal(namespaceK8s)

	fmt.Println(string(body))
	status, result := lib.Sendapi("POST", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces"}, body)
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

// @Title get all namespaces
// @Description get all namespaces
// @Success 200 {string} "get success"
// @router / [get]
func (a *NamespaceController) GetAll() {

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces"}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)
	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	//var namespace_get models.NamespaceGetResponse
	var namespaceListK8s models.NamespaceList
	var namespaceList models.NamespaceGetAllResponse
	var namespace models.NamespaceGetAllResponseItem

	namespaceList.Items = make([]models.NamespaceGetAllResponseItem, 0, 60)

	json.Unmarshal([]byte(responsebodyK8s), &namespaceListK8s)

	for index := 0; index < len(namespaceListK8s.Items); index++ {

		namespace = models.NamespaceGetAllResponseItem{
			Name: namespaceListK8s.Items[index].ObjectMeta.Name,
		}
		namespaceList.Items = append(namespaceList.Items, namespace)
	}

	namespaceList.Kind = namespaceListK8s.TypeMeta.Kind

	responsebody, _ := json.Marshal(namespaceList)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title Get
// @Description get namespace by name
// @Param	name		path 	string	true		"The key for staticblock"
// @Success 200 {string} "get success"
// @router /:name [get]
func (a *NamespaceController) Get() {
	name := a.Ctx.Input.Param(":name")
	fmt.Println("name is")

	status, result := lib.Sendapi("GET", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", name}, []byte{})
	responsebodyK8s, _ := json.Marshal(result)

	if status != 200 {
		a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
		a.Ctx.ResponseWriter.WriteHeader(status)
		fmt.Println("status is ")
		fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebodyK8s))
		return
	}

	var namespaceK8s models.Namespace
	json.Unmarshal([]byte(responsebodyK8s), &namespaceK8s)

	var namespace = models.NamespaceGetResponse{
		Name:              namespaceK8s.ObjectMeta.Name,
		CreationTimestamp: namespaceK8s.ObjectMeta.CreationTimestamp,
		Status:            namespaceK8s.Status,
	}
	responsebody, _ := json.Marshal(namespace)

	a.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	a.Ctx.ResponseWriter.WriteHeader(status)
	fmt.Fprintln(a.Ctx.ResponseWriter, string(responsebody))
}

// @Title delete
// @Description delete the namespace
// @Param	name		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @router /:name [delete]
func (a *NamespaceController) Delete() {

	name := a.Ctx.Input.Param(":name")
	fmt.Println("delete namespace" + name)

	status, result := lib.Sendapi("DELETE", models.KubernetesIp, "8080", "v1beta3", []string{"namespaces", name}, []byte{})
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
