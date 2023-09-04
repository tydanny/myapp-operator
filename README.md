# myapp-operator

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Compiling the code
Run the following command to build the code into an executable binary named `manager` in the bin directory.

```sh
make build
```

### Running tests
The following command will download the [ginkgo](https://onsi.github.io/ginkgo/) cli and use it to run the test suite. This test suite uses a mock kube-api and etcd store to simulate the kubernetes API.

```sh
make test
```

### Deploying to a cluster
1. Install the CRD:

```sh
make install
```

2. Build and push your image to the location specified by `IMG` (if left blank `controller:latest` is used):

```sh
make docker-build docker-push IMG=<some-registry>/myapp-operator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/myapp-operator:tag
```

### Creating the sample resource
Run the following command to create the sample resource as well as a namepsace.

```sh
make sample
```

### Connecting to the pod info endpoint
Run the following command to setup a port forward to the pod info API
```sh
kubectl port-forward service/podinfo 8080:8080 -n myapp
```

>**Note:** The command above assumes you're using the sample namespace `myapp`

You should now be able to hit the pod info endpoint(s) in another terminal. For example:
```sh
curl localhost:8080/
```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

