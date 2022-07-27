package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ci is the type of a ci.
type CiType string

// Ci represents a apps ci.
//
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Ci struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiSpec   `json:"spec,omitempty"`
	Status CiStatus `json:"status,omitempty"`
}

// CiSpec is the spec of a Ci.
type CiSpec struct {
	Model  string `json:"model,omitempty"`
	Repo   string `json:"repo,omitempty"`
	Branch string `json:"branch,omitempty"`
	Term   Term   `json:"term,omitempty"`

	// +listType=map
	// +optional
	Variables map[string]string `json:"variables,omitempty"`
}

type CiStatus struct {
	Histroy []Histroy `json:"history,omitempty"`
}

type Histroy struct {
	CiName  string `json:"ciName"`
	PodName string `json:"podName"`
	Time    string `json:"time"`
	Status  string `json:"status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CiList is a list of Ci resources.
type CiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Ci `json:"items"`
}

type Term struct {
	Schedule string   `json:"schedule"`
	Events   []string `json:"events"`
}
