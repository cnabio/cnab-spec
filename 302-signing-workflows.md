---
title: CNAB Security: Signing workflows
weight: 302
---

# Signing Workflows

* [Signing workflow for the minimum viable product (MVP)](#signing-workflow-for-the-minimum-viable-product-mvp)

This document is a _normative_ part of [CNAB Security](300-CNAB-security.md).

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

## Signing workflow for the minimum viable product (MVP)

This subsection documents how a signing workflow MAY typically work for an [MVP metadata repository](301-metadata-repositories.md). To reiterate, this document is normative, and while a CNAB-Sec-compliant implementation may borrow ideas from here, it MAY also choose to implement different signing workflows, especially for security models different from those described in the MVP.

### Setup (one-time task)

When bundle developers set up an MVP metadata repository for the first time, they SHOULD use one of the [known implementations of CNAB-Sec](304-known-implementations) to set up at least a complete set of the `root`, `timestamp`, `snapshot`, and `targets` TUF metadata.

To do so, developers SHOULD begin by generating the private keys for the `root` and `targets` roles. There are two options for the `timestamp` and `snapshot` roles. One is that developers MAY also generate and manage their private keys. The other is that the metadata repository MAY automatically generate and manage their private keys, as is the default on [Docker Content Trust](https://docs.docker.com/engine/security/trust/trust_key_mng/).

After these private keys have been generated, developers SHOULD generate the initial set of the TUF metadata. The `targets` metadata MAY list no bundles for the moment. Regardless of whether developers or the metadata repository manage the `timestamp` and `snapshot` keys, the `snapshot` metadata SHOULD point to this `targets` metadata file, and the `timestamp` metadata SHOULD point in turn to this `snapshot` metadata file. Finally, the `root` metadata SHOULD list the public keys for all top-level roles, including itself.

If developers wish to use in-toto to provide provenance for their bundles, then they SHOULD also generate and sign the [root layout](https://github.com/in-toto/docs/blob/master/in-toto-spec.md#43-file-formats-layout) for their bundle, which is out of the scope of this document. However, the interested reader MAY consult the [in-toto demo](https://github.com/in-toto/demo) for an example of how to do so.

### New or updated bundles (periodic task)

Whenever developers wish to release a new version of a bundle, they SHOULD perform the following steps.

First, they SHOULD sign new `targets` metadata that points to the new version of the bundle. If developers wish to use in-toto to provide provenance for their bundles, then, as discussed in [metadata repositories](301-metadata-repositories.md), they MAY also list in the custom targets metadata all of the root layout and link metadata required to verify this bundle.

Second, regardless of whether developers or the metadata repository holds these signing keys, the `timestamp` and `snapshot` metadata SHOULD be updated to the point to the new bundle.

### Recovering from a key compromise (exceptional task)

If the private key for the `timestamp`, `snapshot`, or `targets` role has been compromised, then developers SHOULD rotate their keys using the `root` metadata.

If less than a threshold of the `root` keys have been compromised, then developers SHOULD use at least a threshold of the `root` keys to rotate its own keys. Otherwise, if more than a threshold of the `root` keys have been compromised, then it is safest for developers to sign new `root` metadata using no threshold of previous `root` keys, which will require end-users to update `root` metadata using a safe out-of-band mechanism.
