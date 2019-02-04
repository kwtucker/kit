package main

const (
	V1                              APIVersion = "v1"
	Secret                          Kind       = "Secret"
	Lrm                             Namespace  = "lrm"
	KubernetesIoServiceAccountToken Type       = "kubernetes.io/service-account-token"
	Opaque                          Type       = "Opaque"
)

type Kind string
type Namespace string
type Type string
type APIVersion string

type ARG struct {
	Verbose bool
	Objects bool
	Data    bool
	Name    string
	Delete  string
	Secret  string
}

type Key struct {
	Key   string
	Value string
}

type Secrets struct {
	APIVersion APIVersion      `json:"apiVersion,omitempty"`
	Items      []Item          `json:"items,omitempty"`
	Kind       string          `json:"kind,omitempty"`
	Metadata   SecretsMetadata `json:"metadata,omitempty"`
}

type Item struct {
	APIVersion APIVersion        `json:"apiVersion,omitempty"`
	Data       map[string]string `json:"data,omitempty"`
	Kind       Kind              `json:"kind,omitempty"`
	Metadata   Metadata          `json:"metadata,omitempty"`
	Type       Type              `json:"type,omitempty"`
}

type Metadata struct {
	CreationTimestamp string       `json:"creationTimestamp,omitempty"`
	Name              string       `json:"name,omitempty"`
	Namespace         Namespace    `json:"namespace,omitempty"`
	ResourceVersion   string       `json:"resourceVersion,omitempty"`
	SelfLink          string       `json:"selfLink,omitempty"`
	Uid               string       `json:"uid,omitempty"`
	Annotations       *Annotations `json:"annotations,omitempty,omitempty"`
}

type Annotations struct {
	KubernetesIoServiceAccountName string `json:"kubernetes.io/service-account.name,omitempty"`
	KubernetesIoServiceAccountUid  string `json:"kubernetes.io/service-account.uid,omitempty"`
}

type SecretsMetadata struct {
	ResourceVersion string `json:"resourceVersion,omitempty"`
	SelfLink        string `json:"selfLink,omitempty"`
}
