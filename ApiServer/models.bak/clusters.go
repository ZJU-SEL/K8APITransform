package models

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

type UserClient struct {
	UserName string
}

func NewUserClient(userName string) *UserClient {
	return &Client{userName}
}
func (c *UserClient) Envs(cluster string) EnvsInterface {
	return newenvs(c, cluster)
}
func (c *UserClient) Applications(cluster string) AppsInterface {
	//EtcdClient.Get()
	config := &client.Config{
		Host:    "https://" + cluster + "." + c.UserName + PORT,
		Version: apiVersion,
		TLSClientConfig: client.TLSClientConfig{
			// Server requires TLS client certificate authentication
			//CertFile: certDir + "/server.crt",
			// Server requires TLS client certificate authentication
			//KeyFile: certDir + "/server.key",
			// Trusted root certificates for server
			CAFile: "certs/" + c.UserName + "/" + cluster + "/ca.crt",
		},
		BearerToken: "abcdTOKEN1234",
	}

	Client, err := client.New(config)
	return newApps(c, env)
}
