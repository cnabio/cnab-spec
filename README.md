# Cloud Native Application Bundle Specifications

## Abstract

Cloud Native Application Bundles (CNAB) are a package format specification that describes a technology for bundling, installing, and managing distributed applications, that are by design, cloud agnostic.

## CNAB Core 1.0 (Final)

The CNAB Working Group with the joint approval of the Executive Directors has approved the CNAB Core 1.0 specification for publication. CNAB Core 1.0 is complete.

For more information on the approval process, see [the process documentation](901-process.md). Further changes to CNAB Core will be considered for CNAB Core 1.1.

## Table of Contents

- Chapter 1: [Cloud Native Application Bundle Core 1.0.0 (CNAB1)](100-CNAB.md)
  1. [The bundle.json File](101-bundle-json.md)
  1. [The Invocation Image Format](102-invocation-image.md)
  1. [The Bundle Runtime](103-bundle-runtime.md)
  1. [Bundle Formats (Thick and Thin)](104-bundle-formats.md)
- Chapter 2: [Cloud Native Application Bundle Registry 1.0.0 (CNAB-Reg)](200-CNAB-registries.md)
- Chapter 3: [Cloud Native Application Bundle Security 1.0.0 (CNAB-Sec)](300-CNAB-security.md)
- Chapter 4: [Cloud Native Application Bundle Claims 1.0.0 (CNAB-Claims1)](400-claims.md)
- Chapter 5: [Cloud Native Application Bundle Dependencies 1.0.0 (CNAB-Deps)](500-CNAB-dependencies.md)
- Chapter 8: Non-normative Content
  1. [Declarative Invocation Images](801-declarative-images.md)
  1. [Credential Sets](802-credential-sets.md)
  1. [Repositories](803-repositories.md)
  1. [Well known custom actions](804-well-known-custom-actions.md)
  1. [Disconnected Scenarios](805-airgap.md)
- Chapter 9: Appendix
  1. [Appendix A: Preliminary Release Process](901-process.md)

## Contributing

The specification is licensed under [OWF Contributor License Agreement 1.0 - Copyright and Patent](http://www.openwebfoundation.org/legal/the-owf-1-0-agreements/owf-contributor-license-agreement-1-0---copyright-and-patent) in the [LICENSE](./LICENSE) file.

## Communications

### Meetings
* Community Meeting: Every other Wednesday at 09:00 AM US Pacific
  * https://zoom.us/j/653255416
  * [Meeting notes and Agenda](https://aka.ms/cnab/meeting).

### Slack Channel
#cnab Slack channel for related discussion in
[CNCF's Slack workspace](https://slack.cncf.io/).

### Mailing List

We operate a [mailing list](https://lists.jointdevelopment.org/g/CNAB-Main) via the Joint Development Foundation.

## Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119][rfc2119].

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard][c99-unspecified].

An implementation IS compliant if it satisfies all the MUST, REQUIRED, and SHALL requirements.

An implementation IS NOT compliant if it fails to satisfy one or more of the MUST, REQUIRED, or SHALL requirements.

[c99-unspecified]: http://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18
[rfc2119]: http://tools.ietf.org/html/rfc2119

## Document Naming Conventions

- The CNAB Core specification is contained in the 1xx documents.
- The CNAB Registry specification is contained in the 2xx documents.
- The CNAB Security specification reserves 3xx level documents.
- The Claims specification reserves 4xx documents.
- The CNAB Dependencies specification uses 5xx documents.
- The 8xx-level documents are reserved for non-normative guidance.
- The 9xx-level documents are reserved for process documents.
