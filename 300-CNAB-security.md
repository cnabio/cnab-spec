---
title: Bundle Security
weight: 300
---

# Cloud Native Application Bundles Security (CNAB-Sec) 1.0 WD

![cnab-security](https://user-images.githubusercontent.com/686194/61752644-54580b80-ad61-11e9-9518-608534d09bdd.png)

The CNAB Security specification augments the [CNAB Core specification](100-CNAB.md) by standardizing on security mechanisms for signing, verifying, and attesting CNAB packages. It describes both a client/registry security model and a verification chain (provenance) model. By design, this provides security mechanisms that work in air-gapped networks.

This specification is distinct from the CNAB Core specification. An implementation may comply with the CNAB Core specification, and yet not comply with this specification. The use of terms such as MUST and SHOULD in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

> An earlier version of the CNAB Core specification used OpenPGP-based mechanism to sign and verify bundles. This specification supersedes that portion of the specification.

This specification addresses two related problems:

- Securing the delivery of packages from a source to a client
- Ensuring the provenance of a bundle robustly, even across networks and over long periods of time

Industry standards such as [NISTIR-8176](https://csrc.nist.gov/News/2017/NIST-Releases-NISTIR-8176) provide recommendations on how container security may be achieved, and a goal of this specification is to meet such recommendations.

## Key Terms

In addition to the terminology introduced in the CNAB Core specification, this specification uses additional terms:

- *Acceptance Criteria:* The conditions that must obtain before a particular CNAB may be considered acceptable for use. This may be, in some cases, as simple as "the CNAB bundle has the expected digest" or as complicated as a tree-structured diagram of attestations.
- *Attestation:* an assertion, cryptographically verifiable, that a particular identity has performed a particular task. For example, an attestation may capture the assertion that "Developer Padma has approved version 1.2.3 of the WebApp bundle."
- *Bundle:* This term refers to a CNAB package, and is used to clarify that the object in question is a particular package, not the CNAB specification.
- *CNAB Client:* Any CNAB-aware program that operates on a package generated elsewhere.
- *CNAB Registry:* Any repository of CNAB bundles. While there is a specification for how such bundles should be stored (the [CNAB Registry specification](200-CNAB-registries.md)), this specification assumes only that a registry is a remote storage system from whence CNAB bundles can be pulled and (optionally) to which they may be pushed.
- *Layout:* A declaration of a series of attestations that must be made in order to consider a particular bundle to have satisfied the acceptance criteria.
- *Link:* An envelope that binds a particular attestation to a particular step in a layout.
- *Signature:* A cryptographic verification that the possessor of a particular private key has approved a particular body of content. For example, a signature may capture the approval: "Developer Padma has approved this archive file". Note that an attestation always contains a signature, but a signature may be used for things other than attestation.
- *Verification:* The process of applying a public key to a signature and associated body of content, with the purpose of verifying that the body of content is identical to the body that was signed by the private key.

## Securely Moving Bundles To And From Registries

Per the CNAB Core 1.0 specification, bundles may be stored and transmitted in two forms: thin bundles and thick bundles. Both thin and thick bundles can be transmitted over the network. This section describes how CNAB clients (including CNAB runtimes) may safely establish the identity and integrity of a CNAB bundle being fetched from a CNAB registry.

(Per the [CNAB Registry specification](200-CNAB-registries.md), CNABs may be stored in any OCI Registry 1.0-compliant implementation)

Securing a CNAB registry is done by implementing the TUF (The Update Framework) specification on both the client and registry side, and then adhering to the protocol. A registry MUST comply with the [TUF 1.0 specification](https://github.com/theupdateframework/specification/blob/master/tuf-spec.md) to meet this specification. A client also MUST comply with the TUF specification to meet this specification. Furthermore, the client MUST perform additional verification on the content of a bundle.

This specification does not require that a bundle registry support both thick and thin bundles. It prescribes processes for storage of each. Unless otherwise noted, a registry is considered _compliant_ if it implements storage for at least one of the two.

This specification does not specifically deal with network-layer (Layer 7) security, but assumes that implementations will apply correct security constraints for all network traffic. TLS/SSL is likely the method by which such traffic should be secured.

### Overview

There are three major aspects of securely working with registries:

1. Pulling bundles _from_ a registry _to_ a client
2. Pushing bundles _from_ a client _to_ a registry
3. Managing keys between a client and a registry

### Pulling and Verifying Bundles

When a client pulls a bundle from a registry, the client MUST verify the package, where verification consists of the following:

- Fetch (from a TUF endpoint) the signature object for the given bundle
- Fetch the bundle descriptor (optionally by fetching a thick bundle and unpacking the descriptor)
- Verify the requisite signatures for a package, as described in [the TUF specification](https://github.com/theupdateframework/specification/blob/master/tuf-spec.md#5-detailed-workflows).
- Verify that each Invocation Image in the bundle descriptor (a) has a content digest, and (b) the content digest matches the referenced image (by recalculating the hash)
- Verify that each Image referenced in the bundle descriptor (a) has a content digest, and (b) the content digest matches the referenced image (by recalculating the hash)

A client MUST be able to verify SHA-256 and SHA-512 content digests. But this specification accords with the CNAB Core specification in that it does not dictate what are the the contents of the Invocation Images or regular Images. Examples in this document often refer to OCI images as the content of these image fields, but other formats are acceptable.

If any of these steps fail to obtain, the client MUST NOT perform an action on the bundle. Clients SHOULD produce an error.

### Signing and Pushing Bundles

When a client pushes a bundle to a registry, the client MUST prepare the security data on the package, which consists of the following:

- Generate the content digest (SHA-256 or SHA-512) of each invocation image, and write this information into the bundle descriptor
- Generate the content digest of each image, and write this information into the bundle descriptor
- Optionally, push the images to their appropriate remote location (thin bundles)
- Sign the requisite portion of the bundle
	- For thin bundles, this is the bundle descriptor alone
	- For thick bundles, this is the full bundle, in its representation that it will be submitted to the bundle registry
- Send the signature data to the TUF server, per the details of that server's implementation. (The TUF specification does not specify the process for uploading signatures. However, Notary provides one API for doing this.) TODO: Do we want to limit this to the Notary-prescribed protocol for ease of developing compatible clients?
- Send the data that has been signed to the CNAB Registry

In regard to the final two steps, specific implementations may require different ordering of operations. Thus, a CNAB client MAY perform these steps in any order that does not sacrifice the integrity of the cryptographic signatures.

### Key Management

A client SHOULD maintain a cache of trusted public keys. This cache MAY be directly managed, and this cache MAY adopt a strategy of Trust On First Use (TOFU), whereby the first time a package is fetched from a trusted source, the client may also request the public keys from that server. Other than the TOFU case, clients SHOULD NOT attempt to fetch release keys and packages from the same source at the same time. Timestamp, distribution, and root keys (all part of the TUF specification) MUST be managed according to the TUF specification.

## Provenance and Software Supply Chain

The previous section of the specification deals with managing the integrity of the process of transferring a CNAB bundle to and from a CNAB Registry.

In addition to verifying integrity, there is a second related set of verifications that may be performed on a CNAB. These are related to _provenance_.

Provenance captures concern tracing the origin of a package and its components, and ascertaining the integrity of the steps in the process of composing the package from its components. Compare this with the idea of the _software supply chain_, which is a process for composing software in such a way that it preserves provenance in a verifiable way.

This section specifies a provenance model for CNAB based on the [In-Toto specification](https://github.com/in-toto/docs/blob/v0.9/in-toto-spec.md). This specification pays less attention to the software supply chain (e.g. the specific set of steps), focusing instead on how this information is to be stored and verified in a CNAB bundle.

In-Toto, upon which this part of the specification is built, describes two major concepts: A _layout_ that describes the steps which must be performed to establish provenance for a bundle, and multiple _links_, which tie a cryptographic attestation (along with some metadata) to a step in the layout. For example, if a layout has three steps, verification of that layout would require three links -- one for each step. Layouts may be linear, but they also may be tree-shaped, where one layout invokes other "sublayouts", which in turn each may invoke further sublayouts.

For CNAB, bundles may have associated layouts, and such layouts may be satisfied by links.

### Attaching Layouts and Links to Bundles

For layouts and links to function properly, they must be stored in proximity to the bundle.

The "default layout" may be stored inside of the `custom` section of the bundle descriptor.

EXAMPLE

#### Links and Thick Bundles

Links in a thick bundle may be packaged into the tar archive at the following path: `links/`.  The structure of the files and directories inside of that directory are all well-specified in the In-Toto specification.

When a bundle is imported (or extracted), a client SHOULD verify the layout against the contents of `links`.

#### Layouts and Thin Bundles

For thin bundles, links must be stored external to the bundle, since (a) the bundle descriptor is the artifact of record, and (b) a layout will likely make assertions about the bundle descriptor, and therefore cannot be embedded within the descriptor.

An implementation of provenance that uses thin bundles MAY store the links inside of a CNAB registry as layers attached to a bundle, using the `mediaType` `x-application/io.cnab.link` TODO: This should be closer to whatever media types we decide to use for CNAB.

### Key Management

A crucial aspect of verifying layouts is the ability of a client to verify an attestation against the public key used to generate a link. And because layouts are themselves signed, it is implicit in the system that public keys must be shared between the bundle creator(s) and the client.

The specification does not prescribe particular key management steps. However, it explicitly allows for certain key management practices. Clients MAY use the keys managed for TUF transactions in order to verify layouts. Clients MAY include TOFU-acquired public keys. Clients SHOULD NOT use TUF's timestamp, distribution, or root keys to sign or verify layouts, as these keys are ascribed special significance in TUF, and that significance does not align with the provenance function described here.

### Composing Layouts

Layouts are well-defined in the In-Toto specification. However, a certain subset of In-Toto features are considered OPTIONAL in the present spec:

- Inspections: In-Toto provides a method for running arbitrary commands on a machine in order to verify that specific parts of a software supply chain are correct. This specification does not require support for inspections. Moreover, this specification discourages allowing bundles to run external commands even for inspection purposes.

### Verifying Layouts

The In-Toto specification defines the procedure for verifying individual links against the layout. This section adds CNAB-specific context to the more general specification.

If any part of a layout cannot be satisfied, the client SHOULD generate an error and MUST NOT continue installing the bundle.

If a link does not match a step in the layout, the link SHOULD be ignored. But if multiple links apply to the same step of a layout, and at least one can be verified, that layout stem SHOULD be considered verified.

A client MAY choose to apply a non-default (non-package-supplied) layout, provided this is exposed to the user of the client. In this way, a more or less stringent layout may be applied to a package as a method of providing context-specific validations. For example, a bundle may ship with an eight-step layout. However, the destination environment only requires validation of two steps. A client MAY consume a two-step layout and match that against the links provided by the bundle. In this case, the package is considered validated if the newly applied layout passes. External layouts MAY be signed by a key other than the one used on the default layout.


## Signing and Digests for Thick Bundles

When moving thick CNABs between networks, it is often necessary to export a bundle into its thick (single-file) format, transport it, and then import it on the destination host. In this model, the TUF-oriented protocol does very little good, and the In-Toto layouts are stored within the bundle. In such a case, it may be necessary to add a layer of validation around the bundle itself.

To provide this layer, a signed digest of the bundle contents must be prepended to the bundle, creating a new file. The format for this is as follows:

```
<signature-scheme>$<signed-digest>$<gzipped-tar-data>
```

Where `<signature-scheme>` is to be replaced by one of the following schemes, `<signature-digest>` is to be replaced by the signature, and `<gzipped-tar-data>` is to be replaced by the raw gzipped tar data. The `$` characters are literal.

The following signature schemes (as defined in the Section 4.2 of the TUF specification) must be supported:

- "rsassa-pss-sha256" : RSA Probabilistic signature scheme with appendix. The underlying hash function is SHA256.
   https://tools.ietf.org/html/rfc3447#page-29

- "ed25519" : Elliptic curve digital signature algorithm based on Twisted
  Edwards curves.
  https://ed25519.cr.yp.to/

- "ecdsa-sha2-nistp256" : Elliptic Curve Digital Signature Algorithm
   with NIST P-256 curve signing and SHA-256 hashing.
   https://en.wikipedia.org/wiki/Elliptic_Curve_Digital_Signature_Algorithm

The resulting file should have the extension `.cnab`. For example, if `webapp.tgz` is signed and digested, a new file should be created named `webapp.cnab`, with the signature prefix followed by the unaltered data in the original `webapp.tgz` file.
