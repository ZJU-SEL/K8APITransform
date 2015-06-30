package models

import (
	"encoding/json"
)

type AppEnv struct {
	TomcatV string
	JdkV    string
	NodeNum string
	Name    string
	Used    int
}

func (key AppEnv) Validate() error {
	var validationError ValidationError
	if key.TomcatV == "" {
		validationError = validationError.Append(ErrInvalidField{"TomcatV"})
	}

	if key.JdkV == "" {
		validationError = validationError.Append(ErrInvalidField{"JdkV"})
	}
	/*
		if len(key.Ports) == 0 {
			validationError = validationError.Append(ErrInvalidField{"Ports"})
		}
	*/
	if key.Name == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
func AddAppEnv(env *AppEnv) error {
	data, _ := json.Marshal(env)
	_, err := EtcdClient.Create("/envs/"+env.Name, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}
func GetAppEnv(envname string) (*AppEnv, error) {
	response, err := EtcdClient.Get("/envs/"+envname, false, false)
	if err != nil {
		return nil, err
	}
	var env = AppEnv{}
	err = json.Unmarshal([]byte(response.Node.Value), &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
func DeleteAppEnv(envname string) error {
	_, err := EtcdClient.Delete("/envs/"+envname, false)
	return err

}
func GetAllAppEnv() ([]*AppEnv, error) {
	response, err := EtcdClient.Get("/envs/", false, true)
	if err != nil {
		return nil, err
	}
	var envs = []*AppEnv{}
	for _, v := range response.Node.Nodes {
		var env = AppEnv{}
		err = json.Unmarshal([]byte(v.Value), &env)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &env)
	}
	return envs, nil
}
func UpdateAppEnv(envname string, env *AppEnv) error {
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	_, err = EtcdClient.Update("/envs/"+env.Name, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

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
