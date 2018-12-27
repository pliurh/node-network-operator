package v1alpha1

import (
	"fmt"
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var log = logf.Log.WithName("v1alpha1_helper")
const (
	label = "machineconfiguration.openshift.io/role"
)

func MergeNodeNetworkConfigurationPolicies(crs *NodeNetworkConfigurationPolicyList) *NodeNetworkConfigurationPolicy {
	interfaces := []Interface{}
	sort.Slice(crs.Items, func(i, j int) bool { return crs.Items[i].Spec.Priority > crs.Items[j].Spec.Priority })

	for _, cr := range crs.Items {
		for _, i := range cr.Spec.DesiredState.Interfaces {
			if !contains(interfaces, &i) {
				interfaces = append(interfaces, *i.DeepCopy())
			}
		}
	}

	name := "nodenetconf-" + crs.Items[0].ObjectMeta.Labels[label]
	return &NodeNetworkConfigurationPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: crs.Items[0].ObjectMeta.Labels,
		},
		Spec: NodeNetworkConfigurationPolicySpec{
			DesiredState: NodeCfgNetworkState{
				Interfaces: interfaces,
			},
		},
	}
}

func ValidateNodeNetworkConfigurationPolicy(cr *NodeNetworkConfigurationPolicy) error {
	if len(cr.Labels) > 1 || len(cr.Labels) == 0 {
		return fmt.Errorf("One label should be specified")
	}

	if val, ok := cr.Labels[label]; !ok{
		if val == "master" || val == "worker" {
			return nil
		}
	}

	return fmt.Errorf("Label can only be either \"master\" or \"worker\"")
}

func contains(ifaces []Interface, iface *Interface) bool{
	for _, i := range ifaces {
		if i.Name == iface.Name {
			return true
		}
	}
	return false
}