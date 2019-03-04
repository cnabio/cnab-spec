# The bundle.json File

This section describes the format and function of the `bundle.json` document.

The `bundle.json` file is a representation of bundle metadata. It MUST be represented as [Canonical JSON](http://wiki.laptop.org/go/Canonical_JSON). While Canonical JSON is parseable by any JSON parser, its serialized form is consistent. This is a necessity when comparing two textual representations of the same data (such as when hashing).

> JSON data in this document has been formatted for readability using line breaks and indentation. This is not Canonical JSON. These examples were generated using the UNIX command `canonjson SOURCE.json | jq .` Where appropriate, the Canonical JSON text will also be provided. Small snippets of JSON may be shown in the order in which the fields are described (for clarity) rather than in Canonical JSON order.

A `bundle.json` is broken down into the following categories of information:

- The schema version of the bundle, as a string with a `v` prefix. This schema is to be referenced as `v1` or `v1.0.0-WD`
- The top-level package information (`name` and `version`)
  - name: The bundle name, including namespacing. The namespace can have one or more elements separated by a dot (e.g. `acme.tunnels.wordpress`). The left most element of the namespace is the most general moving toward more specific elements on the right.
  - version: Semantic version of the bundle
  - description: Short description of the bundle
- Information on the invocation images, as an array
- A map of images included with this bundle, as a `component name` to `image definition` map
- A specification of which parameters MAY be overridden, and how those are to be validated
- A list of credentials (name and desired location) that the application needs
- An OPTIONAL description of custom actions that this bundle implements

The `bundle.json` is also known as a _thin bundle_. Bundles come in two formats: thick and thin. Read more about thick and thin bundles in the [bundle formats section](104-bundle-formats.md).

For the rest of the documentation, by default we'll be referring to bundles using the "thin" type, but when "thick" bundles become relevant we'll make note that it's a "thick" bundle type.

The following is an example of a `bundle.json` for a bundled distributed as a _thin_ bundle:

```json
{
  "credentials": {
    "hostkey": {
      "env": "HOST_KEY",
      "path": "/etc/hostkey.txt"
    }
  },
  "custom": {
    "com.example.backup-preferences": {
      "frequency": "daily"
    },
    "com.example.duffle-bag": {
      "icon": "https://example.com/icon.png",
      "iconType": "PNG"
    }
  },
  "description": "An example 'thin' helloworld Cloud-Native Application Bundle",
  "images": {
    "my-microservice": {
      "description": "my microservice",
      "digest": "sha256:aaaaaaaaaaaa...",
      "image": "technosophos/microservice:1.2.3",
      "refs": [
        {
          "field": "image.1.field",
          "path": "image1path"
        }
      ]
    }
  },
  "invocationImages": [
    {
      "digest": "sha256:aaaaaaa...",
      "image": "technosophos/helloworld:0.1.0",
      "imageType": "docker"
    }
  ],
  "maintainers": [
    {
      "email": "matt.butcher@microsoft.com",
      "name": "Matt Butcher",
      "url": "https://example.com"
    }
  ],
  "name": "helloworld",
  "parameters": {
    "backend_port": {
      "defaultValue": 80,
      "destination": {
        "env": "BACKEND_PORT"
      },
      "maxValue": 10240,
      "metadata": {
        "description": "The port that the back-end will listen on"
      },
      "minValue": 10,
      "type": "int"
    }
  },
  "schemaVersion": "v1.0.0-WD",
  "version": "0.1.2"
}
```
Source: [101.01-bundle.json](examples/101.01-bundle.json)

The canonical JSON version of the above is:

```json
{"credentials":{"hostkey":{"env":"HOST_KEY","path":"/etc/hostkey.txt"}},"custom":{"com.example.backup-preferences":{"frequency":"daily"},"com.example.duffle-bag":{"icon":"https://example.com/icon.png","iconType":"PNG"}},"description":"An example 'thin' helloworld Cloud-Native Application Bundle","images":{"my-microservice":{"description":"my microservice","digest":"sha256:aaaaaaaaaaaa...","image":"technosophos/microservice:1.2.3","refs":[{"field":"image.1.field","path":"image1path"}]}},"invocationImages":[{"digest":"sha256:aaaaaaa...","image":"technosophos/helloworld:0.1.0","imageType":"docker"}],"maintainers":[{"email":"matt.butcher@microsoft.com","name":"Matt Butcher","url":"https://example.com"}],"name":"helloworld","parameters":{"backend_port":{"defaultValue":80,"destination":{"env":"BACKEND_PORT"},"maxValue":10240,"metadata":{"description":"The port that the back-end will listen on"},"minValue":10,"type":"int"}},"schemaVersion":"v1.0.0-WD","version":"0.1.2"}
```

And here is how a "thick" bundle looks. Notice how the `invocationImage` and `images` fields reference the underlying docker image manifest (`application/vnd.docker.distribution.manifest.v2+json`), which in turn references the underlying images:

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
  },
  "description": "An example 'thick' helloworld Cloud-Native Application Bundle",
  "images": {
    "my-microservice": {
      "description": "helloworld microservice",
      "digest": "sha256:bbbbbbbbbbbb...",
      "image": "technosophos/helloworld:0.1.2",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      },
      "size": 1337
    }
  },
  "invocationImages": [
    {
      "digest": "sha256:aaaaaaaaaaaa...",
      "image": "technosophos/helloworld:1.2.3",
      "imageType": "docker",
      "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
      "platform": {
        "architecture": "amd64",
        "os": "linux"
      },
      "size": 1337
    }
  ],
  "name": "helloworld",
  "parameters": {
    "backend_port": {
      "defaultValue": 80,
      "destination": {
        "path": "/path/to/backend_port"
      },
      "maxValue": 10240,
      "metadata": {
        "description": "The port that the backend will listen on"
      },
      "minValue": 10,
      "type": "int"
    }
  },
  "schemaVersion": "v1.0.0-WD",
  "version": "1.0.0"
}
```
Source: [101.02-bundle.json](examples/101.02-bundle.json)

In descriptions below, fields marked REQUIRED MUST be present in any conformant bundle descriptor, while fields not thusly marked are considered optional.

## Dotted Names

Within this specification, certain user-supplied names SHOULD be expressed in the form of a _dotted name_, which is defined herein as a name composed by concatenating name components (Unicode letters and the dash (`-`) character) together, separated by dot (`.`) characters. Whitespace characters are not allowed within names, nor is there a defined escape sequence for the `.` character.

Dotted names are used to encourage name spacing and reduce the likelihood of naming collisions.

Dotted names SHOULD follow the reverse-DNS pattern used by [Java, C#, and other languages](https://en.wikipedia.org/wiki/Reverse_domain_name_notation).

CNAB tools MUST NOT treat these strings as domain names or domain components, as this specification allows characters that are not legal in DNS addresses.

Examples:

- `com.example.myapp.port`
- `org.example.action.status`
- `example.foo` (This format MAY be used, but the reverse DNS format is preferred)

A string MUST have at least one dot to be considered a _dotted name_.

## Schema Version

Every `bundle.json` MUST have a `schemaVersion` element.

The schema version must reference the version of the schema used for this document. It follows the [SemVer v2 specification](https://semver.org/). The following pre-release markers are recognized:

- `WD` indicates that the document references a working draft of the specification, and is not considered stable.
- `CR` indicates that the document references a candidate recommendation. Stability is not assured.

The current schema version is `v1.0.0-WD`, which is considered unstable. 

## Name and Version: Identifying Metadata

The `name` and `version` fields are used to identify the CNAB bundle. Both fields are REQUIRED.

- `name` MUST contain only characters from the Unicode graph characters
- `version` MUST be a [SemVer2](https://semver.org) string

Fields that do not match this specification SHOULD cause failures.

## Informational Metadata

The following fields are informational pieces of metadata designed to convey additional information about a bundle, but not to be used as identification for a bundle:

- `description`: A short description of a bundle (OPTIONAL)
- `keywords`: A list of keywords (OPTIONAL)
- `license`: The license under which this bundle is covered. This SHOULD use one of the [SPDX License Identifiers](https://spdx.org/licenses/) whenever possible (OPTIONAL)
- `maintainers`: An OPTIONAL list of maintainers, where each maintainer MAY optionally have the following:
  - `name`: Maintainer name
  - `email`: Maintainer's email
  - `url`: URL to relevant maintainer information

## Invocation Images

The `invocationImages` section describes the images that contains the bootstrapping for the image. The appropriate invocation image is selected using the current driver.

A CNAB bundle MUST have at least one invocation image.

```json
"invocationImages": [
    {
        "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
        "image": "technosophos/helloworld:0.1.0",
        "imageType": "docker"
    }
]
```

The `imageType` field MUST describe the format of the image. The list of formats is open-ended, but any CNAB-compliant system MUST implement `docker` and `oci`. The default is `oci`.

> [Duffle](https://github.com/deis/duffle), the reference implementation of a CNAB installer, introduces a layer of user-customizable drivers which are type-aware. Images MAY be delegated to drivers for installation.

The `image` field MUST give a path-like or URI-like representation of the location of the image. It is REQUIRED. The expectation is that an installer should be able to locate the image (given the image type) without additional information.

The `digest` field MUST contain a digest, in [OCI format](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests), to be used to compute the integrity of the image. The calculation of how the image matches the digest is dependent upon image type. (OCI, for example, uses a Merkle tree while VM images are checksums.) If this field is omitted, a runtime is not obligated to validate the image.

The following OPTIONAL fields MAY be attached to an invocation image:

- `size`: The image size in bytes. Implementations SHOULD verify this when a bundle is packaged as a _thick_ bundle, and MAY verify it when the image is part of a thin bundle.
- `platform`: The target platform, as an object with two fields:
  - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
  - `os`: The operating system of the image
- `mediaType`: The media type of the image

## The Image Map

The `bundle.json` maps image metadata (name, origin, tag) to placeholders within the bundle. This allows images to be renamed, relabeled, or replaced during the CNAB bundle build operation. It also specifies the parameters that MAY be overridden in this image, giving tooling the ability to expose configuration options.

The following illustrates an `images` section:

```json
{
"images": {
        "frontend": { 
            "description": "frontend component image",
            "imageType": "docker",
            "image": "gabrtv.azurecr.io/gabrtv/vote-frontend:a5ff67...",
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "refs": [
                {
                    "path": "./charts/azure-voting-app/values.yaml",
                    "field": "AzureVoteFront.deployment.image"
                }
            ]
        },
        "backend": {
            "description": "backend component image",
            "imageType": "docker",
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "image": "gabrtv.azurecr.io/gabrtv/vote-backend:a5ff67...",
            "refs": [
                {
                    "path": "./charts/azure-voting-app/values.yaml",
                    "field": "AzureVoteBack.deployment.image"
                }
            ]
        }
    }
}
```

Fields:

- images: The list of dependent images
  - `description`: The description field provides additional context of the purpose of the image. 
  - `imageType`: The `imageType` field MUST describe the format of the image. The list of formats is open-ended, but any CNAB-compliant system MUST implement `docker` and `oci`. The default is `oci`.
  - `image`: The REQUIRED `image` field provides a valid reference (REGISTRY/NAME:TAG) for the image. Note that SHOULD be a CAS SHA, not a version tag as in the example above.
  - `digest`: MUST contain a digest, in [OCI format](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests), to be used to compute the integrity of the image. The calculation of how the image matches the digest is dependent upon image type. (OCI, for example, uses a Merkle tree while VM images are checksums.)
  - `refs`: An array listing the locations which refer to this image, and whose values should be replaced by the value specified in URI. Each entry contains the following properties:
    - `path`: The path of the file where the value should be replaced
    - `field`: A selector specifying a location (or locations) within that file where the value should be replaced
    - `mediaType`: The media type of the file, which can be used to determine the file type. If unset, tooling may choose any strategy for detecting format
  - `size`: The image size in bytes
  - `platform`: The target platform, as an object with two fields:
    - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
    - `os`: The operating system of the image
  - `mediaType`: The media type of the image

Substitutions MUST be supported for the following formats:

- JSON
- YAML
- XML

In addition to these substitutions, the image map data is also made available to the invocation image at runtime. See [Image map](103-bundle-runtime.md#image-map) for more details.

### Field Selectors

*TODO:* We have multiple competing standards in this space, and those that are popular for JSON are not the same as those popular for XML. This portion is thus not complete.

For fields, the selectors are based on the _de facto_ format used in tools like `jq`, which is a subset of the [CSS selector](https://www.w3.org/TR/selectors-3/) path. Examples:

- `foo.bar.baz` is interpreted as "find element baz whose parent is bar and whose grandparent is foo".
- `#baz` in XML is "the element whose ID attribute is set to "baz"". It is a no-op in YAML and JSON.
- TODO: Will we need to support attribute selectors?

