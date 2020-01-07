---
title: The Bundle Runtime
weight: 103
---

# The Bundle Runtime

This section describes how the invocation image is executed, and how data is injected into the image.

The [Invocation Image definition](102-invocation-image.md) specifies the layout of a CNAB invocation image. This section focuses on how the image is executed, with the goal of managing a cloud application.

## The Run Tool (Main Entry Point)

The main entry point of a CNAB bundle MUST be located at `/cnab/app/run`. When a compliant CNAB runtime executes a bundle, it MUST execute the `/cnab/app/run` tool. In addition, images used as invocation images SHOULD also default to running `/cnab/app/run`. For example, a `Dockerfile`'s `exec` array SHOULD point to this entry point.

> A fixed location for the `run` tool is mandated because not all image formats provide an equivalent method for starting an application. A client implementation of CNAB MAY access the image and directly execute the path `/cnab/app/run`. It is also permissible, given tooling constraints, to set the default entry point to a different path.

The run tool MUST observe standard conventions for executing, exiting, and writing output. On POSIX-based systems, these are:

- The execution mode bit (`x`) MUST be set on the run tool
- Exit codes: Exit code 0 is reserved for the case where the run tool exits with no errors. Non-zero exit codes are considered to be error states. These are interpreted according to [the Open Base Specification](http://pubs.opengroup.org/onlinepubs/9699919799//utilities/V3_chap02.html#tag_18_08_02)
- The special output stream STDERR should be used to write error text

### Bundle Definition

The bundle definition is made accessible from inside the invocation image in order to allow the run tool to reference information in the file. The `bundle.json` MUST be mounted to `/cnab/bundle.json`.

### Injecting Data Into the Invocation Image

CNAB allows injecting data into the invocation image in two ways:

- Environment Variables: This is the preferred method. In this method, data is encoded as a string and placed into the the environment with an associated name.
- Files: Additional files MAY be injected _at known points_ into the invocation image via credentials or parameters.

The spec does not define or constrain any network interactions between the invocation image and external services or sources.

### Environment Variables

When executing an invocation image, a CNAB runtime MUST provide the following three environment variables to `/cnab/app/run`:

```
CNAB_INSTALLATION_NAME=my_installation
CNAB_BUNDLE_NAME=helloworld
CNAB_ACTION=install
```

The _installation name_ is the name of the _instance of_ this application. The value of `CNAB_INSTALLATION_NAME` MUST be the installation name. Consider the situation where an application ("wordpress") is installed multiple times into the same cloud. Each installation MUST have a unique installation name, even though they will be installing the same CNAB bundle. Installation names MUST consist of Graph Unicode characters and MAY be user-readable. The Unicode Graphic characters include letters, marks, numbers, punctuation, symbols, and spaces, from categories L, M, N, P, S, Zs.

The _bundle name_ is the name of the bundle (as represented in `bundle.json`'s `name` field). The specification of this field is in the [bundle definition](101-bundle-json.md). The value of `CNAB_BUNDLE_NAME` MUST be set to the bundle name.

The _action_ is the action name. It MUST be either one of the built-in actions or one of the actions named in the `actions` portion of the bundle descriptor. `CNAB_ACTION` MUST be set to the action name.

### The CNAB Revision Variable

A `CNAB_REVISION` SHOULD be passed into an install operation, and MUST be passed into `upgrade` and `uninstall`, where this is a _unique string_ indicating the current "version" of the _installation_. For example, if the `my_installation` installation is upgraded twice (changing only the parameters), three `CNAB_REVISIONS` should be generated (1. install, 2. upgrade, 3. upgrade).

Revisions are regenerated on destructive operations so that one installation may be tracked over various revisions. CNAB Runtimes MUST generate a new `CNAB_REVISION` for every `install`, `upgrade`, or `uninstall` action. That is, if an application is installed once, upgraded twice, and then uninstalled, four revisions must be generated. For additionally defined targets, a new `CNAB_REVISION` MUST be generated if the target is labeled `"modifies": true`, and MUST NOT be generated if `"modifies": false`.

A `CNAB_REVISION` SHOULD be a [ULID](https://github.com/ulid/spec).

A `CNAB_LAST_REVISION` SHOULD be provided during `upgrade` and `uninstall` operations. It MAY be provided during actions specified in the bundle descriptor. When provided, it MUST be set to the revision from the previous operation. If no previous revision exist, this SHOULD be set to the empty string (`""`). (It SHOULD NOT be set to `0`, as is sometimes the practice in UNIX programming, as `0` is considered a possible, though undesirable, revision ID)

### Parameters as Variables

As specified in the `bundle.json`, some parameters MAY be injected into the environment as environment variables.

A runtime MAY provide other `CNAB_`-prefixed variables. Parameters and credentials SHOULD NOT provide `CNAB_`-prefixed variables.

### Mounting Files

Credentials and parameters MAY be mounted as files within the image's runtime filesystem. This definition does not specify how files are to be attached to an image. However, it specifies the conditions under which the files appear.

Files MUST be attached to the invocation image before the image's `/cnab/app/run` tool is executed. Files MUST NOT be attached to the image when the image is built. That is, files MUST NOT be part of the image itself. This would cause a security violation. Files SHOULD be destroyed immediately following the exit of the invocation image, though secure at-rest encryption MAY be a viable alternative.

### Executing the Run Tool (CNAB Actions)

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

This simple example executes Helm, installing the WordPress chart with the default settings if `install` is sent, or deleting the installation if `uninstall` is sent.

An implementation of a CNAB runtime MUST support sending the following actions to an invocation image:

- `install`
- `upgrade`
- `uninstall`

Invocation images SHOULD implement `install` and `uninstall`. If one of these REQUIRED actions is not implemented, an invocation image MUST NOT generate an error (though it MAY generate a warning). Implementations MAY map the same underlying operations to multiple actions (example: `install` and `upgrade` MAY perform the same action). The runtime MUST NOT perform a [bundle version](101-bundle-json.md#name-and-version-identifying-metadata) comparison when executing an action against an existing installation but the invocation image MAY return an error if the version transition is not supported.

In addition to the default actions, CNAB runtimes MAY support custom actions (as defined in [the bundle definition](101-bundle-json.md)). Any invocation image whose accompanying bundle definition specifies custom actions SHOULD implement those custom actions. A CNAB runtime MAY exit with an error if a custom action is declared in the bundle definition, but cannot be executed by the invocation image.

A bundle MUST exit with an error if the action is executed, but fails to run to completion. A CNAB runtime MUST issue an error if a bundle issues an error. And an error MUST NOT be issued if one of the three built-in actions is requested, but not present in the bundle. Errors are reserved for cases where something has gone wrong.

In the event of an an error, the installation state MUST be considered as undefined. A subsequent execution of the same action or another action MAY resolve the installation state (example: a failed `install` action MAY be fixed by executing the `upgrade` action, a failed `upgrade` action MAY be fixed by executing the `upgrade` action again). A subsequent execution of the `uninstall` action SHOULD resolve the installation state.

## Setting Parameter Values

A CNAB `bundle.json` file MAY specify zero or more parameters whose values MAY be specified by a user.

If the `destination` field contains a key named `env`, values MUST be passed into the container as environment variables, where the value of the `env` field is the name of the environment variable.

```json
{
  "definitions": {
    "greeting": {
      "default": "hello",
      "type": "string"
    }
  },
  "parameters": {
    "greeting": {
      "definition": "greeting",
      "description": "this will be in $GREETING",
      "destination": {
        "env": "GREETING"
      }
    }
  }
}
```

By default (if no override value is provided by the CNAB runtime), the above will set `GREETING=hello`. If the runtime specifies a value `salutations`, then the environment variable would be set to `GREETING=salutations`.

The parameter value is evaluated thus:

- If the CNAB runtime provides a value, that value MAY be sanitized, then validated (as described below), then injected as the parameter value. In the event that sanitization or validation fail, the runtime SHOULD return an error and discontinue the action.
- If the parameter is marked `required` and a value is not supplied, the CNAB Runtime MUST produce an error and discontinue action.
- If the CNAB runtime does not provide a value, but `default` is set, then the default value MUST be used.
- If the parameter is marked `required` and `default` is set, then the requirement is satisfied by the runtime-provided default.
- If no value is provided and `default` is unset, the runtime MUST set the value to an empty string (""), regardless of type.
- Values are encoded as JSON strings.

> Setting the value of other types to a default value based on type, e.g. Boolean to `false` or integer to `0`, is considered _incorrect behavior_. Setting the value to `null`, `nil`, or a related character string is also considered incorrect.

In the case where the `destination` object has a `path` field, the CNAB runtime MUST create a file at that path. The file MUST have sufficient permissions that the effective user ID of the image can read the contents of the file. And the contents of the file MUST be the parameter value (calculated according to the rules above).

```json
{
  "definitions": {
    "greeting": {
      "default": "hello",
      "type": "string"
    }
  },
  "parameters": {
    "greeting": {
      "definition": "greeting",
      "description": "this will be in /var/run/greeting.txt",
      "destination": {
        "path": "/var/run/greeting.txt"
      }
    }
  }
}
```

In the example above, the CNAB runtime creates a file at `/var/run/greeting.txt` whose content (if not overridden) is `hello`. If an empty string is provided as the parameter value, the file must still be created.

A `path` MUST be absolute. But in the event that a CNAB runtime receives a relative path, it MUST treat the file as if the root path were prepended. Thus `var/run/greeting.txt` is treated (on Linux/UNIX) as `/var/run/greeting.txt`. In the cases where operating system pathing types differ, a CNAB runtime MAY freely translate between absolute pathing structures. `c:\foo.txt`, when passed to a Linux/UNIX system, MAY be translated to `/foo.txt`. In this way, multiple invocation images may share parameters regardless of the underlying OS.

If `destination` contains both a `path` and an `env`, the CNAB runtime MUST provide both.

### Validating Parameters

The validation of user-supplied values MUST happen outside of the CNAB bundle. Implementations of CNAB bundle tools MUST validate user-supplied parameter values against the named schema in the `definitions` section of a `bundle.json` before injecting them into the image. The outcome of successful validation MUST be the collection containing all parameters where either the user has supplied a value (that has been validated) or the name definition in the `definitions` section of `bundles.json` contains a `default`.

The resulting calculated values are injected into the bundle before the bundle's `run` is executed (and also in such a way that the `run` has access to these variables.) This works analogously to `CNAB_ACTION` and `CNAB_INSTALLATION_NAME`.

Resolution of conflicts in environment variable names is discussed in [the Bundle.json Description](101-bundle-json.md).

## Credential Files

Credentials MAY be supplied as files on the file system. In such cases, the following rules obtain:

- If a file is specified in the `bundle.json` credentials section, but is not present on the file system, the run tool MAY cause a fatal error
- If a file is NOT specified in the `bundle.json`, and is not present, the run tool SHOULD NOT cause an error (though it MAY emit a warning)
- If a file is present, but not correctly formatted, the run tool MAY cause a fatal error
- If a file's permissions or metadata is incorrect, the run tool MAY try to remediate (e.g. run `chmod`), or MAY cause a fatal error
- The run tool MAY modify credential files. Consequently, any runtime implementation MUST ensure that credentials changed inside of the invocation image will not result in modifications to the source.

## <a name="relocation-mapping">Image Relocation</a>

Images referenced by a CNAB bundle MAY be relocated, for example by copying them to a private registry. A _relocation mapping_ is a JSON map
of original image references to relocated image references. The purpose of a relocation mapping is to enable an invocation image to substitute relocated image references for their original values.

The relocation mapping MUST include in its keys all the image references defined by the CNAB bundle.

Any image references defined by a CNAB bundle which are semantically equivalent MUST be included as separate entries in the map and MUST map to values which are semantically equivalent to each other. For example, "ubuntu" and "library/ubuntu" are semantically equivalent. On the other hand, image references which differ only by tag and/or digest are not semantically equivalent (even though they _could_ refer to the same image).

At runtime a relocation mapping MAY be mounted in the invocation image's container as file `/cnab/app/relocation-mapping.json`. If the file is not mounted, this indicates that images have not been relocated.

For example, if a CNAB bundle with an image `gabrtv/microservice@sha256:cca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120687` and an invocation image `technosophos/helloworld:0.1.0` is relocated to a private registry `my.registry`, a mapping like the following would be mounted as the file `/cnab/app/relocation-mapping.json`:

```json
{
  "gabrtv/microservice@sha256:cca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120687": "my.registry/microservice@sha256:cca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120687",
  "technosophos/helloworld:0.1.0": "my.registry/helloworld:0.1.0"
}
```

Source: [103.01-relocation-mapping.json](examples/103.01-relocation-mapping.json)

The run tool MAY use this file to modify its behavior. For example, a run tool MAY substitute image references using the mapping in this file.
