# The bundle.json File

This section describes the format and function of the `bundle.json` document.

A `bundle.json` is broken down into the following categories of information:

- The schema version of the bundle
- The top-level package information (`name` and `version`)
    - name: The bundle name
    - version: Semantic version of the bundle
    - description: Short description of the bundle
- Information on the invocation images, as an array
- A list of images included with this bundle, as an array
- A specification of which parameters may be overridden, and how those are to be validated
- A list of credentials (name and desired location) that the application needs

There are two formats for a bundle (thin and thick formats). The primary way in which the `bundle.json` file differs is the presence or absence of information in a thick bundle that helps it validate the contents of an image. In a thick bundle, `mediaType` and `size` attributes may assiste the reconstitution of images from the thick format to a runtime format.

For the rest of the documentation, by default we'll be refererring bundles using the "thin" type, but when "thick" bundles become relevant we'll make note that it's a "thick" bundle type.

The following is an example of a `bundle.json` for a bundled distributed as a _thin_ bundle:

```json
{
    "schemaVersion": "v1",
    "name": "helloworld",
    "version": "0.1.2",
    "description": "An example 'thin' helloworld Cloud-Native Application Bundle",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:0.1.0",
            "digest": "sha256:aaaaaaa..."
        }
    ],
    "images": [
        {
            "name": "image1",
            "digest": "sha256:aaaaaaaaaaaa...",
            "uri": "urn:image1uri",
            "refs": [
                {
                    "path": "image1path",
                    "field": "image.1.field"
                }
            ]
        }
    ],
    "parameters": {
        "backend_port" : {
            "type" : "int",
            "defaultValue": 80,
            "minValue": 10,
            "maxValue": 10240,
            "metadata": {
               "description": "The port that the backend will listen on"
            }
        }
    },
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

And here is how a "thick" bundle looks. Notice how the `invocationImage` and `images` fields reference the underlying docker image manifest (`application/vnd.docker.distribution.manifest.v2+json`), which in turn references the underlying images:

```json
{
    "schemaVersion": 1,
    "name": "helloworld",
    "version": "1.0.0",
    "description": "An example 'thick' helloworld Cloud-Native Application Bundle",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:1.2.3",
            "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
            "size": 1337,
            "digest": "sha256:aaaaaaaaaaaa...",
            "platform": {
                "architecture": "amd64",
                "os": "linux"
            }
        }
    ],
    "images": [
        {
            "mediaType": "application/vnd.docker.distribution.manifest.v2+json",
            "size": 1337,
            "digest": "sha256:bbbbbbbbbbbb...",
            "platform": {
                "architecture": "amd64",
                "os": "linux"
            }
        }
    ],
    "parameters": {
        "backend_port" : {
            "type" : "int",
            "defaultValue": 80,
            "minValue": 10,
            "maxValue": 10240,
            "metadata": {
               "description": "The port that the backend will listen on"
            }
        }
    },
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

## Name and Version: Identifying Metadata

The `name` and `version` fields are used to identify the CNAB bundle. Both fields are required.

