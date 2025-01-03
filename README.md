# kzed
Kubernetes operator for integration with IBM Mainframe through ZOWE

## Description
This Kubernetes operator enables zOS management through K8s CRDs in kubernetes-native way.

At this moment the following CRDs are available:

**JCLJob** - enables zOS job management like starting the jobs from JCLs in zOS DataSets or inline JCL definition. JCLJob objects also enables the monitoring of Job execution and insight into job spool outputs.

**PartitionedDataSet** - enables zOS PDS creation and upload of its members. The Members can be defined in the resource itself much like a data in ConfigMap or Secret can be defined.

**SequentialDataSet** - enables zOS sequential data set creation and upload of its data. Data can be defined in the resource itself much like a data in ConfigMap or Secret can be defined.

## Getting Started

### Prerequisites
- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/kzed:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/kzed:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create zowe-config secret**

There are examples of Zowe configuration in ```zowe-config``` folder. Edit ```zove.config.json``` and update the values according your environment. Create secret with the following command: 

```oc -n kzed-system create secret generic zowe-config --from-file=zowe.config.json --from-file=zowe.schema.json --from-literal SYSUID=<SYSUID>```

Change SYSUID value according to your zOS username.

Mount ```zowe-config``` secret to operator deployment to ```/zowe-config```

```oc -n kzed-system set volume deployment/kzed-controller-manager --add --type=secret --secret-name=zowe-config --mount-path=/zowe-config```

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/kzed:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/kzed/<tag or branch>/dist/install.yaml
```
