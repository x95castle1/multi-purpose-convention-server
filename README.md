# Multi-Purpose Convention Server

A sample convention server for adding in liveiness/readiness/startup probes, volumes/volume mounts, container arguments, node affinity, tolerations, and environment variables to a pod spec for a TAP workload.

## Component Overview

This project can be used as a template/exemplar to create your own conventions for a Supply Chain. Boilerplate code with a handler and convention interface has been moved to the convention-server-framework package. To reuse this code you just need to implement your own convention.go logic. 

### server.go

This creates a basic http server to handle webhook calls from the Convention controller. It calls the handler to execute your conventions. 

This component shouldn't need changes (unless you have different logging needs, etc.)

### convention.go 

This contains the logic for your conventions. Each convention is part of variable array that overrides the functions in the convention interface from the framework package. 

## Convention Architecture

[Cartographer Convention Documentation](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/1.6/tap/cartographer-conventions-about.html)

![arch](images/convention-architecture.jpg)

## Prequisites

* [Golang 1.20+](https://go.dev/doc/install)
```shell
brew install go
```
* [Pack CLI](https://buildpacks.io/docs/tools/pack/)
```shell
brew install buildpacks/tap/pack
```
* [Set the default builder](https://buildpacks.io/docs/tools/pack/cli/pack_config_default-builder/)
```shell
pack config default-builder paketobuildpacks/builder-jammy-tiny
```
* [Tanzu CLI](https://docs.vmware.com/en/VMware-Tanzu-Application-Platform/1.6/tap/install-tanzu-cli.html)

* [Kctrl CLI](https://github.com/carvel-dev/carvel) - Needed for bundling and releasing as a Carvel Package

## Available Options

| Annotation | Description | 
| --- | --- |
| `example.com/livenessProbe` | define a liveness probe | 
| `example.com/readinessProbe` | define a readiness probe |
| `example.com/startupProbe` | define a startup probe |
| `example.com/storage` | define volume and volume mounts |
| `example.com/args` | define container args |
| `example.com/tolerations` | define tolerations for a pod |
| `example.com/nodeSelector` | define a node selector for a pod |
| `example.com/affinity` | define scheduling affinity for a pod |

## Example Annotations for a Workload

```yaml
spec:
  params:
  - name: annotations
    value:
      example.com/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
      example.com/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
      example.com/startupProbe: '{"httpGet":{"path":"/healthz","port":"liveness-port"},"failureThreshold":30,"periodSeconds":10}'
      example.com/storage: '{"volumes":[{"name":"config-vol","configMap":{"name":"log-config","items":[{"key":"log_level","path":"log_level"}]}}],"volumeMounts":[{"name":"config-vol","mountPath":"/etc/config"}]}'
      example.com/args: '{["HOSTNAME","KUBERNETES_PORT"]}'
      example.com/tolerations: '[{"key":"rabeyta","operator":"Exists","effect":"NoSchedule"}]'
      example.com/nodeSelector: '{"disktype":"ssd"}'
      example.com/affinity: '{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"topology.kubernetes.io/zone","operator":"In","values":["antarctica-east1","antarctica-west1"]}]}]},"preferredDuringSchedulingIgnoredDuringExecution":[{"weight":1,"preference":{"matchExpressions":[{"key":"another-node-label-key","operator":"In","values":["another-node-label-value"]}]}}]}}'

```

It can sometimes be tricky to convert yaml to json to pass through the annotation. You can use these utilities:

* [Convert Yaml to JSON](https://onlineyamltools.com/convert-yaml-to-json)
* [Compact JSON](https://www.text-utils.com/json-formatter/)

## An example Workload

Below is an example workload that configured two probes.

```yaml
apiVersion: carto.run/v1alpha1
kind: Workload
metadata:
  labels:
    app.kubernetes.io/part-of: app-golang-kpack
    apps.tanzu.vmware.com/workload-type: web
  name: convention-workload
  namespace: jeremy
spec:
  params:
  - name: annotations
    value:
      example.com/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
      example.com/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
  source:
    git:
      ref:
        branch: main
      url: https://github.com/carto-run/app-golang-kpack   
```

You can find more examples in the [workload-examples folder](/workload-examples/) in the repository.

## Example Generated PodSpec with Probes

```yaml
...
spec:
  containers:
  - image: gcr.io/ship-interfaces-dev/supply-chain/app-golang-kpack-dev@sha256:3830de13d0a844420caa3d0a8d77ee1ca5b05897a273465c682032522fc331b5
    livenessProbe:
      exec:
        command:
        - cat
        - /tmp/healthy
      initialDelaySeconds: 5
      periodSeconds: 5
    name: workload
    readinessProbe:
      httpGet:
        path: /healthz
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 5
    resources: {}
...
```

## Install on a Cluster

The multi-purpose-convention-server has been conveniently packaged up via Carvel and can be installed on a TAP cluster via the Tanzu CLI.

### Install via Carvel Package

Run the following command to output a list of available tags.

  ```shell
  imgpkg tag list -i projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo | sort -V
  ```

  For example:

  ```shell
  imgpkg tag list -i projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo | sort -V

  0.1.0
  0.2.0
  0.3.0
  0.4.0
  ```

Use the latest version returned by the command above.

We recommend to relocate the images from VMware Tanzu Network registry to
your own container image registry before installing.

1. Set up environment variables for installation by running:

    ```shell
    export INSTALL_REGISTRY_USERNAME=MY-REGISTRY-USER
    export INSTALL_REGISTRY_PASSWORD=MY-REGISTRY-PASSWORD
    export INSTALL_REGISTRY_HOSTNAME=MY-REGISTRY
    export VERSION=VERSION-NUMBER
    export INSTALL_REPO=TARGET-REPOSITORY
    ```

    Where:

    - `MY-REGISTRY-USER` is the user with write access to MY-REGISTRY.
    - `MY-REGISTRY-PASSWORD` is the password for `MY-REGISTRY-USER`.
    - `MY-REGISTRY` is your own registry.
    - `VERSION` is your Multi-Purpose-Convention-Server version. For example, `0.4.0`.
    - `TARGET-REPOSITORY` is your target repository, a directory or repository on
      `MY-REGISTRY` that serves as the location for the installation files for
      the conventions.

2. Relocate the images with the imgpkg CLI by running:

    ```shell
    imgpkg copy -b projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo:${VERSION} --to-repo ${INSTALL_REGISTRY_HOSTNAME}/${INSTALL_REPO}/multi-purpose-convention-server-bundle-repo
    ```

3. Add Multi-Purpose-Convention-Server package repository to the cluster by running:

    ```shell
    tanzu package repository add multi-purpose-conventions-repository \
      --url ${INSTALL_REGISTRY_HOSTNAME}/${INSTALL_REPO}/multi-purpose-convention-server-bundle-repo:$VERSION \
      --namespace tap-install
    ```

4. Get the status of Multi-Purpose-Convention-Server  package repository, and ensure that the status updates to `Reconcile succeeded` by running:


    ```shell
    tanzu package repository get multi-purpose-conventions-repository --namespace tap-install
    ```

    For example:

    ```console
    tanzu package repository get multi-purpose-conventions-repository --namespace tap-install

    NAMESPACE:               tap-install
    NAME:                    multi-purpose-conventions-repository
    SOURCE:                  (imgpkg) projects.registry.vmware.com/tanzu_practice/conventions/multi-purpose-convention-server-bundle-repo:0.4.0
    STATUS:                  Reconcile succeeded
    CONDITIONS:              - type: ReconcileSucceeded
      status: "True"
      reason: ""
      message: ""
    USEFUL-ERROR-MESSAGE:
    ```

5. List the available packages by running:

    ```shell
    tanzu package available list --namespace tap-install
    ```

    For example:

    ```shell
    $ tanzu package available list --namespace tap-install
    / Retrieving available packages...
      NAME                                                              DISPLAY-NAME                       SHORT-DESCRIPTION
      multi-purpose-convention-server.conventions.tanzu.vmware.com      multi-purpose-convention-server    Set of conventions to enrich pod spec with volumes, probes, affinities
    ```

### Prepare Convention Configuration

You can define the `--values-file` flag to customize the default configuration. You
must define the following fields in the `values.yaml` file for the Convention Server
configuration. You can add fields as needed to activate or deactivate behaviors.
You can append the values to the `values.yaml` file. Create a `values.yaml` file
by using the following configuration:

    ```yaml
    ---
    annotationPrefix: ANNOTATION-PREFIX
    ```

    Where:

    - `ANNOTATION-PREFIX` is the prefix you want to use on your annotation used in the workload. For example: `x95castle1` Defaults to `example.com`.

### Install Multi-Purpose-Convention-Server

Define the `--values-file` flag to customize the default configuration (Optional):

The `values.yaml` file you created earlier is referenced with the `--values-file` flag when running your Tanzu install command:

```shell
tanzu package install REFERENCE-NAME \
  --package SCANNER-NAME \
  --version VERSION \
  --namespace tap-install \
  --values-file PATH-TO-VALUES-YAML
```

Where:

- `ANNOTATION-PREFIX` is the prefix you want to use on your annotation used in the workload. For example: `x95castle1.org` Defaults to `example.com`.

For example:

```console
tanzu package install multi-purpose-convention-server  \
--package multi-purpose-convention-server.conventions.tanzu.vmware.com \
--version 0.4.0 \
--namespace tap-install
```

## Install Locally



### Build Image and Push Image to Repository

To build the image and push it to your repo you need to first set the `DOCKER_ORG` environment variable to the location to push the image and then run the `make image` command. This will build the image using `pack` and then push the image with the `latest` tag to the repo set in the `DOCKER_ORG` environment variable.

```shell
export DOCKER_ORG=registry.harbor.learn.tapsme.org/convention-service

make image
```

### Installation

To install the conventions server onto the Cluster use: 

```shell
make install
```

This will create a new namespace `multi-purpose-convention` and configure cartographer conventions to use this convention provider.


## Build the service using TAP

You can also use TAP to build and deploy the server to make it available as a convention server.

```shell
tanzu apps workload create multi-purpose-convention-server \
  --namespace dev \
  --git-branch main \
  --git-repo https://github.com/x95castle1/multi-purpose-convention-server \
  --label apps.tanzu.vmware.com/has-tests=true \
  --label app.kubernetes.io/part-of=multi-purpose-convention-server \
  --param-yaml testing_pipeline_matching_labels='{"apps.tanzu.vmware.com/pipeline":"golang-pipeline"}' \
  --type web \
  --yes
```
