---
title: CNAB Core
weight: 100
---

# Cloud Native Application Bundle Core 1.0.0 (CNAB1)
*[Working Draft](901-process.md), Nov. 2018*


The Cloud Native Application Bundle (CNAB) is a _standard packaging format_ for multi-component distributed applications. It allows packages to target different runtimes and architectures. Having a common format empowers application distributors to package applications for deployment on a wide variety of cloud platforms, providers, and services. Furthermore, CNAB provides necessary capabilities for delivering multi-container applications in disconnected (airgapped) environments.

CNAB is not a platform-specific tool. While it uses *containers* for encapsulating installation logic, it is unopinionated about what cloud environment the containers run in. CNAB developers can bundle applications targeting environments spanning IaaS (like OpenStack or Azure), container orchestrators (like Kubernetes or Nomad), container runtimes (like local Docker or ACI), and cloud platform services (like object storage or Database as a Service). CNAB can also be used for packaging other distributed applications, such as IoT or edge computing.

## Summary

The current distributed computing landscape involves a combination of executable units and supporting API-based services. Executable units include Virtual Machines (VMs), Containers (e.g. Docker and OCI), Functions-as-a-Service (FaaS), and higher-level Platforms-as-a-Service (PaaS). Along with these executable units, many managed cloud services (from load balancers to databases) are provisioned and interconnected via REST and similar network-accessible APIs. Our goal is to provide a packaging format that can give application providers and developers a way of installing a multi-component application into a distributed computing environment, supporting all of the above types.

The CNAB format is a packaging format for a broad range of distributed applications. A bundle is comprised of a [bundle definition](101-bundle-json.md) and at least one [invocation image](102-invocation-image.md). The invocation image's job is to install zero or more components into the host environment. Such components MAY include (but are not limited to) containers, functions, VMs, IaaS and PaaS layers, and service frameworks.

At run time, the invocation image contains a standardized filesystem layout where metadata, installation data, and the bundle definition are stored in predictable places. A _run tool_ is the executable entry point into a CNAB bundle. Parameterization and credentialing allow injection of configuration data into the invocation image. _Actions_ are sent to the `run` command via environment variables; actions determine whether a bundle is to be installed, upgraded, or uninstalled.

### Bundle Definition

The bundle definition contains the following information:

- Identifying information about the bundle: name and version
- Metadata about the bundle, including: maintainer, description, and keywords
- Information about locating and running invocation image(s)
- A list of executable images that this bundle will install
- A list of user-overridable parameters that this package recognizes
- A list of credential paths or environment variables that this bundle requires to execute
- A list of expected outputs from the invocation image
- A list of schema definitions used to verify the structure and values of the parameters and outputs

The canonical encoding of a bundle definition is a single JSON-formatted file, which MUST be encoded as a Canonical JSON Object stored in a `bundle.json` file, as defined in [the bundle file definition](101-bundle-json.md).

### Thick vs Thin Bundles

As described above, a bundle definition will include information about at least one image (an invocation image). Bundles can decribe their images one of two ways:

* A _thin bundle_ contains references to its images and does *not* contain the images themselves.
* A _thick bundle_ contains an encoded representation of all the invocation and executable images.

In either case, CNAB has the same schema, and this spec refers to this file as the "bundle definition" (or occasionally the "bundle file").

Two adjectives used to describe bundles are _complete_ and _well formed_.

* A thick bundle is _complete_ if all the components are present.
* A thin bundle is _complete_ if all of the references are resolvable.
* A bundle is _well formed_ if its definition follows the CNAB schema and the images are in the correct formats.

Because determining if a bundle is well formed requires checking the format of each image binary, a bundle must first be complete before it can be known to be well formed. The completeness of a thin bundle may vary depending on network access.

### Cryptography

Bundles use cryptographic verification on multiple levels. First, images (Docker, OCI, VM) are digested and the digests are embedded in the bundle definition. [TODO: why are images digested?] Second, the bundle definition is signed a public/private key system to ensure that it is not tampered with between when it is created and when the CNAB runtime receives it. A signed bundle file is named `bundle.cnab`.

Because the goal of signing a bundle definition is to prevent changes (making the bundle immutable), all image references in signed bundle definitions must have a content digest. A bundle is considered _secure_ if the bundle definition contains the correct content digests for all images and the bundle definition is cryptographically signed.

### Invocation Images

This document describes the invocation images, including the file system layout and a functional description of how an invocation image is installed.

Invocation images allow limited configuration as defined in two places in the bundle definition:

- A bundle definition MAY declare zero or more configurable parameters. User-supplied parameters are injected into the invocation image. Parameters MAY be stored.
- A bundle definition MAY declare zero or more credential requirements. This indicates which credentials MUST be passed into the invocation image in order for the invocation image to correctly authenticate to the services used by the bundle. Credentials are injected into the invocation image, but they MUST NOT be stored.


### Key Terms

- Application: the functional unit composed by the components described in a bundle. This MAY be comprised of a mixture of containers, VMs, IaaS and PaaS definitions, and other services, as well as instructions for orchestrators and service frameworks.
- Bundle: the collection of CNAB data and metadata necessary for installing an application on the designated cloud services.
- Bundle definition: the information about a bundle, its parameters, credentials, images, and usage.
- `bundle.json`: the unsigned JSON-encoded representation of a bundle definition.
- `bundle.cnab`: the signed JSON-encoded representation of a bundle definition.
- Image: used generically, a container image (e.g. OCI images) or a VM image.
- Invocation Image: the image that contains the bootstrapping and installation logic for the bundle.

When referencing tooling, we use the following terms:

- `CNAB runtime` or `runtime`: a program capable of reading a CNAB bundle and executing it
- `CNAB builder` or `builder`: a program that can assemble a CNAB bundle
- `bundle tooling`: programs or tooling that generate CNAB bundle contents

Individual tools may meet more than one of the definitions above, but we have chosen to separate them in order to offer guidance such as:

> A runtime MUST support the 'install', 'upgrade', and 'uninstall' actions, while bundle tooling MAY choose not to implement 'upgrade'.

## References To Key Sections

The following subsections define the components of CNAB:

- [The bundle.json File](101-bundle-json.md)
- [The Invocation Image Format](102-invocation-image.md)
- [The Bundle Runtime](103-bundle-runtime.md)

The process for standardization is described in an appendix:

- [Specification Process](901-process.md)

## History

- The `bundle.cnab` is now the name of a signed `bundle.json`.
- The `bundle.json` is now a stand-alone artifact, not part of the invocation image.
- The initial draft of the spec included a `manifest.json`, a `ui.json` and a `parameters.json`. The `bundle.json` is now the only metadata file, containing what was formerly spread across those three.
- The top-level `/cnab` directory was added to the bundle format due to conflicts with file hierarchy.
- The signal handling method was discarded after early research showed its limitations. The replacement uses environment variables to trigger actions.
- The `bundle.json` is now mounted in the invocation image at `/cnab/bundle.json`.
- The generic action `run` has been replaced by specific actions: `install`, `uninstall`, `upgrade`.
- The `status` action has been removed.
- Registries, security, and claims have all be moved to separate specifications.

Next section: [The bundle.json definition](101-bundle-json.md)
