package models

import (
	//"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/serviceaccount"
	"github.com/coreos/go-etcd/etcd"
	"github.com/dgrijalva/jwt-go"
	"log"
	"path"
	"sort"
	"strconv"
	"strings"
)

const (
	Issuer                  = "kubernetes/serviceaccount"
	SubjectClaim            = "sub"
	IssuerClaim             = "iss"
	ServiceAccountNameClaim = "kubernetes.io/serviceaccount/service-account.name"
	ServiceAccountUIDClaim  = "kubernetes.io/serviceaccount/service-account.uid"
	SecretNameClaim         = "kubernetes.io/serviceaccount/secret.name"
	NamespaceClaim          = "kubernetes.io/serviceaccount/namespace"
)

type Env struct {
	ContainerVersion string `json:"containerVersion,omitempty"`
	BuildEnv         string `json:"buildEnv,omitempty"`
	Instance         string `json:"instance,omitempty"`
	Name             string `json:"name,omitempty"`
	Cpu              string `json:"cpu,omitempty"`
	Memory           string `json:"memory,omitempty"`
	Disk             string `json:"disk,omitempty"`
	Target           string `json:"target,omitempty"`
	//UpdateTime       string
	//Version          string
	//NewVersion       string
	//Status           string
	//AppName          string
	//NewAppName       string
	//Address          string
	//NewAddress       string
	//Used             int

}

var (
	EtcdClient *etcd.Client
)
var K8ClientMap = map[string]*client.Client{}

const (
	ApiVersion = "v1beta3"
	AppsRoot   = "applications"
	CertRoot   = "certs"
	HostRoot   = "iptohost"
	IpRoot     = "hosttoip"
	EnvRoot    = "envs"
	PORT       = ":8081"
)

type Envs struct {
	C       *UserClient
	Cluster string
}

func newenvs(client *UserClient, Cluster string) *Envs {
	return &Envs{client, Cluster}
}

