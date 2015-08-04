package models

type InstanceInfo struct {
	Status string `json:"status,omitempty"`
	NodeIp string `json:"nodeIp,omitempty"`
	Id     string `json:"id,omitempty"`
}
type Detail struct {
	Name         string          `json:"name,omitempty"`
	Cpu          string          `json:"cpu,omitempty"`
	Memory       string          `json:"memory,omitempty"`
	Disk         string          `json:"disk,omitempty"`
	InstanceInfo []*InstanceInfo `json:"instanceInfo"`
	Resource     []Detail        `json:"resource,omitempty"`
	PodName      string          `json:"podName,omitempty"`
	Instance     int             `json:"instance"`
	Status       string          `json:"status,omitempty"`
	NodeType     int             `json:"nodeType,omitempty"`
	AppVersion   string          `json:"appVersion,omitempty"`
	OriginalName string          `json:"originalName,omitempty"`
	AppName      string          `json:"appName,omitempty"`
	IP           string          `json:"IP,omitempty"`
	Context      []Detail        `json:"context,omitempty"`
	Value        string          `json:"value,omitempty"`
	Children     []Detail        `json:"children,omitempty"`
}
