---
title: CNAB Security: Signing workflows
weight: 302
---

# Cloud Native Application Bundles Security (CNAB-Sec): Signing workflows 1.0 WD

* [Signing workflows](#signing-workflows)
* [References](#references)

This document is a _normative_ part of [CNAB Security](300-CNAB-security.md).

The keywords MUST, MUST NOT, REQUIRED, SHALL, SHALL NOT, SHOULD, SHOULD NOT, RECOMMENDED, MAY, and OPTIONAL in this document are to be interpreted as described in [RFC 2119](https://tools.ietf.org/html/rfc2119). The use of these keywords in this document are statements about how a CNAB implementation may fulfill the CNAB Security specification _only_.

## Signing workflows

This subsection documents signing workflows for images and bundles.

### Images

The signing workflows for traditional and community image repositories are generally the same. Assuming that images are signed as [RECOMMENDED](#image-repositories), the following SHOULD recommendations provide level 2 security. The additional MAY recommendations provide level 3 security.
The signing workflows for all metadata repositories are generally the same. Assuming that images are signed as [RECOMMENDED](301-CNAB-metadata-repositories.md), the following SHOULD recommendations provide level 2 security. The additional MAY recommendations provide level 3 security.

In the beginning, administrators set up the top-level TUF roles once as described above. Every year or so, they SHOULD revoke and replace keys for the top-level roles. Less frequently, every month or so, or in the event of a key compromise, they SHOULD add, update, or remove delegations of images to their respective developers.

After administrators have delegated an image to its respective developers, the latter SHOULD add, update, or remove TUF metadata about tags or delegations whenever they like. Administrators SHOULD make available authenticated Application Programming Interfaces (APIs) for updating these metadata.
If end-to-end authenticity and integrity of tags is desired, the developers MAY also sign targets metadata about the in-toto root layout for this image, as well as the public keys used to verify this root layout. They MAY also delegate the signing of links for these tags to automation. They MAY do this every year or so, or in the event of a key compromise.
If developers have delegated the signing of tags to automation, then the latter SHOULD add, update, or remove TUF metadata whenever tags have been updated. The automation SHOULD update signatures on metadata periodically (say, every day) in order to prevent expired metadata from being accidentally served to users.

If end-to-end authenticity and integrity of tags is desired, the automation MAY include metadata about links for tags. In this case, developers MAY also sign some of these links. How links are propagated from developers to the automation is up to the implementor, and thus out of the scope of this document.

Administrators SHOULD refresh the `timestamp` and `snapshot` metadata whenever developers or automation upload new metadata to the repository.
Administrators SHOULD refresh the `timestamp` and `snapshot` metadata whenever developers or automation upload new metadata to the metadata repository.

It is RECOMMENDED for administrators to turn on [consistent snapshots](https://github.com/theupdateframework/specification/blob/master/tuf-spec.md#7-consistent-snapshots) so that new versions of metadata can be safely written while old versions are being consumed at the same time. It is RECOMMENDED that administrators, developers, and automation garbage-collect obsolete or expired metadata.

### Thin bundles

[**TODO**]

### Thick bundles

[**TODO**]

## References

1. [The TUF specification](https://github.com/theupdateframework/specification)
2. [The in-toto specification](https://github.com/in-toto/docs)
3. [ITE-2: A general overview of combining TUF and in-toto to build compromise-resilient CI/CD
](https://github.com/in-toto/ITE/pull/4)
4. [Secure Publication of Datadog Agent Integrations with TUF and in-toto](https://www.datadoghq.com/blog/engineering/secure-publication-of-datadog-agent-integrations-with-tuf-and-in-toto/)
5. [ITE-3: Real-world example of combining TUF and in-toto for packaging Datadog Agent integrations
](https://github.com/in-toto/ITE/pull/5)
6. [PEP 458 -- Surviving a Compromise of PyPI](https://www.python.org/dev/peps/pep-0458/)
7. [PEP 480 -- Surviving a Compromise of PyPI: The Maximum Security Model](https://www.python.org/dev/peps/pep-0480/)