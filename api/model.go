package api

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type Model struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelSpec   `json:"spec"`
	Status ModelStatus `json:"status,omitempty"`
}

type ModelSpec struct {
	Tasks     []Task            `json:"tasks,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

type ModelStatus struct {
	StartTime  string
	Completime string
	Succeeded  bool
}

type ModelList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Model `json:"items"`
}

type Task struct {
	Name    string   `json:"name,omitempty"`
	Image   string   `json:"image,omitempty"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}
