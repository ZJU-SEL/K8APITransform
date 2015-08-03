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
type AppConfig struct {
	Name           string
	Version        string
	Ports          []Port
	Replicas       int
	Containername  string
	Containerimage string
	Cpu            string
	Memory         string
	Storage        string
	Warpath        string
	ContainerPort  []Containerport
	Volumes        []api.Volume
}

func (key AppConfig) Validate() error {
	var validationError ValidationError
	if key.Name == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}

	if key.Containerimage == "" {
		validationError = validationError.Append(ErrInvalidField{"Containerimage"})
	}

	if len(key.ContainerPort) == 0 {
		validationError = validationError.Append(ErrInvalidField{"ContainerPort"})
	}

	if !validationError.Empty() {
		return validationError
	}

	return nil
}
