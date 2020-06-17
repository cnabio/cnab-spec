---
title: Disconnected Scenario
weight: 805
---

# Disconnected Scenarios

This section is non-normative and describes the disconnected (air gapped) scenario mentioned in the introduction of [CNAB-Core](100-CNAB.md).

## Summary

Using CNAB in a disconnected environment involves transferring a bundle and its images into the 
environment such that the bundle can be installed and the bundled software executed successfully. 

## Internet Access

A typical disconnected scenario will have limited, intermittent or no internet access, whether by design or by situation.
To install a bundle directly in a disconnected environment, the bundle and its images need to be included in a [CNAB thick bundle](104-bundle-formats.md)
and transferred into the disconnected environment, for instance on a USB stick.

Alternatively, if a DMZ is available, it may be possible to read the bundle and/or its images from an external network
and write the bundle and/or its images to a registry inside the disconnected environment. 

## Private Registries

Common cloud patterns today reference artifacts and code from multiple sources as described in
[CNAB-Sec](300-CNAB-security.md).
Industry best practices, particularly [NISTIR-8176](https://csrc.nist.gov/News/2017/NIST-Releases-NISTIR-8176), emphasise the importance of
private registries and of avoiding the use of uncontrolled artifacts in production. CNAB aims to enable adherence to these best practices.

A CNAB runtime is NOT required to provide registry functionality to a CNAB in a disconnected environment.
It is assumed that an OCI compliant registry is available in the disconnected environment for hosting the bundle
and/or its images.

A private registry:
* Can be hosted in a disconnected environment for security, auditability, or other reasons.
* Provides complete control over when a bundle or image is updated or deleted:
    * This provides isolation from unwanted updates or deletion of the original bundle or image.
    * If the bundle or image becomes stale, for instance when it has known vulnerabilities, it can be deleted.

## CNAB Thick Bundles

A CNAB thick bundle provides a convenient archive format for transferring a bundle and its images into a 
disconnected environment. But thick bundles have other benefits.

A CNAB bundle may reference artifacts that are hosted in different repositories or registries.
These remote artifacts may change over time without changing the references.
If the digests of artifacts are provided in the bundle, the content of the artifacts cannot change without
changing the digests, but even then the artifacts, or their repositories or registries, may be deleted.

Archiving a CNAB and its images at a point of time as a CNAB thick bundle
provides protection against modification or deletion of images and also provides a central location for code
auditing and digital forensics of all code and references used in the CNAB.

## Image Relocation

When the images of a CNAB are _relocated_ to (that is, stored in), a private registry, the images should be loaded from the private registry when they are run.
This ensures that CNAB operations can function properly even if the original image repositories are unavailable.

The runtime uses the relocated reference of the invocation image so that the image is loaded from the private registry.

A [relocation mapping](103-bundle-runtime.md#relocation-mapping) is mounted so that the invocation
image is aware of the original and new values of image references and can replace the original image references with their relocated counterparts.
Thus the images referenced by the CNAB are also loaded from the private registry.

## CNAB Security

When copying images from public registries to an airgapped environment, the user has at least two options.

The first option is to use one of the [known
implementations](304-known-implementations) of CNAB Security to verify
signatures on bundles (and possibly images, although this is out of scope for
    CNAB Security) from public registries, but push unsigned bundles and images
to their own private, airgapped registries. The pros here are that the user
does verify the authenticity and integrity of public bundles, and does not need
to maintain any private signing and verification infrastructure, but the con is
that the user is not protected from internal attacks between registries and
consumers.

The second option is to use one of the [known implementations] of CNAB Security
to not only verify signatures on bundles (and optionally images) from public
registries, but also _independently_ sign and push bundles (and images) to
their own private, airgapped registries. The signing keys are independent from
the public, upstream registries, because: (1) the known implementation may not
support copying signatures, (2) the original signatures may expire within the
airgapped environment, and (3) there may be transformation of bundles (such as
producing private thick bundles from public thin bundles). This option is
more complicated than the first one, but does not suffer from the same cons.
