package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

type TororuResourceSpec struct {
	Kind   string `json:"kind"`
	Rotate int    `json:"rotate"`
	Config string `json:"config"`
}

//go:generate controller-gen object paths=$GOFILE

// +k8s:deepcopy-gen=true
type TororuResourceConsumers struct {
	RW string   `json:"rw"`
	RO []string `json:"ro"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TororuResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec                TororuResourceSpec      `json:"spec"`
	LastUpdated         string                  `json:"lastUpdated"`
	Approved            bool                    `json:"approved"`
	Consumers           TororuResourceConsumers `json:"consumers"`
	Secret              string                  `json:"secret"`
	PodsRestartRequired bool                    `json:"podsRestartRequired"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type TororuResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []TororuResource `json:"items"`
}

func (t *TororuResource) ToUnstructured() (*unstructured.Unstructured, error) {
	unstructuredObject, err := runtime.DefaultUnstructuredConverter.ToUnstructured(t)
	if err != nil {
		fmt.Printf("Error converting Go struct to unstructured object: %v\n", err)
	}

	return &unstructured.Unstructured{Object: unstructuredObject}, err
}

func FromUnstructured(obj *unstructured.Unstructured) (*TororuResource, error) {
	var tRes TororuResource
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.Object, &tRes)
	if err != nil {
		fmt.Printf("Error converting unstructured object to Go struct: %v\n", err)
	}

	return &tRes, err
}
