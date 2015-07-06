package models

import (
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
)

type Port struct {
	Port       int
	TargetPort int
	Protocol   string
}
type Containerport struct {
	Port     int
	Protocol string
}
type AppCreateRequest struct {
	Name           string
	Version        string
	Ports          []Port
	Replicas       int
	Runlocal       bool
	Containername  string
	Containerimage string
	Cpu            string
	Memery         string
	Storage        string
	Warpath        string
	ContainerPort  []Containerport
	Volumes        []api.Volume
	PublicIPs      []string
}

func (key AppCreateRequest) Validate() error {
	var validationError ValidationError
	if key.Name == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}

	if key.Containerimage == "" {
		validationError = validationError.Append(ErrInvalidField{"Containerimage"})
	}
	/*
		if len(key.Ports) == 0 {
			validationError = validationError.Append(ErrInvalidField{"Ports"})
		}
	*/

	if len(key.ContainerPort) == 0 {
		validationError = validationError.Append(ErrInvalidField{"ContainerPort"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}
