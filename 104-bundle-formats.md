---
title: Bundle Formats
weight: 104
---

# Bundle Formats

This section of the specification addresses how bundles are to be represented. The CNAB Core specification defines two representations (thick and thin), and any CNAB compliant bundle MUST be representable in both formats. CNAB runtimes MUST support thin bundles, and SHOULD support thick bundles. Various other CNAB tools MAY support only one or the other. For example, a bundle builder could support generating thick bundles, but not thin bundles, and yet still meet the requirements of the specification.

Thick bundles do impose a slightly heavier burden on the runtime, and thus the specification allows that some runtimes might have pragmatic reasons for not supporting a thick bundle. Such runtimes MUST produce an error when given a thick bundle.

## Thick and Thin Bundles

The thick and thin bundle formats refer to how much information must be transmitted in a package. Due to the nature of cloud technologies, it is sometimes possible (and desirable) to have a small artifact that references other external artifacts. In other cases, it is desirable to have one large self-contained artifact that depends on no external artifacts. CNAB represents the former as a thin bundle and the later as a thick bundle.

A thin bundle contains only one object: The bundle descriptor. Thus, the format for a thin bundle is a JSON file.

A thick bundle contains multiple objects:

- The bundle descriptor (`bundle.json`)
- One or more invocation images
- Zero or more images

As defined in this specification, objects of these three types are collected into a single archive file encoded as a gzipped tar archive.

## Differences in the Content of Bundle Descriptors

Bundle descriptors have slightly different fields for thick and thin bundles. These are described in [the Bundle Descriptor description](101-bundle-json.md).

## Formatting and Transmitting Thin Bundles

Thin bundles MUST be represented as Canonical JSON as specified in [the Bundle Descriptor description](101-bundle-json.md), and MUST conform to the schema provided in this specification.

Thin bundles MAY then be stored and transmitted as JSON data.

## Formatting and Transmitting Thick Bundles

Thick bundles MUST contain the bundle descriptor. In addition, a thick bundle MUST contain:

- _All_ of the images referenced in the `invocationImages` section of the bundle descriptor
- _All_ of the images referenced in the `images` section of the bundle descriptor

### File Format for Thick Bundles

A thick bundle SHOULD be encoded as a gzipped TAR. This specification is neutral as to what compression ratio is used.

The internal layout of the TAR SHOULD be as follows:
```
├── artifacts
│   └── layout
│       └── ...
└── bundle.json
```

The `bundle.json` MUST always be located at the root of the archive.

CNAB implementations MAY create other directories at the root of the archive.

All images MUST be located inside of the `artifacts/layout` directory, the contents of which SHOULD be an image layout conforming to the [OCI Image Layout Specification](https://github.com/opencontainers/image-spec/blob/master/image-layout.md), for example:
```
├── artifacts
│   └── layout
│       ├── blobs
│       │   └── sha256
│       │       ├── 3588d02542238316759cbf24502f4344ffcc8a60c803870022f335d1390c13b4
│       │       ├── 4b0bc1c4050b03c95ef2a8e36e25feac42fd31283e8c30b3ee5df6b043155d3c
│       │       └── 7968321274dc6b6171697c33df7815310468e694ac5be0ec03ff053bb135e768
│       ├── index.json
│       └── oci-layout
└── bundle.json
```


### Transmitting Thick Bundles

Thick bundles are represented as one large gzipped TAR file. As such, a thick bundle MAY be transmitted as a single unit. Some transmission implementations MAY decompose thick bundles for the purpose of efficiently transmitting them. In such cases, the transmission MUST NOT alter the contentDigest values of the artifacts, and MUST NOT alter the content of the `bundle.json`. Furthermore, they MUST use the same compression level when recomposing the bundle. To this end, the digest calculated on the source bundle file MUST be valid when the bundle is recovered.
