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
	Containername  string
	Containerimage string
	Warpath        string
	ContainerPort  []Containerport
	PublicIPs      []string
}
