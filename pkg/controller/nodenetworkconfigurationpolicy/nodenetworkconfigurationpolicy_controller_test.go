package nodenetworkconfigurationpolicy

import (
	"strings"
	"testing"

	nodenetwork "github.com/pliurh/node-network-operator/pkg/apis/nodenetwork/v1alpha1"
	"github.com/vincent-petithory/dataurl"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRenderMachineConfigWithVfNum(t *testing.T) {
	numVfs := uint(2)
	policy := &nodenetwork.NodeNetworkConfigurationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-policy",
			Namespace: "node-network-operator",
		},
		Spec: nodenetwork.NodeNetworkConfigurationPolicySpec{
			Priority: 99,
			DesiredState: nodenetwork.NodeCfgNetworkState{
				Interfaces: []nodenetwork.Interface{
					{
						Name:   "eth2",
						NumVfs: &numVfs,
					},
				},
			},
		},
	}

	machineConfig, _ := renderMachineConfig(policy)

	if machineConfig.Spec.Config.Storage.Files[0].Path != "/etc/udev/rules.d/99-sriov.rules" {
		t.Error("Exp. the udev rule for sriov to be created.")
	}

	dataURL, err := dataurl.DecodeString(machineConfig.Spec.Config.Storage.Files[0].Contents.Source)
	if err != nil {
		t.Error("Exp. the file content to be encoded in data-url format")
	}

	content := strings.TrimRight(string(dataURL.Data), "\n")
	rules := strings.Split(content, "\n")

	if len(rules) != 1 {
		t.Error("Exp. one rule to be created.")
	}

	if strings.Compare(rules[0], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth2", ATTR{device/sriov_numvfs}="2"`) != 0 {
		t.Error("Exp. correct rule to be created")
	}
}

func TestRenderMachineConfigWithVfNumMutiIf(t *testing.T) {
	numVfs := uint(2)
	policy := &nodenetwork.NodeNetworkConfigurationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-policy",
			Namespace: "node-network-operator",
		},
		Spec: nodenetwork.NodeNetworkConfigurationPolicySpec{
			Priority: 99,
			DesiredState: nodenetwork.NodeCfgNetworkState{
				Interfaces: []nodenetwork.Interface{
					{
						Name:   "eth2",
						NumVfs: &numVfs,
					},
					{
						Name:   "eth1",
						NumVfs: &numVfs,
					},
				},
			},
		},
	}

	machineConfig, _ := renderMachineConfig(policy)

	if machineConfig.Spec.Config.Storage.Files[0].Path != "/etc/udev/rules.d/99-sriov.rules" {
		t.Error("Exp. the udev rule for sriov to be created.")
	}

	dataURL, err := dataurl.DecodeString(machineConfig.Spec.Config.Storage.Files[0].Contents.Source)
	if err != nil {
		t.Error("Exp. the file content to be encoded in data-url format")
	}

	content := strings.TrimRight(string(dataURL.Data), "\n")
	rules := strings.Split(content, "\n")

	if len(rules) != 2 {
		t.Error("Exp. two rules to be created.")
	}

	if strings.Compare(rules[0], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth2", ATTR{device/sriov_numvfs}="2"`) != 0 {
		t.Error("Exp. correct rule to be created")
	}

	if strings.Compare(rules[1], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth1", ATTR{device/sriov_numvfs}="2"`) != 0 {
		t.Error("Exp. correct rule to be created")
	}
}

func TestRenderMachineConfigWithMtu(t *testing.T) {
	mtu := uint(1200)
	policy := &nodenetwork.NodeNetworkConfigurationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-policy",
			Namespace: "node-network-operator",
		},
		Spec: nodenetwork.NodeNetworkConfigurationPolicySpec{
			Priority: 99,
			DesiredState: nodenetwork.NodeCfgNetworkState{
				Interfaces: []nodenetwork.Interface{
					{
						Name: "eth2",
						Mtu:  &mtu,
					},
				},
			},
		},
	}

	machineConfig, _ := renderMachineConfig(policy)

	if machineConfig.Spec.Config.Storage.Files[0].Path != "/etc/udev/rules.d/99-mtu.rules" {
		t.Error("Exp. the udev rule for mtu to be created.")
	}

	dataURL, err := dataurl.DecodeString(machineConfig.Spec.Config.Storage.Files[0].Contents.Source)
	if err != nil {
		t.Error("Exp. the file content to be encoded in data-url format")
	}

	content := strings.TrimRight(string(dataURL.Data), "\n")
	rules := strings.Split(content, "\n")

	if len(rules) != 1 {
		t.Error("Exp. one rule to be created.")
	}

	if strings.Compare(rules[0], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth2", RUN+="/sbin/ip link set mtu 1200 dev '%k'"`) != 0 {
		t.Error("Exp. correct rule to be created")
	}
}

func TestRenderMachineConfigWithMtuAndVfNum(t *testing.T) {
	numVfs := uint(2)
	mtu := uint(1200)
	policy := &nodenetwork.NodeNetworkConfigurationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-policy",
			Namespace: "node-network-operator",
		},
		Spec: nodenetwork.NodeNetworkConfigurationPolicySpec{
			Priority: 99,
			DesiredState: nodenetwork.NodeCfgNetworkState{
				Interfaces: []nodenetwork.Interface{
					{
						Name:   "eth2",
						NumVfs: &numVfs,
						Mtu:    &mtu,
					},
				},
			},
		},
	}

	machineConfig, _ := renderMachineConfig(policy)

	for i := range machineConfig.Spec.Config.Storage.Files {
		dataURL, err := dataurl.DecodeString(machineConfig.Spec.Config.Storage.Files[i].Contents.Source)
		if err != nil {
			t.Error("Exp. the file content to be encoded in data-url format")
		}
		content := strings.TrimRight(string(dataURL.Data), "\n")
		rules := strings.Split(content, "\n")
		if len(rules) != 1 {
			t.Error("Exp. one rule to be created.")
		}

		if machineConfig.Spec.Config.Storage.Files[i].Path == "/etc/udev/rules.d/99-mtu.rules" {
			if strings.Compare(rules[0], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth2", RUN+="/sbin/ip link set mtu 1200 dev '%k'"`) != 0 {
				t.Error("Exp. correct rule to be created")
			}
		} else if machineConfig.Spec.Config.Storage.Files[i].Path == "/etc/udev/rules.d/99-sriov.rules" {
			if strings.Compare(rules[0], `ACTION=="add", SUBSYSTEM=="net", KERNEL=="eth2", ATTR{device/sriov_numvfs}="2"`) != 0 {
				t.Error("Exp. correct rule to be created")
			}
		} else {
			t.Error("Exp. the udev rule for mtu and sriov vf to be created.")
		}
	}
}
