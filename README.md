# Vault K8S Canister

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![GithubActions](https://github.com/Zondax/vault-k8s-canister/actions/workflows/checks.golem.yml/badge.svg)](https://github.com/Zondax/vault-k8s-canister/blob/master/.github/workflows/checks.golem.yml)

---

![zondax_light](docs/assets/zondax_light.png#gh-light-mode-only)
![zondax_dark](docs/assets/zondax_dark.png#gh-dark-mode-only)

_Please visit our website at [zondax.ch](https://www.zondax.ch)_

---

This repository is a Proof of Concept of a decentralized secret management solution that leverages Internet Computer (ICP) technology to compete with existing applications such as 1Password, Doppler or Hashicorp Vault.
Our project aims to provide means for services to share secrets in a flexible, transparent and secure way. It simplifies the flow of secret management between consumers in the cluster as well as rotate secrets based on config for added security.

In the future we want to keep building on top of what we have now by adding support for various secret types, adding more specialized sidecars as well as improve the distribution methods.


## About the project :book::book:

Please visit the folder `docs` or the [documentation site](https://docs.zondax.ch) for more information! 


## How to try by yourself :gear:
### Pre-requisites

1. [Docker](https://docs.docker.com/engine/install/)
1. Local k8 cluster using any of the following tools:
   - [minikube](https://minikube.sigs.k8s.io/docs/start/) (preferred)
   - [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
   - [k3d](https://k3d.io/v5.6.0/)
1. [cloudflared](https://github.com/cloudflare/cloudflared) for port forwarding
1. [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to manage the local cluster.
1. [Make](https://www.gnu.org/software/make/) the build automation tool - most likely you will already have it but just in case.
1. [Golang](https://go.dev/doc/install)
1. [Earthly](https://docs.earthly.dev/)

### Steps

1. Local cluster: minikube/k3d/kind
   - If required to start a new cluster using config, you can do it using the `k8s/mock2/clusterWithWebhookRegistration/kinDCluster.yaml`
1. Build the postgres sidecar image by running `make build-sidecars`

#### Running local
1. Store the cluster config at `~/.kube/config`
1. Start cloudflared tunnel `make tunnel-adm-controller` and get the tunnel url from output file
1. Update the URL in `k8s/mock2/mutating-webhook.yaml`
1. Apply `k8s/mock2/mutating-webhook.yaml` using kubectl on your local cluster
1. Start the operators and adm_controller through `make run`
1. At this point, you should be able to see these components created on the cluster (you can use lens for that)
1. Try experimenting and user flows by applying manifests in `k8s/mock2` to see things in action :)

#### Running as chart

1. Start cloudflared tunnel `make tunnel-icp` and get the tunnel url:
1. Update the URL in `tororu-operator` helm chart values, under `config.icpNodeUrl`
1. Update the canister id in `tororu-operator` helm chart values, under `config.canisterId`
1. Start the operators and adm_controller through `make install-chart`
1. Try experimenting and user flows by applying manifests in `k8s/mock2` to see things in action :)

### Demo

To run a complete demo, please follow these steps after you finish the previous setup.

#### Option A: Using manifests
- Create a tororu resource:
  1. Run `kubectl apply -f k8s/mock2/tororu-api/tororu-crd.yaml`
  1. Check the result under Custom Resources on lens
- Create two new CRDs:
  1. Prateek's secret: `kubectl apply -f k8s/mock2/postgres-crd-1.yaml`
  1. Juan's secret: `kubectl apply -f k8s/mock2/postgres-crd-2.yaml`

#### Option B:  Using helm charts
- Create a tororu resource, deploying two new CRDs:
    1. Run `make install-crds`
    1. Check the result under Custom Resources on lens


#### Creating RW and RO consumers
- Create a pvc `kubectl apply -f k8s/mock2/persistantVolume.yaml`
- Deploy the postgres server `kubectl apply -f k8s/mock2/postgres-server.yaml`
- Deploy the postgres client `kubectl apply -f k8s/mock2/postgres-client.yaml`
