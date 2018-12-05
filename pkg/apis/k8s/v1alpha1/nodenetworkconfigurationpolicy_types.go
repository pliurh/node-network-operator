package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MatchCondition defines the interface match condition
type MatchCondition struct {
	Name string    `json:"name"`
}

// NodeNetworkConfigurationPolicySpec defines the desired state of NodeNetworkConfigurationPolicy
type NodeNetworkConfigurationPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Priority     uint                  `json:"priority"`
	// Match        []MatchCondition      `json:"match"`
	DesiredState NodeCfgNetworkState   `json:"desiredState"`

	// Label selector for Machines.
	NodeSelector *metav1.LabelSelector `json:"nodeSelector,omitempty"`
	//TODO: add AutoConfig support
}

// NodeNetworkConfigurationPolicyStatus defines the observed state of NodeNetworkConfigurationPolicy
type NodeNetworkConfigurationPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeNetworkConfigurationPolicy is the Schema for the nodenetworkconfigurationpolicies API
// +k8s:openapi-gen=true
type NodeNetworkConfigurationPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeNetworkConfigurationPolicySpec   `json:"spec,omitempty"`
	Status NodeNetworkConfigurationPolicyStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeNetworkConfigurationPolicyList contains a list of NodeNetworkConfigurationPolicy
type NodeNetworkConfigurationPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeNetworkConfigurationPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeNetworkConfigurationPolicy{}, &NodeNetworkConfigurationPolicyList{})
}
