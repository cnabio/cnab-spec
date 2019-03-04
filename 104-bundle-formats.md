# Bundle Formats

This section of the specification addresses how bundles are to be represented. The CNAB Core specification defines two representations ("thick" and "thin"), and any CNAB compliant bundle MUST be representable in both formats. CNAB runtimes MUST support thin bundles, and SHOULD support thick bundles. Various other CNAB tools MAY support only one or the other. For example, a bundle builder could support generating thick bundles, but not thin bundles, and yet still meet the requirements of the specification.

Thick bundles do impose a slightly heavier burden on the runtime, and thus the specification allows that some runtimes might have pragmatic reasons for not supporting a thick bundle. Such runtimes MUST produce an error when given a thick bundle.

## Thick and Thin Bundles

The thick and thin bundle formats refer to how much information must be transmitted in a package. Due to the nature of cloud technologies, it is sometimes possible (and desirable) to move a small artifact that references other external artifacts. In other cases, it is desirable to have one large self-contained artifact that depends on no external artifacts. CNAB represents the former as a _thin bundle_ and the later as a _thick bundle_.

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

The internal layout of the TAR should be as follows:

```
├── artifacts
│   ├── example.com-technosophos-invocation-image-1.2.3.tar
│   └── example.com-technosophos-microservice-latest.tar
└── bundle.json
```

The `bundle.json` MUST always be located at the root of the archive.

All images MUST be located inside of the `artifacts` directory, with no subdirectories.

While CNAB implementations MAY create other directories at the root of the archive, they MUST NOT create subdirectories in `artifacts`.

### Transmitting Thick Bundles

Thick bundles are represented as one large gzipped TAR file. As such, a thick bundle MAY be transmitted as a single unit. Some transmission implementations MAY decompose thick bundles for the purpose of efficiently transmitting them. In such cases, the transmission MUST NOT alter the digest values of the artifacts, and MUST NOT alter the content of the `bundle.json`. Furthermore, they MUST use the same compression level when recomposing the bundle. To this end, the digest calculated on the source bundle file MUST be valid when the bundle is recovered.