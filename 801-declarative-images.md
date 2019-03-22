---
title: Declarative Invocation Images
weight: 801
---

# Declarative Invocation Images

This section is non-normative. A CNAB implementation MAY implement this portion to be considered conforming. However, this section maps out what the authors consider a best practice.

## Declarative Infrastructure and Declarative Installers

In declarative models, authors express the desired entities to be produced during an action. Declarative infrastructure, for example, is the practice of declaring which objects should exist in an infrastructure layer. Likewise, declarative installers provide authors with tools for describing what the installed program should look like. It is then left to the tooling to realize the authors' expression.

This section of the documentation explains how CNAB can be used to create declarative CNAB bundles by leveraging _invocation image middleware_.

## Invocation Image Middleware

As described in [the Invocation Image definition](102-invocation-image.md) and [the Bundle Runtime definition](103-bundle-runtime.md), the invocation image is responsible for executing an _action_ (install, upgrade, uninstall) inside of the image.

The aforementioned documents show examples of building CNAB bundles from non-CNAB base images. This document introduces the idea of building CNAB bundles from _middleware invocation images_, which provide underlying CNAB functionality that can be used to simplify the construction of a CNAB bundle.

Specifically, this section exhibits how a middleware package can provide the necessary tooling to make CNAB bundle construction a declarative matter.

By providing a run tool (`/cnab/app/run`), a middleware image can remove the necessity to write the imperative portions of a CNAB bundle, essentially allowing construction of declarative CNAB bundles. In this model, the middleware image provides the tooling necessary for handling CNAB actions. Images layered on top of this middleware merely need to describe what entities are being installed, uninstalled, or upgraded.

Here is an example of a declarative CNAB bundle that uses the Azure Resource Manager (ARM) templates to orchestrate an installation of an Azure application:

```text
arm-aci
└── cnab
    ├── build/Dockerfile
    ├── bundle.json
    └── app
        └── arm
            ├── parameters.json
            └── template.json
```

Note that this bundle is composed only of the following:

- Dockerfile: The Dockerfile
- bundle.json: the bundle file
- template.json: an ARM template
- parameters.json: an ARM parameters file

The `Dockerfile` begins by importing a CNAB middleware image:

```Dockerfile
FROM cnab/armbase:0.1.0

COPY Dockerfile /cnab
COPY cnab/app/arm /cnab/app/arm
COPY bundle.json /cnab/bundle.json
```

The `cnab/armbase` middleware provides the tools necessary for executing Azure Resource Manager (ARM) templates. Consequently, the present CNAB bundle MAY have imperative components.

The middleware image (`cnab/armbase`) contains tooling that looks in predefined locations for ARM templates, and understands how to install, upgrade, and uninstall those resources.

## Why Is This Non-normative?

While declarative invocation images are considered the best practice, they are non-normative because CNAB does not require specific images to be used as base images. The CNAB definition is focused on describing the conditions under which a bundle MAY be correctly packaged and executed. We have chosen, however, to not prescribe the shape of the CNAB executable.

Next section: [Credential Sets](802-credential-sets.md)
