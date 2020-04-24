---
title: Well known custom extensions
weight: 810
---

# Well known custom actions

This section is non-normative, but is here to propose a common set of custom extensions a bundle MAY define, and that CNAB tools MAY implement.

A bundle indicates that it uses a custom extension by including them in its custom section (as defined in [Custom Extensions](101-bundle-json.md#custom-extensions)).

## Dependencies

See the [Dependencies Specification](500-CNAB-dependencies.md).

## Parameter Sources

A custom bundle extension, `io.cnab.parameterSources`, MAY be defined that provides additional guidance for how a runtime MAY determine the default value of a [parameter][parameter]. The parameter source extension does not cover how the runtime should set the value of a parameter because each implemenation may have different information to take into account when determining the final value of a parameter.

In the example below, the [output][output] `tfstate` is initially generated during the install action, and the required parameter `tfstate` is used by the upgrade and uninstall actions. The parameter source specifies that the `tfstate` output can be used to set the `tfstate` parameter. This enables a bundle to pass data, in this case state, between actions.

```json
{
  "custom": {
    "io.cnab.parameterSources": {
      "tfstate": {
        "output": {
          "name": "tfstate"
        }
      }
    }
  },
  "outputs": {
    "tfstate": {
      "applyTo": [ "install", "upgrade", "uninstall" ],
      "definition": "tfstate",
      "path": "/cnab/app/outputs/tfstate"
    }
  },
  "parameters": {
    "tfstate": {
      "applyTo": [ "upgrade", "uninstall" ],
      "definition": "tfstate",
      "required": true
    }
  }
}
```

- `io.cnab.parameterSources`: Defines the parameter sources extension. A map of key/value pairs of a parameter name and a source for the parameter value.
  - `<parameterName>`: The name of the destination parameter.
    - `output`: Specifies that the source of the parameter is an output.
      - `name`: The name of the source output.

[parameter]: 101-bundle-json.md#parameters
[output]: 101-bundle-json.md#outputs