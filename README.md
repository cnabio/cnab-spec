# Cloud Native Application Bundle Specifications

## Abstract

Cloud Native Application Bundles (CNAB) are a package format specification that describes a technology for bundling, installing, and managing distributed applications, that are by design, cloud agnostic.

## CNAB Core 1.0 (Final)

The CNAB Working Group with the joint approval of the Executive Directors has approved the CNAB Core 1.0 specification for publication. CNAB Core 1.0 is complete.

For more information on the approval process, see [the process documentation](901-process.md). Further changes to CNAB Core will be considered for CNAB Core 1.1.

### Branch/Tag Structure

- [master](https://github.com/cnabio/cnab-spec) is the current working draft of all specs
- [cnab-core-1.0.1](https://github.com/cnabio/cnab-spec/tree/cnab-core-1.0.1) is the CNAB Core 1.0.1 final draft
- [cnab-core-1.0](https://github.com/cnabio/cnab-spec/tree/cnab-core-1.0) is the CNAB Core 1.0.0 final draft

## Table of Contents

- Chapter 1: [Cloud Native Application Bundle Core 1.0.0 (CNAB1)](100-CNAB.md)
  1. [The bundle.json File](101-bundle-json.md)
  1. [The Invocation Image Format](102-invocation-image.md)
  1. [The Bundle Runtime](103-bundle-runtime.md)
  1. [Bundle Formats (Thick and Thin)](104-bundle-formats.md)
- Chapter 2: [Cloud Native Application Bundle Registry 1.0.0 (CNAB-Reg)](200-CNAB-registries.md)
- Chapter 3: [# Cloud Native Application Bundles Security (CNAB-Sec) 1.0.0 WD](300-CNAB-security.md)
  1. [Cloud Native Application Bundles Security (CNAB-Sec): Metadata repositories 1.0.0 WD](301-metadata-repositories.md)
  1. [Cloud Native Application Bundles Security (CNAB-Sec): Signing workflows 1.0.0 WD](302-signing-workflows.md)
  1. [Cloud Native Application Bundles Security (CNAB-Sec): Verification workflows 1.0.0 WD](303-verification-workflows.md)
- Chapter 4: [Cloud Native Application Bundle Claims 1.0.0 (CNAB-Claims1)](400-claims.md)
- Chapter 5: [Cloud Native Application Bundle Dependencies 1.0.0 (CNAB-Deps)](500-CNAB-dependencies.md)
- Chapter 8: Non-normative Content
  1. [Declarative Invocation Images](801-declarative-images.md)
  1. [Credential Sets](802-credential-sets.md)
  1. [Repositories](803-repositories.md)
  1. [Well known custom actions](804-well-known-custom-actions.md)
  1. [Disconnected Scenarios](805-airgap.md)
  1. [Known implementations for CNAB Security 1.0.0 (CNAB-Sec 1.0.0)](806-security-known-implementations.md)
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

## Versioned Schema URLs

CNAB Spec schema defined in this repo (under `/schema`) will be hosted for all tagged versions.  This way, implementations can require specific schema versions for validation and assert compatibility with the corresponding versions.

Note that some tagged versions don't directly map to official schema versions.  For instance, a Git SHA may be appended if the spec is still in a Draft state, e.g. `cnab-claim-1.0.0-DRAFT+abc1234`.  Again, this facilitates the ability for implementations to pin to a certain tag whilst a spec is under heavy development with many breaking changes.

The schema files are hosted via `https://cnab.io/schema/<VERSION>/<SCHEMA>.schema.json`, e.g. https://cnab.io/schema/cnab-core-1.0/bundle.schema.json

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

The top-level description of each specification is found in the `x00` sections (e.g. 100 is the top of the CNAB Core specification, while 200 is the top of the CNAB Registry specification). Specifications may be broken down into subsections, which are numbered incrementally according to their section. Thus, the CNAB Core specification has a 100 (top-level) with four subsections (101, 102, 103, and 104). The 8XX section is entirely composed of non-normative documents, and each 8XX document stands alone.
