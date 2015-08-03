package models

import (
	"crypto/tls"
	"net/http"
)

var monitClient = map[string]*http.Client{}

func NewMonitClient(username string, cluster string) *http.Client {
	if client, exist := monitClient[cluster+"."+username]; exist {
		return client
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	monitClient[cluster+"."+username] = client
	return client
}