func (key Env) Validate() error {
	var validationError ValidationError
	if key.ContainerVersion == "" {
		validationError = validationError.Append(ErrInvalidField{"TomcatV"})
	}

	if key.BuildEnv == "" {
		validationError = validationError.Append(ErrInvalidField{"BuildEnv"})
	}
	if key.Instance == "" {
		validationError = validationError.Append(ErrInvalidField{"Instance"})
	}
	if replcas, _ := strconv.Atoi(key.Instance); replcas == 0 {
		validationError = validationError.Append(ErrInvalidField{"Instance is zero"})
	}
	if key.Name == "" {
		validationError = validationError.Append(ErrInvalidField{"Name"})
	}
	if !validationError.Empty() {
		return validationError
	}

	return nil
}
func (e *Envs) NewClient() (*client.Client, error) {
	if client, exist := K8ClientMap[e.Cluster+"."+e.C.UserName]; exist {
		return client, nil
	}
	token := jwt.New(jwt.SigningMethodRS256)
	// Set some claims
	token.Claims[SubjectClaim] = "zjusel"
	token.Claims[IssuerClaim] = Issuer
	token.Claims[ServiceAccountNameClaim] = "zjusel"
	token.Claims[ServiceAccountUIDClaim] = "123456789sel"
	token.Claims[SecretNameClaim] = "zjusel"
	token.Claims[NamespaceClaim] = "zjusel"

	// Sign and get the complete encoded token as a string

	//publicKey, err := serviceaccount.ReadPublicKey("ca.key")

	privateKey, err := serviceaccount.ReadPrivateKey("certs/sa.key")
	//fmt.Println(privateKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	//fmt.Println(tokenString)
	config := &client.Config{
		Host:     "https://" + e.Cluster + "." + e.C.UserName + PORT,
		Version:  ApiVersion,
		Insecure: true,
		//TLSClientConfig: client.TLSClientConfig{
		//	// Server requires TLS client certificate authentication
		//	//CertFile: certDir + "/server.crt",
		//	// Server requires TLS client certificate authentication
		//	//KeyFile: certDir + "/server.key",
		//	// Trusted root certificates for server
		//	CAFile: path.Join(CertRoot, e.C.UserName, e.Cluster, "ca.crt"),
		//},
		BearerToken: "abcdTOKEN1234",
	}

	Client, err := client.New(config)
	if err == nil {
		K8ClientMap[e.Cluster+"."+e.C.UserName] = Client
	}
	return Client, err
}
func Ip2UC(ip string) (string, string, error) {
	resp, err := EtcdClient.Get(path.Join(HostRoot, ip), false, false)
	if err != nil {
		return "", "", err
	}
	value := strings.Split(resp.Node.Value, ".")
	if len(value) != 2 {
		return "", "", fmt.Errorf("ip's host not right :%s", resp.Node.Value)
	}
	//clusterpath := path.Join(value[1], value[0])
	return value[1], value[0], nil

}
func Host2Ip(userName string, cluster string) (string, error) {
	resp, err := EtcdClient.Get(path.Join(IpRoot, userName, cluster), false, false)
	if err != nil {
		return "", err
	}
	value := resp.Node.Value
	return value, nil
}
func (e *Envs) ToMapArray(envs []*Env) []map[string]string {
	data, _ := json.Marshal(envs)
	result := []map[string]string{}
	json.Unmarshal(data, &result)
	fmt.Println(result)
	return result
}
func (e *Envs) ToMapMap(envs map[string]*Env) map[string]map[string]string {
	data, _ := json.Marshal(envs)
	result := map[string]map[string]string{}
	json.Unmarshal(data, &result)
	return result
}
func (e *Envs) ToMap(env *Env) map[string]string {
	data, _ := json.Marshal(env)
	result := map[string]string{}
	json.Unmarshal(data, &result)
	fmt.Println(result)
	return result
}
func (e *Envs) Apps(envName string) AppsInterface {
	Client, err := e.NewClient()
	if err != nil {
		log.Println(err.Error())

	}
	return newApps(e, Client, envName)
}
func (e *Envs) EnvPath(envName string) string {
	return path.Join(EnvRoot, e.C.UserName, e.Cluster, envName)
}
func (e *Envs) Create(env *Env) (err error) {
	data, _ := json.Marshal(env)
	_, err = EtcdClient.Create(e.EnvPath(env.Name), string(data), 0)
	if err != nil {
		return err
	}
	err = e.C.IdPools(e.Cluster).Create(env.Name)
	if err != nil {
		EtcdClient.Delete(e.EnvPath(env.Name), false)
		return err
	}
	return nil
}
func (e *Envs) Get(envName string) (*Env, error) {
	response, err := EtcdClient.Get(e.EnvPath(envName), false, false)
	if err != nil {
		return nil, err
	}
	var env = Env{}
	err = json.Unmarshal([]byte(response.Node.Value), &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
func (e *Envs) Delete(envName string) error {
	//err :=
	err := e.Apps(envName).DeleteAll()
	if err != nil {
		return err
	}
	_, err = EtcdClient.Delete(e.EnvPath(envName), false)
	if err != nil {
		return err
	}
	err = e.C.IdPools(e.Cluster).Delete(envName)
	if err != nil {
		return err
	}
	return nil

}
func (e *Envs) DeleteAll() error {
	envs, err := e.List()
	if err != nil {
		return err
	}
	for _, env := range envs {
		e.Delete(env.Name)
	}
	return err
}
func (e *Envs) Update(envName string, env *Env) error {
	data, _ := json.Marshal(env)
	_, err := EtcdClient.Update(e.EnvPath(envName), string(data), 0)
	if err != nil {
		return err
	}
	return nil
}
func (e *Envs) List() ([]*Env, error) {
	response, err := EtcdClient.Get(e.EnvPath(""), false, true)
	if err != nil {
		return nil, err
	}
	var envs = []*Env{}
	for _, v := range response.Node.Nodes {
		var env = Env{}
		err = json.Unmarshal([]byte(v.Value), &env)
		if err != nil {
			return nil, err
		}
		envs = append(envs, &env)
	}
	return envs, nil
}
func (e *Envs) ListMap() (map[string]*Env, error) {
	response, err := EtcdClient.Get(e.EnvPath(""), false, true)
	if err != nil {
		return nil, err
	}
	var envs = map[string]*Env{}
	for _, v := range response.Node.Nodes {
		var env = Env{}
		err = json.Unmarshal([]byte(v.Value), &env)
		if err != nil {
			return nil, err
		}
		envs[path.Base(v.Key)] = &env
	}
	return envs, nil
}
func (e *Envs) GetInfo(envName string) (map[string]string, error) {
	info, err := e.Get(envName)
	if err != nil {
		return nil, err
	}
	Client, err := e.NewClient()
	if err != nil {
		log.Println(err.Error())
	}
	env := e.ToMap(info)
	sevicelist, err := Client.Services("default").List(labels.SelectorFromSet(map[string]string{"env": env["name"]}))
	sort.Sort(&Services{*sevicelist})
	if err != nil {
		return nil, err
	}
	env["target"] = e.Cluster
	//sevicelist, err := Client.Services("default").List(labels.SelectorFromSet(map[string]string{"env": env["name"]}))
	rclist, err := Client.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": env["name"]}))
	if err != nil {
		return nil, err
	}
	fmt.Println(sevicelist.Items)
	switch len(sevicelist.Items) {
	case 0:
		env["updateTime"] = ""
		env["version"] = ""
		env["newVersion"] = ""
		env["status"] = ""
		env["appName"] = ""
		env["newAppName"] = ""
		env["address"] = ""
		env["newAddress"] = ""
		break
	case 1:
		env["appName"] = sevicelist.Items[0].ObjectMeta.Labels["name"]
		data := strings.Split(env["appName"], "-")
		env["version"] = data[len(data)-1]
		env["address"] = sevicelist.Items[0].Spec.ClusterIP + ":8080"
		env["updateTime"] = fmt.Sprintf("%d", sevicelist.Items[0].ObjectMeta.CreationTimestamp.Unix())
		env["newVersion"] = ""
		env["newAppName"] = ""
		env["newAddress"] = ""
		details, err := newApps(e, Client, env["name"]).Get(sevicelist.Items[0].ObjectMeta.Labels["name"])
		if err != nil {
			return nil, err
		}

		if details.Status == "Stopped" {
			env["status"] = "Stopped"
			break
		}
		count := 0
		for _, v := range details.InstanceInfo {
			if v.Status == "Running" {
				count++
			}
		}
		if count == 0 {
			env["status"] = "Unavailable"
		} else if count != rclist.Items[0].Spec.Replicas {
			env["status"] = "Pending"
		} else {
			env["status"] = "Running"
		}

		env["instance"] = fmt.Sprintf("%d", rclist.Items[0].Spec.Replicas)
		break
	case 2:
		env["appName"] = sevicelist.Items[0].ObjectMeta.Labels["name"]
		data := strings.Split(env["appName"], "-")
		env["version"] = data[len(data)-1]
		env["address"] = sevicelist.Items[0].Spec.ClusterIP + ":8080"

		env["newAppName"] = sevicelist.Items[1].ObjectMeta.Labels["name"]
		data = strings.Split(env["newAppName"], "-")
		env["newVersion"] = data[len(data)-1]
		env["newAddress"] = sevicelist.Items[1].Spec.ClusterIP + ":8080"
		env["updateTime"] = fmt.Sprintf("%d", sevicelist.Items[1].ObjectMeta.CreationTimestamp.Unix())
		stopped := 0
		count := 0
		for i := 0; i < 2; i++ {
			details, err := newApps(e, Client, env["name"]).Get(sevicelist.Items[i].ObjectMeta.Labels["name"])
			if err != nil {
				return nil, err
			}
			if details.Status == "Stopped" {
				stopped++
			}
			for _, v := range details.InstanceInfo {
				if v.Status == "Running" {
					count++
				}
			}
		}
		if stopped == 2 {
			env["status"] = "Stopped"
		} else if count == 0 {
			env["status"] = "Unavailable"
		} else if count != rclist.Items[0].Spec.Replicas+rclist.Items[1].Spec.Replicas {
			env["status"] = "Pending"
		} else {
			env["status"] = "Running"
		}
		env["instance"] = fmt.Sprintf("%d", rclist.Items[0].Spec.Replicas+rclist.Items[1].Spec.Replicas)
		break
	default:
	}
	return env, nil
}
func (e *Envs) ListInfo() ([]map[string]string, error) {
	infos, err := e.List()
	if err != nil {
		return nil, err
	}
	Client, err := e.NewClient()
	if err != nil {
		log.Println(err.Error())
	}
	envs := e.ToMapArray(infos)
	for _, env := range envs {
		env["target"] = e.Cluster
		sevicelist, err := Client.Services("default").List(labels.SelectorFromSet(map[string]string{"env": env["name"]}))
		sort.Sort(&Services{*sevicelist})
		rclist, err := Client.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"env": env["name"]}))
		if err != nil {
			return nil, err
		}
		switch len(sevicelist.Items) {
		case 0:
			env["updateTime"] = ""
			env["version"] = ""
			env["newVersion"] = ""
			env["status"] = ""
			env["appName"] = ""
			env["newAppName"] = ""
			env["address"] = ""
			env["newAddress"] = ""
			break
		case 1:
			env["appName"] = sevicelist.Items[0].ObjectMeta.Labels["name"]
			data := strings.Split(env["appName"], "-")
			env["version"] = data[len(data)-1]
			env["address"] = sevicelist.Items[0].Spec.ClusterIP + ":8080"
			env["updateTime"] = fmt.Sprintf("%d", sevicelist.Items[0].ObjectMeta.CreationTimestamp.Unix())
			env["newVersion"] = ""
			env["newAppName"] = ""
			env["newAddress"] = ""
			details, err := newApps(e, Client, env["name"]).Get(sevicelist.Items[0].ObjectMeta.Labels["name"])
			if err != nil {
				return nil, err
			}

			if details.Status == "Stopped" {
				env["status"] = "Stopped"
				break
			}
			count := 0
			for _, v := range details.InstanceInfo {
				if v.Status == "Running" {
					count++
				}
			}
			if count == 0 {
				env["status"] = "Unavailable"
			} else if count != rclist.Items[0].Spec.Replicas {
				env["status"] = "Pending"
			} else {
				env["status"] = "Running"
			}

			env["instance"] = fmt.Sprintf("%d", rclist.Items[0].Spec.Replicas)
			break
		case 2:
			env["appName"] = sevicelist.Items[0].ObjectMeta.Labels["name"]
			data := strings.Split(env["appName"], "-")
			env["version"] = data[len(data)-1]
			env["address"] = sevicelist.Items[0].Spec.ClusterIP + ":8080"

			env["newAppName"] = sevicelist.Items[1].ObjectMeta.Labels["name"]
			data = strings.Split(env["newAppName"], "-")
			env["newVersion"] = data[len(data)-1]
			env["newAddress"] = sevicelist.Items[1].Spec.ClusterIP + ":8080"
			env["updateTime"] = fmt.Sprintf("%d", sevicelist.Items[1].ObjectMeta.CreationTimestamp.Unix())
			stopped := 0
			count := 0
			for i := 0; i < 2; i++ {
				details, err := newApps(e, Client, env["name"]).Get(sevicelist.Items[i].ObjectMeta.Labels["name"])
				if err != nil {
					return nil, err
				}
				if details.Status == "Stopped" {
					stopped++
				}
				for _, v := range details.InstanceInfo {
					if v.Status == "Running" {
						count++
					}
				}
			}
			if stopped == 2 {
				env["status"] = "Stopped"
			} else if count == 0 {
				env["status"] = "Unavailable"
			} else if count != rclist.Items[0].Spec.Replicas+rclist.Items[1].Spec.Replicas {
				env["status"] = "Pending"
			} else {
				env["status"] = "Running"
			}
			env["instance"] = fmt.Sprintf("%d", rclist.Items[0].Spec.Replicas+rclist.Items[1].Spec.Replicas)
			break
		default:
		}
	}
	return envs, nil
}
func (e *Envs) ListInfo2() ([]map[string]string, error) {
	infos, err := e.ListMap()
	if err != nil {
		return nil, err
	}
	Client, err := e.NewClient()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	envs := e.ToMapMap(infos)
	sevicelist, err := Client.Services("default").List(nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	sort.Sort(&Services{*sevicelist})
	rclist, err := Client.ReplicationControllers("default").List(nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	podlist, err := Client.Pods("default").List(nil, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	envstatus := map[string]map[string]int{}
	for k, env := range envs {
		env["updateTime"] = ""
		env["version"] = ""
		env["newVersion"] = ""
		env["status"] = ""
		env["appName"] = ""
		env["newAppName"] = ""
		env["address"] = ""
		env["newAddress"] = ""
		envstatus[k] = map[string]int{}
	}
	for _, v := range sevicelist.Items {
		name := v.ObjectMeta.Labels["env"]
		if env, exist := envs[name]; exist {
			if r, exist := env["hasone"]; exist {
				if r == "1" {
					env["newAppName"] = v.ObjectMeta.Labels["name"]
					data := strings.Split(v.ObjectMeta.Labels["name"], "-")
					env["newVersion"] = data[len(data)-1]
					env["newAddress"] = v.Spec.ClusterIP + ":8080"
					env["updateTime"] = fmt.Sprintf("%d", v.ObjectMeta.CreationTimestamp.Unix())
					env["hasone"] = "2"
				}
			} else {
				env["appName"] = v.ObjectMeta.Labels["name"]
				data := strings.Split(v.ObjectMeta.Labels["name"], "-")
				env["version"] = data[len(data)-1]
				env["address"] = v.Spec.ClusterIP + ":8080"
				env["updateTime"] = fmt.Sprintf("%d", v.ObjectMeta.CreationTimestamp.Unix())
				env["hasone"] = "1"
			}
		}
	}
	//fmt.Println(envs)
	for _, v := range rclist.Items {
		name := v.ObjectMeta.Labels["env"]
		if _, exist := envs[name]; exist {
			if v.ObjectMeta.Labels["stopped"] != "" {
				stopped, _ := strconv.Atoi(v.ObjectMeta.Labels["stopped"])
				envstatus[name]["stopped"] += stopped
				envstatus[name]["replicas"] += stopped
			} else {
				envstatus[name]["replicas"] += v.Spec.Replicas
			}
		}
	}
	for _, v := range podlist.Items {
		name := v.ObjectMeta.Labels["env"]
		if _, exist := envs[name]; exist {
			if string(v.Status.Phase) == "Running" {
				envstatus[name]["running"] += 1
			} else {
				envstatus[name]["pending"] += 1
			}
		}
	}
	result := []map[string]string{}
	for k, env := range envs {
		env["target"] = e.Cluster
		fmt.Println(env)
		if _, exist := env["hasone"]; !exist {
			fmt.Println("not exist")
			env["updateTime"] = ""
			env["version"] = ""
			env["newVersion"] = ""
			env["status"] = ""
			env["appName"] = ""
			env["newAppName"] = ""
			env["address"] = ""
			env["newAddress"] = ""
		} else {
			env["instance"] = fmt.Sprintf("%d", envstatus[k]["replicas"])
			if exist && envstatus[k]["replicas"] == envstatus[k]["stopped"] && envstatus[k]["stopped"] != 0 {
				env["status"] = "Stopped"

			} else {
				if envstatus[k]["running"] == 0 {
					env["status"] = "Unavailable"
				} else {
					if envstatus[k]["running"]+envstatus[k]["stopped"] >= envstatus[k]["replicas"] {
						env["status"] = "Running"
					} else {
						env["status"] = "Pending"
					}
				}
			}
			delete(env, "hasone")
		}
		result = append(result, env)
	}
	return result, nil
}
