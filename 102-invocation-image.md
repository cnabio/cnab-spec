---
title: The Invocation Images
weight: 102
---

# The Invocation Images

The `invocationImages` section of a `bundle.json` MUST contain at least one image (the invocation image). This image MUST be formatted according to the specification laid out in the present document.
The appropriate invocation image is selected using the current driver.

When a bundle is executed, the invocation image will be retrieved (if necessary) and started. Credential, parameter and image map data is passed to it, and then its `run` tool is executed. (See [The Bundle Runtime](103-bundle-runtime.md) for details).

This section describes the layout of an invocation image.

## Components of an Invocation Image

An invocation image is composed of the following:

- A file system hierarchy following a defined pattern (below)
- A main entry point, called the _run tool_, which is an executable (often a script) responsible for translating action requests (`install`, `upgrade`,...) to a sequence of tasks
- Runtime metadata (Helm charts, Terraform templates, etc)
- The material necessary for reproducing the invocation image (`Dockerfile` and `packer.json` are two examples)

Note that the bundle definition itself is not stored inside of the invocation image.

### The File System Layout

The following exhibits the filesystem layout:

```yaml
cnab/                  # REQUIRED top-level directory
└── build/
    │   └──Dockerfile​  # OPTIONAL
└── app​                # REQUIRED
    ├── run​            # REQUIRED: This is the main entrypoint, and MUST be executable
    ├── charts​         # Example: Helm charts might go here
    │   └── azure-voting-app​
    │       ├── Chart.yaml​
    │       ├── templates​​
    │       │   └── ...
    │       └── values.yaml​
    └── sfmesh​         # Example: Service Fabric definitions might go here
        └── sfmesh-deploy.json
```

### The `/cnab` Directory

An invocation image MUST have a directory named `cnab` placed directly under the root of the file system hierarchy inside of an image.

This directory MUST have a subdirectory named `app`.

This directory MAY have any of the following:

- `build/`: A directory containing files used in the construction of this image
    - `Dockerfile`: A valid Dockerfile used for constructing this image
    - Files for the form `Dockerfile.$INFO`, where `$INFO` is a further specification of that Dockerfile (e.g. `Dockerfile.arm64`)
    - `packer.json`: A valid Packer configuration file
    - Other build-related files
- `README.txt` or `README.md`: A text file containing information about this bundle
- `LICENSE`: A text file containing the license(s) for this image

This directory MUST NOT have any files or directories not explicitly named in the present document. The `/cnab` directory is considered a reserved namespace where future CNAB revisions MAY place new files or directories.

### The `/cnab/app` Directory

The `app/` directory contains subdirectories, each of which stores configuration for a particular target environment. The `app/run` file _MUST be an executable file_ that will act as the "main" installer for this CNAB bundle. This is the only file that is REQUIRED in this directory.

The contents beneath `/cnab/app/SUBDIRECTORY` are undefined by the spec. `run` is considered the only reserved word underneath `/cnab/app/`

### The OPTIONAL `/cnab/build` Directory

The directory `/cnab/build` MAY be present within the CNAB hierarchy. This directory houses files used in the construction of the invocation image, and are provided to make it easier to rebuild the image during rewrite operations. If a `Dockerfile` was used to build the image, a `Dockerfile` SHOULD be included. Other files MAY be included.

Examples:
    - `packer.json`
    - `Dockerfile.arm32`

#### Dockerfiles for Constructing Invocation Images

This subsection is non-normative. Images may be built using any suitable tooling, but this section describes the process using a `Dockerfile`.

The `Dockerfile` used to build the invocation image MAY be stored inside of the invocation image. This is to ensure reproducibility, and in order to allow rename operations that require a rebuild. (Likewise, if a build tool like Packer is used, this tool's configuration MAY be placed in the bundle.)

This is a normal Dockerfile, and MAY derive from any base image.

Example:

```Dockerfile
FROM ubuntu:latest

COPY ./Dockerfile /cnab/Dockerfile
COPY ./bundle.json /cnab/manifest.json
COPY ./parameters.json /cnab/parameters.json

RUN curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash
RUN helm init --client-only
RUN helm repo add example-stable https://repo.example.com/stable

CMD /cnab/app/run
```

The above example installs and configures Helm inside of a base Ubuntu image. Note that there are no restrictions on what tools MAY be installed.

## The Run Tool

The run tool MUST be located at the path `/cnab/app/run`. It MUST be executable. It MUST react to the `CNAB_ACTION` provided to it.

The specification does not define what language(s) the tool must be written in, or any details about how it processes the information. However, the following is a non-normative example:

```bash
#!/bin/sh

action=$CNAB_ACTION
name=$CNAB_INSTALLATION_NAME 

case $action in
    install)
    echo "Install action"
    ;;
    uninstall)
    echo "uninstall action"
    ;;
    upgrade)
    echo "Upgrade action"
    ;;
    *)
    echo "No action for $action"
    ;;
esac
echo "Action $action complete for $name"
```

The run tool above is implemented as a shell script, and merely reacts to each given `CNAB_ACTION` by printing a message.

See [The Bundle Runtime](103-bundle-runtime.md) for a description on how this tool is used.

Next section: [The Bundle Runtime](103-bundle-runtime.md)
