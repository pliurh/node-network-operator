package nodenetworkconfigurationpolicy

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"regexp"

	ignv2_2types "github.com/coreos/ignition/config/v2_2/types"
	k8sv1alpha1 "github.com/pliurh/node-network-operator/pkg/apis/k8s/v1alpha1"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"

	"github.com/vincent-petithory/dataurl"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_nodenetworkconfigurationpolicy")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NodeNetworkConfigurationPolicy Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNodeNetworkConfigurationPolicy{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nodenetworkconfigurationpolicy-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NodeNetworkConfigurationPolicy
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.NodeNetworkConfigurationPolicy{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NodeNetworkConfigurationPolicy
	err = c.Watch(&source.Kind{Type: &k8sv1alpha1.NodeNetworkConfigurationPolicy{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &k8sv1alpha1.NodeNetworkConfigurationPolicy{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileNodeNetworkConfigurationPolicy{}

// ReconcileNodeNetworkConfigurationPolicy reconciles a NodeNetworkConfigurationPolicy object
type ReconcileNodeNetworkConfigurationPolicy struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NodeNetworkConfigurationPolicy object and makes changes based on the state read
// and what is in the NodeNetworkConfigurationPolicy.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNodeNetworkConfigurationPolicy) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NodeNetworkConfigurationPolicy")

	// Fetch the NodeNetworkConfigurationPolicy instance
	instance := &k8sv1alpha1.NodeNetworkConfigurationPolicy{}
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

	k8sv1alpha1.ValidateNodeNetworkConfigurationPolicy(instance)

	b := new(bytes.Buffer)
    for key, value := range instance.Labels {
        fmt.Fprintf(b, "%s=%s", key, value)
    }
	reqLogger.Info("Instance","Labels", b.String())

	// Fetch the NodeNetworkConfigurationPolicy instances with the same label.
	policies := &k8sv1alpha1.NodeNetworkConfigurationPolicyList{}
	listOpts := &client.ListOptions{}
	listOpts.SetLabelSelector(b.String())
	err = r.client.List(context.TODO(), listOpts, policies)
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
	policy := k8sv1alpha1.MergeNodeNetworkConfigurationPolicies(policies)

	// Render MachineConfig based on policies
	machineConfig, err:= renderMachineConfig(policy)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Check if the MachineConfig already exists.
	found := &mcfgv1.MachineConfig{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: machineConfig.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new machineConfig", "name", machineConfig.Name)
		// Create MachineConfig if not exist
		err = r.client.Create(context.TODO(), machineConfig)
		if err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	} else if err != nil {
		reqLogger.Info("Error when finding machineConfig ")
		return reconcile.Result{}, err
	}

	// MachineConfig already exists, update it if the hash of spec is different
	foundHash, err := getMachineConfigHash(&found.Spec)
	if err != nil {
		return reconcile.Result{}, err
	}
	hash, err := getMachineConfigHash(&machineConfig.Spec)
	if err != nil {
		return reconcile.Result{}, err
	}
	if foundHash == hash {
		//No need to update
		return reconcile.Result{}, nil
	}

	reqLogger.Info("Update MachineConfig", "name", found.Name)
	found.Spec.Config = machineConfig.Spec.Config
	err = r.client.Update(context.TODO(), found)
	if err != nil {
		return reconcile.Result{}, err
	}

	err = r.updateNodeNetworkState(policy)
	return reconcile.Result{}, nil
}

func newNodeNetworkState(node *corev1.Node)*k8sv1alpha1.NodeNetworkState{
	return &k8sv1alpha1.NodeNetworkState{
		ObjectMeta: metav1.ObjectMeta{
			Name: node.Name,
		},
		Spec: k8sv1alpha1.NodeNetworkStateSpec{
			Managed: true,
		},
	}
}

func (r *ReconcileNodeNetworkConfigurationPolicy)updateNodeNetworkState(cr *k8sv1alpha1.NodeNetworkConfigurationPolicy) error{
	listOpts := &client.ListOptions{}
	nodes := &corev1.NodeList{}
	err := r.client.List(context.TODO(), listOpts, nodes)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return nil
		}
		// Error reading the object - requeue the request.
		return err
	}

	for _, node := range nodes.Items {
		cfg := &k8sv1alpha1.NodeNetworkState{}
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: node.Name, Namespace: ""}, cfg)
		if err != nil {
			if errors.IsNotFound(err) {
				// Request object not found, create it
				log.Info("NodeNetworkState is not found, create it", "node", node.Name)
				cfg := newNodeNetworkState(&node)
				err = r.client.Create(context.TODO(), cfg)
				if err != nil {
					return err
				}
			} else {
				// Error reading the object - requeue the request.
				return err
			}
		}
		
		// update node network desired config
		err = r.client.Get(context.TODO(), types.NamespacedName{Name: node.Name, Namespace: ""}, cfg)
		if err != nil {
			return err
		}
		log.Info("Update node network config", "desired state", cr.Spec.DesiredState)
		cfg.Status.DesiredState = *cr.Spec.DesiredState.DeepCopy()
		err = r.client.Status().Update(context.TODO(), cfg)
		if err != nil {
			// Error reading the object - requeue the request.
			return err
		}
	}
	return nil
}

