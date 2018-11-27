# Cloud Native Application Bundle Specifications

# Abstract

Cloud Native Application Bundles (CNAB) are a package format specification that describes a technology for bundling, installing, and managing distributed applications, that are by design, cloud agnostic.

# Table of Contents

1. [Cloud Native Application Bundle Core 1.0.0 (CNAB1)](./100-CNAB.md)
    1. [The bundle.json File](101-bundle-json.md)
    2. [The Invocation Image Format](102-invocation-image.md)
    3. [The Bundle Runtime](103-bundle-runtime.md)
    4. [The Claims System](104-claims.md)
    5. [Signing and Provenance](105-signing.md)
    6. [Declarative Invocation Images](106-declarative-images.md)
    7. [CNAB Repositories](107-repositories.md)

## Contributing

The specification and code is licensed under the Apache 2.0 license found in the [LICENSE](./LICENSE) file.

# Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119][rfc2119].

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard][c99-unspecified].

An implementation IS compliant if it satisfies all the MUST, REQUIRED, and SHALL requirements.

An implementation IS NOT compliant if it fails to satisfy one or more of the MUST, REQUIRED, or SHALL requirements.

[c99-unspecified]: http://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18
[rfc2119]: http://tools.ietf.org/html/rfc2119

# Document Naming Conventions

During the draft period of the specification, the documents are named according to the following convention:

- `000-099` contains front matter
- `100`-`799` contain the specification proper
    - The first digit is the chapter number
    - The following two digits are the section numbers
- `800`-`899` is reserved for non-normative content (examples, patterns, best practices)
- `900`-`999` is reserved for appendices 