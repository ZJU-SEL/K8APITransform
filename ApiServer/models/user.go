package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	//"strconv"
	//"time"
)

var (
	EtcdClient *etcd.Client
	UserList   map[string]*User
)

type User struct {
	Id       string
	Username string
	Company  string
	Password string
	Ip       string
}

func (u User) Validate() error {
	var validationError ValidationError
	if u.Username == "" {
		validationError = validationError.Append(ErrInvalidField{"username"})
	}

	if u.Password == "" {
		validationError = validationError.Append(ErrInvalidField{"password"})
	}
	if u.Company == "" {
		validationError = validationError.Append(ErrInvalidField{"Company"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
func AddUser(u User) (string, bool) {
	u.Id = "user_" + u.Username
	_, err := EtcdClient.Get("/users/"+u.Id, false, false)
	if err == nil {
		return "", true
	}
	//UserList[u.Id] = &u
	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = EtcdClient.Create("/users/"+u.Id, string(data), 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	//_, err = EtcdClient.Create("/ips/"+u.Ip, string(data), 0)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	return u.Id, false
}

func GetUser(uid string) (u *User, err error) {
	response, err := EtcdClient.Get("/users/"+uid, false, false)
	if err == nil {
		var user = User{}
		data := []byte(response.Node.Value)
		json.Unmarshal(data, &user)
		return &user, nil
	}
	return nil, errors.New("User not exists")
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	u, err := GetUser(uid)
	if err != nil {
		return nil, errors.New("User Not Exist")
	}
	if uu.Username != "" {
		u.Username = uu.Username
	}
	if uu.Password != "" {
		u.Password = uu.Password
	}
	if uu.Ip != "" {
		u.Ip = uu.Ip
	}
	if uu.Company != "" {
		u.Company = uu.Company
	}
	data, err := json.Marshal(u)
	//fmt.Println(u.Username, u.Password, u.Ip)
	response, err := EtcdClient.Update("/users/"+uid, string(data), 0)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(response.Node.Value), u)
	if err != nil {
		return nil, err
	}
	return u, nil

}

func Login(username, password string) (a *User, b bool) {
	uid := "user_" + username
	u, err := GetUser(uid)
	fmt.Println(u)
	fmt.Println(username, password)
	if err != nil {
		return nil, false
	}
	if u.Username == username && u.Password == password {
		return u, true
	}
	return nil, false
}

func Auth(userid string) bool {
	_, err := GetUser(userid)
	if err != nil {
		return false
	}
	return true
}

func DeleteUser(uid string) {
	EtcdClient.Delete("/users/"+uid, false)
	//delete(UserList, uid)
}
