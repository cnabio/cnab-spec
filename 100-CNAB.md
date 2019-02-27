# Cloud Native Application Bundle Core 1.0.0 (CNAB1)
*[Working Draft](901-process.md), Nov. 2018*


The Cloud Native Application Bundle (CNAB) is a _standard packaging format_ for multi-component distributed applications. It allows packages to target different runtimes and architectures. It empowers application distributors to package applications for deployment on a wide variety of cloud platforms, providers, and services. Furthermore, it provides necessary capabilities for delivering multi-container applications in disconnected (airgapped) environments.

CNAB is not a platform-specific tool. While it uses *containers* for encapsulating installation logic, it remains un-opinionated about what cloud environment it runs in. CNAB developers can bundle applications targeting environments spanning IaaS (like OpenStack or Azure), container orchestrators (like Kubernetes or Nomad), container runtimes (like local Docker or ACI), and cloud platform services (like object storage or Database as a Service). 

CNAB can also be used for packaging other distributed applications, such as IoT or edge computing.

## Summary

The CNAB format is a packaging format for a broad range of distributed applications. It specifies a pairing of a _bundle definition_ [(`bundle.json`)](101-bundle-json.md) to define the app, and an _invocation image_ to install the app.

The _bundle definition_ is a single file that contains the following information:

- Information about the bundle, such as name, bundle version, description, and keywords
- Information about locating and running the _invocation image_ (the installer program)
- A list of user-overridable parameters that this package recognizes
- The list of executable images that this bundle will install
- A list of credential paths or environment variables that this bundle requires to execute

The canonical encoding of a bundle definition is a JSON-formatted file, which MAY be presented in either signed or unsigned format.

- If unsigned, the CNAB is encoded as a JSON Object stored in a `bundle.json` file, as defined in [the bundle file definition](101-bundle-json.md)
- If signed, the CNAB is encoded as a JSON object stored in a `bundle.cnab` file, as described in [the signature definition](105-signing.md)

In either case, CNAB has the same schema, and this spec refers to this file as the "bundle definition" (or occasionally "bundle file"). 

However, as a signed bundle definition represents an immutable bundle, all invocation images and images references must have a digest.

The bundle definition can be stored on its own, or as part of a _packed archive_, which is a CNAB bundle that includes the JSON file and exported images (including the [invocation image](102-invocation-image.md)).

- A _thin bundle_ consists of just a bundle definition.
- A _thick bundle_ consists of a packaged archive that contains both the bundle definition and an encoded representation of all of the invocation images and images.

When thin bundles are processed, the referenced content (such as invocation images and other images) are retrieved from their respective storage repositories and registries. A bundle is considered to be _well formed_ if it's definition follows the CNAB schema and the images are in the correct formats. A bundle is considered _complete_ if it is packaged as a thick bundle, and all the components are present OR if it is packaged as a thin bundle and all of the references are resolvable. Completeness is thus in some cases contingent upon external factors such as network access.

Bundles use cryptographic verification on multiple levels. Images (Docker, OCI, VM) are digested, and their cryptographic digest is then embedded into the `bundle.json`. The `bundle.json` is then signed using a public/private key system to ensure that it has not been tampered with. A signed bundle is named `bundle.cnab`. A bundle is considered _secure_ if the bundle definition contains the correct digests for all images, and the bundle definition is cryptographically signed.

Finally, this document describes a format for invocation images, including file system layout and a functional description of how an invocation image is installed.

## Approach

The current distributed computing landscape involves a combination of executable units and supporting API-based services. Executable units include Virtual Machines (VMs), Containers (e.g. Docker and OCI) and Functions-as-a-Service (FaaS), as well as higher-level PaaS services. Along with these executable units, many managed cloud services (from load balancers to databases) are provisioned and interconnected via REST (and similar network-accessible) APIs. Our goal is to provide a packaging format that can enable application providers and developers with a way of installing a multi-component application into a distributed computing environment, supporting all of the above types.

A bundle is comprised of a bundle definition and at least one _invocation image_. The invocation image's job is to install zero or more components into the host environment. Such components MAY include (but are not limited to) containers, functions, VMs, IaaS and PaaS layers, and service frameworks.