TODO: How do we specify multiple replacements within a single file?

TODO: How do we specify URI is a VM image (or Jar or other) instead of a Docker-style image? Or do we? And if not, why not?

## Parameters

The `parameters` section of the `bundle.json` defines which parameters a user (person installing a CNAB bundle) MAY configure on an invocation image. Parameters represent information about the application configuration, and may be persisted by the runtime.

Parameter specifications are flat (not tree-like), consisting of name/value pairs. The name is fixed, but the value MAY be overridden by the user. The parameter definition includes a specification on how to constrain the values submitted by the user.

```json
"parameters": {
    "backend_port" : {
        "type" : "int",
        "defaultValue": 80,
        "minValue": 10,
        "maxValue": 10240,
        "metadata": {
            "description": "The port that the backend will listen on"
        },
        "destination": {
            "env": "MY_ENV_VAR",
            "path": "/my/destination/path"
        },
        "apply-to": ["install", "action1", "action2"]
    }
}
```

- parameters: name/value pairs describing a user-overridable parameter
  - `<name>`: The name of the parameter. This is REQUIRED. In the example above, this is `backend_port`. This
    is mapped to a value definition, which contains the following fields:
    - type: "null", "boolean", "string", "number", or "integer". These types correspond to the main primitive types of JSON (with _integer_ as a special case). Objects and arrays are not supported, as the serialization of these is underdetermined. (REQUIRED)
    - required: if this is set to true, a value MUST be specified (OPTIONAL, not shown)
    - defaultValue: The default value (OPTIONAL)
    - enum (alias allowedValues): an array of allowed values (OPTIONAL)
    - minimum (alias minValue): Minimum value (for numeric and integer) (OPTIONAL)
    - exclusiveMinimum: Exclusive minimum value (for numeric and integer) (OPTIONAL) 
    - maximum (alias maxValue): Maximum value (for numeric and integer) (OPTIONAL)
    - exclusiveMaximum: Exclusive maximum value (for numeric and integer) (OPTIONAL)
    - minLength: Minimum number of characters allowed in the field (for strings) (OPTIONAL)
    - maxLength: Maximum number of characters allowed in the field (for strings) (OPTIONAL)
    - pattern: An ECMA 262 regular expression that must match the string (OPTIONAL)
    - metadata: Holds fields that are not used in validation (OPTIONAL)
      - description: A user-friendly description of the parameter
    - destination: Indicates where (in the invocation image) the parameter is to be written (REQUIRED)
      - env: The name of an environment variable
      - path: The fully qualified path to a file that will be created
    - apply-to: restricts this parameter to a given list of actions. If empty or missing, applies to all actions (OPTIONAL)

