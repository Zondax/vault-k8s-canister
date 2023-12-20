# tororu-operator

For development purposes, you will need to have Docker Desktop or similar
and run a local Kubernetes cluster.

Then you can run

```bash
make helm-install
```

## admission controller local development

### Notes

Then you can run a tunnel using:

```bash
ngrok http 8282
```

Download the pem file and convert it to base64:

```bash
echo | openssl s_client -showcerts -servername 93a2-2a09-bac5-1e7b-282-00-40-3d.eu.ngrok.io -connect 93a2-2a09-bac5-1e7b-282-00-40-3d.eu.ngrok.io:443 2>/dev/null | openssl x509 -inform pem | base64
```

Otherwise, let's encrypt
Download the pem file and convert it to base64:

```bash
curl -s https://letsencrypt.org/certs/isrgrootx1.pem | base64
```

If you want to debug the webhook, you can find the info in grafana using something like:

```text
{component="kube-apiserver"}
```

Clear up events

```bash
kubectl delete events --all-namespaces --all
```

List webhooks

```bash
kubectl get mutatingwebhookconfigurations
```

List webhooks

```bash
kubectl get mutatingwebhookconfigurations webhook.tororu.io -o yaml
```

Get ngrok url

```bash
curl -s http://localhost:4040/api/tunnels | jq '.tunnels[0].public_url'
```
