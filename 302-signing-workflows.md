---
title: CNAB Security: Signing workflows
weight: 302
---

# Cloud Native Application Bundles Security (CNAB-Sec): Signing workflows 1.0 WD

* [Signing workflow for the minimum viable product (MVP)](#signing-workflow-for-the-minimum-viable-product-mvp)

This document is a _normative_ part of [CNAB Security](300-CNAB-security.md).

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

## Signing workflow for the minimum viable product (MVP)

This subsection documents how a signing workflow MAY typically work for an [MVP metadata repository](301-metadata-repositories.md).

### Setup (one-time task)

When bundle developers set up an MVP metadata repository for the first time, they SHOULD use one of the [known implementations of CNAB-Sec](304-known-implementations) to set up at least a complete set of the `root`, `timestamp`, `snapshot`, and `targets` TUF metadata.

There are two options for the TUF `timestamp` and `snapshot` roles. One is that the private keys for these roles MAY be generated automatically and managed by the MVP metadata repository, if it supports it, as is the default on [Docker Content Trust](https://docs.docker.com/engine/security/trust/trust_key_mng/). The other is that, just as with the `root` and `targets` roles, the private keys for the `timestamp` and `snapshot` roles MAY be generated and managed by the bundle developers.

in-toto requires bundle developers to generate and sign the [root layout](https://github.com/in-toto/docs/blob/master/in-toto-spec.md#43-file-formats-layout) for their bundle, which is out of the scope of this document. However, the interested reader MAY consult the [in-toto demo](https://github.com/in-toto/demo) for an example of how to do so.

### New or updated bundles (periodic task)

Administrators SHOULD refresh the `timestamp` and `snapshot` metadata whenever developers or automation upload new metadata to the metadata repository.

### Recovering from a key compromise (exception task)

