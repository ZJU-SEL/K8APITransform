package backend

import (
	"fmt"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
)

// PodInterface has methods to work with Pod resources.
type ApplicationInterface interface {
	List() (*models.Detail, error)
	Get(name string) (*models.Detail, error)
	Delete(name string) error
	Create(app models.AppCreateRequest) (*models.Detail, error)
	Update(name string, replicas int) error
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
	return &pods{
		b:   b,
		env: env,
	}
}
func (a *applications) Create(app models.AppCreateRequest) (*models.Detail, error) {
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
	id := IdPools.GetId(a.env)
	labels := map[string]string{"env": a.env, "name": a.env + "-" + app.Name + "-" + app.Version}
	containers := []api.Container{
		models.Container{
			Name:         id,
			Image:        app.Containerimage,
			Ports:        containerports,
			VolumeMounts: volumemount,
		},
	}
	var rc = &api.ReplicationController{
		TypeMeta: api.TypeMeta{
			Kind:       "ReplicationController",
			APIVersion: "v1",
		},
		ObjectMeta: api.ObjectMeta{
			Name:   id,
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
	rc, err := a.b.ReplicationControllers("default").Create(rc)
	if err != nil {
		return nil, err
	}
	var Ports = []api.ServicePort{}
	for k, v := range app.Ports {
		Ports = append(Ports, models.ServicePort{
			Name:     "default" + strconv.Itoa(k),
			Port:     v.Port,
			Protocol: models.Protocol(v.Protocol),
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
	service, err := a.b.Services("default").Create(service)
	if err != nil {
		return nil, err
	}
	detail, err := a.servicetodetail(labels["name"])
	return detail, nil
}
func (a *applications) List() (*models.Detail, error) {
	lables := map[string]string{"env": a.env}
	serviceslist, err := b.Services("default").List(lables, nil)
	if err != nil {
		return nil, err
	}
	podslist, err := a.b.Pods("default").List(lables, nil)
	if err != nil {
		return nil, err
	}
	detail := &models.Detail{Name: env.Name, Status: 1, NodeType: 1, Context: []models.Detail{}, Children: []models.Detail{}}
	detail.Children = append(detail.Children, models.Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []models.Detail{
			models.Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := models.Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []models.Detail{}, Children: []models.Detail{}}
	if len(podslist.Items) == 0 {
		num, _ := strconv.Atoi(env.NodeNum)
		for k := 0; k < num; k++ {
			//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
			tomcat.Context = append(tomcat.Context, models.Detail{
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
			tomcat.Context = append(tomcat.Context, models.Detail{
				Name:       "Node" + strconv.Itoa(k+1),
				AppVersion: v.ObjectMeta.Labels["name"],
				Status:     status,
				NodeType:   3,
			})
		}
	}
	apps := []models.Detail{}
	for _, v := range serviceslist.Items {
		//names := strings.Split(v.ObjectMeta.Labels["name"], "-")
		apps = append(apps, models.Detail{
			Name:     v.ObjectMeta.Labels["name"],
			NodeType: 4,
			Status:   1,
			Resource: []models.Detail{models.Detail{Name: "IP", Value: v.Spec.ClusterIP + ":8080"}},
		})
	}
	tomcat.Children = append(tomcat.Children, models.Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []models.Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, apps...)
	detail.Children = append(detail.Children, tomcat)
	return detail, nil
}
func (a *applications) Get(name string) (*models.Detail, error) {
	sevicelist, err := a.b.Services("default").List(map[string]string{"name": name}, nil)
	if err != nil {
		return nil, err
	}
	if len(sevicelist.Item) != 1 {
		return nil, models.ErrResponse(fmt.Sprintf("a app with %d services", len(sevicelist.Item)))
	}
	service := sevicelist.Item[0]
	podslist, err := a.b.Pods("default").List(service.ObjectMeta.Labels, nil)
	if err != nil {
		return nil, err
	}
	detail := &models.Detail{Name: env.Name, Status: 1, NodeType: 1, Context: []models.Detail{}, Children: []models.Detail{}}
	detail.Children = append(detail.Children, models.Detail{
		Name:     "Nginx",
		Status:   1,
		NodeType: 2,
		Context: []models.Detail{
			models.Detail{
				Name:     "Node1",
				NodeType: 2,
			},
		},
	})
	tomcat := models.Detail{Name: "tomcat", Status: 1, NodeType: 2, Context: []models.Detail{}, Children: []models.Detail{}}
	for k, v := range podslist.Items {
		status := 0
		if v.Status.Phase == api.PodRunning {
			status = 1
		}

		tomcat.Context = append(tomcat.Context, models.Detail{
			Name: "Node" + strconv.Itoa(k+1),
			//AppVersion: v.ObjectMeta.Labels["name"],
			AppVersion: deployReq.AppVersion,
			Status:     status,
			NodeType:   3,
		})
	}
	tomcat.Children = append(tomcat.Children, models.Detail{
		Name:     "应用",
		NodeType: 3,
		Context:  []models.Detail{},
	})
	tomcat.Children[0].Context = append(tomcat.Children[0].Context, models.Detail{
		Name:         service.ObjectMeta.Labels["name"],
		NodeType:     4,
		Status:       1,
		Resource:     []models.Detail{models.Detail{Name: "IP", Value: service.Spec.ClusterIP + ":8080"}},
		OriginalName: deployReq.WarName,
	})
	detail.Children = append(detail.Children, tomcat)
}
func (a *applications) Delete(name string) error {
	sevicelist, err := a.b.Services("default").List(map[string]string{"name": name}, nil)
	if err != nil {
		return err
	}
	for _, v := range sevicelist.Item {
		a.b.Services("default").Delete(v.ObjectMeta.Name)
	}
	rclist, err := a.b.ReplicationControllers("default").List(map[string]string{"name": name}, nil)
	if err != nil {
		return err
	}
	for _, v := range rclist.Item {
		a.b.ReplicationControllers("default").Delete(v.ObjectMeta.Name)
	}
	return nil
}
func (a *applications) Update(name string, replcas int) error {
	rclist, err := a.b.ReplicationControllers("default").List(map[string]string{"name": name}, nil)
	if err != nil {
		return err
	}
	if len(rclist.Item) != 1 {
		return nil, models.ErrResponse(fmt.Sprintf("a app with %d services", len(rclist.Item)))
	}
	rc := rclist.Item[0]
	rc.TypeMeta.Kind = "ReplicationController"
	rc.TypeMeta.APIVersion = a.b.Client.RESTClient.APIVersion()
	_, err = a.b.ReplicationControllers("default").Update(rc)
	if err != nil {
		return err
	}
	return nil
}
