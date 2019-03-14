FROM registry.svc.ci.openshift.org/openshift/release:golang-1.10 AS builder
WORKDIR /go/src/github.com/pliurh/node-network-operator
COPY . .
RUN make

FROM registry.svc.ci.openshift.org/openshift/origin-v4.0:base
ADD build/_output/bin/node-network-operator /usr/bin/node-network-operator
WORKDIR /usr/bin
ENTRYPOINT ["/usr/bin/node-network-operator"]
LABEL io.k8s.display-name="OpenShift node-network-operator" \
      io.k8s.description="This is a component of OpenShift Container Platform that manages the SR-IOV NICs on hosts." \
      io.openshift.tags="openshift,sriov" \
      maintainer="Peng Liu <pliu@redhat.com>"
