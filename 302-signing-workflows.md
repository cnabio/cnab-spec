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

When bundle developers set up an MVP metadata repository for the first time, they SHOULD use a tool like [`signy`](https://github.com/engineerd/signy) to set up encrypted keys for the TUF `root` and `targets` roles. The encrypted keys for the TUF `timestamp` and `snapshot` roles MAY be generated automatically and managed by the MVP metadata repository, as is the default on [Docker Content Trust](https://docs.docker.com/engine/security/trust/trust_key_mng/).

The generation of encrypted keys for the in-toto root layout and its associated functionaries as well as the root layout itself are out of the scope of this document as well as `signy`. However, the interested reader MAY consult the [in-toto demo](https://github.com/in-toto/demo) for an example of how to do so.

### New or updated bundles (periodic task)

Administrators SHOULD refresh the `timestamp` and `snapshot` metadata whenever developers or automation upload new metadata to the metadata repository.

### Recovering from a key compromise (exception task)

