# Dev Setup

Please have a look the sequence diagrams in `/docs` to learn more about the high level design of this project.

## Pre-requisites

1. [Docker](https://docs.docker.com/engine/install/)
1. Local k8 cluster using any of the following tools:
   - [minikube](https://minikube.sigs.k8s.io/docs/start/) (preferred)
   - [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
   - [k3d](https://k3d.io/v5.6.0/)
1. [ngrok](https://ngrok.com/download) for port forwarding
1. [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl) to manage the local cluster.
1. [Make](https://www.gnu.org/software/make/) the build automation tool - most likely you will already have it but just in case.
1. [Golang](https://go.dev/doc/install)
1. [Earthly](https://docs.earthly.dev/)

## Steps

1. Local cluster: minikube/k3d/kind

   - If required to start a new cluster using config, you can do it using the `k8s/mock2/clusterWithWebhookRegistration/kinDCluster.yaml`

### Running local
1. Store the cluster config at `~/.kube/config`
1. Start ngrok tunnel `make tunnel-adm-controller` and get the tunnel url:

   ```bash
   curl -s http://localhost:4040/api/tunnels | jq '.tunnels[0].public_url'
   ```

1. Update the URL in `k8s/mock2/mutating-webhook.yaml`
1. Apply `k8s/mock2/mutating-webhook.yaml` using kubectl on your local cluster
1. Start the operators and adm_controller through `make run`
1. At this point, you should be able to see these components created on the cluster (you can use lens for that)
1. Try experimenting and user flows by applying manifests in `k8s/mock2` to see things in action :)

### Running as chart

1. Start ngrok tunnel `make tunnel-icp` and get the tunnel url:

   ```bash
   curl -s http://localhost:4040/api/tunnels | jq '.tunnels[0].public_url'
   ```
1. Update the URL in `tororu-operator` helm chart values, under `config.icpNodeUrl`
1. Update the canister id in `tororu-operator` helm chart values, under `config.canisterId`
1. Start the operators and adm_controller through `make install-chart`
1. Try experimenting and user flows by applying manifests in `k8s/mock2` to see things in action :)

## Demo

To run a complete demo, please follow these steps after you finish the previous setup.

- Create a tororu resource:
  1. Run `kubectl apply -f k8s/mock2/tororu-api/tororu-crd.yaml`
  1. Check the result under Custom Resources on lens
- Create two new CRDs:
  1. Prateek's secret: `kubectl apply -f k8s/mock2/postgres-crd-1.yaml`
  1. Juan's secret: `kubectl apply -f k8s/mock2/postgres-crd-2.yaml`
- Build the postgres sidecar image by running `make build-sidecars`
- Create a pvc `kubectl apply -f k8s/mock2/persistantVolume.yaml`
- Deploy the postgres server `kubectl apply -f k8s/mock2/postgres-server.yaml`
- Deploy the postgres client `kubectl apply -f k8s/mock2/postgres-client.yaml`
