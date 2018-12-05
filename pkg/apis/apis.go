package apis

import (
	"k8s.io/apimachinery/pkg/runtime"
	mcfgv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
)

// AddToSchemes may be used to add all resources defined in the project to a Scheme
var AddToSchemes runtime.SchemeBuilder

// AddToScheme adds all Resources to the Scheme
func AddToScheme(s *runtime.Scheme) error {
	mcfgv1.AddToScheme(s)
	return AddToSchemes.AddToScheme(s)
}
