package models

//type AppViewRequest struct {
//	Namespace string
//}

//type AppCreateRequest struct {
//	// Name is unique within a service.
//	Name string `json:"name,omitempty"`
//}

type AppGetResponse struct {
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

	// Labels are key value pairs that may be used to scope and select individual resources.
	// Label keys are of the form:
	//     label-key ::= prefixed-name | name
	//     prefixed-name ::= prefix '/' name
	//     prefix ::= DNS_SUBDOMAIN
	//     name ::= DNS_LABEL
	// The prefix is optional.  If the prefix is not specified, the key is assumed to be private
	// to the user.  Other system components that wish to use labels must specify a prefix.  The
	// "kubernetes.io/" prefix is reserved for use by kubernetes components.
	// TODO: replace map[string]string with labels.LabelSet type
	Labels map[string]string `json:"labels,omitempty"`

	// Spec defines the desired quota
	Spec ServiceSpec `json:"spec,omitempty"`

	// Status defines the actual enforced quota and its current usage
	Status ServiceStatus `json:"status,omitempty"`
}

type AppGetAllResponseItem struct {
	Name string `json:"name,omitempty"`
	//CreationTimestamp *Time `json:"creationTimestamp,omitempty"`
	//Status ServiceStatus `json:"status,omitempty"`
}

type AppGetAllResponse struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind,omitempty"`

	Items []AppGetAllResponseItem `json:"items"`
}
