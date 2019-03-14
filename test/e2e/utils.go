package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	framework "github.com/operator-framework/operator-sdk/pkg/test"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	// "k8s.io/client-go/kubernetes"
	// dynclient "sigs.k8s.io/controller-runtime/pkg/client"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	// nnfgv1 "github.com/pliurh/node-network-operator/pkg/apis/nodenetwork/v1alpha1"
	// policy "github.com/pliurh/node-network-operator/pkg/controller/nodenetworkconfigurationpolicy"
)

var (
	mcName = "99-worker-nodenetconf"
)

// WaitForMachineConfig for customer resource to be created
func WaitForMachineConfig(t *testing.T, client framework.FrameworkClient, path string, retryInterval, timeout time.Duration) error {
	found := &mcfgv1.MachineConfig{}
	err := wait.Poll(retryInterval, timeout, func() (done bool, err error) {
		err = client.Get(context.Background(), types.NamespacedName{Name: mcName, Namespace: ""}, found)
		if err != nil {
			if apierrors.IsNotFound(err) {
				t.Logf("Waiting for availability of %s MachineConfig\n", mcName)
				return false, nil
			}
			return false, err
		}

		if found.Spec.Config.Storage.Files[0].Path == path {
			return true, nil
		}

		return false, fmt.Errorf("MachineConfig %s content is incorrect", mcName)
	})

	if err != nil {
		return err
	}
	t.Logf("MachineConfig available.\n")
	return nil
}
