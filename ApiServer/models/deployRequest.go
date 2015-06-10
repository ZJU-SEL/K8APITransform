package models

type DeployRequest struct {
	EnvName    string
	WarName    string
	AppVersion string
}

func (key DeployRequest) Validate() error {
	var validationError ValidationError
	if key.EnvName == "" {
		validationError = validationError.Append(ErrInvalidField{"TomcatV"})
	}

	if key.WarName == "" {
		validationError = validationError.Append(ErrInvalidField{"JdkV"})
	}
	/*
		if len(key.Ports) == 0 {
			validationError = validationError.Append(ErrInvalidField{"Ports"})
		}
	*/

	if key.AppVersion == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
