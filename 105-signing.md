# Calculating Digests in Bundle Descriptors

This section describes how image digests are calculated and verified.

## Summary

- The JSON data MUST be represented as [Canonical JSON](http://wiki.laptop.org/go/Canonical_JSON)
- Every image and invocation image field in the `bundle.json` MUST have a `digest:` field that MUST contain a digest of the image.
- Digests are computed in accordance with the underlying image type (e.g. OCI bundles are validated by computing the top hash of a Merkle tree, VM images are computed by digest of the image)

## Image Integrity with Digests

CNAB correlates a number of images, of varying types, together in a single bundle. This section of the specification defines how image integrity is to be tested via digests and checksumming.

### Digests, OCI CAS, and Check Summing

A frequently used tool for validating that a file has not been changed between time T and time T+n is _checksumming_. In this strategy, the creator of the file runs a cryptographic hashing function (such as a SHA-512) on a file, and generates a _digest_ of the file's content. The digest can then be transmitted to a recipient separately from the file. The recipient can then re-run the same cryptographic hashing function on the file, and verify that the two functions are identical. There are more elaborate strategies for digesting complex objects, such as a [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree), which is what the OCI specification uses. In any event, the output _digest_ can be used to later verify the integrity of an object.

The OCI specification contains a [standard for representing digests](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests). In its simplest case, it looks like this:

```text
ALGO:DIGEST
```

Where `ALGO` is the name of the cryptographic hashing function (`sha512`, `md5`, `blake2`...) plus some OPTIONAL metadata, and DIGEST is the ASCII representation of the hash (typically as a hexadecimal number).

> Note: The OCI specification only allows `sha256` and `sha512`. This is not a restriction we make here.

For example:

```text
sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b
```

### Digesting Objects in the `bundle.json`

CNAB is composed of a `bundle.json` and a number of supporting images. Those images are referenced by the `bundle.json`. Thus, digesting those artifacts and including their digest in the `bundle.json` provides a convenient way to store (and locate) digests.

To that end, in a signed bundle, anything that shows up in the `invocationImages` or `images` section of the `bundle.json` MUST have a `digest`:

```json
{
    "name": "technosophos.helloworld",
    "version": "0.1.2",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:0.1.0@sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b",
            "digest": "sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b"
        }
    ],
    "images": {
        "image1": {
            "image": "image1name:image1tag@sha256:aaaa624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b",
            "digest": "sha256:aaaa624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b",
            "refs": [
                {
                    "path": "image1path",
                    "field": "image.1.field"
                }
            ],
        }
    }
}
```

Objects MUST contain a `digest` field even if the digest is present in another field. This is done to provide _access uniformity_, as well as to support image types that do not provide a built-in mechanism for embedding digest information. When (as is seen in the example above) a particular image type admits more than one way of attaching a digest, implementations SHOULD provide the digest in all applicable places.

> OCI images, for example, MAY embed a digest in the image's _version_ field. According to this specification, while this is allowed, it does not remove the requirement that the `digest` field be present and filled.

Different formats (viz. OCI) provide definitions for validating a digest. When possible, images should be validated using these definitions, according to their `imageType`. If a particular image type does not already define what it means to have a digest verified, the default method is to retrieve the object as-is, and checksum it in the format in which it was delivered when accessed.

Drivers MAY choose to accept the digesting by another trusted agent in lieu of performing the digest algorithm themselves. For example, if a driver requests that a remote agent install an image on its behalf, it MAY trust that the image digest given by that remote agent is indeed the digest of the object in question. And it MAY then compare that digest to the `bundle.json`'s digest. In such cases, a driver SHOULD ensure that the channel between the driver itself and the trusted remote agent is itself secured (for example, via TLS). Failure to do so will invalidate the integrity of the check.

As an example, the digest of Docker and OCI images in a thin bundle are not validated by the CNAB runtime, but by the container engine and the registry hosting those images. In the example above, when the docker driver creates the container with the specified invocation image, it specifies `technosophos/helloworld:0.1.0@sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b` as the container image. It is the responsibility of the container runtime and registry to guarantee immutability of the referenced image.

Signing and provenance operations on a bundle are covered in the [CNAB-Sec specification](300-security.md), which is a separate specification from the present.

Next section: [Declarative Images](801-declarative-images.md)
