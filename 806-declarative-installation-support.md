---
title: Declarative Installation Support
weight: 806
---

# Declarative Installation Support

This section is non-normative and describes how bundles with no invocation image are intended to be used.

## Summary

A bundle with no invocation images cannot be installed using a standard CNAB runtime. A standard runtime should fail when asked to perform an action on such a bundle. Such a bundle may, however, be installed by a CNAB runtime which supports the bundle's extension(s) and/or other metadata included with the bundle.

A bundle with no invocation images should omit the parameters, credentials, and outputs sections since these are intended to pass information to/from invocation images.

## Background

In environments, notably Kubernetes, where applications define their installations declaratively, the presence of an invocation image in a bundle was seen by some as a potential security exposure (similar to ActiveX). This was likely to limit the adoption of CNAB in those environments.

The invocation image was therefore made optional. A bundle without an invocation image can be installed only by runtimes that support the bundle's extension(s) and/or other metadata included with the bundle.

This makes CNAB more broadly applicable, notably within the Kubernetes community. Features of the CNAB specification other than invocation images, including image relocation, air gap support, registry representation, and supply chain security, can then be exploited without the necessity of using invocation images.

Optional invocation images enable CNAB to support the following experiences:

* Provide a purely declarative install mechanism. All necessary artifacts (with metadata about those artifacts) is included in the bundle. External tools are required to interpret and use those artifacts to install something, or do anything else that is needed.

* Move the install tool chain into a container (or set of containers) that are left as a suggestion (e.g. because the invocation image is omitted) in the metadata on how to use the artifacts. It would be up to the user to use (or not use) those containers.

* Have everything bundled into the invocation image. That image would be opaque and must be used to install the application as some of the required artifacts may be baked into that image.
