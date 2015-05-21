package main

import (
	//	"ContainerCloudCli/Sendreq/Apitool"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func sendGettest(host string, port string, version string, getcommands []string) ([]byte, int) {
	url := "http://" + host + ":" + port + "/" + version

	for _, str := range getcommands {
		url = url + "/" + str

	}

	fmt.Println("send request:" + url)

	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(10 * time.Second)
				c, err := net.DialTimeout(netw, addr, time.Second*5)
				if err != nil {
					fmt.Println("time out , checkout your net connection")
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}
	reqest, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}

	response, _ := client.Do(reqest)

	//body为 []byte类型
	body, _ := ioutil.ReadAll(response.Body)
	status := response.StatusCode

	return body, status
}

//func main() {
//	//fmt.Println("test search")
//	serverip := "10.10.105.204"
//	//	namespace := "localnamespace"
//	getcommands := []string{"baseimage", "search"}
//	//Apitool.Sendapi("GET", serverip, "8080", "v1", getcommands, "")
//	responsebody, status := sendGettest(serverip, "8080", "v1", getcommands)
//	fmt.Println(string(responsebody), status)

//}
