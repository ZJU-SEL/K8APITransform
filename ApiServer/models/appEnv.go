package models

import (
	"encoding/json"
)

type AppEnv struct {
	TomcatV string
	JdkV    string
	NodeNum int
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
