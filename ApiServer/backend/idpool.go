package backend

import (
	"K8APITransform/ApiServer/models"
	"github.com/coreos/go-etcd/etcd"
	"sync"
)

var mutex = &sync.Mutex{}
var IdPools IdPoolsInterface

type IdPoolsInterface interface {
	GetId(env string) string
	CreateIdPool(env string) error
}

func NewIdPools() IdPoolsInterface {
	response, err := models.EtcdClient.Get("/idpools", false, false)
	if err != nil {
		models.EtcdClient.CreateDir("/idpools", 0)
	}
	return &idpool{}
}

type idpools struct {
}

func (pools *idpools) CreateIdPool(env string) error {
	response, err := models.EtcdClient.Get("/idpools/"+env, false, false)
	if err != nil {
		_, err := models.EtcdClient.Create("/idpools"+env, "aaaaaaaaaaaaa", false)
		return err
	}
	return nil
}
func (pools *idpools) GetId(env string) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()
	response, err := models.EtcdClient.Get("/idpools/"+env, false, false)
	if err != nil {
		return "", err
	}
	id := pool.next(response.Node.Value)
	_, err = models.EtcdClient.Update("/idpools/"+env, id, false)
	if err != nil {
		return "", err
	}
	return response.Node.Value, nil
}
func (pool *idpool) next(Id string) string {
	id := []byte(Id)
	t := 0
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
