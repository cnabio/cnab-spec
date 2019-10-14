---
title: CNAB Security: Verification workflows
weight: 303
---

# Cloud Native Application Bundles Security (CNAB-Sec): Verification workflows 1.0 WD

- [Verification workflows](#verification-workflows)
- [References](#references)

This section describes the verification workflow that runtimes compliant with the CNAB Security Specifications must perform in order to verify the provenance and attestation of bundles. Regardless of [the bundle format][bundle-formats], the verification workflows is the same, with the signature of a thin bundle being the content digest of the `bundle.json` file, while the signature of a thick bundle being the content digest of the bundle archive.

End-users SHOULD see no difference in their experience, unless an attack is caught. In that case, the installation of the bundle in question MUST be denied, and users SHOULD see an error message indicating why the verification workflow failed.

## Verification workflows

1. The runtime obtains the bundle file. This operation can follow the [CNAB Registries Specification][registry-spec] for thin bundles and pull the bundle from an OCI registry, or use other protocols of obtaining it.
1. The content digest of the bundle is computed locally.
1. Follow [the TUF workflow][tuf-detailed-workflow] to download the bundle signature (as defined by) from the metadata repository.
1. The bundle signature and the locally computed digest are compared. If they are not equal, runtimes MUST stop the execution and not use the bundle for any subsequent operation.
1. If present in the metadata repository, runtimes MUST download the in-toto metadata associated with the bundle (as described in [the signing workflows document][signing-workflow]).
1. The runtime fetches the public signing key of the in-toto root layout from TUF.
1. The runtime validates all layouts structurally, and runs all in-toto verifications. If any validation or verification fails, the runtime MUST stop the execution and not use the bundle for any subsequent operations.

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

[bundle-formats]: 104-bundle-formats.md
[registry-spec]: 200-CNAB-registries.md
[tuf-detailed-workflow]: https://github.com/theupdateframework/specification/blob/master/tuf-spec.md#5-detailed-workflows
[metadata-repos]: 301-metadata-repositories.md
[signing-workflow]: 302-signing-workflows.md
