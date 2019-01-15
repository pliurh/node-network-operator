package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Interface defines the state of network interface
type Interface struct {
	Name    string    `json:"name"`
	NumVfs  *uint      `json:"numVfs,omitempty"`
	Mtu     *uint      `json:"mtu,omitempty"`
	Promisc *bool      `json:"promisc,omitempty"`
}

// NodeCfgNetworkState defines the configuration state of node network
type NodeCfgNetworkState struct {
	Interfaces []Interface       `json:"interfaces"`
}

// NodeOpNetworkState defines the operational state of node network configuration, which is a superset of configuration state
type NodeOpNetworkState struct {
	NodeCfgNetworkState
	// TODO: Capabilities string     `json:"capabilities"`
}

// NodeNetworkStateSpec defines the desired state of NodeNetworkState
type NodeNetworkStateSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	Managed bool    `json:"managed"`
}

// NodeNetworkStateStatus defines the observed state of NodeNetworkState
type NodeNetworkStateStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
//	NodeName         string               `json:"nodeName"`
	DesiredState     NodeCfgNetworkState  `json:"desiredState"`
	OperationalState NodeOpNetworkState   `json:"operationalState"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeNetworkState is the Schema for the nodenetworkstates API
// +k8s:openapi-gen=true
type NodeNetworkState struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodeNetworkStateSpec   `json:"spec,omitempty"`
	Status NodeNetworkStateStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeNetworkStateList contains a list of NodeNetworkState
type NodeNetworkStateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeNetworkState `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeNetworkState{}, &NodeNetworkStateList{})
}