// Render MachineConfig based on policies
func renderMachineConfig(cr *k8sv1alpha1.NodeNetworkConfigurationPolicy) (*mcfgv1.MachineConfig, error){	
	config, err := generateIgnConfig(cr)
	if err != nil {
		return nil,err
	}

	return &mcfgv1.MachineConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name: cr.Name,
			Labels: cr.ObjectMeta.Labels,
		},
		Spec: mcfgv1.MachineConfigSpec{
			Config: *config,
		},
	}, nil
}

func generateIgnConfig(cr *k8sv1alpha1.NodeNetworkConfigurationPolicy) (*ignv2_2types.Config, error) {
	contents := make(map[string]string)
	for _, iface := range cr.Spec.DesiredState.Interfaces {
		parseInterface(iface, contents)
	}

	config := ignv2_2types.Config{
		Storage: ignv2_2types.Storage{
			Files: generateFiles(contents),
		},
	}

	return &config, nil
}

func parseInterface(i k8sv1alpha1.Interface, contents map[string]string) {
	if i.Mtu != nil && *i.Mtu != 0 {
		contents["mtu"] += fmt.Sprintf("ACTION==\"add\", SUBSYSTEM==\"net\", KERNEL==\"%s\", RUN+=\"/sbin/ip link set mtu %d dev '%%k'\"\n", i.Name, *i.Mtu)
	}
	if i.NumVfs != nil && *i.NumVfs >= 0 {
		contents["sriov"] += (fmt.Sprintf("[device-%s]\n", i.Name) + fmt.Sprintf("match-device=interface-name:%v\nsriov-num-vfs=%v\n\n", i.Name, *i.NumVfs))
	}
	if i.Promisc != nil {
		filename := "ifcfg-" + i.Name
		contents[filename] += fmt.Sprintf("DEVICE=%s\n", i.Name)
		contents[filename] += fmt.Sprintf("ONBOOT=%s\n", "yes")
		contents[filename] += "NM_CONTROLLED=yes\n"
		if *i.Promisc {
			contents[filename] += fmt.Sprintf("PROMISC=%s\n", "yes")
		} else {
			contents[filename] += fmt.Sprintf("PROMISC=%s\n", "no")
		}
	}
}

func getEncodedContent(inp string) string {
	return (&url.URL{
		Scheme: "data",
		Opaque: "," + dataurl.Escape([]byte(inp)),
	}).String()
}

func generateFiles(contents map[string]string) []ignv2_2types.File {
	var files []ignv2_2types.File
	fileMode := int(420)
	r := regexp.MustCompile(`^ifcfg-.*`)

	for k, v := range contents {
		switch k {
		case "mtu":
			log.Info("file content", "mtu.conf", v)
			files = append (files, ignv2_2types.File{
				Node: ignv2_2types.Node{
					Path: "/etc/udev/rules.d/99-mtu.rules",
				},
				FileEmbedded1: ignv2_2types.FileEmbedded1{
					Contents: ignv2_2types.FileContents{
						Source: getEncodedContent(v),
					},
					Mode: &fileMode,
				},
			})
		case "sriov":
			log.Info("file content", "sriov.conf", v)
			files = append (files, ignv2_2types.File{
				Node: ignv2_2types.Node{
					Path: "/etc/NetworkManager/conf.d/sriov.conf",
				},
				FileEmbedded1: ignv2_2types.FileEmbedded1{
					Contents: ignv2_2types.FileContents{
						Source: getEncodedContent(v),
					},
					Mode: &fileMode,
				},
			})
		default:
			if r.MatchString(k) {
				log.Info("file content", k, v)
				files = append (files, ignv2_2types.File{
					Node: ignv2_2types.Node{
						Path: fmt.Sprintf("/etc/sysconfig/network-scripts/%s", k),
					},
					FileEmbedded1: ignv2_2types.FileEmbedded1{
						Contents: ignv2_2types.FileContents{
							Source: getEncodedContent(v),
						},
						Mode: &fileMode,
					},
				})
			}
		}
	}

	return files
}
