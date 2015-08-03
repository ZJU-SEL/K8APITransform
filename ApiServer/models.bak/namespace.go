package models

type NamespaceCreateRequest struct {
	// Name is unique within a namespace.
	Name string `json:"name,omitempty"`
}

type NamespaceGetResponse struct {
	Name string `json:"name,omitempty"`

	// CreationTimestamp is a timestamp representing the server time when this object was created.
	CreationTimestamp *Time `json:"creationTimestamp,omitempty"`

	// Status describes the current status of a Namespace
	Status NamespaceStatus `json:"status,omitempty"`
}

type NamespaceGetAllResponseItem struct {
	Name string `json:"name,omitempty"`
	//CreationTimestamp *Time `json:"creationTimestamp,omitempty"`
	//Status NamespaceStatus `json:"status,omitempty"`
}

type NamespaceGetAllResponse struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	Kind string `json:"kind,omitempty"`

	Items []NamespaceGetAllResponseItem `json:"items"`
}
