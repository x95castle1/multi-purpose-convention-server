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
```
brew install go
```
* [Pack CLI](https://buildpacks.io/docs/tools/pack/)
```
brew install buildpacks/tap/pack
```
* [Set the default builder](https://buildpacks.io/docs/tools/pack/cli/pack_config_default-builder/)
```
pack config default-builder paketobuildpacks/builder-jammy-tiny
```

## Build Image and Push Image to Repository

To build the image and push it to your repo you need to first set the `DOCKER_ORG` environment variable to the location to push the image and then run the `make image` command. This will build the image using `pack` and then push the image with the `latest` tag to the repo set in the `DOCKER_ORG` environment variable.

```
export DOCKER_ORG=registry.harbor.learn.tapsme.org/convention-service

make image

```

## Installation

To install the conventions server onto the Cluster use: 

```
make install
```

This will create a new namespace `simple-convention` and configure cartographer conventions to use this convention provider.

## Available Options

| Annotation | Description | 
| --- | --- |
| `x95castle1.org/livenessProbe` | define a liveness probe | 
| `x95castle1.org/readinessProbe` | define a readiness probe |
| `x95castle1.org/startupProbe` | define a startup probe |
| `x95castle1.org/storage` | define volume and volume mounts |
| `x95castle1.org/args` | define container args |
| `x95castle1.org/tolerations` | define tolerations for a pod |
| `x95castle1.org/nodeSelector` | define a node selector for a pod |
| `x95castle1.org/affinity` | define scheduling affinity for a pod |

## Example Annotations for a Workload

```
spec:
  params:
  - name: annotations
    value:
      x95castle1.org/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
      x95castle1.org/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
      x95castle1.org/startupProbe: '{"httpGet":{"path":"/healthz","port":"liveness-port"},"failureThreshold":30,"periodSeconds":10}'
      x95castle1.org/storage: '{"volumes":[{"name":"config-vol","configMap":{"name":"log-config","items":[{"key":"log_level","path":"log_level"}]}}],"volumeMounts":[{"name":"config-vol","mountPath":"/etc/config"}]}'
      x95castle1.org/args: '{["HOSTNAME","KUBERNETES_PORT"]}'
      x95castle1.org/tolerations: '[{"key":"rabeyta","operator":"Exists","effect":"NoSchedule"}]'
      x95castle1.org/nodeSelector: '{"disktype":"ssd"}'
      x95castle1.org/affinity: '{"nodeAffinity":{"requiredDuringSchedulingIgnoredDuringExecution":{"nodeSelectorTerms":[{"matchExpressions":[{"key":"topology.kubernetes.io/zone","operator":"In","values":["antarctica-east1","antarctica-west1"]}]}]},"preferredDuringSchedulingIgnoredDuringExecution":[{"weight":1,"preference":{"matchExpressions":[{"key":"another-node-label-key","operator":"In","values":["another-node-label-value"]}]}}]}}'

```

It can sometimes be tricky to convert yaml to json to pass through the annotation. You can use these utilities:

* [Convert Yaml to JSON](https://onlineyamltools.com/convert-yaml-to-json)
* [Compact JSON](https://www.text-utils.com/json-formatter/)

## An example Workload

Below is an example workload that configured two probes.

```
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
      x95castle1.org/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
      x95castle1.org/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
  source:
    git:
      ref:
        branch: main
      url: https://github.com/carto-run/app-golang-kpack   
```

You can find more examples in the [workload-examples folder](/workload-examples/) in the repository.

## Example Generated PodSpec with Probes

```
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

### Build the service using TAP

You can also use TAP to build and deploy the server to make it available as a convention server.

```
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
