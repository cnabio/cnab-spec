---
title: The bundle.json File
weight: 101
---

# The bundle.json File

This section describes the format and function of the `bundle.json` document.

The `bundle.json` file is a representation of bundle metadata. It MUST be represented as [Canonical JSON](http://wiki.laptop.org/go/Canonical_JSON). While Canonical JSON is parseable by any JSON parser, its serialized form is consistent. This is a necessity when comparing two textual representations of the same data (such as when hashing).

> JSON data in this document has been formatted for readability using line breaks and indentation. This is not Canonical JSON. These examples were generated using the UNIX command `canonjson SOURCE.json | jq .` Where appropriate, the Canonical JSON text will also be provided. Small snippets of JSON may be shown in the order in which the fields are described (for clarity) rather than in Canonical JSON order.

A `bundle.json` is broken down into the following categories of information:

- The schema version of the bundle, as a string with a `v` prefix. This schema is to be referenced as `v1` or `v1.0.0-WD`
- The top-level package information (`name` and `version`)
  - `name`: The bundle name, including namespacing. The namespace can have one or more elements separated by a dot (e.g. `acme.tunnels.wordpress`). The left most element of the namespace is the most general moving toward more specific elements on the right.
  - `version`: Semantic version of the bundle
  - `description`: Short description of the bundle
- Information on the invocation images, as an array
- A map of images included with this bundle, as a `component name` to `image definition` map
- A specification of which parameters MAY be overridden, and how those are to be validated
- A list of credentials (name and desired location) that the application needs
- An OPTIONAL description of custom actions that this bundle implements
- A list of outputs (name, type and location) that the application produces

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
      "contentDigest": "sha256:aaaaaaaaaaaa...",
      "image": "technosophos/microservice:1.2.3"
    }
  },
  "invocationImages": [
    {
      "contentDigest": "sha256:aaaaaaa...",
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
  "outputs": {
    "clientCert" : {
      "contentEncoding" : "base64",
      "contentMediaType" : "application/x-x509-user-cert",
      "path" : "/cnab/app/outputs/clientCert",
      "sensitive" : true,
      "type" : "file"
    },
    "hostName" : {
        "applyTo" : ["install"],
        "description" : "the hostname produced installing the bundle",
        "path" : "/cnab/app/outputs/hostname",
        "type" : "string"
     },
    "port" : {
      "path" : "/cnab/app/outputs/port",
      "type" : "integer"
    }
  },
  "parameters": {
    "fields": {
      "backend_port": {
        "default": 80,
        "destination": {
          "env": "BACKEND_PORT"
        },
        "maximum": 10240,
        "description": "The port that the back-end will listen on",
        "minimum": 10,
        "type": "integer"
      }
    }
  },
  "schemaVersion": "v1.0.0-WD",
  "version": "0.1.2"
}
```

Source: [101.01-bundle.json](examples/101.01-bundle.json)

The canonical JSON version of the above is:

```json
{"credentials":{"hostkey":{"env":"HOST_KEY","path":"/etc/hostkey.txt"}},"custom":{"com.example.backup-preferences":{"frequency":"daily"},"com.example.duffle-bag":{"icon":"https://example.com/icon.png","iconType":"PNG"}},"description":"An example 'thin' helloworld Cloud-Native Application Bundle","images":{"my-microservice":{"description":"my microservice","contentDigest":"sha256:aaaaaaaaaaaa...","image":"technosophos/microservice:1.2.3"}},"invocationImages":[{"contentDigest":"sha256:aaaaaaa...","image":"technosophos/helloworld:0.1.0","imageType":"docker"}],"maintainers":[{"email":"matt.butcher@microsoft.com","name":"Matt Butcher","url":"https://example.com"}],"name":"helloworld","outputs":{"clientCert":{"contentEncoding":"base64","contentMediaType":"application/x-x509-user-cert","path":"/cnab/app/outputs/clientCert","sensitive":true,"type":"string"},"hostName":{"applyTo":["install"],"description":"the hostname produced installing the bundle","path":"/cnab/app/outputs/hostname","type":"string"},"port":{"path":"/cnab/app/outputs/port","type":"integer"}},"parameters":{"fields":{"backend_port":{"default":80,"description":"The port that the back-end will listen on","destination":{"env":"BACKEND_PORT"},"maximum":10240,"minimum":10,"type":"integer"}}},"schemaVersion":"v1.0.0-WD","version":"0.1.2"}
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
      "contentDigest": "sha256:bbbbbbbbbbbb...",
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
      "contentDigest": "sha256:aaaaaaaaaaaa...",
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
  "outputs": {
    "clientCert" : {
      "contentEncoding" : "base64",
      "contentMediaType" : "application/x-x509-user-cert",
      "path" : "/cnab/app/outputs/clientCert",
      "sensitive" : true,
      "type" : "file"
    },
    "hostName" : {
        "applyTo" : ["install"],
        "description" : "the hostname produced installing the bundle",
        "path" : "/cnab/app/outputs/hostname",
        "type" : "string"
     },
    "port" : {
      "path" : "/cnab/app/outputs/port",
      "type" : "integer"
    }
  },
  "parameters": {
    "fields" : {
      "backend_port": {
        "default": 80,
        "description": "The port that the backend will listen on",
        "destination": {
          "path": "/path/to/backend_port"
        },
        "maximum": 10240,
        "minimum": 10,
        "type": "integer"
      }
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

The `invocationImages` section describes the images that are responsible for bootstrapping the installation. The appropriate invocation image is selected by the CNAB runtime, typically by considering the runtime requirements of the bundle. For example, both a Windows and a Linux version of the invocation image may be included in the list. It is up to the CNAB runtime to determine which one to use. If no sufficient image is found, the CNAB runtime MUST emit an error and stop processing.

A CNAB bundle MUST have at least one invocation image.

```json
"invocationImages": [
    {
        "contentDigest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
        "image": "technosophos/helloworld:0.1.0",
        "imageType": "docker"
    }
]
```

The `imageType` field MUST describe the format of the image. The list of formats is open-ended, but any CNAB-compliant system MUST implement `docker` and `oci`. The default is `oci`.

The `image` field MUST give a path-like or URI-like representation of the location of the image. It is REQUIRED. The expectation is that an installer should be able to locate the image (given the image type) without additional information.

The `contentDigest` field MUST contain a digest, in [OCI format](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests), to be used to compute the integrity of the image. The calculation of how the image matches the contentDigest is dependent upon image type. (OCI, for example, uses a Merkle tree while VM images are checksums.) If this field is omitted, a runtime is not obligated to validate the image. During bundle development, it may be ideal to omit the contentDigest field and/or skip validation. Once a bundle is ready to be transmitted as a thick or thin bundle, it must have a contentDigest field.

The following OPTIONAL fields MAY be attached to an invocation image:

- `size`: The image size in bytes. Implementations SHOULD verify this when a bundle is packaged as a _thick_ bundle, and MAY verify it when the image is part of a thin bundle.
- `platform`: The target platform, as an object with two fields:
  - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
  - `os`: The operating system of the image
- `mediaType`: The media type of the image
- `originalImage`: provides a path-like or URI-like representation of the original location of the image. If the `originalImage` field is omitted, as in the example above, its value defaults to that of the `image` field.

## The Image Map

The `bundle.json` maps image metadata (name, origin, tag) to placeholders within the bundle. This allows images to be renamed, relabeled, or replaced during the CNAB bundle build operation.

The following illustrates an `images` section:

```json
{
"images": {
        "frontend": { 
            "description": "frontend component image",
            "imageType": "docker",
            "image": "example.com/gabrtv/vote-frontend@sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "originalImage": "example.com/opendeis/vote-frontend@sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685",
            "contentDigest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685"
        },
        "backend": {
            "description": "backend component image",
            "imageType": "docker",
            "image": "example.com/gabrtv/vote-backend@sha256:bca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120686",
            "originalImage": "example.com/opendeis/vote-backend@sha256:bca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120686",
            "contentDigest": "sha256:bca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120686"
        }
    }
  }
}
```

Fields:

- `images`: The list of dependent images
  - `description`: The description field provides additional context of the purpose of the image.
  - `imageType`: The `imageType` field MUST describe the format of the image. The list of formats is open-ended, but any CNAB-compliant system MUST implement `docker` and `oci`. The default is `oci`.
  - `image`: The REQUIRED `image` field provides a valid reference for the image. Note that SHOULD be a CAS SHA, as in the example above, not a version tag.
  - `originalImage`: provides a valid reference for the image. Note that SHOULD be a CAS SHA, as in the example above, not a version tag. If the `originalImage` field is omitted, its value defaults to that of the `image` field.
  - `contentDigest`: MUST contain a digest of the contents of the image, in [OCI format](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests), to be used to compute the integrity of the image. The calculation of how the image matches the contentDigest is dependent upon image type. (OCI, for example, uses a Merkle tree while VM images use checksums.)
  - `size`: The image size in bytes
  - `platform`: The target platform, as an object with two fields:
    - `architecture`: The architecture of the image (`i386`, `amd64`, `arm32`...)
    - `os`: The operating system of the image
  - `mediaType`: The media type of the image

The image map data is made available to the invocation image at runtime. This allows invocation images to perform various substitutions during installation (for example, moving images to different storage mechanisms or registries, and renaming appropriately). See [Image map](103-bundle-runtime.md#image-map) for more details.

## Parameters

The `parameters` section of the `bundle.json` defines which parameters a user (person installing a CNAB bundle) MAY configure on an invocation image. Parameters represent information about the application configuration, and may be persisted by the runtime.

Parameter specifications consist of name/value pairs. The name is fixed, but the value MAY be overridden by the user. The parameter definition includes a specification of how to constrain the values submitted by the user.

```json
{
  "parameters": {
    "fields" : {
      "backend_port": {
        "applyTo": [
          "install",
          "action1",
          "action2"
        ],
        "default": 80,
        "description": "The port that the backend will listen on",
        "destination": {
          "env": "MY_ENV_VAR",
          "path": "/my/destination/path"
        },
        "maximum": 10240,
        "minimum": 10,
        "type": "integer"
      }
    },
    "required": ["backend_port"]
}
```

- `parameters`: A collection of parameter definitions and a list of those parameters that are required.
  - `fields`: name/value pairs describing a user-overridable parameter:
    - `<name>`: The name of the parameter. In the example above, this is `backend_port`. This
    is mapped to a value definition, which contains the following fields (REQUIRED):
      - `$comment`: Reserved for comments from bundle authors to readers or maintainers of the bundle. This MUST be a string (OPTIONAL)
      - `$id`: A URI for the schema resolved against the base URI of its parent schema. MUST be a uri-reference string in accordance with [RFC3986](https://tools.ietf.org/html/rfc3986) (OPTIONAL)
      - `$ref`: A URI reference used to resolve a schema located elsewhere. This MUST be a uri-reference string in accordance with [RFC3986](https://tools.ietf.org/html/rfc3986) (OPTIONAL)
      - `additionalItems`: Parameter validation requiring that any additional items included in a user-provided array must conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
      - `additionalProperties`: Parameter validation requiring that any additional properties in the user-provided object conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
      - `allOf`: Parameter validation requiring that the user-provided value match ALL of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
      - `anyOf`: Parameter validation requiring that the user-provided value match ANY of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
      - `applyTo`: restricts this parameter to a given list of actions. If empty or missing, applies to all actions (OPTIONAL)
      - `const`: Parameter validation requiring that the user-provided value matches exactly the specified const. MAY be of any type, including null. (OPTIONAL)
      - `contains`: Parameter validation requiring at least one item included in the user-provided array conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
      - `contentEncoding`: Indicates that the user-provided content should interpreted as binary data and decoded using the encoding named by this property. MUST be a string in accordance with [RFC2045, Sec 6.1](https://json-schema.org/latest/json-schema-validation.html#RFC2045). (OPTIONAL)
      - `contentMediaType`: MIME type indicating the media type of the user-provided content. MUST be a string in accordance with [RFC2046](https://json-schema.org/latest/json-schema-validation.html#RFC2046). (OPTIONAL)
      - `default`: A default JSON value associated with a particular schema. RECOMMENDED that a default value be valid against the associated schema. (OPTIONAL)
      - `definitions`: Provides a standardized location for bundle authors to inline re-usable JSON Schemas into a more general schema. MUST be an object where each named property contains a JSON schema. (OPTIONAL)
      - `dependencies`: Specifies rules that are evaluated if the parameter type is an object and contains a certain property. MUST be an object where each named dependency is either an array of unique strings or a JSON schema. (OPTIONAL)
      - `description`: Descriptive text for the field. Can be used to decorate a user interface. MUST be a string. (OPTIONAL)
      - `destination`: Indicates where (in the invocation image) the parameter is to be written (REQUIRED)
        - `env`: The name of an environment variable
        - `path`: The fully qualified path to a file that will be created. Specified path MUST NOT be a subpath of `/cnab/app/outputs`.
      - `else`: Parameter validation requiring that the user-provided value match the specified schema. Only matches if the user-provided value does NOT match the schema provided in the `if` property. MUST be a JSON schema. (OPTIONAL)
      - `enum`: Parameter validation requiring that the user-provided value is one of the specified items in the specified array. MUST be a non-empty array of unique elements that can be of any type. (OPTIONAL)
      - `examples`: Sample JSON values associated with a particular schema. MUST be an array. (OPTIONAL)
      - `exclusiveMaximum`: Parameter validation requiring that the user-provided number be less than the number specified. MUST be a number. (OPTIONAL)
      - `exclusiveMinimum`: Parameter validation requiring that the user-provided number be greater than the number specified. MUST be a number. (OPTIONAL)
      - `format`: Parameter validation requiring that the user-provided value adhere to the specified format. MUST be a string. (OPTIONAL)
      - `if`: Provides a method to conditionally validate user-provided values against a schema. MUST be a JSON schema. (OPTIONAL)
      - `items`: Parameter validation requiring the items included in a user-provided array must conform to the specified schema(s). MUST be either a JSON schema or an array of JSON schemas. (OPTIONAL)
      - `maxItems`: Parameter validation requiring the length of the user-provided array be less than or equal to the number specified. MUST be a non-negative number. (OPTIONAL)
      - `maxLength`: Parameter validation requiring that the length of the user-provided string be less than or equal to the number specified. MUST be a non-negative integer. (OPTIONAL)
      - `maxProperties`: Parameter validation requiring the number of properties included in the user-provided object be less than or equal to the specified number. MUST be a non-negative integer. (OPTIONAL)
      - `maximum`: Parameter validation requiring that the user-provided number be less than or equal to the number specified. MUST be a number. (OPTIONAL)
      - `minItems`: Parameter validation requiring the length of the user-provided array be greater than or equal to the number specified. MUST be a non-negative number. (OPTIONAL)
      - `minLength`: Parameter validation requiring that the length of the user-provided string be greater than or equal to the number specified. MUST be a non-negative integer. (OPTIONAL)
      - `minProperties`: Parameter validation requiring the number of properties included in the user-provided object be greater than or equal to the specified number. MUST be a non-negative integer. (OPTIONAL)
      - `minimum`: Parameter validation requiring that the user-provided number be greater than or equal to the number specified. MUST be a number. (OPTIONAL)
      - `multipleOf`: Parameter validation requiring that the user-provided number be wholly divisible by the number specified. MUST be a number strictly greater than zero. (OPTIONAL)
      - `not`: Parameter validation requiring that the user-provided value NOT match the specified schema. MUST be a JSON schema. (OPTIONAL)
      - `oneOf`: Parameter validation requiring that the user-provided value match ONE of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
      - `patternProperties`: The set of matching properties and schemas for their values included in an object type parameter. MUST be an object where each named property is a regular expression with a JSON schema as the value. (OPTIONAL)
      - `pattern`: Parameter validation requiring that the user-provided string match the regular expression specified. MUST be a string representation of a valid ECMA 262 regular expression. (OPTIONAL)
      - `properties`: The set of named properties and schemas for their values included in an object type parameter. MUST be an object where each named property contains a JSON schema. (OPTIONAL)
      - `propertyNames`: Parameter validation requiring that each property name in an object match the specified schema. MUST be a JSON schema. (OPTIONAL)
      - `readOnly`: Indicates that the value of the parameter cannot be modified. MUST be a boolean. (OPTIONAL)
      - `required`: Parameter validation requiring the properties named in the user-provided object include the specified list of properties. MUST be an array of strings. (OPTIONAL)
      - `then`: Parameter validation requiring that the user-provided value match the specified schema. Only matches if the user-provided value matches the schema provided in the `if` property. MUST be a JSON schema. (OPTIONAL)
      - `title`: Short, human-readable descriptive name for the field. Can be used to decorate a user interface. MUST be a string. (OPTIONAL)
      - `type`: Parameter validation requiring that the user-provided value is either a "null", "boolean", "object", "array", "number", "string", or "integer". MUST be a string or an array of strings with unique elements. (OPTIONAL)
      - `uniqueItems`: Parameter validation requiring the items included in the user-provided array be unique. MUST be a boolean. (OPTIONAL)
      - `writeOnly`: Indicates that the value of the parameter is sensitive and cannot be viewed once set or updated. MUST be a boolean. (OPTIONAL)
  - `required`: A list of required parameters. MUST be an array of strings.(OPTIONAL)

For more information on the supported parameter properties, visit the [JSON Schema documentation](https://json-schema.org/)

Parameter names (the keys in `fields`) ought to conform to the [Open Group Base Specification Issue 6, Section 8.1, paragraph 4](http://pubs.opengroup.org/onlinepubs/000095399/basedefs/xbd_chap08.html) definition of environment variable names with one exception: parameter names MAY begin with a digit (approximately `[A-Z0-9_]+`).

Evaluation of the validation keywords should conform to the applicable sections of [Section 6 of the JSONSchema specification](https://tools.ietf.org/html/draft-handrews-json-schema-validation-01#section-6).

Some validation terms have aliased keywords for backward compatibility with ARM template validators. A CNAB runtime SHOULD implement all aliases. A CNAB builder SHOULD NOT prefer the aliases over the regular names. In this way, aliases will be phased out. The present syntax is based on a subset of JSONSchema's validators.

> The term _parameters_ indicates the present specification of what can be provided to a bundle. The term _values_ is frequently used to indicate the user-supplied values which are tested against the parameter definitions.

### Format of Parameter Specification

The structure of a parameters section looks like the section below.

```
{
  "parameters": {
    "fields" : {
      "<parameter-name>": {
        "$comment": <string>,
        "$id": <uri-reference>,
        "$ref": <uri-reference>,
        "additionalItems": <json-schema>,
        "additionalProperties": <json-schema>,
        "allOf": [ <json-schema> ],
        "anyOf": [ <json-schema> ],
        "applyTo": [ <string> ],
        "const": <any-value>,
        "contains": <json-schema>,
        "contentEncoding": <string>,
        "contentMediaType": <string>,
        "default": <any-value>,
        "definitions": {
          "<definition-name>": <json-schema>
        },
        "dependencies": {
          "<first-property-name>": <json-schema>,
          "<second-property-name>": [ <string> ]
        },
        "description": <string>,
        "destination": {
          "env": <string>,
          "path": <string>
        },
        "else": <json-schema>,
        "enum": [ <any-value> ],
        "examples": [ <any-value> ],
        "exclusiveMaximum": <number>,
        "exclusiveMinimum": <number>,
        "format": <string>,
        "if": <json-schema>,
        "items": <json-schema> | [ <json-schema> ],
        "maxItems": <integer>,
        "maxLength": <integer>,
        "maxProperties": <integer>,
        "maximum": <number>,
        "minItems": <integer>,
        "minLength": <integer>,
        "minProperties": <integer>,
        "minimum": <integer>,
        "multipleOf": <number>,
        "not": <json-schema>,
        "oneOf": [ <json-schema> ],
        "pattern": <string>,
        "patternProperties": {
          "<regular-expression-for-property-name>": <json-schema>
        },
        "properties": {
          "<property-name>": <json-schema>
        },
        "propertyNames": <json-schema>,
        "readOnly": <boolean>,
        "required": [ <string> ],
        "then": <json-schema>,
        "title": <string>,
        "type": <string> | [ <string> ],
        "uniqueItems": <boolean>,
        "writeOnly": <boolean>
      }
    },
    "required": [ <string> ]
  }
}
```

See [The Bundle Runtime](103-bundle-runtime.md) for details of how parameters are injected into the invocation image.

### Examples

Below are a few more complicated examples that outline some of the possible parameter configurations that a bundle author can expose.

Check out the [JSON Schema specification](https://json-schema.org/) for more examples and further documentation.

```json
{
  "parameters": {
    "fields" : {
      "address": {
        "destination": {
          "path": "/tmp/address.adr"
        },
        "properties": {
          "country_name": {
            "type": "string"
          },
          "extended_street_address": {
            "type": "string"
          },
          "locality": {
            "type": "string"
          },
          "postal_code": {
            "type": "string"
          },
          "region": {
            "type": "string"
          },
          "street_address": {
            "type": "string"
          }
        },
        "required": [
          "country_name",
          "locality",
          "postal_code",
          "region",
          "street_address"
        ],
        "type": "object"
      },
      "email": {
        "format": "idn-email",
        "type": "string"
      },
      "greetings": {
        "default": [
          "Hello"
        ],
        "description": "a list of greetings",
        "destination": {
          "env": "GREETINGS"
        },
        "items": {
          "examples": [
            "Bonjour",
            "Aloha",
            "こんにちは"
          ],
          "type": "string"
        },
        "title": "Greetings for new users",
        "type": "array"
      },
      "image": {
        "contentEncoding": "base64",
        "contentMediaType": "image/jpeg",
        "destination": {
          "path": "/tmp/user.jpg"
        },
        "type": "string"
      }
    }
  },
  "required" : [ "address" ]
}
```

### Resolving Destinations

When resolving destinations, there are two ways a particular parameter value MAY be placed into the invocation image. Here is an example illustrating both:

```json
{
  "parameters": {
    "fields" : {
      "config": {
        "default": "",
        "description": "this will be located in a file",
        "destination": {
          "path": "/opt/example-parameters/config.txt"
        },
        "type": "string"
      },
      "greeting": {
        "default": "hello",
        "description": "this will be in $GREETING",
        "destination": {
          "env": "GREETING"
        },
        "type": "string"
      }
    }
  }
}
```

For the (REQUIRED) `destination` field, at least one of `env` or `path` MUST be specified. (Both MAY be provided).

If `env` is set, the value of the parameter will be assigned to the given environment variable name. In the example in the previous section, `GREETING` is set to `hello`.

If `path` is set, the value of the parameter will be written into a file at the specified location on the invocation image's filesystem. This file name MUST NOT be present already on the invocation image's filesystem.

If both `env` and `path` are specified, implementations MUST put a copy of the data in each destination.

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
        "required": true,
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
    - `path` describes the _absolute path within the invocation image_ where the invocation image expects to find the credential. Specified path MUST NOT be a subpath of `/cnab/app/outputs`.
    - `env` contains _the name of an environment variable_ that the invocation image expects to have available when executing the CNAB `run` tool (covered in the next section).
    - `description` contains a user-friendly description of the credential.
    - `required` indicates whether this credential MUST be supplied. By default, it is `false`, which means the credential is optional. When `true`, a runtime MUST fail if the credential is not provided.

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
  - `EXTENSION NAME`: a unique name for an extension. Names SHOULD follow the dotted name format described earlier in this section.
  - The value of the extension must be valid JSON, but is otherwise undefined.

The usage of extensions is undefined. However, bundles SHOULD be installable by runtimes that do not understand the extensions.

## Outputs

The `outputs` section of the `bundle.json` defines which outputs an application will produce during the course of executing a bundle. Outputs are expected to be written to one or more files on the file system of the invocation image. The location of this file MUST be provided in the output definition.

Output specifications are flat (not tree-like), consisting of name/value pairs. The output definition includes a destination the output will be written to, along with a type to help bundle users identify how to consume them.

```json
"outputs" : {
  "clientCert" : {
        "contentEncoding" : "base64",
        "contentMediaType" : "application/x-x509-user-cert",
        "type" : "string",
        "path" : "/cnab/app/outputs/clientCert",
        "applyTo": ["install", "action2"],
        "sensitive" : true,
    },
    "hostName" : {
        "type" : "string",
        "path" : "/cnab/app/outputs/hostname"
    },
    "port" : {
        "type" : "integer",
        "path" : "/cnab/app/outputs/port"
    }
}
```

- `outputs`: name/value pairs describing an application output
  - `<name>`: The name of the output. In the example above, this is `clientCert`, `hostName`, and `port`. This is mapped to a value definition, which contains the following fields (REQUIRED):
    - `$comment`: Reserved for comments from bundle authors to readers or maintainers of the bundle. This MUST be a string (OPTIONAL)
    - `$id`: A URI for the schema resolved against the base URI of its parent schema. MUST be a uri-reference string in accordance with [RFC3986](https://tools.ietf.org/html/rfc3986) (OPTIONAL)
    - `$ref`: A URI reference used to resolve a schema located elsewhere. This MUST be a uri-reference string in accordance with [RFC3986](https://tools.ietf.org/html/rfc3986) (OPTIONAL)
    - `additionalItems`: Output validation requiring that any additional items included in a resulting array must conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
    - `additionalProperties`: Output validation requiring that any additional properties in the resulting object conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
    - `allOf`: Output validation requiring that the resulting value match ALL of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
    - `anyOf`: Output validation requiring that the resulting value match ANY of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
    - `applyTo`: restricts this output to a given list of actions. If empty or missing, applies to all actions (OPTIONAL)
    - `const`: Output validation requiring that the resulting value matches exactly the specified const. MAY be of any type, including null. (OPTIONAL)
    - `contains`: Output validation requiring at least one item included in the resulting array conform to the specified schema. MUST be a JSON schema. (OPTIONAL)
    - `contentEncoding`: Indicates that the resulting content should interpreted as binary data and decoded using the encoding named by this property. MUST be a string in accordance with [RFC2045, Sec 6.1](https://json-schema.org/latest/json-schema-validation.html#RFC2045). (OPTIONAL)
    - `contentMediaType`: MIME type indicating the media type of the resulting content. MUST be a string in accordance with [RFC2046](https://json-schema.org/latest/json-schema-validation.html#RFC2046). (OPTIONAL)
    - `default`: A default JSON value associated with a particular schema. RECOMMENDED that a default value be valid against the associated schema. (OPTIONAL)
    - `definitions`: Provides a standardized location for bundle authors to inline re-usable JSON Schemas into a more general schema. MUST be an object where each named property contains a JSON schema. (OPTIONAL)
    - `dependencies`: Specifies rules that are evaluated if the output type is an object and contains a certain property. MUST be an object where each named dependency is either an array of unique strings or a JSON schema. (OPTIONAL)
    - `description`: Descriptive text for the field. Can be used to decorate a user interface. MUST be a string. (OPTIONAL)
    - `else`: Output validation requiring that the resulting value match the specified schema. Only matches if the resulting value does NOT match the schema provided in the `if` property. MUST be a JSON schema. (OPTIONAL)
    - `enum`: Output validation requiring that the resulting value is one of the specified items in the specified array. MUST be a non-empty array of unique elements that can be of any type. (OPTIONAL)
    - `examples`: Sample JSON values associated with a particular schema. MUST be an array. (OPTIONAL)
    - `exclusiveMaximum`: Output validation requiring that the resulting number be less than the number specified. MUST be a number. (OPTIONAL)
    - `exclusiveMinimum`: Output validation requiring that the resulting number be greater than the number specified. MUST be a number. (OPTIONAL)
    - `format`: Output validation requiring that the resulting value adhere to the specified format. MUST be a string. (OPTIONAL)
    - `if`: Provides a method to conditionally validate resulting values against a schema. MUST be a JSON schema. (OPTIONAL)
    - `items`: Output validation requiring the items included in a resulting array must conform to the specified schema(s). MUST be either a JSON schema or an array of JSON schemas. (OPTIONAL)
    - `maxItems`: Output validation requiring the length of the resulting array be less than or equal to the number specified. MUST be a non-negative number. (OPTIONAL)
    - `maxLength`: Output validation requiring that the length of the resulting string be less than or equal to the number specified. MUST be a non-negative integer. (OPTIONAL)
    - `maxProperties`: Output validation requiring the number of properties included in the resulting object be less than or equal to the specified number. MUST be a non-negative integer. (OPTIONAL)
    - `maximum`: Output validation requiring that the resulting number be less than or equal to the number specified. MUST be a number. (OPTIONAL)
    - `minItems`: Output validation requiring the length of the resulting array be greater than or equal to the number specified. MUST be a non-negative number. (OPTIONAL)
    - `minLength`: Output validation requiring that the length of the resulting string be greater than or equal to the number specified. MUST be a non-negative integer. (OPTIONAL)
    - `minProperties`: Output validation requiring the number of properties included in the resulting object be greater than or equal to the specified number. MUST be a non-negative integer. (OPTIONAL)
    - `minimum`: Output validation requiring that the resulting number be greater than or equal to the number specified. MUST be a number. (OPTIONAL)
    - `multipleOf`: Output validation requiring that the resulting number be wholly divisible by the number specified. MUST be a number strictly greater than zero. (OPTIONAL)
    - `not`: Output validation requiring that the resulting value NOT match the specified schema. MUST be a JSON schema. (OPTIONAL)
    - `oneOf`: Output validation requiring that the resulting value match ONE of the specified schemas. MUST be a non-empty array of JSON schemas. (OPTIONAL)
    - `path`: The fully qualified path to a file that will be created (REQUIRED). The path specified MUST be a _strict_ subpath of `/cnab/app/outputs` and MUST be distinct from the paths for all other outputs in this bundle.
    - `patternProperties`: The set of matching properties and schemas for their values included in an object type output. MUST be an object where each named property is a regular expression with a JSON schema as the value. (OPTIONAL)
    - `pattern`: Output validation requiring that the resulting string match the regular expression specified. MUST be a string representation of a valid ECMA 262 regular expression. (OPTIONAL)
    - `properties`: The set of named properties and schemas for their values included in an object type output. MUST be an object where each named property contains a JSON schema. (OPTIONAL)
    - `propertyNames`: Output validation requiring that each property name in an object match the specified schema. MUST be a JSON schema. (OPTIONAL)
    - `sensitive`: Indicates that a runtime should treat the output as a sensitive value. Defaults to false (OPTIONAL)
    - `then`: Output validation requiring that the resulting value match the specified schema. Only matches if the resulting value matches the schema provided in the `if` property. MUST be a JSON schema. (OPTIONAL)
    - `title`: Short, human-readable descriptive name for the field. Can be used to decorate a user interface. MUST be a string. (OPTIONAL)
    - `type`: Output validation requiring that the resulting value is either a "null", "boolean", "object", "array", "number", "string", or "integer". MUST be a string or an array of strings with unique elements. (OPTIONAL)
    - `uniqueItems`: Output validation requiring the items included in the resulting array be unique. MUST be a boolean. (OPTIONAL)

An invocation image should write outputs to a file specified by the `path` attribute for each output. A bundle runtime can then extract values from the specified path and present them to a user. A runtime can leverage appropriate [in-memory](https://docs.docker.com/v17.09/engine/admin/volumes/tmpfs/#choosing-the-tmpfs-or-mount-flag) volume mounted at the `path` location for storing these outputs.

For more information on the supported output properties, visit the [JSON Schema documentation](https://json-schema.org/)

A runtime may validate outputs. Evaluation of the validation keywords should conform to the applicable sections of [Section 6 of the JSONSchema specification](https://tools.ietf.org/html/draft-handrews-json-schema-validation-01#section-6).

Next section: [The invocation image definition](102-invocation-image.md)
