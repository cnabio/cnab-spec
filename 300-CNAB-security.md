---
title: Bundle Security
weight: 300
---

# Cloud Native Application Bundles Security (CNAB-Sec) 1.0.0 GA
*[Working Group Approval](901-process.md), Jun. 2020*

The CNAB Working Group has approved the CNAB Security 1.0 specification. This specification is now pending approval from the CNAB Foundation executive committee, after which it will be finalized as CNAB Security 1.0 (AD). For more information on the approval process, see [the process documentation](901-process.md). Further changes to CNAB Security will be considered for CNAB Security 1.1.


* [Introduction](#introduction)
* [Roadmap](#roadmap)
  * [Gradual security](#gradual-security)
* [References](#references)

![cnab-security](https://user-images.githubusercontent.com/686194/61752644-54580b80-ad61-11e9-9518-608534d09bdd.png)

This specification is distinct from the CNAB Core specification. An implementation may comply with the CNAB Core specification, and yet not comply with this specification.

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

The reader is assumed to be thoroughly familiar with the [TUF][tuf-spec] and [in-toto][in-toto-spec] specifications, and some deployment models such as [PEP 458][pep-458], [PEP 480][pep-480], [ITE-4][ite-4], [ITE-5][ite-5], and [Datadog Agent integrations][datadog-agent-integrations].

The CNAB Security specification augments the [CNAB Core specification](100-CNAB.md) by standardizing on security mechanisms for signing, verifying, and attesting CNAB packages. It describes both a client/registry security model and a verification chain (provenance) model.

This specification is distinct from the CNAB Core specification. An implementation may comply with the CNAB Core specification, and yet not comply with this specification. The use of terms such as MUST and SHOULD in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

> An earlier version of the CNAB Core specification used OpenPGP-based mechanism to sign and verify bundles. This specification supersedes that portion of the specification.

## Introduction

It is important to sign bundles and images so that attackers cannot tamper with them without being detected. To enable faster deployments, they are likely to be built and signed by continuous integration / continuous delivery (CI/CD) instead of developers. While CI/CD provides benefits such as safer handling of signing keys, it also has a severe drawback. Since signing keys must be kept _online_ so that CI/CD can build and sign new bundles or images on demand, attackers who compromise the infrastructure can also abuse these keys to build, sign, and distribute malicious versions. The goal of this document is to discuss how to achieve _compromise-resilience_ in this setting: that is, even if the infrastructure is compromised anywhere between developers and end-users, attackers should not be able to cause end-users to install malicious versions of bundles or images that were not released by developers.

While SSL / TLS protects users from man-in-the-middle (MitM) attacks, it is _not_ compromise-resilient, because attackers who compromise the infrastructure can abuse the online SSL / TLS key to replace files signed in transit but not at rest. Likewise, other solutions, such as GPG or X.509, do not support necessary features such as in-band key revocation or rotation.

We propose using [The Update Framework (TUF)](https://theupdateframework.com) and [in-toto](https://in-toto.io) to solve this problem. TUF is a framework that provides security between _registries_, or the servers used to distribute bundles or images, and end-users. Specifically, TUF ensures the authenticity and integrity of bundles and images downloaded from registries. in-toto is a framework that provides security between developers and registries. Specifically, in-toto ensures the end-to-end integrity of the _software supply chain_ used to build and sign bundles or images. By combining both, we get a system that provides compromise-resilience between developers and end-users.

 Signing and verifying [claims](400-claims.md) are out of the scope of this document.

## Roadmap

The basic idea is to use TUF as a transport protocol that distributes several files in a compromise-resilient manner:

* The root of trust for all bundles and images, as well as TUF and in-toto metadata.
* The software supply chains defined using in-toto.
* The public keys used to verify these supply chains.

While images and bundles are published on registries, TUF and in-toto metadata about images and bundles are published on [_metadata repositories_](301-metadata-repositories.md). Metadata repositories are conceptually, but not necessarily physically, distinct from [CNAB](200-CNAB-registries.md) or [OCI](https://github.com/opencontainers/distribution-spec/blob/master/spec.md) registries. Furthermore, the key management models on metadata repositories are largely prescriptive. For these reasons and more, metadata repositories are out of the scope of this document.

Unless specified otherwise, metadata are TUF and in-toto metadata, which may be produced using the suggested [signing workflows](302-signing-workflows.md).

When a [bundle runtime](103-bundle-runtime.md) installs a bundle, it should first verify images and bundles using the suggested [verification workflows](303-verification-workflows.md).

Considerations for airgapped environments are discussed [here](805-airgap.md#cnab-security).

### Gradual security

Implementors / adopters MAY wish to implement security gradually instead of an all-or-nothing approach.

Level 0: neither images nor bundles are signed.
   * Pros: no security to maintain.
   * Cons: completely lacking compromise-resilience.

Level 1a: images are signed using TUF, but bundles are unsigned.
   * Pros: easiest level of security to achieve, e.g. using Docker Content Trust (metadata repository) and Notary (client).
   * Cons: only partly compromise-resilient. If a [CNAB registry](200-CNAB-registries.md) is compromised, a bundle could be easily replaced with malicious one.

Level 1b: bundles are signed using TUF, but images are unsigned.
   * Pros: easiest level of security to achieve, e.g. using Docker Content Trust (metadata repository) and Notary (client).
   * Cons: only partly compromise-resilient. If an [OCI registry](https://github.com/opencontainers/distribution-spec/blob/master/spec.md) is compromised, a bundle runtime may be vulnerable to accidentally bundling malicious images.
  
Level 2: both bundles and images are signed using TUF.
   * Pros: resilient against compromise of bundle and image registries.
   * Cons: not resilient against compromise of CI for bundles or images.

Level 3a: images are signed using TUF and in-toto. Bundles are signed using TUF.
   * Pros: resilient against compromise of CI for images.
   * Cons: adds work to developers and administrators. Not resilient against compromise of CI for bundles.

Level 3b: bundles are signed using TUF and in-toto. Images are signed using TUF.
   * Pros: resilient against compromise of CI for bundles.
   * Cons: adds work to developers and administrators. Not resilient against compromise of CI for images.

Level 4: both bundles and images are signed using TUF and in-toto.
   * Pros: resilient against compromise of CI for bundles and images.
   * Cons: greatest amount of work in implementation and adoption.

[tuf-spec]: https://github.com/theupdateframework/specification
[in-toto-spec]: https://github.com/in-toto/docs
[ite-4]: https://github.com/in-toto/ITE/pull/4
[datadog-agent-integrations]: https://www.datadoghq.com/blog/engineering/secure-publication-of-datadog-agent-integrations-with-tuf-and-in-toto/
[ite-5]: https://github.com/in-toto/ITE/pull/5
[pep-458]: https://www.python.org/dev/peps/pep-0458/
[pep-480]: https://www.python.org/dev/peps/pep-0480/
