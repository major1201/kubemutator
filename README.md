# k8s-mutator

A Kubernetes resource mutator

[![GoDoc](https://godoc.org/github.com/major1201/k8s-mutator?status.svg)](https://godoc.org/github.com/major1201/k8s-mutator)
<!--[![Go Report Card](https://goreportcard.com/badge/github.com/major1201/k8s-mutator)](https://goreportcard.com/report/github.com/major1201/k8s-mutator)-->

## Get start

### Set up a Kubernetes cluster

You should set up your self

### Generate a certification pair

Modify `examples/tls/ca.conf`, `examples/tls/csr-prod.conf`.

Generate certifate,

```bash
cd examples/tls
DEPLOYMENT=us-east-1 CLUSTER=PRODUCTION ./new-cluster-injector-cert.rb
```

### Apply MutatingWebhookConfiguration

1. Encode ca.crt with base64

```bash
cat examples/tls/us-east-1/PRODUCTION/ca.crt | base64
```

2. First replace `examples/kubernetes/mutating-webhook-configuration.yaml` `webhooks[0].clientConfig.caBundle` with your ca.crt base64 generated before

3. Apply the mutating webhook configuration

```bash
kubectl -n kube-system apply -f examples/kubernetes/mutating-webhook-configuration.yaml
```

### Apply other kubernetes configurations

```bash
kubectl -n kube-system apply -f serviceaccount.yaml
kubectl -n kube-system apply -f clusterrole.yaml
kubectl -n kube-system apply -f clusterrolebinding.yaml
kubectl -n kube-system apply -f configmap.yaml  # !! rewrite configmap with your config file and sidecar injector cert and key
kubectl -n kube-system apply -f deployment.yaml
kubectl -n kube-system apply -f service.yaml

# if you have prometheus operator deployed, you can add service monitor below
kubectl -n kube-system apply -f service-monitor.yaml
```

## Configuration

An example configuration is in examples/conf/config.yml

1. For each rule in `rules`, match `namespace` and `selector`.
2. If match failed, match next.
3. If match succeeded, append the strageties to the strategy list.
4. Merge all the strategy patches and response the JSONPatch object to the Kubernetes server.

## License

MIT
