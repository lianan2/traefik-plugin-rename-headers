# Rename Header

Traefik custom headers plugin is a middleware plugin for [Traefik](https://traefik.io) which renames headers in the response, while keeping their values.

## Configuration

### Static

```yaml
# File(YAML)
pilot:
  token: "xxxx"

experimental:
  plugins:
    renameHeaders:
      modulename: "gitlab.com/lianan2/traefik-plugin-rename-headers"
      version: "v0.0.1"
```

### Dynamic

To configure the Rename Headers plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). 
The following example creates and uses the renameHeaders middleware plugin to rename the "custom_id" header

```yaml
# K8S CRD
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: rename-headers
  namespace: my-namespace
spec:
  plugin:
    rename-headers-traefik-plugin:
      rename:
        - headerName: "X-Custom-Header1"
          newHeaderName: "X-Custom-Header2"
```
