# Project Modules Architecture

Primarily there's two major kind of components to this part of the project.

## Admission Controller

This module listens for specific resources being created on the cluster and act accordingly to update the [manifests](https://faun.pub/understanding-the-kubernetes-manifest-e96d680f2a11) of the resources when necessary.

Service exposes endpoint and waits for kubernetes cluster to invoke it therefore it's important to register this through `MutatingWebhookConfiguration`. More on this is explained in choronological order in `DEV.md`.

## Operators

Operators in general are specialized state managers for specific kind of resources. The goal of the operator is to bring it's respective resources in a determinitic state by performing operations on it.

The operators in this project follow a pattern i.e. they work in conjunction with admission controller, the admission controller updates the manifests of resources while the operators are incharge of restarting the pods when required.

### Sidecar operator

The sidecar operator listens for pods being created with specific annotation and restarts them when they do not have a sidecar attached to them.

#### Specialized Sidecars

The code for the sidecars lives in `/sidecars`. So far there's postgres-sidecar implemented. These apps are packed as containers and will end up being attached to a pod based on the declaration in CRD.

### CRD operator

[CRD](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) operator listens for resources being created `kind: TororuResource`.

---

#### Highly Recommended

There's component level sequence diagrams in dir `/docs/component\ diagrams`. Please have a look for high level description of each flow.
