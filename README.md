# Probes Convention Service

A sample conventions server for adding in liveiness/readiness/startup probes, volumes/volume mounts, container arguments, and environment variables.

## Component Overview

ADD IN HOW THIS THING WORKS!!!

## Prequisites

* Golang 1.20+ 
```
brew install go
```
* Pack CLI 
```
brew install buildpacks/tap/pack
```
* Set the default builder 
```
pack config default-builder paketobuildpacks/builder-jammy-tiny
```

## Build Image and Push Image to Repository

To build the image and push it to your repo you need to first set the ```DOCKER_ORG``` environment variable to the location to push the image and then run the ```make image``` command. This will build the image using ```pack``` and then push the image with the ```latest``` tag to the repo set in the ```DOCKER_ORG``` environment variable.

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

## An example Workload

Below is an example workload that configured two probes.

```
      1 + |---
      2 + |apiVersion: carto.run/v1alpha1
      3 + |kind: Workload
      4 + |metadata:
      5 + |  labels:
      6 + |    app.kubernetes.io/part-of: app-golang-kpack
      7 + |    apps.tanzu.vmware.com/has-tests: "true"
      8 + |    apps.tanzu.vmware.com/workload-type: web
      9 + |  name: app-golang-kpack
     10 + |  namespace: dev
     11 + |spec:
     12 + |  params:
     13 + |  - name: annotations
     14 + |    value:
     15 + |      x95castle1.org/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
     16 + |      x95castle1.org/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
     17 + |  - name: testing_pipeline_matching_labels
     18 + |    value:
     19 + |      apps.tanzu.vmware.com/pipeline: golang-pipeline
     20 + |  source:
     21 + |    git:
     22 + |      ref:
     23 + |        branch: main
     24 + |      url: https://github.com/carto-run/app-golang-kpack
```

## Generated PodSpec

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

```
tanzu apps workload create simple-conventions \
  --namespace dev \
  --git-branch main \
  --git-repo https://github.com/x95castle1/probes-convention-service \
  --label apps.tanzu.vmware.com/has-tests=true \
  --label app.kubernetes.io/part-of=simple-conventions \
  --param-yaml testing_pipeline_matching_labels='{"apps.tanzu.vmware.com/pipeline":"golang-pipeline"}' \
  --type web \
  --yes
```
