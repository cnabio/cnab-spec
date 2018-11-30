# Credential Sets

This is a non-normative section that describes how credentials can be passed into an invocation image. This strategy is implemented by `duffle`.

## Credential Mapping

In the bundle descriptor, credentials are declared like this:

```json
{
    "credentials": {
        "kubeconfig": {
            "path": "/home/.kube/config",
        },
        "image_token": {
            "env": "AZ_IMAGE_TOKEN",
        },
        "hostkey": {
            "path": "/etc/hostkey.txt",
            "env": "HOST_KEY"
        }
    }
}
```

The description above indicates that there are three credentials required to execute this bundle.

The names on each credential are unique within the bundle:

- kubeconfig
- image_token
- hostkey

That name can be used to reference the credential from tooling. Consequently, tooling may choose to construct a map from an external value to the credential name.

Duffle, for example, maps local credentials to bundle credentials by _credential sets_. Here is an example:

```json
name: test_credentials
credentials:
- name: kubeconfig
  source:
    path: $HOME/.kube/config
- name: image_token
  source:
    value: "1234aaaaa"
- name: hostkey
  source:
    env: HOSTKEY
```

This credential set tells the client application (`duffle`) to map CNAB bundle credential sections to certain local values:

- When a bundle requests `kubeconfig`, a file will be loaded from the local file system and injected into the container at `/home/.kube/config`
- When a bundle requests `image_token`, the literal value `1234aaaaa` will be loaded into the environment variable `AZ_IMAGE_TOKEN
- When a bundle requests `hostkey`, the local environment variable `$HOSTKEY` will be dereferenced, and its value injected into the container's `HOST_KEY` variable as well as the container's filesystem path `/etc/hostkey.txt`.

Similar tooling could choose to load the values by name from a database, vault, or file.

## Credential Injection

Credentials must be injected into the runtime of the invocation image. The following strategies are known to work:

- Environment Variables:
  - Injected via environment variable services (Docker and Kubernetes)
  - Loaded into VMs via cloud-init
- Paths:
  - Mounted as volumes (or wrappers, such as Kubernetes secrets)
  - Injected at runtime into the top layer of the container
  - Loaded into VMs via cloud-init