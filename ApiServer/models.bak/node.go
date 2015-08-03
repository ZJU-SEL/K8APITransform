package models

type NodePhase string

// These are the valid phases of node.
const (
	// NodePending means the node has been created/added by the system, but not configured.
	NodePending NodePhase = "Pending"
	// NodeRunning means the node has been configured and has Kubernetes components running.
	NodeRunning NodePhase = "Running"
	// NodeTerminated means the node has been removed from the cluster.
	NodeTerminated NodePhase = "Terminated"
)

type NodeConditionType string

// These are valid conditions of node. Currently, we don't have enough information to decide
// node condition. In the future, we will add more. The proposed set of conditions are:
// NodeReady, NodeReachable
const (
	// NodeReady means kubelet is healthy and ready to accept pods.
	NodeReady NodeConditionType = "Ready"
)

type NodeCondition struct {
	Type               NodeConditionType `json:"type"`
	Status             ConditionStatus   `json:"status"`
	LastHeartbeatTime  Time              `json:"lastHeartbeatTime,omitempty"`
	LastTransitionTime Time              `json:"lastTransitionTime,omitempty"`
	Reason             string            `json:"reason,omitempty"`
	Message            string            `json:"message,omitempty"`
}
type NodeAddressType string

// These are valid address types of node. NodeLegacyHostIP is used to transit
// from out-dated HostIP field to NodeAddress.
const (
	NodeLegacyHostIP NodeAddressType = "LegacyHostIP"
	NodeHostName     NodeAddressType = "Hostname"
	NodeExternalIP   NodeAddressType = "ExternalIP"
	NodeInternalIP   NodeAddressType = "InternalIP"
)

type NodeAddress struct {
	Type    NodeAddressType `json:"type"`
	Address string          `json:"address"`
}

// NodeSystemInfo is a set of ids/uuids to uniquely identify the node.
type NodeSystemInfo struct {
	// MachineID is the machine-id reported by the node
	MachineID string `json:"machineID"`
	// SystemUUID is the system-uuid reported by the node
	SystemUUID string `json:"systemUUID"`
	// BootID is the boot-id reported by the node
	BootID string `json:"bootID"`
	// Kernel version reported by the node
	KernelVersion string `json:"kernelVersion""`
	// OS image used reported by the node
	OsImage string `json:"osImage"`
	// Container runtime version reported by the node
	ContainerRuntimeVersion string `json:"containerRuntimeVersion"`
	// Kubelet version reported by the node
	KubeletVersion string `json:"kubeletVersion"`
	// Kube-proxy version reported by the node
	KubeProxyVersion string `json:"kubeProxyVersion"`
}

// NodeStatus is information about the current status of a node.
type NodeStatus struct {
	// Capacity represents the available resources of a node.
	Capacity ResourceList `json:"capacity,omitempty"`
	// NodePhase is the current lifecycle phase of the node.
	Phase NodePhase `json:"phase,omitempty"`
	// Conditions is an array of current node conditions.
	Conditions []NodeCondition `json:"conditions,omitempty"`
	// Queried from cloud provider, if available.
	Addresses []NodeAddress `json:"addresses,omitempty"`
	// NodeSystemInfo is a set of ids/uuids to uniquely identify the node
	NodeInfo NodeSystemInfo `json:"nodeInfo,omitempty"`
}

// NodeSpec describes the attributes that a node is created with.
type NodeSpec struct {
	// PodCIDR represents the pod IP range assigned to the node
	// Note: assigning IP ranges to nodes might need to be revisited when we support migratable IPs.
	PodCIDR string `json:"podCIDR,omitempty"`

	// External ID of the node assigned by some machine database (e.g. a cloud provider)
	ExternalID string `json:"externalID,omitempty"`

	// Unschedulable controls node schedulability of new pods. By default node is schedulable.
	Unschedulable bool `json:"unschedulable,omitempty"`
}

// Node is a worker node in Kubernetes
// The name of the node according to etcd is in ObjectMeta.Name.
type Node struct {
	TypeMeta   `json:",inline"`
	ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the behavior of a node.
	Spec NodeSpec `json:"spec,omitempty"`

	// Status describes the current status of a Node
	Status NodeStatus `json:"status,omitempty"`
}

// NodeList is a list of nodes.
type NodeList struct {
	TypeMeta `json:",inline"`
	ListMeta `json:"metadata,omitempty"`

	Items []Node `json:"items"`
}
