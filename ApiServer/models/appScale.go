package models

type AppScale struct {
	EnvName string `json:"envName,omitempty"`
	Name    string `json:"appName,omitempty"`
	Num     string `json:"insNum,omitempty"`
}
