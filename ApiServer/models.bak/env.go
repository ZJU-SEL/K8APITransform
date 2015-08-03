package models

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

type Env struct {
	ContainerVersion string
	BuildEnv         string
	Target           string
	Instance         string
	Name             string
	Cpu              string
	Memory           string
	Storage          string
	//UpdateTime       string
	//Version          string
	//NewVersion       string
	//Status           string
	//AppName          string
	//NewAppName       string
	//Address          string
	//NewAddress       string
	//	Used             int
}

const (
	HostRoot = "/iptohost"
	EnvRoot  = "/envs"
)

type EnvsInterface interface {
	List() (*Detail, error)
	Get(name string) (*Detail, error)
	Delete(name string) error
	DeleteAll() error
	Create(app AppCreateRequest) (*Detail, error)
	Update(name string, replicas int) (*Detail, error)
}
type envs struct {
	C       *UserClient
	Cluster string
}

func newenvs(client *UserClient, cluster string) *envs {
	return &envs{client, cluster}
}
func (key Env) Validate() error {
	var validationError ValidationError
	if key.ContainerVersion == "" {
		validationError = validationError.Append(ErrInvalidField{"TomcatV"})
	}

	if key.BuildEnv == "" {
		validationError = validationError.Append(ErrInvalidField{"JdkV"})
	}
	if key.Target == "" {
		validationError = validationError.Append(ErrInvalidField{"JdkV"})
	}
	if key.Instance == "" {
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
func IptoPath(ip string) (string, error) {
	resp, err := EtcdClient.Get(path.Join(HostRoot, ip), false, false)
	if err != nil {
		return "", err
	}
	value := strings.Split(resp.Node.Value, ".")
	if len(value) != 2 {
		return "", fmt.Errorf("ip's host not right :%s", resp.Node.Value)
	}
	clusterpath := path.Join(value[1], value[0])
	return clusterpath, nil

}
func (e *envs) EnvPath(envName string) string {
	return path.Join(e.C.UserName, e.Cluster, envName)
}
func (e *envs) Create(env *Env) error {
	data, _ := json.Marshal(env)
	_, err = EtcdClient.Create(e.EnvPath(env.Name), string(data), 0)
	if err != nil {
		return err
	}
	err = IdPools.CreateIdPool(ip, env.Name)
	if err != nil {
		return err
	}
	return nil
}
func GetEnv(ip string, envName string) (*Env, error) {
	clusterpath, err := IptoPath(ip)
	if err != nil {
		return nil, err
	}
	response, err := EtcdClient.Get(EnvPath(clusterpath, envName), false, false)
	if err != nil {
		return nil, err
	}
	var env = Env{}
	err = json.Unmarshal([]byte(response.Node.Value), &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
func DeleteEnv(ip string, envName string) error {
	clusterpath, err := IptoPath(ip)
	if err != nil {
		return err
	}
	_, err = EtcdClient.Delete(EnvPath(clusterpath, envName), false)
	err = IdPools.DeleteIdPool(ip, envName)
	return err

}
func ListEnv(ip string) ([]*Env, error) {
	clusterpath, err := IptoPath(ip)
	if err != nil {
		return nil, err
	}
	response, err := EtcdClient.Get(EnvPath(clusterpath, ""), false, true)
	if err != nil {
		return nil, err
	}
	var envs = []*Env{}
	for _, v := range response.Node.Nodes {
		var env = Env{}
		err = json.Unmarshal([]byte(v.Value), &env)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &env)
	}
	return envs, nil
}
func UpdateEnv(ip string, envname string, env *Env) error {
	data, _ := json.Marshal(env)
	clusterpath, err := IptoPath(ip)
	if err != nil {
		return err
	}
	_, err = EtcdClient.Update(EnvPath(clusterpath, env.Name), string(data), 0)
	if err != nil {
		return err
	}
	return nil
}
