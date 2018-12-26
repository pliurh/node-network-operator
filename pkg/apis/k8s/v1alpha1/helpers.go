package v1alpha1

import (
	"sort"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	return &NodeNetworkConfigurationPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodenetconf",
			Labels: crs.Items[0].ObjectMeta.Labels,
		},
		Spec: NodeNetworkConfigurationPolicySpec{
			DesiredState: NodeCfgNetworkState{
				Interfaces: interfaces,
			},
		},
	}
}

func contains(ifaces []Interface, iface *Interface) bool{
	for _, i := range ifaces {
		if i.Name == iface.Name {
			return true
		}
	}
	return false
}