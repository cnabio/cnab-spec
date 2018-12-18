# Cloud Native Application Bundle Specifications

## Abstract

Cloud Native Application Bundles (CNAB) are a package format specification that describes a technology for bundling, installing, and managing distributed applications, that are by design, cloud agnostic.

## Table of Contents

- Chapter 1: [Cloud Native Application Bundle Core 1.0.0 (CNAB1)](100-CNAB.md)
  1. [The bundle.json File](101-bundle-json.md)
  1. [The Invocation Image Format](102-invocation-image.md)
  1. [The Bundle Runtime](103-bundle-runtime.md)
  1. [The Claims System](104-claims.md)
  1. [Signing and Provenance](105-signing.md)
- Chapter 8: Non-normative Content
  1. [Declarative Invocation Images](801-declarative-images.md)
  1. [Credential Sets](802-credential-sets.md)
  1. [The Base Bundle Pattern](803-base-bundles.md)
  1. [Repositories](804-repositories.md)
- Chapter 9: Appendix
  1. [Appendix A: Preliminary Release Process](901-process.md)

## Contributing

The specification and code is licensed under the Apache 2.0 license found in the [LICENSE](./LICENSE) file.

## Communications

The mailing list for e-mail communications:
- Send emails to: [dev@opencontainers.org](mailto:dev@opencontainers.org)
- To subscribe see: https://groups.google.com/a/opencontainers.org/forum/#!forum/dev

Temporary #cnab Slack channel for related discussion in
[CNCF's Slack workspace](https://slack.cncf.io/).

## Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119][rfc2119].

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard][c99-unspecified].

An implementation IS compliant if it satisfies all the MUST, REQUIRED, and SHALL requirements.

An implementation IS NOT compliant if it fails to satisfy one or more of the MUST, REQUIRED, or SHALL requirements.

[c99-unspecified]: http://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18
[rfc2119]: http://tools.ietf.org/html/rfc2119

## Document Naming Conventions

During the draft period of the specification, the documents are named according to the following convention:

- `000-099` contains front matter
- `100`-`799` contain the specification proper
  - The first digit is the chapter number
  - The following two digits are the section numbers
- `800`-`899` is reserved for non-normative content (examples, patterns, best practices)
- `900`-`999` is reserved for appendices
