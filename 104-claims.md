# Claims: Tracking an Installation

A _claim_ (or _claim receipt_) is a record of a CNAB installation. This document describes how the claim system works.

CNAB implementations MAY implement claims as an external storage mechanism. However, they MUST inject information into an invocation image as explained in this document.

## Concepts of Package Management

A _package_ is a discrete data chunk that can be moved from location to location, and can be unpacked and installed onto a system. Typically, a package contains an application or application description. All package managers provide some explicit definition of a package and a package format.

When a package is installed, the contents of a package are extracted and placed into the appropriate spaces on the target system, thus becoming an _installation_ (or _instance_) of the package.

There are three core feature categories of a package manager system:

- It can _install_ packages (initially put something onto a system)
- It can _query_ installations (to see what is installed)
- It can _upgrade_ and _delete_ packages (in other words, it can perform additional mutations on an existing installation)

Package managers provide a wealth of other features, but the above are standard across all package managers. (For example, most package managers also provide a way to query what packages are available for installation.)

This proposal explains how CNAB records are generated such that continuity can be established across applications. In other words, this describes how CNAB bundles can be treated analogously to traditional packages.

## Managing State

Fundamentally, package managers provide a state management layer to keep records of what was installed. For example, [homebrew](http://homebrew.sh), a popular macOS package manager, stores records for all installed software in `/usr/local/Cellar`. Helm, the package manager for Kubernetes, stores state records in Kubernetes ConfigMaps located in the system namespace. The Debian Apt system stores state in `/var/run`. In all of these cases, the stored state allows the package managing system to be able to answer (quickly) the question of whether a given package is installed.

```console
$ brew info cscope
cscope: stable 15.8b (bottled)
Tool for browsing source code
https://cscope.sourceforge.io/
/usr/local/Cellar/cscope/15.8b (10 files, 714.2KB) *
  Poured from bottle on 2017-05-15 at 09:24:58
From: https://github.com/Homebrew/homebrew-core/blob/master/Formula/cscope.rb
```

CNAB does not define where or how records are stored, nor how these records may be used by an implementation. However, it does describe how a CNAB-based system MUST emit the record to an invocation image, and provides some guidance on maintaining integrity of the system.

This is done so that implementors can standardize on a way of relating a release claim (the record of a release) to release operations like `install`, `upgrade`, or `delete`. This, in turn, is necessary if CNAB bundles are expected to be executable by different implementations.

### Anatomy of a Claim

While implementors are not REQUIRED to implement claims, this is that standard format for claims-supporting systems.

The CNAB claim is defined as a JSON document.

```json
{
  "name": "hellohelm",
  "revision": "01CP6XM0KVB9V1BQDZ9NK8VP29",
  "created": "2018-08-30T20:39:55.549002887-06:00",
  "modified": "2018-08-30T20:39:59.611068556-06:00",
  "bundle": {
    "name": "hellohelm",
    "version": "0.1.0",
    "invocationImages": [
      {
        "imageType": "docker",
        "image": "technosophos/demo2:0.2.0"
      }
    ],
    "images": [],
    "parameters": {},
    "credentials": {}
  },
  "result": {
    "message": "",
    "action": "install",
    "status": "success"
  },
  "parameters": {}
}
```

- name: The name of the _installation_. This can be automatically generated, though humans may need to interact with it. It MUST be unique within the installation environment, though that constraint MUST be imposed externally. Elsewhere, this field is referenced as the _installation name_.
- revision: An [ULID](https://github.com/ulid/spec) that MUST change each time the release is modified.
- bundle: The bundle, as defined in [the Bundle Definition](101-bundle-json.md)
- created: A timestamp indicating when this release claim was first created. This MUST not be changed after initial creation.
- updated: A timestamp indicating the last time this release claim was modified
- result: The outcome of the bundle's last action (e.g. if action is install, this indicates the outcome of the installation.). It is an object with the following fields:
  - message: A human-readable string that communicates the outcome. Error messages may be included in `failure` conditions.
  - action: Indicates the action that the current bundle is in. Valid actions are:
    - install
    - upgrade
    - delete
    - downgrade
    - status
    - unknown
  - status: Indicates the status of the last phase transition. Valid statuses are:
    - success: completed successfully
    - failure: failed before completion
    - underway: in progress. This should only be used if the invocation container MUST exit before it can determine whether all operations are complete. Note that `underway` is a _long term status_ that indicates that the installation's final state cannot be determined by the system. For this reason, it should be avoided.
    - unknown: the state is unknown. This is an error condition.
- parameters: Key/value pairs that were passed in during the operation. These are stored so that the operation can be re-run. Some implementations may choose not to store these for security or portability reasons.

> Note that credential data is _never_ stored in a claim. For this reason, a claim is not considered _trivially repeatable_. Credentials MUST be re-supplied.

TODO: What is the best timestamp format to use? Does JSON have a preference?

### ULIDs for Revisions

ULIDs have two properties that are desirable:

- High probability of [uniqueness](https://github.com/ulid/javascript)
- Sortable by time. The first 48 bits contain a timestamp

Compared to a monotonic increment, this has strong advantages when it cannot be assumed that only one actor will be acting upon the CNAB claim record. While other unique IDs are not meaningfully sortable, ULIDs are. Thus, even unordered claim storage records can be sorted.

### Parameters

The parameter data stored in a claim data is _the resolved key/value pairs_ that result from the following transformation:

- The values supplied by the user are validated by the rules specified in the `bundle.json` file
- The output of this operation is a set of key/value pairs in which:
  - Valid user-supplied values are presented
  - Default values are supplied for all parameters where `defaultValue` is provided and no user-supplied value overrides this

### How is the Claim Used

The claim is used to inform any CNAB tooling about how to address a particular installation. For example, given the claim record, a package manager that implements CNAB should be able to:

- List the _names_ of the installations, given a _bundle name_
- Given an installation's _name_, return the _bundle info_ that is installed under that name
- Given an installation _name_ and a _bundle_, generate a _bundle info_.
  - This is accompanied by running the `install` path in the bundle
- Given an installation's _name_, replace the _bundle info_ with updated _bundle info_, and update the revision with a new ULID, and the modified timestamp with the current time. This is an upgrade operation.
  - This is accompanied by running the `upgrade` path in the bundle
- Given an installation's name, mark the claim as deleted.
  - This is accompanied by running the `uninstall` path in the bundle
  - XXX: Do we want to allow the implementing system to remove the claim from its database (e.g. helm delete --purge) or remain silent on this matter?

To satisfy these requirements, implementations of a CNAB package manager are expected to be able to store and retrieve state information. However, note that nothing in the CNAB specification tells _how or where_ this state information is to be stored. It is _not a requirement_ to store that state information inside of the invocation image. (In fact, this is discouraged.)

## Injecting Claim Data into an Invocation Image

Complaint CNAB implementations MUST conform to this section.

The claim is produced outside of the CNAB package. The following claim data is injected
into the invocation container at runtime:

- `$CNAB_INSTALLATION_NAME`: The value of `claim.name`.
- `$CNAB_BUNDLE_NAME`: The name of the present bundle.
- `$CNAB_ACTION`: The action to be performed (install, upgrade, ...)
- `$CNAB_REVISION`: The ULID for the present release revision. (On upgrade, this is the _new_ revision)

> Credential data, which is also injected into the invocation image, is _not_ managed by the claim system. Rules for injecting credentials are found in [the bundle runtime definition](103-bundle-runtime.md).

The parameters passed in by the user are vetted against `parameters.json` outside of the container, and then injected into the container as environment variables of the form: `$CNAB_P_{parameterName.toUpper}="{parameterValue}"`.

For example, the parameter `hello_world` in the claim is presented to the invocation image as the environment variables `CNAB_P_HELLO_WORLD`.

## Calculating the Result

The `result` object is populated by the result of the invocation image's action. For example, consider the case where an invocation image executes an installation action. The action is represented by the following shell script, and `$CNAB_INSTALLATION_NAME` is set to `my_first_install`:

```bash
#!/bin/bash

set -eo pipefail

helm install stable/wordpress -n $CNAB_INSTALLATION_NAME > /dev/null
kubectl create pod $CNAB_INSTALLATION_NAME > /dev/null
echo "yay!"
```

(Note that the above script redirects data to `/dev/null` just to make the example easier. A production CNAB bundle might choose to include more verbose output.)

If both commands exit with code `0`, then the resulting claim will look like this:

```json
{
    "name": "my_first_install",
    "revision": "01CN530TF9Q095VTRYP1M8797C",
    "bundle": {
        "uri": "hub.docker.com/technosophos/example_cnab",
        "name": "example_cnab",
        "version": "0.1.0"
    },
    "created": "TIMESTAMP",
    "modified": "TIMESTAMP",
    "result": {
        "message": "yay!",    // From STDOUT (echo)
        "action": "install",  // Determined by the action
        "status": "success"   // if exit code == 0, success, else failure
    }
}
```

Tools that implement claims may then present `result` info to end users to show the result of running an invocation image.

## TODO

- Define how action is determined, as this is beyond merely running an executable

Next section: [signing and verifying bundles](105-signing.md)
