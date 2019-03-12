// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Interface) DeepCopyInto(out *Interface) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Interface.
func (in *Interface) DeepCopy() *Interface {
	if in == nil {
		return nil
	}
	out := new(Interface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MatchCondition) DeepCopyInto(out *MatchCondition) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MatchCondition.
func (in *MatchCondition) DeepCopy() *MatchCondition {
	if in == nil {
		return nil
	}
	out := new(MatchCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeCfgNetworkState) DeepCopyInto(out *NodeCfgNetworkState) {
	*out = *in
	if in.Interfaces != nil {
		in, out := &in.Interfaces, &out.Interfaces
		*out = make([]Interface, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeCfgNetworkState.
func (in *NodeCfgNetworkState) DeepCopy() *NodeCfgNetworkState {
	if in == nil {
		return nil
	}
	out := new(NodeCfgNetworkState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicy) DeepCopyInto(out *NodeNetworkConfigurationPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicy.
func (in *NodeNetworkConfigurationPolicy) DeepCopy() *NodeNetworkConfigurationPolicy {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkConfigurationPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicyList) DeepCopyInto(out *NodeNetworkConfigurationPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeNetworkConfigurationPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicyList.
func (in *NodeNetworkConfigurationPolicyList) DeepCopy() *NodeNetworkConfigurationPolicyList {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkConfigurationPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicySpec) DeepCopyInto(out *NodeNetworkConfigurationPolicySpec) {
	*out = *in
	in.DesiredState.DeepCopyInto(&out.DesiredState)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicySpec.
func (in *NodeNetworkConfigurationPolicySpec) DeepCopy() *NodeNetworkConfigurationPolicySpec {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkConfigurationPolicyStatus) DeepCopyInto(out *NodeNetworkConfigurationPolicyStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkConfigurationPolicyStatus.
func (in *NodeNetworkConfigurationPolicyStatus) DeepCopy() *NodeNetworkConfigurationPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkConfigurationPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkState) DeepCopyInto(out *NodeNetworkState) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkState.
func (in *NodeNetworkState) DeepCopy() *NodeNetworkState {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkState) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateList) DeepCopyInto(out *NodeNetworkStateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]NodeNetworkState, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateList.
func (in *NodeNetworkStateList) DeepCopy() *NodeNetworkStateList {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NodeNetworkStateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateSpec) DeepCopyInto(out *NodeNetworkStateSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateSpec.
func (in *NodeNetworkStateSpec) DeepCopy() *NodeNetworkStateSpec {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeNetworkStateStatus) DeepCopyInto(out *NodeNetworkStateStatus) {
	*out = *in
	in.DesiredState.DeepCopyInto(&out.DesiredState)
	in.OperationalState.DeepCopyInto(&out.OperationalState)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeNetworkStateStatus.
func (in *NodeNetworkStateStatus) DeepCopy() *NodeNetworkStateStatus {
	if in == nil {
		return nil
	}
	out := new(NodeNetworkStateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeOpNetworkState) DeepCopyInto(out *NodeOpNetworkState) {
	*out = *in
	in.NodeCfgNetworkState.DeepCopyInto(&out.NodeCfgNetworkState)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeOpNetworkState.
func (in *NodeOpNetworkState) DeepCopy() *NodeOpNetworkState {
	if in == nil {
		return nil
	}
	out := new(NodeOpNetworkState)
	in.DeepCopyInto(out)
	return out
}