# image-transformer

# Prerequisite

## Deploy cert-manager

cert-manager is used by TLS certificates for admission webhooks.

```
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
```

## Install webhook

If no custom configuration is provided, the system defaults to using DaoCloud as the registry. Alternatively, you can deploy [crproxy](https://github.com/DaoCloud/crproxy) to set up your own registry.

You can use `ORIGINAL_REPO` to set the registry you want to modify.

The default configuration is `docker.io,gcr.io,ghcr.io,registry.k8s.io`

```
kubeclt apply -f deploy
```

