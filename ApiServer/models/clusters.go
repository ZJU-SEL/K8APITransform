package models

import (
//"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

type UserClient struct {
	UserName string
}

func NewUserClient(userName string) *UserClient {
	return &UserClient{userName}
}
func (c *UserClient) Envs(cluster string) *Envs {
	return newenvs(c, cluster)
}
func (c *UserClient) IdPools(cluster string) IdPoolsInterface {
	return newIdPools(c, cluster)
}
