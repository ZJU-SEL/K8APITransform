package models

import (
	//"K8APITransform/ApiServer/models"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	//"github.com/coreos/go-etcd/etcd"
	"path"
	"sync"
)

const (
	IdpoolRoot = "/idpools"
)

var mutex = &sync.Mutex{}
var IdPools IdPoolsInterface

type IdPoolsInterface interface {
	GetId(env string) (string, error)
	Create(env string) error
	DeleteIdPool(env string) error
}

func NewIdPools(C *UserClient, cluster string) IdPoolsInterface {
	return &idpools{C, cluster}
}

type idpools struct {
	c       *UserClient
	cluster string
}

func (pools *idpools) IppoolPath(env string) string {
	return path.Join(IdpoolRoot, pools.c.UserName, pools.cluster, env)
}
func (pools *idpools) Create(env string) error {
	_, err := EtcdClient.Create(pools.IppoolPath(env), "aaaaaaaaaaaaa", 0)
	return err

}
func (pools *idpools) Delete(ip string, env string) error {
	_, err := EtcdClient.Delete(pools.IppoolPath(env), false)
	return err

}
func (pools *idpools) GetId(env string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	response, err := EtcdClient.Get(pools.IppoolPath(env), false, false)
	if err != nil {
		return "", err
	}
	id := pools.next(response.Node.Value)
	_, err = EtcdClient.Update(pools.IppoolPath(env), id, 0)
	if err != nil {
		return "", err
	}
	return response.Node.Value, nil
}
func (pool *idpools) next(Id string) string {
	id := []byte(Id)
	//t := 0
	for i := 12; i >= 0; i-- {
		if id[i] == byte('z') {
			id[i] = byte('a')
		} else {
			id[i]++
			break
		}
	}
	return string(id)
}
