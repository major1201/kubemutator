# k8s-mutator

A Kubernetes resource mutator

[![GoDoc](https://godoc.org/github.com/major1201/k8s-mutator?status.svg)](https://godoc.org/github.com/major1201/k8s-mutator)
[![Go Report Card](https://goreportcard.com/badge/github.com/major1201/k8s-mutator)](https://goreportcard.com/report/github.com/major1201/k8s-mutator)

## Get start

### Set up a Kubernetes cluster

You should set up your self

### Generate a certification pair

Modify `examples/tls/ca.conf`, `examples/tls/csr-prod.conf`.

Generate certifate,

```bash
cd examples/tls
DEPLOYMENT=us-east-1 CLUSTER=PRODUCTION ./new-k8s-mutator-cert.rb
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
kubectl -n kube-system apply -f configmap.yaml  # !! rewrite configmap with your config file and mutator cert and key
kubectl -n kube-system apply -f deployment.yaml
kubectl -n kube-system apply -f service.yaml

# if you have prometheus operator deployed, you can add service monitor below
kubectl -n kube-system apply -f service-monitor.yaml
```

### Or you can deploy k8s-mutator with helm

```bash
cd examples/chart

# generate you custom values
helm inspect values k8s-mutator > custom.yaml

# make some changes to custom.yaml

# show what would happen next
helm template --name k8s-mutator --namespace kube-system -f custom.yaml k8s-mutator

# install to your Kubernetes cluster
helm install --name k8s-mutator --namespace kube-system -f custom.yaml k8s-mutator
```

## Configuration

An example configuration is in examples/conf/config.yml

```yaml
annotationKey: k8s-mutator.example.com/requests

strategies:
  - name: filebeat
    patches:
      # add filebeat sidecar
      - isTemplate: true
        data: |
          op: add
          path: /spec/containers/-
          value:
            name: filebeat
            image: myrepo/filebeat
            resources:
              requests:
                cpu: 100m
                memory: 128Mi
              limits:
                cpu: 100m
                memory: 128Mi
            volumeMounts:
              - name: logs
                mountPath: /var/log/{{ .Labels.k8s-app }}

rules:
  - namespace:
      - default
    selector:
      matchLabels:
        k8s-app: myapp
    strategies:
      - filebeat
```

*patch data see: <https://tools.ietf.org/html/rfc6902>*

1. For each rule in `rules`, match `namespace` and `selector`.
2. If match failed, match next.
3. If match succeeded, append the strategies to the strategy list.
4. Read the pod annotation prefixed by `annotationKey`, append the strategies joined by comma to the strategy list.
5. Merge all the strategy patches and response the JSONPatch object to the Kubernetes server.

## License

MIT