- Name should be human-readable (TODO: Define allowed format)
- Version MUST be a [SemVer2](https://semver.org) string

Fields that do not match this specification _should_ cause failures.

## Informational Metadata

The following fields are informational pieces of metadata designed to convey additional information about a bundle, but not to be used as identification for a bundle:

- `description`: A short description of a bundle
- `keywords`: A list of keywords
- `maintainers`: A list of maintainers, where each maintainer may have the following:
  - `name`: Maintainer name
  - `email`: Matainer's email
  - `url`: URL to relevant maintainer information

*TODO:* `bundle.json` probably requires a few more top-level fields, such as something about who published it, and something about the license, as well as a bundle api version. A decision on this is deferred until after the PoC

## Invocation Images

The `invocationImages` section describes the images that contains the bootstrapping for the image. The appropriate invocation
image is selected using the current driver.

```json
"invocationImages": [
    {
        "imageType": "docker",
        "image": "technosophos/helloworld:0.1.0",
        "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685"
    }
]
```

The `imageType` field is required, and must describe the format of the image. The list of formats is open-ended, but any CNAB-compliant system MUST implement `docker` and `oci`.

> Duffle, the reference implementation of a CNAB installer, introduces a layer of user-customizable drivers which are type-aware. Images may be delegated to drivers for installation.

The `image` field must give a path-like or URI-like representation of the location of the image. The expectation is that an installer should be able to locate the image (given the image type) without additional information.

The `digest` field _must_ contain a digest, in [OCI format](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests), to be used to compute the integrity of the image. The calculation of how the image matches the digest is dependent upon image type. (OCI, for example, uses a Merkle tree while VM images are checksums.)

The following optional fields may be attached to an invocation image:

- `size`: The image size in bytes. Implementations _should_ verify this when a bundle is packaged as a _thick_ bundle, and _may_ verify it when the image is part of a thin bundle.
- `platform`: The target platform, as an object with two fields:
  - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
  - `os`: The operating system of the image
- `mediaType`: The media type of the image

## The Image List

The `bundle.json` maps image metadata (name, origin, tag) to placeholders within the bundle. This allows images to be renamed, relabeled, or replaced during the CNAB bundle build operation. It also specifies the parameters that may be overridden in this image, giving tooling the ability to expose configuration options.

The following illustrates an `images` section:

```json
{ ​
"images": [​
        { ​
            "name": "frontend",​
            "uri": "gabrtv.azurecr.io/gabrtv/vote-frontend:a5ff67...",​
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "refs": [​
                {​
                    "path": "./charts/azure-voting-app/values.yaml",​
                    "field": "AzureVoteFront.deployment.image"​
                }​
            ]​
        },​
        { ​
            "name": "backend",​
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "uri": "gabrtv.azurecr.io/gabrtv/vote-backend:a5ff67...",​
            "refs": [​
                {​
                    "path": "./charts/azure-voting-app/values.yaml",​
                    "field": "AzureVoteBack.deployment.image"​
                }​
            ]​
        }​
    ]
}
```

Fields:

- images: The list of dependent images
  - name: The image name
  - URI: The image reference (REGISTRY/NAME:TAG). Note that _should_ be a CAS SHA, not a version tag as in the example above.
  - digest: the cryptographically hashed digest of the image. The implementation of hash verification depends on image type.
  - refs: An array listing the locations which refer to this image, and whose values should be replaced by the value specified in URI. Each entry contains the following properties:
    - path: the path of the file where the value should be replaced
    - field:a selector specifying a location (or locations) within that file where the value should be replaced
  - `size`: The image size in bytes
  - `platform`: The target platform, as an object with two fields:
    - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
    - `os`: The operating system of the image
  - `mediaType`: The media type of the image

Substitutions _must_ be supported for the following formats:

- JSON
- YAML
- XML

### Field Selectors

*TODO:* We have multiple competing standards in this space, and those that are popular for JSON are not the same as those popular for XML. This portion is thus not complete.

For fields, the selectors are based on the _de facto_ format used in tools like `jq`, which is a subset of the [CSS selector](https://www.w3.org/TR/selectors-3/) path. Examples:

- `foo.bar.baz` is interpreted as "find element baz whose parent is bar and whose grandparent is foo".
- `#baz` in XML is "the element whose ID attribute is set to "baz"". It is a no-op in YAML and JSON.
- TODO: Will we need to support attribute selectors?

TODO: How do we specify multiple replacements within a single file?

TODO: How do we specify URI is a VM image (or Jar or other) instead of a Docker-style image? Or do we? And if not, why not?

## Parameters

The `parameters` section of the `bundle.json` defines which parameters a user (person installing a CNAB bundle) may _override_. Parameter specifications are flat (not tree-like), consisting of name/value pairs. The name is fixed, but the value may be overridden by the user. The parameter definition includes a specification on how to constrain the values submitted by the user.

> The parameters definition is a subset of the ARM template laguage.

```json
"parameters": {
    "backend_port" : {
        "type" : "int",
        "defaultValue": 80,
        "minValue": 10,
        "maxValue": 10240,
        "metadata": {
            "description": "The port that the backend will listen on"
        }
    }
}
```

- parameters: name/value pairs describing a user-overridable parameter
  - `<name>`: The name of the parameter. In the example above, this is `backend_port`. This
    is mapped to a value definition, which contains the following fields:
    - type: one of string, int, boolean
    - defaultValue: The default value (optional)
    - allowedValues: an array of allowed values (optional)
    - minValue: Minimum value (for ints) (optional)
    - maxValue: Maximum value (for ints) (optional)
    - minLength: Minimum number of characters allowed in the field (for strings) (optional)
    - maxLength: Maximum number of characters allowed in the field (for strings) (optional)
    - metadata: Holds fields that are not used in validation
      - description: A user-friendly description of the parameter

Parameter names (the keys in `parameters`) ought to conform to the [Open Group Base Specification Issue 6, Section 8.1, paragraph 4](http://pubs.opengroup.org/onlinepubs/000095399/basedefs/xbd_chap08.html) definition of environment variable names with one exception: parameter names _may_ begin with a digit (approximately `[A-Z0-9_]+`).

For convenience, if lowercase characters are used in parameter names, they will be automatically capitalized. This effectively makes parameter names case-insensitive.

> The term _parameters_ indicates the present specification of what can be provided to a bundle. The term _values_ is frequently used to indicate the user-supplied values which are tested against the parameter definitions.

### Format of Parameter Specification

The structure of a parameters section looks like this:

```json
"parameters": {
    "<parameter-name>" : {
        "type" : "<type-of-parameter-value>",
        "defaultValue": "<default-value-of-parameter>",
        "allowedValues": [ "<array-of-allowed-values>" ],
        "minValue": <minimum-value-for-int>,
        "maxValue": <maximum-value-for-int>,
        "minLength": <minimum-length-for-string-or-array>,
        "maxLength": <maximum-length-for-string-or-array-parameters>,
        "metadata": {
            "description": "<description-of-the parameter>"
        }
    }
}
```

See [The Bundle Runtime](103-bundle-runtime.md) for details of how parameters are injected into the invocation image.

## Credentials

A `bundle.json` may optionally contain a section that describes which credentials the bundle expects to have access to in the invocation image. This information is provided so that users can be informed about the credentials that must be provided.

```json
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
```

- The `credentials` container is a map of human-friendly credential names to a description of where the invocation image expects to find them.
  - The name key must be human-readable
    - `path` describes the _absolute path within the invocation image_ where the invocation image expects to find the credential
    - `env` contains _the name of an environment variable_ that the invocation image expects to have available when executing the CNAB `run` tool (covered in the next section).

When _both a path and an env_ are specified, _only one is required_ (properties are disjunctive). To require two presentations of the same material, two separate entries must be made.

Next section: [The invocation image definition](102-invocation-image.md)
