---
title: Registries
weight: 200
---

# CNAB Registries 1.0

Draft, Feb. 2020

![cnab-registry](https://user-images.githubusercontent.com/686194/61753147-2b387a80-ad63-11e9-8a63-f250bcdf06b0.png)

---

This specification describes how CNAB bundles can be stored inside of OCI Registries.

This specification, the CNAB Registries specification, is not part of the CNAB Core specification. An implementation of CNAB Core MAY NOT implement this specification, yet still claim compliance with CNAB Core. A CNAB implementation MUST NOT claim to be CNAB Registries-compliant unless it meets this specification.

> Previous versions of the CNAB Core draft specification included a repository protocol. That work is completely superseded by this specification.

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://www.ietf.org/rfc/rfc2119.txt). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB registry storage specification only.

The CNAB Core 1.0 specification (CNAB1) does not dictate how bundles should be distributed. This is intentional, so that organizations that already have a way of distributing artifacts may continue to use it.
This document, the CNAB Registries specification, is not part of the core specification, and an implementation of CNAB Core MAY NOT implement it, yet still claim compliance with CNAB Core. A CNAB implementation MUST NOT claim to be CNAB Registries compliant unless it meets this specification.

This specification proposes the use of [OCI registries][oci-org] to distribute CNAB artifacts.

## Approach

Container registries provide a reliable and highly scalable distribution for container images. Implementations of container registries can be found in cloud provider services, vendor services, and open source projects. More and more artifact types are distributed with OCI registries, with the process standardized by the [OCI Artifacts][artifacts] project.

A Cloud Native Application Bundle is a collection of metadata and container images that are needed in order to successfully deploy an application - as such, it is not a _single_ new artifact, but represents a _collection of multiple artifacts_.

An [OCI image index (or simply OCI index)][oci-index] represents a collection of container images stored in a repository - so rather than using a new artifact, this specification proposes reusing the OCI index to represent Cloud Native Application Bundles, and use it as a natural distribution mechanism for the bundle file, invocation image, and images used by the application.

The approach of this specification is to reuse the existing distribution infrastructure for OCI artifacts when distributing CNAB bundles, to ensure compatibility with the CNAB Security specification, and to enable workflows such as the relocation of bundles and images.

In the following sections, we will explore [the representation of bundles in OCI registries](201-representing-CNAB-in-OCI.md) and [on disk representations](210-on-disk-representation.md).

[oci-org]: https://github.com/opencontainers/
[artifacts]: https://github.com/opencontainers/artifacts
[oci-index]: https://github.com/opencontainers/image-spec/blob/master/image-index.md
[oci-manifest]: https://github.com/opencontainers/image-spec/blob/master/manifest.md
[docker-manifest]: https://docs.docker.com/registry/spec/manifest-v2-2/
