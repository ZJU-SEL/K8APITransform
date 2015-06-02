package models

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
	Ports          []Port
	Replicas       int
	Runlocal       bool
	Containername  string
	Containerimage string
	Warpath        string
	ContainerPort  []Containerport
	Volumes        []Volume
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
