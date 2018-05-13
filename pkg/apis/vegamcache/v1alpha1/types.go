package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VegamCache describes a VegamCache.
type VegamCache struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec VegamCacheSpec `json:"spec"`
}

// VegamCacheSpec is the spec
type VegamCacheSpec struct {
	DeploymentName  string `json:"deploymentname"`
	VegamPort       uint   `json:"vegamport"`
	VegamServerPort uint   `json:"vegamserverport"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VegamCacheList is a list of VegamCache resources
type VegamCacheList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []VegamCache `json:"items"`
}
