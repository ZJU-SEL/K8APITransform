package models

var Appinfo = AppsInfo{}

type AppMetaInfo struct {
	Name     string
	Replicas int
	Status   int
}
type NamespaceInfo map[string]*AppMetaInfo
type AppsInfo map[string]NamespaceInfo

func init() {
	Appinfo["default"] = NamespaceInfo{}
}
