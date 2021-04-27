---
title: Credential and Parameter Sets
weight: 802
---

# Credential and Parameter Sets

This is a non-normative section that describes how credentials and parameters MAY be passed into an invocation image. This strategy is implemented by [Duffle] and [Porter].

[Duffle]: https://duffle.sh
[Porter]: https://porter.sh

## Credential Mapping

In the bundle descriptor, credentials are declared like this:

```json
{
  "credentials": {
    "hostkey": {
      "env": "HOST_KEY",
      "path": "/etc/hostkey.txt"
    },
    "image_token": {
      "env": "AZ_IMAGE_TOKEN"
    },
    "kubeconfig": {
      "path": "/home/.kube/config"
    }
  }
}
```

The description above indicates that there are three credentials required to execute this bundle.

The names on each credential are unique within the bundle:

- `kubeconfig`
- `image_token`
- `hostkey`

That name MAY be used to reference the credential from tooling. Consequently, tooling MAY choose to construct a map from an external value to the credential name.

Porter, for example, maps local credentials to bundle credentials by _credential sets_. Here is an example:

```json
{
   "name": "test_credentials",
   "namespace": "test",
   "created": "2020-01-20T12:00:00.000Z",
   "modified": "2020-01-28T18:30:00.111Z",
   "labels": {
     "env": "test"
   },
   "credentials": [
      {
         "name": "kubeconfig",
         "source": {
            "path": "$HOME/.kube/config"
         }
      },
      {
         "name": "image_token",
         "source": {
            "secret": "deploy_token"
         }
      },
      {
         "name": "hostkey",
         "source": {
            "env": "HOSTKEY"
         }
      }
   ]
}
```
Source: [802.01-credential-set.json](/examples/802.01-credential-set.json)

This credential set tells the client application (`porter`) to map CNAB bundle credential sections to certain local values:

- When a bundle requests `kubeconfig`, a file will be loaded from the local file system and injected into the container at `/home/.kube/config`
- When a bundle requests `image_token`, the literal value `1234aaaaa` will be loaded into the environment variable `AZ_IMAGE_TOKEN`
- When a bundle requests `hostkey`, the local environment variable `$HOSTKEY` will be dereferenced, and its value injected into the container's `HOST_KEY` variable as well as the container's filesystem path `/etc/hostkey.txt`.

Similar tooling could choose to load the values by name from a database, vault, or file.

The created and modified timestamps are in [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format.

## Parameter Mapping

Parameter values MAY be mapped by the runtime using parameter sets, which work the same as credential sets.

In the bundle.json parameters are defined:

```json
{
 "parameters": {
    "backend_port": {
      "definition": "http_port",
      "description": "The port that the back-end will listen on",
      "destination": {
        "env": "BACKEND_PORT"
      }
    }
  }
}
```

The following parameter set specifies that the hard-coded value "8080" should be used as the value for the backend_port parameter:

```json
{
   "name": "local-dev-parameters",
   "namespace": "dev",
   "created": "2020-01-20T12:00:00.000Z",
   "modified": "2020-01-28T18:30:00.111Z",
   "labels": {
     "env": "dev"
   },
   "parameters": [
      {
         "name": "backend_port",
         "source": {
            "value": "8080"
         }
      }
   ]
}
```
Source: [802.01-parameter-set.json](/examples/802.01-parameter-set.json)

### Namespaces

Credential and Parameter Sets MAY be scoped to a namespace.
The combination of namespace and name MUST be unique.
When a namespace is unset or empty, the document is considered to be global.
Documents in a namespace MAY reference a global document, but global documents MUST NOT reference namespaced documents and namespaced documents MUST NOT reference documents in a different namespace.

Namespaces MUST follow the [CNAB Namespace Format].

### Labels

Credential and Parameter Sets MAY define labels which MAY be used by storage providers to query for the document.
For example, retrieving credential sets intended for use within a particular environment such as test or production.

How labels are used for querying is out-of-scope of this spec and is up to the implementing storage provider.
Labels MUST follow the [CNAB Label Format].

## Well-known Value Sources

This document does not cover all possible sources for values and the supported set of sources is up to the implementing runtime.
Below are some well-known sources of values that a runtime MAY support:

* value: A hard-coded value stored in the credential or parameter set.
* env: The value is stored in the specified environment variable on the host.
* path: The value is stored in a file on the host at the specified path.
* secret: The value is stored in a secret store under the specified key.

## Value Injection

Credential and parameters MUST be injected into the runtime of the invocation image. The following strategies are known to work:

- Environment Variables:
  - Injected via environment variable services (Docker and Kubernetes)
  - Loaded into VMs via cloud-init
- Paths:
  - Mounted as volumes (or wrappers, such as Kubernetes secrets)
  - Injected at runtime into the top layer of the container
  - Loaded into VMs via cloud-init

Next section: [Well known custom actions](804-well-known-custom-actions.md)

[CNAB Label Format]: /105-labels.md
[CNAB Namespace Format]: /106-namespaces.md
