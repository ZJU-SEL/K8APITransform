package controllers

import (
	"K8APITransform/ApiServer/lib"
	"K8APITransform/ApiServer/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Operations about Users
type UserController struct {
	beego.Controller
}

// @Title createUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *UserController) Post() {
	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	uid, exist := models.AddUser(user)
	if exist {
		u.Ctx.ResponseWriter.Header().Set("Content-Type", "text/plain; charset=utf-8")
		u.Ctx.ResponseWriter.WriteHeader(406)
		fmt.Fprintln(u.Ctx.ResponseWriter, "user is exist")
		return
	}
	body := `{"name":"` + user.Username + `"}`
	lib.Sendapi("POST", "127.0.0.1", "8081", "", []string{"namespaces"}, []byte(body))
	u.Data["json"] = map[string]string{"uid": uid}
	u.ServeJson()
}

// @Title Get
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (u *UserController) GetAll() {
	users := models.GetAllUsers()
	u.Data["json"] = users
	u.ServeJson()
}

// @Title Get
// @Description get user by uid
// @Param	uid		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :uid is empty
// @router /:uid [get]
func (u *UserController) Get() {
	uid := u.GetString(":uid")
	if uid != "" {
		user, err := models.GetUser(uid)
		if err != nil {
			u.Data["json"] = err
		} else {
			u.Data["json"] = user
		}
	}
	u.ServeJson()
}

// @Title update
// @Description update the user
// @Param	uid		path 	string	true		"The uid you want to update"
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {object} models.User
// @Failure 403 :uid is not int
// @router /:uid [put]
func (u *UserController) Put() {
	uid := u.GetString(":uid")
	if uid != "" {
		var user models.User
		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
		uu, err := models.UpdateUser(uid, &user)
		if err != nil {
			u.Data["json"] = err
		} else {
			u.Data["json"] = uu
		}
	}
	u.ServeJson()
}

// @Title delete
// @Description delete the user
// @Param	uid		path 	string	true		"The uid you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 uid is empty
// @router /:uid [delete]
func (u *UserController) Delete() {
	uid := u.GetString(":uid")
	models.DeleteUser(uid)
	u.Data["json"] = "delete success!"
	u.ServeJson()
}

// @Title login
// @Description Logs user into the system
// @Param	username		query 	string	true		"The username for login"
// @Param	password		query 	string	true		"The password for login"
// @Success 200 {string} login success
// @Failure 500 user not exist
// @router /login [post]
func (u *UserController) Login() {

	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	fmt.Println(string(u.Ctx.Input.RequestBody))
	username := user.Username
	password := user.Password
	//username := u.GetString("username")
	//password := u.GetString("password")
	fmt.Println(username)
	fmt.Println(password)
	ip := strings.Split(u.Ctx.Request.RemoteAddr, ":")[0]
	if resuser, err := models.Login(username, password); err {
		response, _ := json.Marshal(resuser)
		fmt.Println(string(response))
		resp, err := http.Get("http://" + models.KubernetesIp + ":8080/api/v1beta3/nodes/" + ip)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			var node = models.Node{}
			json.Unmarshal(body, &node)
			node.ObjectMeta.Labels = map[string]string{"namespace": username, "ip": ip}
			body, _ = json.Marshal(node)
			status, _ := lib.Sendapi("PUT", models.KubernetesIp, "8080", "v1beta3", []string{"nodes", ip}, body)
			fmt.Println("add label status:" + strconv.Itoa(status))
		}
		//u.Data["json"] = response
		http.Error(u.Ctx.ResponseWriter, string(response)+"@login successful", 200)
		return
		//u.Data["json"] = "login successful"
	} else {
		http.Error(u.Ctx.ResponseWriter, "user not exist", 500)
		return
	}
	//u.ServeJson()
}

// @Title auth
// @Description auth user into the system
// @Param	userid		query 	string	true		"The userid for user"
// @Success 200 {string} auth success
// @Failure 500 user not exist
// @router /auth [post]
func (u *UserController) Auth() {

	var user models.User
	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
	fmt.Println(string(u.Ctx.Input.RequestBody))
	userid := user.Id

	fmt.Println(userid)

	if models.Auth(userid) {

		u.Data["json"] = "Auth successful"
	} else {
		http.Error(u.Ctx.ResponseWriter, "user not exist", 500)
		return
	}
	u.ServeJson()
}

// @Title logout
// @Description Logs out current logged in user session
// @Success 200 {string} logout success
// @router /logout [get]
func (u *UserController) Logout() {
	u.Data["json"] = "logout success"
	u.ServeJson()
}
