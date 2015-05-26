package models

type ResourceQuotaCreateRequest struct {
	// Name is unique within a resourcequota.
	Name string `json:"name,omitempty"`
}

type ResourceQuotaGetResponse struct {
	Kind string `json:"kind,omitempty"`

	Name string `json:"name,omitempty"`

	// Namespace defines the space within which name must be unique. An empty namespace is
	// equivalent to the "default" namespace, but "default" is the canonical representation.
	// Not all objects are required to be scoped to a namespace - the value of this field for
	// those objects will be empty.
	Namespace string `json:"namespace,omitempty"`

	// CreationTimestamp is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	CreationTimestamp *Time `json:"creationTimestamp,omitempty"`

	// Spec defines the desired quota
	Spec ResourceQuotaSpec `json:"spec,omitempty"`

	// Status defines the actual enforced quota and its current usage
	Status ResourceQuotaStatus `json:"status,omitempty"`
}

type ResourceQuotaGetAllResponseItem struct {
	Name string `json:"name,omitempty"`
	//CreationTimestamp *Time `json:"creationTimestamp,omitempty"`
	//Status ResourceQuotaStatus `json:"status,omitempty"`
}

type ResourceQuotaGetAllResponse struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind,omitempty"`

	Items []ResourceQuotaGetAllResponseItem `json:"items"`
}
