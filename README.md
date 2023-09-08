# simple-conventions

a sample conventions server for adding in liveiness/readiness probes.

## Installation

To install the conventions server use: 

```
make install
```

This will create a new namespace `simple-convention` and configure cartographer conventions to use this convention provider.

## Available Options

| Annotation | Description |
| --- | --- |
| `garethjevans.org/livenessProbe` | define a liveness probe |
| `garethjevans.org/readinessProbe` | define a readiness probe |
| `garethjevans.org/startupProbe` | define a startup probe |

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
     15 + |      garethjevans.org/livenessProbe: '{"exec":{"command":["cat","/tmp/healthy"]},"initialDelaySeconds":5,"periodSeconds":5}'
     16 + |      garethjevans.org/readinessProbe: '{"httpGet":{"path":"/healthz","port":8080},"initialDelaySeconds":5,"periodSeconds":5}'
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

### Build this in TAP

```
tanzu apps workload create simple-conventions \
  --namespace dev \
  --git-branch main \
  --git-repo https://github.com/garethjevans/simple-conventions \
  --label apps.tanzu.vmware.com/has-tests=true \
  --label app.kubernetes.io/part-of=simple-conventions \
  --param-yaml testing_pipeline_matching_labels='{"apps.tanzu.vmware.com/pipeline":"golang-pipeline"}' \
  --type web \
  --yes
```
