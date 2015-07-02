package models

import (
//"K8APITransform/ApiServer/backend"
//"encoding/json"
)

func GetPodtoSe(podip string) (string, error) {
	response, err := EtcdClient.Get("/potose/"+podip, false, false)
	if err != nil {
		return "", err
	}

	seip := response.Node.Value
	//if get the value err is nil
	return seip, nil
}

func UpdatePodtoSe(podip string, seip string) error {
	_, err := EtcdClient.Update("/potose/"+podip, seip, 0)
	if err != nil {
		return err
	}
	return nil
}

func AddPodtoSe(podip string, seip string) error {

	_, err := GetPodtoSe(podip)
	//using err to adjust if get value
	if err != nil {
		_, err := EtcdClient.Create("/potose/"+podip, seip, 0)
		if err != nil {
			return err
		}
	} else {
		err := UpdatePodtoSe(podip, seip)
		return err
	}
	return nil

}
