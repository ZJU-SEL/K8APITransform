package models

import (
	"encoding/json"
	"fmt"
	api "github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"strconv"
	//"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
	"strings"
)

// PodInterface has methods to work with Pod resources.
type ApplicationInterface interface {
	List() (*Detail, error)
	Get(name string) (*Detail, error)
	Delete(name string) error
	Create(app AppCreateRequest) (*Detail, error)
	Update(name string, replicas int) error
	Restart(name string) error
	//Watch(label labels.Selector, field fields.Selector, resourceVersion string) (watch.Interface, error)
	//Bind(binding *api.Binding) error
	//UpdateStatus(pod *api.Pod) (*api.Pod, error)
}

// pods implements PodsNamespacer interface
type applications struct {
	b   *Backend
	env string
}

// newPods returns a pods
func newApplications(b *Backend, env string) *applications {
	return &applications{
		b:   b,
		env: env,
	}
}
func (a *applications) Create(app AppCreateRequest) (*Detail, error) {
	containerports := []api.ContainerPort{}
	for _, v := range app.ContainerPort {
		containerports = append(containerports, api.ContainerPort{
			ContainerPort: v.Port, //
			Protocol:      api.Protocol(v.Protocol),
		})
	}
	//name := app.Name + "-" + app.Version
	//name = strings.ToLower(name)
	//name = strings.Replace(name, ".", "", -1)
	id, err := IdPools.GetId(a.env)
	if err != nil {
		return nil, err
	}
	labels := map[string]string{"env": a.env, "name": a.env + "-" + app.Name + "-" + app.Version}
	containers := []api.Container{
		api.Container{
			Name:  a.env + "-" + id,
			Image: app.Containerimage,
			Ports: containerports,
		},
	}
	var rc = &api.ReplicationController{
		TypeMeta: api.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   a.env + "-" + id,
			Labels: labels,
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: app.Replicas,
			Selector: labels,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: labels,
				},
				Spec: api.PodSpec{
					Containers: containers,
					Volumes:    app.Volumes,
				},
			},
		},
	}
	rc, err = a.b.ReplicationControllers("default").Create(rc)
	if err != nil {
		return nil, err
	}
	var Ports = []api.ServicePort{}
	for k, v := range app.Ports {
		Ports = append(Ports, api.ServicePort{
			Name:     "default" + strconv.Itoa(k),
			Port:     v.Port,
			Protocol: api.Protocol(v.Protocol),
		})
	}
	var service = &api.Service{
		TypeMeta: api.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   id,
			Labels: labels,
		},
		Spec: api.ServiceSpec{
			Selector: labels,
			Ports:    Ports,
		},
	}
	service, err = a.b.Services("default").Create(service)
	if err != nil {
		return nil, err
	}
	detail, err := a.Get(labels["name"])
	return detail, nil
}
func (a *applications) List() (*Detail, error) {
	label := map[string]string{"env": a.env}
	serviceslist, err := a.b.Services("default").List(labels.SelectorFromSet(label))
	fmt.Println("get servicelists")
	if err != nil {
		return nil, err
	}
	podslist, err := a.b.Pods("default").List(labels.SelectorFromSet(label), nil)
	if err != nil {
		return nil, err
	}
	detail := &Detail{Name: a.env, Status: 1, NodeType: 1, Context: []Detail{}, Children: []Detail{}}
	detail.Children = append(detail.Children, Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []Detail{
			Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []Detail{}, Children: []Detail{}}
	if len(podslist.Items) == 0 {
		e, _ := EtcdClient.Get("envs/"+a.env, false, false)
		e_tmp := AppEnv{}
		json.Unmarshal([]byte(e.Node.Value), &e_tmp)
		num, _ := strconv.Atoi(e_tmp.NodeNum)
		for k := 0; k < num; k++ {
			//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
			tomcat.Context = append(tomcat.Context, Detail{
				Name:     "Node" + strconv.Itoa(k+1),
				NodeType: 3,
			})
		}
	} else {
		for k, v := range podslist.Items {
			status := 0
			if v.Status.Phase == api.PodRunning {
				status = 1
			}
			//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
			tomcat.Context = append(tomcat.Context, Detail{
				Name:       "Node" + strconv.Itoa(k+1),
				AppVersion: v.ObjectMeta.Labels["name"],
				Status:     status,
				NodeType:   3,
			})
		}
	}
	apps := []Detail{}
	for _, v := range serviceslist.Items {
		//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
		apps = append(apps, Detail{
			Name:     v.ObjectMeta.Labels["name"],
			NodeType: 4,
			Status:   1,
			Resource: []Detail{Detail{Name: "IP", Value: v.Spec.ClusterIP + ":8080"}},
		})
	}
	tomcat.Children = append(tomcat.Children, Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, apps...)
	detail.Children = append(detail.Children, tomcat)
	return detail, nil
}
func (a *applications) Get(name string) (*Detail, error) {
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return nil, err
	}
	if len(sevicelist.Items) != 1 {
		return nil, ErrResponse{fmt.Sprintf("a app with %d services", len(sevicelist.Items))}
	}
	service := sevicelist.Items[0]
	podslist, err := a.b.Pods("default").List(labels.SelectorFromSet(service.ObjectMeta.Labels), nil)
	if err != nil {
		return nil, err
	}
	detail := &Detail{Name: a.env, Status: 1, NodeType: 1, Context: []Detail{}, Children: []Detail{}}
	detail.Children = append(detail.Children, Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []Detail{
			Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []Detail{}, Children: []Detail{}}
	warName := ""
	for k, v := range podslist.Items {
		status := 0
		if v.Status.Phase == api.PodRunning {
			status = 1
		}
		info := strings.Split(v.ObjectMeta.Labels["name"], "-")
		warName = info[len(info)-2]
		tomcat.Context = append(tomcat.Context, Detail{
			Name: "Node" + strconv.Itoa(k+1),
			//AppVersion: v.ObjectMeta.Labels["name"],
			AppVersion: info[len(info)-1],
			Status:     status,
			NodeType:   3,
		})
	}
	tomcat.Children = append(tomcat.Children, Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, Detail{
		Name:         service.ObjectMeta.Labels["name"],
		NodeType:     4,
		Status:       1,
		Resource:     []Detail{Detail{Name: "IP", Value: service.Spec.ClusterIP + ":8080"}},
		OriginalName: warName,
	})
	detail.Children = append(detail.Children, tomcat)
	return detail, nil
}
func (a *applications) Delete(name string) error {
	sevicelist, err := a.b.Services("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}
	for _, v := range sevicelist.Items {
		a.b.Services("default").Delete(v.ObjectMeta.Name)
	}
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}
	for _, v := range rclist.Items {
		a.b.ReplicationControllers("default").Delete(v.ObjectMeta.Name)
	}
	return nil
}
func (a *applications) Update(name string, replicas int) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	rc.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		return err
	}
	return nil
}
func (a *applications) Restart(name string) error {
	rclist, err := a.b.ReplicationControllers("default").List(labels.SelectorFromSet(map[string]string{"name": name}))
	if err != nil {
		return err
	}
	if len(rclist.Items) != 1 {
		return ErrResponse{fmt.Sprintf("a app with %d services", len(rclist.Items))}
	}
	rc := rclist.Items[0]
	replicas := rc.Spec.Replicas
	rc.Spec.Replicas = 0
	rcnew, err := a.b.ReplicationControllers("default").Update(&rc)
	if err != nil {
		return err
	}
	rcnew.Spec.Replicas = replicas
	_, err = a.b.ReplicationControllers("default").Update(rcnew)
	if err != nil {
		return err
	}
	return nil
}
