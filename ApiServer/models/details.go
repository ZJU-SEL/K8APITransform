package models

type Detail struct {
	Name         string   `json:"name,omitempty"`
	Cpu          string   `json:"cpu,omitempty"`
	Memory       string   `json:"memory,omitempty"`
	Storage      string   `json:"storage,omitempty"`
	Resource     []Detail `json:"resource,omitempty"`
	PodName      string   `json:"podName,omitempty"`
	Status       int      `json:"status,omitempty"`
	NodeType     int      `json:"nodeType,omitempty"`
	AppVersion   string   `json:"appVersion,omitempty"`
	OriginalName string   `json:"originalName,omitempty"`
	AppName      string   `json:"appName,omitempty"`
	IP           string   `json:"IP,omitempty"`
	Context      []Detail `json:"context,omitempty"`
	Value        string   `json:"value,omitempty"`
	Children     []Detail `json:"children,omitempty"`
}
