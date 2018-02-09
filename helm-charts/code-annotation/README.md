# Code Annotation Tool

This chart deploys source{d} [code annotation tool](https://github.com/src-d/code-annotation)

## Pre-requisites

* Kubernetes 1.4+ with Beta APIs enabled

## Installing the chart

All parameters under `settings` in [values.yaml](values.yaml) must be provided.

```
helm install -n <release-name> <chart path> --set \
secrets.jwt_signing_key=<signing-key>,\
secrets.github_client=<client-id>,\
secrets.github_secret=<client-secret>,\
ingress.hostname=<hostname>,\
ingress.globalStaticIpName=<static-ip>,\
image.tag=<your-favourite.tag>
```

These are the mandatory parameters that need to be provided or installation will fail.
Other parameters can be provided too but, if not, a default value will be used.

It's also possible to set name of separately deployed secrets using parameter `secretName`.

Full command for deployment:

```
helm install -n <release-name> <chart path> --set \
secretName=<name-of-secret-object>,\
ingress.hostname=<hostname>,\
ingress.globalStaticIpName=<static-ip>,\
image.tag=<your-favourite.tag>
```

# Configuration

Please refer to [values.yaml](values.yaml) for the full run-down on defaults.

To override any of those default values,
specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided
while installing the chart.
For example,

```bash
$ helm install --name <release-name> -f values.yaml <path-to-chart>
```

> **Tip**: You can use the default [values.yaml](values.yaml)
