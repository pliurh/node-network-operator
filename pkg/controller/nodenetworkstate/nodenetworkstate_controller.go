package nodenetworkstate

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	ignv2_2types "github.com/coreos/ignition/config/v2_2/types"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	k8sv1alpha1 "github.com/pliurh/node-network-operator/pkg/apis/k8s/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"github.com/vincent-petithory/dataurl"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
 	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nodenetworkstate")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NodeNetworkState Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNodeNetworkState{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nodenetworkstate-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NodeNetworkState
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.NodeNetworkState{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NodeNetworkState
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NodeNetworkState{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNodeNetworkState{}

// ReconcileNodeNetworkState reconciles a NodeNetworkState object
type ReconcileNodeNetworkState struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NodeNetworkState object and makes changes based on the state read
// and what is in the NodeNetworkState.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNodeNetworkState) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NodeNetworkState")

	// Fetch the NodeNetworkState instance
	instance := &k8sv1alpha1.NodeNetworkState{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new MachineConfig object if node is managed.
	if !instance.Spec.Managed {
		reqLogger.Info("Node", instance.Name, "is not configured to be managed by operator")
		return reconcile.Result{}, nil
	}

	machineConfig := newMachineConfig(instance)

	// Set NodeNetworkState instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, machineConfig, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the MachineConfig already exists. 
	found := &mcfgv1.MachineConfig{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: machineConfig.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new machineConfig", "name", machineConfig.Name)
		err = r.client.Create(context.TODO(), machineConfig)
		if err != nil {
			return reconcile.Result{}, err
		}
		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Info("Error when finding machineConfig ")
		return reconcile.Result{}, err
	}

	// MachineConfig already exists, update it
	reqLogger.Info("Update MachineConfig", "name", found.Name)
	found.Spec.Config = machineConfig.Spec.Config
	err = r.client.Update(context.TODO(), found)
	if err != nil {
		return reconcile.Result{}, err
	}	
	return reconcile.Result{}, nil
}

func newMachineConfig(nns *k8sv1alpha1.NodeNetworkState) *mcfgv1.MachineConfig{
	labels := map[string]string{
		"machineconfiguration.openshift.io/role": "worker",
	}
	config, err := generateIgnConfig(nns)
	if err != nil {
		return nil
	} 

	return &mcfgv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nodenetconf", 
			Labels: labels,
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: *config,
		},
	}
}

func generateIgnConfig(nns *k8sv1alpha1.NodeNetworkState) (*ignv2_2types.Config, error) {

	fileMode := int(420)
	contents, err := generateNetConfigFileContent(nns)
	if err != nil {
		return nil, err
	}

	file := ignv2_2types.File{
		Node: ignv2_2types.Node{
			Path: "/etc/NetworkManager/conf.d/sriov.conf",
		},
		FileEmbedded1: ignv2_2types.FileEmbedded1{
			Contents: ignv2_2types.FileContents{
				Source: getEncodedContent(contents),
			},
			Mode: &fileMode,
		},
	}
	config := ignv2_2types.Config{
		Storage: ignv2_2types.Storage{
			Files: []ignv2_2types.File{},
		},
	}
	config.Storage.Files = append(config.Storage.Files, file)
	return &config, nil
}

func getEncodedContent(inp string) string {
	return (&url.URL{
		Scheme: "data",
		Opaque: "," + dataurl.Escape([]byte(inp)),
	}).String()
}

func generateNetConfigFileContent(nns *k8sv1alpha1.NodeNetworkState) (string, error) {
	var content strings.Builder
	var i k8sv1alpha1.Interface

	fmt.Fprintf(&content, "[devices]\n")

	for _, i = range nns.Status.DesiredState.Interfaces {
		if i.NumVfs > 0 {
			fmt.Fprintf(&content, "match-device=interface-name:%v\nsriov-num-vfs=%v\n", i.Name, i.NumVfs)
		}
	}
	// log.Printf("NetConfig file content is %s", content.String())
	
	return content.String(), nil
}

