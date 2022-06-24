package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ci is the type of a ci.
type CiType string

const (
	// CiTypeUnknown is a super ci unknown.
	CiTypeUnknown = CiType("unknown")
	// CiTypeSuperCi is a super ci e.g Batman, Spiderman...
	CiTypeSuperCi = CiType("superci")
	// CiTypeAntiCi is a anti ci e.g Punisher, Deadpool...
	CiTypeAntiCi = CiType("antici")
	// CiTypeVillain is a Villain e.g Fisk, Joker...
	CiTypeVillain = CiType("villain")
)

// Ci represents a apps ci.
//
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Ci struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CiSpec   `json:"spec,omitempty"`
	Status CiStatus `json:"status,omitempty"`
}

// CiSpec is the spec of a Ci.
type CiSpec struct{
	Model     string            `json:"model,omitempty"`
	Repo      string            `json:"repo,omitempty"`
	Branch    string            `json:"branch,omitempty"`
	On        On                `json:"on,omitempty"`

	// +listType=map
	// +optional
	Variables map[string]string `json:"variables,omitempty"`
}

type CiStatus struct {
	Histroy []string `json:"history,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// CiList is a list of Ci resources.
type CiList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Ci `json:"items"`
}

type On struct {
	Schedule string   `json:"schedule"`
	Events   []string `json:"events"`
}
