package e2e

import (
	goctx "context"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	apis "github.com/pliurh/node-network-operator/pkg/apis"
	nodenetwork "github.com/pliurh/node-network-operator/pkg/apis/nodenetwork/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	namespace            = "node-network-operator"
	retryInterval        = time.Second * 5
	timeout              = time.Second * 120
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)

func TestNodeNetworkConfigurationPolicy(t *testing.T) {
	nodeNetworkConfigurationPolicyList := &nodenetwork.NodeNetworkConfigurationPolicyList{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
	}
	err := framework.AddToFrameworkScheme(apis.AddToScheme, nodeNetworkConfigurationPolicyList)
	if err != nil {
		t.Fatalf("failed to add custom resource scheme to framework: %v", err)
	}

	// run subtests
	t.Run("NodeNetworkConfigurationPolicy-group", func(t *testing.T) {

		t.Run("Node network with SR-IOV", NodeNetworkConfigurationPolicyVfnum)
		// time.Sleep(time.Minute * 1) // wait for objects to be deleted/cleaned up
	})
}

func NodeNetworkConfigurationPolicyVfnum(t *testing.T) {
	t.Parallel()
	ctx := framework.NewTestCtx(t)
	defer ctx.Cleanup()
	numVfs := uint(2)

	// create NodeNetworkConfigurationPolicy custom resource
	policy := &nodenetwork.NodeNetworkConfigurationPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "NodeNetworkConfigurationPolicy",
			APIVersion: nodenetwork.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "example-policy",
			Namespace: namespace,
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

	// get global framework variables
	f := framework.Global
	err := f.Client.Create(goctx.TODO(), policy, &framework.CleanupOptions{TestContext: ctx, Timeout: time.Second * 5, RetryInterval: time.Second * 1})
	if err != nil {
		return
	}

	err = WaitForMachineConfig(t, f.Client, "/etc/udev/rules.d/99-sriov.rules", retryInterval, timeout)
	if err != nil {
		return
	}
}
