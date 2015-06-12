package models

type Detail struct {
	Name       string   `json:"name,omitempty"`
	PodName    string   `json:"podName,omitempty"`
	Status     int      `json:"status,omitempty"`
	NodeType   int      `json:"nodeType,omitempty"`
	AppVersion string   `json:"appVersion,omitempty"`
	AppName    string   `json:"appName,omitempty"`
	IP         string   `json:"IP,omitempty"`
	Context    []Detail `json:"context,omitempty"`
	Children   []Detail `json:"children,omitempty"`
}
