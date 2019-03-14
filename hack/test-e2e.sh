#!/bin/bash -x
set -e

if [ -n "${IMAGE_CLUSTER_LOGGING_OPERATOR:-}" ] ; then
  source "$(dirname $0)/common"
fi

# STABLE_IMAGE_CLUSTER_LOGGING_OPERATOR=$(echo $IMAGE_FORMAT | sed 's,${component},cluster-logging-operator,')
# IMAGE_CLUSTER_LOGGING_OPERATOR=${IMAGE_CLUSTER_LOGGING_OPERATOR:-quay.io/openshift/origin-cluster-logging-operator:latest}
# IMAGE_MANIFEST_CLUSTER_LOGGING_OPERATOR=${STABLE_IMAGE_CLUSTER_LOGGING_OPERATOR:-$IMAGE_CLUSTER_LOGGING_OPERATOR}

repo_dir="$(dirname $0)/.."
if ! oc get namespace node-network-operator > /dev/null 2>&1 ; then
    oc create -f deploy/namespace.yaml
fi

manifest=$(mktemp)
files="service_account.yaml role.yaml role_binding.yaml operator.yaml"
pushd deploy;
  for f in ${files}; do
     cat ${f} >> ${manifest};
  done;
popd

global_manifest=$(mktemp)
global_files="nodenetwork_v1alpha1_nodenetworkconfigurationpolicy_crd.yaml nodenetwork_v1alpha1_nodenetworkstate_crd.yaml"

pushd deploy/crds;
  for f in ${global_files}; do
    cat ${f} >> ${global_manifest}
  done;
popd

TEST_NAMESPACE=${NAMESPACE} go test ./test/e2e/... \
  -root=$(pwd) \
  -kubeconfig=${KUBECONFIG} \
  -globalMan ${global_manifest} \
  -namespacedMan ${manifest} \
  -v \
  -parallel=1 \
  -singleNamespace
