package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Model is the type of a model.
type ModelType string

const (
	// ModelTypeUnknown is a super model unknown.
	ModelTypeUnknown = ModelType("unknown")
	// ModelTypeSuperModel is a super model e.g Batman, Spiderman...
	ModelTypeSuperModel = ModelType("supermodel")
	// ModelTypeAntiModel is a anti model e.g Punisher, Deadpool...
	ModelTypeAntiModel = ModelType("antimodel")
	// ModelTypeVillain is a Villain e.g Fisk, Joker...
	ModelTypeVillain = ModelType("villain")
)

// Model represents aapps model.
//
// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec,omitempty"`
	Status ModelStatus `json:"status,omitempty"`
}

// ModelSpec is the spec of a Model.
type ModelSpec struct{
	Tasks     []Task  `json:"tasks,omitempty"`

	// +listType=map
	// +optional
	Variables map[string]string `json:"variables,omitempty"`
}

type ModelStatus struct {
	StartTime  string
	Completime string
	Succeeded  bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ModelList is a list of Model resources.
type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Model `json:"items"`
}

type Task struct {
	Name    string   `json:"name,omitempty"`
	Image   string   `json:"image,omitempty"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

