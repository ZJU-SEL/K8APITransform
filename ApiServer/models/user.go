package models

import (
	"errors"
	"strconv"
	"time"
)

var (
	UserList map[string]*User
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

func AddUser(u User) string {
	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	UserList[u.Id] = &u
	return u.Id
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
