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
	UserList   map[string]*User
	EtcdClient *etcd.Client
)

func init() {
	UserList = make(map[string]*User)
	u := User{"user_abcd", "abcd", "123456"}
	UserList["user_abcd"] = &u
	UserList["abcd"] = &u
}

type User struct {
	Id       string
	Username string
	Password string
}

func AddUser(u User) (string, bool) {
	u.Id = "user_" + u.Username
	response, err := EtcdClient.Get("/users/"+u.Id, false, false)
	if err != nil {
		fmt.Println(err.Error())
	}
	if _, exist := UserList[u.Id]; exist {
		return "", true
	}
	UserList[u.Id] = &u
	data, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err := EtcdClient.Create("/users/"+u.Id, data, 0)
	if err != nil {
		fmt.Println(err.Error())
	}
	return u.Id, false
}

func GetUser(uid string) (u *User, err error) {
	if u, ok := UserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	if u, ok := UserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}

		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func Login(username, password string) (a *User, b bool) {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return u, true
		}
	}
	return nil, false
}

func Auth(userid string) bool {
	for _, u := range UserList {
		if u.Id == userid {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}
