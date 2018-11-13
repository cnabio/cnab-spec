# The Invocation Images

The `invocationImages` section of a `bundle.json` must contain at least one image (the invocation image). This image must be formatted according to the specification laid out in the present document.
The appropriate invocation image is selected using the current driver.

When a bundle is executed, the invocation image will be retrieved (if necessary) and started. Credential and parameter data is passed to it, and then its `run` tool is executed. (See [The Bundle Runtime](103-bundle-runtime.md) for details).

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
cnab/                  # Required top-level directory
└── Dockerfile​         # Optional
└── app​                # Required
    ├── run​            # Required: This is the main entrypoint, and must be executable
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

- `Dockerfile`: A valid Dockerfile used for constructing this image
- Files for the form `Dockerfile.$INFO`, where `$INFO` is a further specification of that Dockerfile (e.g. `Dockerfile.arm64`)
- `packer.json`: A valid Packer configuration file
- `bundle.json`: A valid bundle definition, possibly without digests
- `build/`: A directory containing other files used in the construction of this image
- `README.txt` or `README.md`: A text file containing information about this bundle
- `LICENSE`: A text file containing the license(s) for this image

This directory MUST NOT have any files or directories not explicitly named in the present document. The `/cnab` directory is considered a reserved namespace where future CNAB revisions may place new files or directories.

### The `/cnab/app` Directory

The `app/` directory contains subdirectories, each of which stores configuration for a particular target environment. The `app/run` file _must be an executable file_ that will act as the "main" installer for this CNAB bundle. This is the only file that is required in this directory.

The contents beneath `/cnab/app/SUBDIRECTORY` are undefined by the spec. `run` is considered the only reserved word underneath `/cnab/app/`

### The Optional `/cnab/build` Directory

The directory `/cnab/build` MAY be present within the CNAB hierarchy. The contents of this directory are undefined by the spec. However, the intention is to provide space for artifacts that are used during the assembly (or re-assembly) of the bundle, but are not part of the execution of a bundle action.

For example, _base bundles_ may use this directory to store utilities that derived bundles may use during bundle building.

### Base Bundles

This section of the specification is non-normative. It describes a pattern for using the contents of one bundle in another bundle as a form of extension or inheritance. While this is non-normative to the specification, facilitating this pattern was a goal of CNAB's design.

The _base bundle pattern_ is a pattern for sharing common bundle tooling across multiple bundles. Succinctly expressed, a _base bundle_ is a bundle whose contents are inherited by another bundle. The most frequent way in which is occurs is when one bundle uses a `Dockerfile` to import another bundle's invocation image in its `FROM` line.

A bundle that uses a base bundle as a source is referred to here as an _extending bundle_ or _extension bundle_, drawing from the object oriented terminology of base classes and extension classes.

For example, a base CNAB image may provide something like this:

```Dockerfile
FROM ubuntu:latest

COPY ./some-chart /cnab/app/charts/some-chart
```

If the above bundle is built as `base-bundle:latest`, then it may be referenced by other images.

The `Dockerfile` for an extending bundle can import `base-bundle:latest` as a starting point:

```Dockerfile
FROM base-bundle:latest

RUN helm inspect values /cnab/app/charts/some-chart > ./myvals.yaml &&  sed ...
```

The example above is simply intended to show how by reserving the `/cnab` directory, we can make images extensible, while not worrying about non-CNAB images putting data in places CNAB treats as special.

The mechanisms for taking a base bundle and creating an extension bundle are not formalized in the definition, as tooling may implement this pattern as fit for the tooling's domain. However, there are particular points of consideration that implementations may wish to address:

- Run tools make good candidates for reuse, and may then be well suited for base bundles
- A base bundle MAY define parameters and credentials in its bundle definition. In such cases, implementations of the base bundle pattern may wish to bubble this configuration up to the extension bundle via tooling.
- The `/cnab/build` directory may be used to store tooling in a base bundle that extension bundles can use during bundle assembly.
- Docker and other OCI tooling may provide mechanisms for creating a bundle by composing via more than one base image. The present specification does not prohibit composition, and this section of the specification is non-normative.

## Image Construction Files

Including a Dockerfile is RECOMMENDED for all images built with Docker. It is useful for reproducing a bundle. For other build tools, the build tool's definition may be included instead (e.g. `packer.json` for VM images built with Packer). Any image construction artifacts that are not explicitly allowed in the `/cnbab` directory may be placed in the `/cnab/build` directory.

The remainder of this subsection is non-normative.

The `Dockerfile` used to build the invocation image MAY be stored inside of the invocation image. This is to ensure reproducibility, and in order to allow rename operations that require a rebuild. (Likewise, if a build tool like Packer is used, this tool's configuration MAY be placed in the bundle.)

This is a normal Dockerfile, and may derive from any base image.

Example:

```Dockerfile
FROM ubuntu:latest

COPY ./Dockerfile /cnab/Dockerfile
COPY ./bundle.json /cnab/manfiest.json
COPY ./parameters.json /cnab/parameters.json

RUN curl https://raw.githubusercontent.com/kubernetes/helm/master/scripts/get | bash
RUN helm init --client-only
RUN helm repo add example-stable https://repo.example.com/stable

CMD /cnab/app/run
```

The above example installs and configures Helm inside of a base Ubuntu image. Note that there are no restrictions on what tools may be installed.

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
