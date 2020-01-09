---
title: Well known custom actions
weight: 805
---

# Well known custom actions

This section is non-normative, but is here to propose a common set of optional actions a CNAB package MAY implement, and that CNAB tools can understand.

A CNAB indicates that it supports those actions by including them in its custom action list (as defined in [the bundle definition](101-bundle-json.md)).
- `io.cnab.dry-run` (with `stateless`: true and `modifies`: false): execute the installation in a dry-run mode, allowing to see what would happen with the given set of parameter values.
- `io.cnab.help` (with `stateless`: true and `modifies`: false): print an help message to the standard output. Implementations MAY print different messages depending on the parameters values passed to the invocation image.
- `io.cnab.log` (with `stateless`: false and `modifies`: false): print logs of the installed system to the standard output.
- `io.cnab.status` (with `stateless`: false and `modifies`: false): print a human readable status message to the standard output. This action also produces an output file named `status` describing the detailed status as in the example below ([here](schema/status.schema.json) is the JSON schema for this):
	```json
    {
      "components": {
        "backend": {
          "status": "Ready"
        },
        "front-end": {
          "message": "Failed to pull Docker image mginx:latest",
          "status": "Failed"
        }
      },
      "message": "Component front-end failed to deploy properly",
      "status": "Failed"
    }
    ```
    Source: [804.01-status.json](examples/804.01-status.json)
  - `status` is a value indicating the deployment status of a component or of the whole bundle. Well known values are `Failed`, `Ready`, `Pending`. An invocation image MAY use other custom-values, but tools must at least understand those values.
  - `message` optional field giving details about the current bundle or component status.
  - `components`: map of components/statuses allowing a more detailed status. A component can itself have subcomponents, to enable scenarios like aggregating statuses in CNABs composed of other CNABs. e.g.:
	```json
    {
      "components": {
        "gke": {
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
          },
          "status": "Ready"
        },
        "wordpress": {
          "status": "Ready"
        }
      },
      "status": "Ready"
    }
    ```
    Source: [804.02-status.json](examples/804.02-status.json)
  - Additionaly to the fields defined above, an invocation image MAY add custom fields as they want. To avoid ambiguous naming, those fields names MAY be namespaced. An example of that could be a CNAB bundle deploying Kubernetes workloads and exposing scale details:
    ```json
    {
      "components": {
        "backend": {
          "com.example.scale": {
            "actual": 3,
            "desired": 3,
            "failed": 0,
            "pending": 0
          },
          "status": "Ready"
        },
        "front-end": {
          "com.example.scale": {
            "actual": 2,
            "desired": 3,
            "failed": 0,
            "pending": 1
          },
          "status": "Pending"
        }
      },
      "status": "Pending"
    }
    ```
    Source: [804.03-status.json](examples/804.03-status.json)

Next section: [Disconnected Scenarios](805-airgap.md)
