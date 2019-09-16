---
title: CNAB Security: Verification workflows
weight: 303
---

# Cloud Native Application Bundles Security (CNAB-Sec): Verification workflows 1.0 WD

* [Verification workflows](#verification-workflows)
  * [Thin bundles](#thin-bundles)
  * [Images](#images)
  * [Thick bundles](#thick-bundles)
* [References](#references)

This document is a _prescriptive_ part of [CNAB Security](300-CNAB-security.md).

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

## Verification workflows

This section describes the general workflows that any bundle runtime SHOULD follow in order to verify bundles and images.

End-users SHOULD see no difference in their experience, unless an attack is caught by TUF or in-toto. In that case, installation of the bundle or image in question SHOULD be denied, and users SHOULD see an error message indicating why TUF or in-toto failed to verify the bundle or image.

### Thin bundles

The download and verification workflow for a thin bundle is roughly as follows:

1. Use TUF to download and verify the wheel for a given integration name and version number.
1. Use TUF to download and verify the in-toto metadata for the given wheel.
1. Use TUF to download and verify public keys for the in-toto root layout.
1. Use in-toto to verify that the given wheel matches the rules specified in the in-toto root layout.
1. If all of the above checks pass, return the wheel to the Datadog agent.

[**TODO**]

### Images

 When interpreting metadata from a metadata repository, it is up to a bundle runtime to determine how to download an image depending on its name and type. The signed bundle specifies where to download the image from.

 [**TODO**]

### Thick bundles

A [thick bundle](104-bundle-formats.md) contains all of the files needed to transmit a complete bundle in an airgapped / offline manner. It is simply a gzipped TAR containing at least:

* The bundle descriptor (`bundle.json`)
* One or more invocation images
* Zero or more images

The download and verification workflow for a thick bundle is the same as.

## References

1. [The TUF specification](https://github.com/theupdateframework/specification)
2. [The in-toto specification](https://github.com/in-toto/docs)
3. [ITE-2: A general overview of combining TUF and in-toto to build compromise-resilient CI/CD
](https://github.com/in-toto/ITE/pull/4)
4. [Secure Publication of Datadog Agent Integrations with TUF and in-toto](https://www.datadoghq.com/blog/engineering/secure-publication-of-datadog-agent-integrations-with-tuf-and-in-toto/)
5. [ITE-3: Real-world example of combining TUF and in-toto for packaging Datadog Agent integrations
](https://github.com/in-toto/ITE/pull/5)
6. [PEP 458 -- Surviving a Compromise of PyPI](https://www.python.org/dev/peps/pep-0458/)
7. [PEP 480 -- Surviving a Compromise of PyPI: The Maximum Security Model](https://www.python.org/dev/peps/pep-0480/)