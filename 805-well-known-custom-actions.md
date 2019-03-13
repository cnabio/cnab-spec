---
title: Well known custom actions
---

# Well known custom actions

This section is non-normative, but is here to propose a common set of optional actions a CNAB package MAY implement, and that CNAB tools can understand.

A CNAB indicates that it supports those actions by including them in its custom action list (as defined in [the bundle definition](101-bundle-json.md)).
- `io.cnab.dry-run` (with `stateless`: true and `modifies`: false): execute the installation in a dry-run mode, allowing to see what would happen with the given set of parameter values.
- `io.cnab.help` (with `stateless`: true and `modifies`: false): print an help message to the standard output. Implementations MAY print different messages depending on the parameters values passed to the invocation image.
- `io.cnab.log` (with `stateless`: false and `modifies`: false): print logs of the installed system to the standard output.
- `io.cnab.status` (with `stateless`: false and `modifies`: false): print a human readable status message to the standard output.
- `io.cnab.status+json` (with `stateless`: false and `modifies`: false): print a json payload describing the detailed status with the following schema ([here](schema/status.schema.json) is the json schema for this):

```json
{
    "status": "Failed",
    "message": "Component front-end failed to deploy properly",
    "components":{
        "backend": {
            "status": "Ready"
        },
        "front-end": {
            "status": "Failed",
            "message": "Failed to pull Docker image mginx:latest",
        },
    }
}
```
  - `status` is a value indicating the deployment status of a component or of the whole bundle. Well known values are `Failed`, `Ready`, `Pending`. An invocation image MAY use other custom-values, but tools must at least understand those values.
  - `message` optional field giving details about the current bundle or component status.
  - `components`: map of components/statuses allowing a more detailed status. A component can itself have subcomponent, to enable scenarios like aggregating statuses in CNABs composed of other CNABs. e.g.:
```json
{
    "status": "Ready",
    "components": {
        "gke": {
            "status" : "Ready",
            "components": {
                "kube-api": {
                    "status": "Ready"
                },
                "network": {
                    "status": "Ready"
                },
                "scheduler": {
                    "status": "Ready"
                }
            }
        },
        "wordpress": {
            "status": "Ready"
        }
    }
}
```
  - Additionaly to the fields defined above, an invocation image MAY add custom fields as they want. To avoid ambiguous naming, those fields names MAY be namespaced. An example of that could be a CNAB bundle deploying Kubernetes workloads and exposing scale details:
```json
{
    "status": "Pending",
    "components":{
        "backend": {
            "status": "Ready",
            "com.example.scale": {
                "desired": 3,
                "actual": 3,
                "failed": 0,
                "pending": 0
            }
        },
        "front-end": {
            "status": "Pending",
            "com.example.scale": {
                "desired": 3,
                "actual": 2,
                "failed": 0,
                "pending": 1
            }
        },
    }
}
```


Next section: [Standardization Process](901-process.md)