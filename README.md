# Cloud Native Application Bundle Specification

This document describes the Cloud Native Application Bundle (CNAB) specification.

1. [CNAB Specification](./100-CNAB.md)
    1. [The bundle.json File](101-bundle-json.md)
    2. [The Invocation Image Format](102-invocation-image.md)
    3. [The Bundle Runtime](103-bundle-runtime.md)
    4. [The Claims System](104-claims.md)
    5. [Signing and Provenance](105-signing.md)
    6. [Declarative Invocation Images](106-declarative-images.md)
    7. [CNAB Repositories](107-repositories.md)

## Contributing

The specification and code is licensed under the Apache 2.0 license found in the [LICENSE](./LICENSE) file.

### Git commit

#### Sign your work

A DCO is lightweight way for a developer to certify that they wrote or otherwise have the right to submit code or documentation to a project. The way a developer does this is by adding a Signed-off-by line to a commit. When they do this they are agreeing to the DCO.

The full text of the DCO can be found at https://developercertificate.org. It reads:

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```
An example signed commit message might look like:

    Signed-off-by: Some Developer somedev@example.com

Git has a flag that can sign a commit for you. An example using it is:

```bash
$ git commit -s -m 'An example commit message'
```

# Notational Conventions

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and "OPTIONAL" are to be interpreted as described in [RFC 2119][rfc2119].

The key words "unspecified", "undefined", and "implementation-defined" are to be interpreted as described in the [rationale for the C99 standard][c99-unspecified].

An implementation IS compliant if it satisfies all the MUST, REQUIRED, and SHALL requirements.

An implementation IS NOT compliant if it fails to satisfy one or more of the MUST, REQUIRED, or SHALL requirements.

[c99-unspecified]: http://www.open-std.org/jtc1/sc22/wg14/www/C99RationaleV5.10.pdf#page=18
[rfc2119]: http://tools.ietf.org/html/rfc2119
