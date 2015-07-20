package models

import (
	//"K8APITransform/ApiServer/backend"
	"encoding/json"
)

type AppEnv struct {
	TomcatV string
	JdkV    string
	NodeNum string
	Name    string
	Cpu     string
	Memory  string
	Storage string
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
	if key.Name == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
func AddAppEnv(ip string, env *AppEnv) error {
	data, _ := json.Marshal(env)
	_, err := EtcdClient.Create("/envs/"+ip+"/"+env.Name, string(data), 0)
	//IdPools.CreateIdPool(env.Name)
	if err != nil {
		return err
	}
	err = IdPools.CreateIdPool(ip, env.Name)
	if err != nil {
		return err
	}
	return nil
}
func GetAppEnv(ip string, envname string) (*AppEnv, error) {
	response, err := EtcdClient.Get("/envs/"+ip+"/"+envname, false, false)
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
func DeleteAppEnv(ip string, envname string) error {
	_, err := EtcdClient.Delete("/envs/"+ip+"/"+envname, false)
	err = IdPools.DeleteIdPool(ip, envname)
	return err

}
func GetAllAppEnv(ip string) ([]*AppEnv, error) {
	response, err := EtcdClient.Get("/envs/"+ip, false, true)
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
func UpdateAppEnv(ip string, envname string, env *AppEnv) error {
	data, err := json.Marshal(env)
	if err != nil {
		return err
	}
	_, err = EtcdClient.Update("/envs/"+ip+"/"+env.Name, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}
