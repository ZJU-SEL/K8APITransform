package models

import (
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	//"encoding/json"
	//"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
)

func GetPodtoSeFromEtcd(clusterip string, podip string) (string, error) {
	username, cluster, err := Ip2UC(clusterip)
	if err != nil {
		fmt.Println(1)
		return "", err
	}
	client := NewMonitClient(username, cluster)
	request, err := http.NewRequest("GET", "https://"+clusterip+":50000/get?key="+"podiptose/"+podip, nil)
	if err != nil {
		fmt.Println(1)
		return "", err
	}
	request.Header.Set("Authorization", "qwertyuiopasdfghjklzxcvbnm1234567890")
	form := url.Values{}
	form.Add("key", "podiptose/"+podip)
	request.Form = form
	response, err := client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(2)
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(3)
		return "", err
	}
	fmt.Println(string(body))
	resp := &etcd.RawResponse{response.StatusCode, body, response.Header}
	etcdres, err := resp.Unmarshal()
	if err != nil {
		fmt.Println(4)
		return "", err
	}
	sename := etcdres.Node.Value
	fmt.Println(sename)
	request, err = http.NewRequest("GET", "https://"+clusterip+":50000/get?key="+"setoip/"+sename, nil)
	if err != nil {
		fmt.Println(5)
		return "", err
	}
	request.Header.Set("Authorization", "qwertyuiopasdfghjklzxcvbnm1234567890")
	//request.Form.Add("key", "setoip/"+sename)
	response, err = client.Do(request)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(6)
		return "", err
	}
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(7)
		return "", err
	}
	resp = &etcd.RawResponse{response.StatusCode, body, response.Header}
	etcdres, err = resp.Unmarshal()
	if err != nil {
		fmt.Println(8)
		return "", err
	}
	//sename := etcdres.Node.Value
	seip := etcdres.Node.Value
	fmt.Println(seip)

	//seip := response.Node.Value
	//if get the value err is nil
	return seip, nil
}
func GetPodtoSe(clusterip string, podip string) (string, error) {
	fmt.Println("/potose/" + clusterip + "/" + podip)
	response, err := EtcdClient.Get("/potose/"+clusterip+"/"+podip, false, false)
	if err != nil {
		return "", err
	}
	seip := response.Node.Value
	//if get the value err is nil
	return seip, nil
}

func UpdatePodtoSe(clusterip string, podip string, seip string) error {
	_, err := EtcdClient.Update("/potose/"+clusterip+"/"+podip, seip, 0)
	if err != nil {
		return err
	}
	return nil
}

func AddPodtoSe(clusterip string, podip string, seip string) error {
	_, err := EtcdClient.Set("/potose/"+clusterip+"/"+podip, seip, 0)
	if err != nil {
		return err
	}

	return nil

}
