#!/bin/bash

repo_dir="$(dirname $0)/.."

pushd ${repo_dir}/deploy
files="service_account.yaml role.yaml role_binding.yaml operator.yaml crds/nodenetwork_v1alpha1_nodenetworkconfigurationpolicy_crd.yaml crds/nodenetwork_v1alpha1_nodenetworkstate_crd.yaml"
for file in ${files}; do
  oc delete -f $file --ignore-not-found
done
popd