Parameter names (the keys in `parameters`) ought to conform to the [Open Group Base Specification Issue 6, Section 8.1, paragraph 4](http://pubs.opengroup.org/onlinepubs/000095399/basedefs/xbd_chap08.html) definition of environment variable names with one exception: parameter names MAY begin with a digit (approximately `[A-Z0-9_]+`).

Evaluation of the validation keywords should conform to the applicable sections of [Section 6 of the JSONSchema specification](https://tools.ietf.org/html/draft-handrews-json-schema-validation-01#section-6).

Some validation terms have aliased keywords for backward compatibility with ARM template validators. A CNAB runtime SHOULD implement all aliases. A CNAB builder SHOULD NOT prefer the aliases over the regular names. In this way, aliases will be phased out. The present syntax is based on a subset of JSONSchema's validators.

> The term _parameters_ indicates the present specification of what can be provided to a bundle. The term _values_ is frequently used to indicate the user-supplied values which are tested against the parameter definitions.

### Resolving Destinations

When resolving destinations, there are two ways a particular parameter value MAY be placed into the invocation image. Here is an example illustrating both:

```json
"parameters": {
    "greeting": {
        "defaultValue": "hello",
        "type": "string",
        "destination": {
            "env": "GREETING"
        },
        "metadata":{
            "description": "this will be in $GREETING"
        }
    },
    "config": {
        "defaultValue": "",
        "type": "string",
        "destination": {
            "path": "/opt/example-parameters/config.txt"
        },
        "metadata": {
            "description": "this will be located in a file"
        }
    }
}
```

For the (REQUIRED) `destination` field, at least one of `env` or `path` MUST be specified. (Both MAY be provided).

If `env` is set, the value of the parameter will be assigned to the given environment variable name. In the example in the previous section, `GREETING` is set to `hello`.

If `path` is set, the value of the parameter will be written into a file at the specified location on the invocation image's filesystem. This file name MUST NOT be present already on the invocation image's filesystem.

If both `env` and `path` are specified, implementations MUST put a copy of the data in each destination.

### Format of Parameter Specification

The structure of a parameters section looks like this:

```json
"parameters": {
    "<parameter-name>" : {
        "type" : "<type-of-parameter-value>",
        "required": true|false
        "defaultValue": "<default-value-of-parameter>",
        "allowedValues": [ "<array-of-allowed-values>" ],
        "minValue": <minimum-value-for-int>,
        "maxValue": <maximum-value-for-int>,
        "minimum": <minimum-length-for-string-or-array>,
        "maximum": <maximum-length-for-string-or-array-parameters>,
        "metadata": {
            "description": "<description-of-the parameter>"
        },
        "destination": {
            "env": "<name-of-env-var>",
            "path": "<fully-qualified-path>"
        },
        "apply-to": ["action1", "action2"]
    }
}
```

See [The Bundle Runtime](103-bundle-runtime.md) for details of how parameters are injected into the invocation image.

## Credentials

A `bundle.json` MAY contain a section that describes which credentials the bundle expects to have access to in the invocation image. This information is provided so that users can be informed about the credentials that are required by the invocation image in order for the invocation image to perform its tasks.

Credentials differ from parameters in intention: A _parameter_ represents application configuration. A _credential_ represents the _identity of the agent performing the action_. For example, the lifecycle of one bundle installation may be managed by several individuals. When interacting with that bundle, individuals may use their own credentials. However, the bundle's _parameters_ are assumed to be attached to the bundle itself, regardless of which individual is presently acting on that bundle.

What about parameters such as database passwords used by the application? Properly speaking, these are _parameters_ if the application uses them.

> Note that the CNAB specification does not mandate specific security implementations on storage of either parameters or credentials. CNAB Runtimes ought to consider that both parameters and credentials may contain secret data, and OUGHT to secure both parameter and credential data in ways appropriate to the underlying platform. It is recommended that credentials are not persisted between actions, and that parameters _are_ persisted between actions.

```json
"credentials": {
    "kubeconfig": {
        "path": "/home/.kube/config"
    },
    "image_token": {
        "env": "AZ_IMAGE_TOKEN"
    },
    "hostkey": {
        "path": "/etc/hostkey.txt",
        "env": "HOST_KEY"
    }
}
```

- The `credentials` container is a map of human-friendly credential names to a description of where the invocation image expects to find them.
  - The name key MUST be human-readable
    - `path` describes the _absolute path within the invocation image_ where the invocation image expects to find the credential
    - `env` contains _the name of an environment variable_ that the invocation image expects to have available when executing the CNAB `run` tool (covered in the next section).
    - `description` contains a user-friendly description of the credential.

When _both a path and an env_ are specified, _only one is REQUIRED_ (properties are disjunctive). To require two presentations of the same material, two separate entries MUST be made.

### Resolving Destination Conflicts in Environment Variables and Paths

Parameters and credentials may specify environment variables or paths as destinations.

- Implementations SHOULD produce an error if a parameter and a credential use the same destination environment variable
- Implementations SHOULD produce an error if a parameter and a credential use the same destination path
- Implementations MUST NOT override a credential value with a parameter value
- Implementations SHOULD NOT allow any parameter or credential to declare an environment variable with the prefix `CNAB_`
- Implementations MUST NOT allow a parameter or credential to override any environment variable with the `CNAB_` prefix
    - The `CNAB_` variables are defined in the [Bundle Runtime Description](./103-bundle-runtime.md) of this specification

## Custom Actions

Every implementation of a CNAB tool MUST support three built-in actions:

- `install`
- `upgrade`
- `uninstall`

Implementations MAY support user-defined additional actions as well. Such actions are exposed via the `bundle` definition file. An action definition contains an action _name_ followed by a description of that action:

```json
"actions": {
    "io.cnab.status":{
        "modifies": false,
        "description": "retrieves the status of an installation"
    },
    "io.cnab.migrate":{
        "modifies": false
    },
    "io.cnab.dry-run":{
        "modifies": false,
        "stateless": true,
        "description": "prints what install would do with the given parameters values"
    }
}
```

The action _name_ SHOULD use _dotted name_ syntax as defined earlier in this section.

The above declares three actions: `io.cnab.status`, `io.cnab.migrate` and `io.cnab.dry-run`. This means that the associated invocation images can handle requests for `io.cnab.status`, `io.cnab.migrate` and `io.cnab.dry-run` in addition to `install`, `upgrade`, and `uninstall`.

Each action is accompanied by a description, which contains the following fields:

- `modifies`: Indicates whether the given action will _modify resources_ in any way. If not provided, it will be assumed `false`.
- `description`: A human readable description of the action (OPTIONAL)
- `stateless`: The action does not act on a claim, and does not require credentials. This is useful for exposing dry-run actions, printing documentation, etc. (OPTIONAL)

The `modifies` field MUST be set to `true` if any resource that is managed by the bundle is changed in any way. The `modifies` field assists CNAB implementations in tracking history of changes over time. An implementation of CNAB MAY use this information when describing history or managing releases.

The `stateless` field indicates that the runtime bypass credentials validation (the user MAY or MAY NOT pass any credentials), and will not keep track of this action. Primary intent is to allow invocation bundles to provide dry-run or detailed help functionalities.
Stateless actions can be invoked on a non-existing installation.

An invocation image _ought_ to handle all custom targets declared in the `actions` section. An invocation image SHOULD NOT handle actions that are not included by the default list (`install`, `upgrade`, `uninstall`) and the custom actions section.

The built-in actions (`install`, `upgrade`, `uninstall`) MUST NOT appear in the `actions` section, and an implementation MUST NOT allow custom actions named `install`, `upgrade`, or `uninstall`.

Implementations that do not support custom actions MUST NOT emit errors (either runtime or validation) if a bundle defines custom actions. That is, even if an implementation cannot execute custom actions, it MUST NOT fail to operate on bundles that declare custom actions.

## Custom Extensions

In many cases, the bundle descriptor is a sufficient artifact for delivering a CNAB bundle, since invocation images and other images may be retrieved from registries and repositories. However, it is important to provide an extension mechanism. A _custom extension_ is a named collection of auxiliary data whose meaning is defined outside of this specification.

Tools MAY define and declare additional fields inside of the `custom` section. Tools MUST NOT define additional fields anywhere else in the bundle descriptor. Implementations MAY produce an error or MAY ignore additional fields outside of the extension, but MUST be consistent in either ignoring or producing an error. However, implementations SHOULD preserve the data inside of the `custom` section even when that information is not understood by the implementation.

The `custom` object is used as follows:

```json
{
    "custom": {
        "com.example.duffle-bag": {
            "icon": "https://example.com/icon.png",
            "iconType": "PNG"
        },
        "com.example.backup-preferences": {
            "frequency": "daily"
        }
    }
}
```

The format is:

```json
{
    "custom": {
        "EXTENSION NAME": "ARBITRARY JSON DATA"
    }
}
```

The fields are defined as follows:

- `custom` defines the wrapper object for extensions
  - EXTENSION NAME: a unique name for an extension. Names SHOULD follow the dotted name format described earlier in this section.
  - The value of the extension must be valid JSON, but is otherwise undefined.

The usage of extensions is undefined. However, bundles SHOULD be installable by runtimes that do not understand the extensions.


Next section: [The invocation image definition](102-invocation-image.md)
