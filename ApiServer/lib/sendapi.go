package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func sendbase(client *http.Client, apitype string, url string, body []byte) (int, map[string]interface{}) {
	var finalreqest *http.Request
	reqest, _ := http.NewRequest(apitype, url, bytes.NewBuffer(body))
	reqest.Header.Set("Content-Type", "application/json")
	finalreqest = reqest

	response, _ := client.Do(finalreqest)

	//body为 []byte类型
	body1, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body1))
	//decoding the []body into the map
	var result map[string]interface{}
	if err := json.Unmarshal(body1, &result); err != nil {
		panic(err)
	}
	//fmt.Println(result)
	status := response.StatusCode
	return status, result

}

func sendGet(client *http.Client, url string) (int, map[string]interface{}) {
	var result map[string]interface{}
	var status int
	fmt.Println("sent get request ")
	//using "" to represent the nil of string
	status, result = sendbase(client, "GET", url, []byte{})

	return status, result

}

func sendPost(client *http.Client, url string, body []byte) (int, map[string]interface{}) {

	var result map[string]interface{}
	var status int
	fmt.Println("sent post request ")

	status, result = sendbase(client, "POST", url, body)

	return status, result

}

func sendDelete(client *http.Client, url string) (int, map[string]interface{}) {
	var result map[string]interface{}
	var status int
	fmt.Println("sent delete request ")
	status, result = sendbase(client, "DELETE", url, []byte{})
	return status, result

}

func sendPut(client *http.Client, url string, body []byte) (int, map[string]interface{}) {
	var result map[string]interface{}
	var status int
	fmt.Println("sent put request ")
	status, result = sendbase(client, "PUT", url, body)
	return status, result

}

//problems in using patch
func sendPatch(client *http.Client, url string, body []byte) (int, map[string]interface{}) {

	var result map[string]interface{}
	var status int
	fmt.Println("sent patch request ")
	status, result = sendbase(client, "PATCH", url, []byte{})

	return status, result

}

func Sendapi(apitype string, host string, port string, version string, commands []string, body []byte) (int, map[string]interface{}) {

	client := &http.Client{}
	fmt.Println(reflect.TypeOf(client))
	url := ""
	//注意前面要加上http://
	if version == "" {
		url = "http://" + host + ":" + port + "/v1"
	} else {
		url = "http://" + host + ":" + port + "/api" + "/" + version
	}
	for _, str := range commands {
		url = url + "/" + str

	}
	fmt.Println(url)

	var result map[string]interface{}
	var status int
	if apitype == "GET" {
		status, result = sendGet(client, url)
	} else if apitype == "POST" {
		status, result = sendPost(client, url, body)

	} else if apitype == "DELETE" {
		status, result = sendDelete(client, url)

	} else if apitype == "PUT" {
		status, result = sendPut(client, url, body)

	} else if apitype == "PATCH" {
		status, result = sendPatch(client, url, body)

	} else {
		panic("error api type")

	}

	return status, result
}
