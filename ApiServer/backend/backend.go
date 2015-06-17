package backend

import (
	"K8APITransform/ApiServer/models"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"net/url"
	"path"
)

type Backend struct {
	*client.Client
}

func NewBackend(host string, apiVersion string) (*Backend, error) {
	Client, err := client.New(&Config{Host: host, Version: apiVersion})
	if err != nil {
		return nil, err
	}
	return &Backend{Client}, nil
}
func (c *Backend) Applications(env string) ApplicationInterface {
	return newApplications(c, env)
}
