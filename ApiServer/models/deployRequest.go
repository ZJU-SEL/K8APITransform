package models

type DeployRequest struct {
	EnvName        string
	WarName        string
	Version        string
	IsGreyUpdating string
}

func (key DeployRequest) Validate() error {
	var validationError ValidationError
	if key.EnvName == "" {
		validationError = validationError.Append(ErrInvalidField{"EnvName"})
	}

	if key.WarName == "" {
		validationError = validationError.Append(ErrInvalidField{"WarName"})
	}
	/*
		if len(key.Ports) == 0 {
			validationError = validationError.Append(ErrInvalidField{"Ports"})
		}
	*/

	if key.Version == "" {
		validationError = validationError.Append(ErrInvalidField{"Version"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
