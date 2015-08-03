package models

import (
	//"K8APITransform/ApiServer/backend"
	"fmt"
	//"encoding/json"
)

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

	_, err := GetPodtoSe(clusterip, podip)
	//using err to adjust if get value
	if err != nil {
		_, err := EtcdClient.Create("/potose/"+clusterip+"/"+podip, seip, 0)
		if err != nil {
			return err
		}
	} else {
		err := UpdatePodtoSe(clusterip, podip, seip)
		return err
	}
	return nil

}