The invocation image contains a standardized filesystem layout where metadata and installation data is stored in predictable places. A _run tool_ is the executable entry point into a CNAB bundle. Parameterization and credentialing allow injection of configuration data into the invocation image. The invocation image is described in detail in [the invocation image definition](102-invocation-image.md).

_Actions_ are sent to the `run` command via environment variables. Actions determine whether a bundle is to be installed, upgraded, downgraded, or uninstalled.

Invocation images allow limited configuration, as defined in two places in the bundle definition:

- A bundle definition MAY declare zero or more configurable parameters. User-supplied parameters are injected into the invocation image. Parameters MAY be stored.
- A bundle definition MAY declare zero or more credential requirements. This indicates which credentials MUST be passed into the invocation image in order for the invocation image to correctly authenticate to the services used by the bundle. Credentials are injected into the invocation image, but they MUST NOT be stored.

Additionally, this document defines two auxiliary components of the CNAB system: claims and repositories.

A _claim_ is a record of an installation of a bundle. When a bundle is installed into a host environment, it MAY be useful to track that bundle's history. The claims system is an OPTIONAL part of the specification which describes a standard format for storing bundle processing history.

Bundles are stored in _bundle repositories_. A bundle repository is a network-accessible service for uploading and distributing CNAB objects.

### Key Terms

- Application: The functional unit composed by the components described in a bundle. This MAY be comprised of a mixture of containers, VMs, IaaS and PaaS definitions, and other services, as well as instructions for orchestrators and service frameworks.
- Bundle: the collection of CNAB data and metadata necessary for installing an application on the designated cloud services.
- Bundle definition: The information about a bundle, its parameters, credentials, images, and usage
- `bundle.json`: The unsigned JSON-encoded representation of a bundle definition.
- `bundle.cnab`: The signed JSON-encoded representation of a bundle definition.
- Claim: A record of a particular installation of a bundle.
- Image: Used generically, a container image (e.g. OCI images) or a VM image.
- Invocation Image: The image that contains the bootstrapping and installation logic for the bundle
- Registry: A storage and retrieval service for CNAB objects.

Also, when referencing tooling, we use the following terms:

- `CNAB runtime` or `runtime`: A program capable of reading a CNAB bundle and executing it
- `CNAB builder` or `builder`: A program that can assemble a CNAB bundle
- `bundle tooling`: Programs or tooling that generate CNAB bundle contents

Individual tools may meet more than one of the definitions above, but we have chosen to separate them in order to offer guidance such as:

> A runtime MUST support the 'install', 'upgrade', and 'uninstall' actions, while bundle tooling MAY choose not to implement 'upgrade'.

### The Definitions

The following subsections define the components of CNAB:

- [The bundle.json File](101-bundle-json.md)
- [The Invocation Image Format](102-invocation-image.md)
- [The Bundle Runtime](103-bundle-runtime.md)
- [The Claims System](104-claims.md)
- [Signing and Provenance](105-signing.md)

The following sections contain non-normative guidance

- [Declarative Bundles](801-declarative-images.md)
- [Credential Sets](802-credential-sets.md)
- [Base Bundles](803-base-bundles.md)
- [Storing CNAB Bundles](804-repositories.md)

The process for standardization is described in an appendix:

- [Specification Process](901-process.md)

## History

- The `bundle.cnab` is now the name of a signed `bundle.json`.
- The `bundle.json` is now a stand-alone artifact, not part of the invocation image.
- The initial draft of the spec included a `manifest.json`, a `ui.json` and a `parameters.json`. The `bundle.json` is now the only metadata file, containing what was formerly spread across those three.
- The top-level `/cnab` directory was added to the bundle format due to conflicts with file hierarchy.
- The signal handling method was discarded after early research showed its limitations. The replacement uses environment variables to trigger actions.
- The credential set and claims concepts were introduced to cover areas upon which the original spec was silent.
- The generic action `run` has been replaced by specific actions: `install`, `uninstall`, `upgrade`, `status`.

Next section: [The bundle.json definition](101-bundle-json.md)
