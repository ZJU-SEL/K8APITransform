package models

type LimitRangeCreateRequest struct {
	// Name is unique within a limitrange.
	Name string `json:"name,omitempty"`
}

type LimitRangeGetResponse struct {
	Kind string `json:"kind,omitempty"`

	Name string `json:"name,omitempty"`

	// Namespace defines the space within which name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	Namespace string `json:"namespace,omitempty"`

	// CreationTimestamp is a timestamp representing the server time when this object was created.
	CreationTimestamp *Time `json:"creationTimestamp,omitempty"`

	// Spec defines the limits enforced
	Spec LimitRangeSpec `json:"spec,omitempty"`
}

type LimitRangeGetAllResponseItem struct {
	Name string `json:"name,omitempty"`
	//CreationTimestamp *Time `json:"creationTimestamp,omitempty"`
	//Status LimitRangeStatus `json:"status,omitempty"`
}

type LimitRangeGetAllResponse struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind,omitempty"`

	Items []LimitRangeGetAllResponseItem `json:"items"`
}
