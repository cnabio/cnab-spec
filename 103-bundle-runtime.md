# The Bundle Runtime

This section describes how the invocation image is executed, and how data is injected into the image.

The [Invocation Image definition](102-invocation-image.md) specifies the layout of a CNAB invocation image. This section focuses on how the image is executed, with the goal of managing a cloud application.

## The Run Tool (Main Entry Point)

The main entry point of a CNAB bundle _must_ be located at `/cnab/app/run`. When a compliant driver executes a CNAB bundle, it _must_ execute the `/cnab/app/run` tool. In addition, images used as invocation images _should_ also default to running `/cnab/app/run`. For example, a `Dockerfile`'s `exec` array must point to this entry point.

> A fixed location for the `run` tool is mandated because not all image formats provide an equivalent method for starting an application. A client implementation of CNAB may access the image and directly execute the path `/cnab/app/run`. It is also permissible, given tooling constraints, to set the default entry point to a different path.

The run tool must observe standard conventions for executing, exiting, and writing output. On POSIX-based systems, these are:

- The execution mode bit (`x`) must be set on the run tool
- Exit codes: Exit code 0 is reserved for the case where the run tool exits with no errors. Non-zero exit codes are considered to be error states. These are interpreted according to [the Open Base Specification](http://pubs.opengroup.org/onlinepubs/9699919799//utilities/V3_chap02.html#tag_18_08_02)
- The special output stream STDERR should be used for outputing error text

### Injecting Data Into the Invocation Image

CNAB allows injecting data into the invocation image in two ways:

- Environment Variables: This is the preferred method. In this method, data is encoded as a string and placed into the the environment with an associated name.
- Files: Additional files may be injected _at known points_ into the invocation image. In the current specification, only credentials may be injected this way.

The spec does not define or constrain any network interactions between the invocation image and external services or sources.

### Environment Variables

When executing an invocation image, a CNAB runtime _must_ provide the following three environment variables to `/cnab/app/run`:

```
CNAB_INSTALLATION_NAME=my_installation
CNAB_BUNDLE_NAME=helloworld
CNAB_ACTION=install
```

The _installation name_ is the name of the _instance of_ this application. Consider the situation where an application ("wordpress") is installed multiple times into the same cloud. Each installation _must_ have a unique installation name, even though they will be installing the same CNAB bundle. Instance names may consist of Graph Unicode characters and need not be user-readable. The Unicode Graphic characters include letters, marks, numbers, punctuation, symbols, and spaces, from categories L, M, N, P, S, Zs.

The _bundle name_ is the name of the bundle (as represented in `bundle.json`'s `name` field). The specification of this field is in the [101-bundle-json.md](bundle definition).

The _action_ is one of the action verbs defined in the section below.

Optionally, `CNAB_REVISION` _may_ be passed, where this is a _unique string_ indicating the current "version" of the _installation_. For example, if the `my_installation` installation is upgraded twice (changing only the parameters), three `CNAB_REVISIONS` should be generated (1. install, 2. upgrade, 3. upgrade). See [the Claims definition](104-claims.md) for details on revision ids. That `status` action _must not_ increment the revision.

### Mounting Files

Credentials may be mounted as files within the image's runtime filesystem. This definition does not specify how files are to be attached to an image. However, it specifies the conditions under which the files appear.

Files _must_ be attached to the invocation image before the image's `/cnab/app/run` tool is executed. Files _must not_ be attached to the image when the image is built. That is, files _must not_ be part of the image itself. This would cause a security violation. Files _should_ be destroyed immediately following the exit of the invocation image, though secure at-rest encryption may be a viable alternative.

### Executing the Run Tool

The environment will provide the name of the current installation as `$CNAB_INSTALLATION_NAME` and the name of the action will be passed as `$CNAB_ACTION`.

Example:

```bash
#!/bin/bash
action=$CNAB_ACTION

if [[ action == "install" ]]; then
  helm install example-stable/wordpress -n $CNAB_INSTALLATION_NAME
elif [[ action == "uninstall" ]]; then
  helm delete $CNAB_INSTALLATION_NAME
fi
```

This simple example executes Helm, installing the Wordpress chart with the default settings if `install` is sent, or deleting the installation if `delete` is sent.

None of the actions are required to be implemented. Bundles _ought not_ return errors simply because an action is not implemented. Errors are reserved for cases where the bundle's action failed to run correctly.

## Overriding Parameter Values

A CNAB `bundle.json` file may specify zero or more parameters whose values may be specified by a user.

Values that are passed into the container are passed in as environment variables, where each environment variable begins with the prefix `CNAB_P_` and to which the uppercased parameter name is appended. For example `backend_port` will be exposed inside the container as `CNAB_P_BACKEND_PORT`, and thus can be accessed inside of the `run` script:

```bash
#!/bin/sh

echo $CNAB_P_BACKEND_PORT
```

> Implementations _may_ also assign another variable to the same value. In other words, every parameter _must_ have a `CNAB_P_` version, and _may_ be duplicated under different names.

The validation of user-supplied values _must_ happen outside of the CNAB bundle. Implementations of CNAB bundle tools _must_ validate user-supplied values against the `parameters` section of a `bundle.json` before injecting them into the image. The outcome of successful validation _must_ be the collection containing all parameters where either the user has supplied a value (that has been validated) or the `parameters` section of `bundles.json` contains a `defaultValue`.

The resulting calculated values are injected into the bundle before the bundle's `run` is executed (and also in such a way that the `run` has access to these variables.) This works analogously to `CNAB_ACTION` and `CNAB_INSTALLATION_NAME`.

## Credential Files

Credentials may be supplied as files on the file system. In such cases, the following rules obtain:

- If a file is specified in the `bundle.json` credentials section, but is not present on the file system, the run tool _may_ cause a fatal error
- If a file is NOT specified in the `bundle.json`, and is not present, the run tool _should not_ cause an error (though it may emit a warning)
- If a file is present, but not correctly formatted, the run tool _may_ cause a fatal error
- If a file's permissions or metadata is incorrect, the run tool _may_ try to remediate (e.g. run `chmod`), or _may_ cause a fatal error
- The run tool _may_ modify credential files. Consequently, any runtime implementation _must_ ensure that credentials changed inside of the invocation image will not result in modifications to the source.

Next Section: [The claims definition](104-claims.md)
