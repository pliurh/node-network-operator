# node-network-operator

Node network operator is designed to work in a Openshift 4.0+ cluster. It relies on the machine-config-operator to configure the number of VFs and MTU for worker nodes.

## Quick Start

Install node-network-operator on a Openshift cluster:

```
 $ make deploy-setup
```

Node network config operator defines a new CRD NodeNetworkConfigurationPolicy. With creating a customer resource of it, you are able to configure the number of VFs on each worker node.

Here comes an example:

```
cat <<EOF | oc create -f -
|apiVersion: k8s.cni.cncf.io/v1alpha1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: policy
  labels:
    machineconfiguration.openshift.io/role: worker
spec:
  priority: 99
  desiredState:
    interfaces:
    - name: eth0
      numVfs: 2
    - name: eth1
      mtu: 1400
```

In this example, we set NIC 'eth0' with 2 VFs, and NIC 'eth1' with MTU 1400. Node network config operator will create a MachineConfig customer resource, which will trigger the machine config operator to apply the configuration files of SR-IOV to each work node, then reboot the node. After rebooting, you should be able to see the VFs are provisioned.

## Hacking

### Running the operator

Running locally outside an OKD cluster:

```
 $ make run
```


## Testing

### E2E Testing

Run e2e test locally

```
$ make test-e2e
```

## Unit Testing

Run unit test locally

```
$ make test-unit
```
