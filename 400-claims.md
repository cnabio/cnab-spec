---
title: Claims System
weight: 400
---

# CNAB Claims 1.0 (CNAB-Claims1)
*[Working Draft](901-process.md), Feb. 2019*

This specification describes the CNAB Claims system. This is not part of the [CNAB Core](100-CNAB.md) specification, but is a specification describing how records of CNAB installations may be formatted for storage. Implementations of CNAB Core can meet the CNAB Core specification without implementing this specification.

A _claim_ (or _claim receipt_) is a record of a CNAB installation. This specification describes how the claim system works.

The claims system is designed to satisfy the requirements of the [Bundle Runtime specification](103-bundle-runtime.md) regarding the tracking `CNAB_REVISION` and `CNAB_LAST_REVISION`. It also provides a description of how the state of a bundle installation may be represented.

This specification uses the terminology defined in [the CNAB Core specification](100-CNAB.md).

## Concepts of Package Management

A _package_ is a discrete data chunk that can be moved from location to location, and can be unpacked and installed onto a system. Typically, a package contains an application or application description. All package managers provide some explicit definition of a package and a package format.

When a package is installed, the contents of a package are extracted and placed into the appropriate spaces on the target system, thus becoming an _installation_ (or _instance_) of the package.

There are three core feature categories of a package manager system:

- It can _install_ packages (initially put something onto a system)
- It can _query_ installations (to see what is installed)
- It can _upgrade_ and _delete_ packages (in other words, it can perform additional mutations on an existing installation)

Package managers provide a wealth of other features, but the above are standard across all package managers. (For example, most package managers also provide a way to query what packages are available for installation.)

This specification explains how CNAB records are generated such that continuity can be established across applications. In other words, this describes how CNAB bundles can be treated analogously to traditional packages.

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

The CNAB Claims specification does not define where or how records are stored, nor how these records MAY be used by CNAB tooling. However, this specification describes how a CNAB-based system MUST emit the record to an invocation image, and provides some guidance on maintaining integrity of the system.

This is done so that implementors can standardize on a way of relating a release claim (the record of a release) to release operations like `install`, `upgrade`, or `uninstall`. This, in turn, is necessary if CNAB bundles are expected to be executable by different implementations.

### Anatomy of a Claim

While implementors are not REQUIRED to implement claims, this is the standard format for claims-supporting systems.

The CNAB claim is defined as a JSON document. The specification currently does not require that claims be formatted as Canonical JSON.

```json
{
  "bundle": {
    "credentials": {},
    "images": {},
    "invocationImages": [
      {
        "image": "technosophos/demo2:0.2.0",
        "imageType": "docker"
      }
    ],
    "name": "technosophos.hellohelm",
    "outputs": {},
    "parameters": {},
    "schemaVersion": "v1.0.0-WD",
    "version": "0.1.0"
  },
  "created": "2018-08-30T20:39:55.549002887-06:00",
  "modified": "2018-08-30T20:39:59.611068556-06:00",
  "name": "technosophos.hellohelm",
  "outputs": {},
  "parameters": {},
  "result": {
    "action": "install",
    "message": "",
    "status": "success"
  },
  "revision": "01CP6XM0KVB9V1BQDZ9NK8VP29"
}
```
Source: [400.01-claim.json](examples/400.01-claim.json)

- bundle: The bundle, as defined in [the Bundle Definition](101-bundle-json.md).
- created: A timestamp indicating when this release claim was first created. This MUST not be changed after initial creation.
- modified: A timestamp indicating the last time this release claim was modified.
- name: The name of the _installation_. This can be automatically generated, though humans may need to interact with it. It MUST be unique within the installation environment, though that constraint MUST be imposed externally. Elsewhere, this field is referenced as the _installation name_.
- outputs: Key/value pairs that were created by the operation. These are stored so that the user can access them after the operation completes. Some implementations MAY choose not to store these for security or portability reasons.
- parameters: Key/value pairs that were passed in during the operation. These are stored so that the operation can be re-run. Some implementations MAY choose not to store these for security or portability reasons.
- result: The outcome of the bundle's last action (e.g. if action is install, this indicates the outcome of the installation.). It is an object with the following fields:
  - action: Indicates the action that the current bundle is in. This may be any of the built-in actions (`install`, `upgrade`, `uninstall`) as well as any custom actions as defined in the bundle descriptor. The special name `unknown` MAY be used in the case where the CNAB Runtime cannot determine the action name of a claim.
  - message: A human-readable string that communicates the outcome. Error messages MAY be included in `failure` conditions.
  - status: Indicates the status of the last phase transition. Valid statuses are:
    - failure: failed before completion
    - underway: in progress. This should only be used if the invocation container MUST exit before it can determine whether all operations are complete. Note that `underway` is a _long term status_ that indicates that the installation's final state cannot be determined by the system. For this reason, it should be avoided.
    - unknown: the state is unknown. This is an error condition.
    - success: completed successfully
- revision: An [ULID](https://github.com/ulid/spec) that MUST change each time the release is modified.

> Note that credential data is _never_ stored in a claim. For this reason, a claim is not considered _trivially repeatable_. Credentials MUST be supplied on each action.

Timestamps in JSON are defined in the [ECMAScript specification](https://www.ecma-international.org/ecma-262/9.0/index.html#sec-date-time-string-format), which matches the [ISO-8601 Extended Format](https://www.iso.org/iso-8601-date-and-time-format.html).

### ULIDs for Revisions

ULIDs have two properties that are desirable:

- High probability of [uniqueness](https://github.com/ulid/javascript)
- Sortable by time. The first 48 bits contain a timestamp

Compared to a monotonic increment, ULID has strong advantages when it cannot be assumed that only one actor will be acting upon the CNAB claim record. While other unique IDs are not meaningfully sortable, ULIDs are. Thus, even unordered claim storage records can be sorted.

### Parameters

The parameter data stored in a claim data is _the resolved key/value pairs_ that result from the following transformation:

- The values supplied by the user are validated by the rules specified in the `bundle.json` file
- The output of this operation is a set of key/value pairs in which:
  - Valid user-supplied values are presented
  - Default values are supplied for all parameters where `default` is provided and no user-supplied value overrides this

### Outputs

The output data saves the contents of each output file at the end of the action. The data is a map from output name as defined in the bundle to the content of the file located at the path defined for that output.

Below, you can see an example of a claim for a bundle that included a single output file, `clientCert`. The value in the claim's output data is the contents of `/cnab/app/outputs/clientCert` at the end of the `install` action.

```json
{
  "bundle": {
    "credentials": {},
    "definitions": {
      "x509Certificate": {
         "contentEncoding": "base64",
         "contentMediaType": "application/x-x509-user-cert",
         "type": "string",
         "writeOnly": true
      }
    },
    "images": {},
    "invocationImages": [
      {
        "image": "technosophos/demo2:0.2.0",
        "imageType": "docker"
      }
    ],
    "name": "technosophos.hellohelm",
    "outputs": {
      "fields": {
        "clientCert": {
          "definition": "x509Certificate",
          "path": "/cnab/app/outputs/clientCert"
        }
      }
    },
    "parameters": {},
    "schemaVersion": "v1.0.0-WD",
    "version": "0.1.0"
  },
  "created": "2018-08-30T20:39:55.549002887-06:00",
  "modified": "2018-08-30T20:39:59.611068556-06:00",
  "name": "technosophos.hellohelm",
  "outputs": {
    "clientCert": "-----BEGIN CERTIFICATE-----\nMIIDBzCCAe+gAwIBAgIJAL2nOwEePOPvMA0GCSqGSIb3DQEBBQUAMBoxGDAWBgNV\nBAMMD3d3dy5leGFtcGxlLmNvbTAeFw0xOTA3MTUxNjQwMzdaFw0yOTA3MTIxNjQw\nMzdaMBoxGDAWBgNVBAMMD3d3dy5leGFtcGxlLmNvbTCCASIwDQYJKoZIhvcNAQEB\nBQADggEPADCCAQoCggEBAJ7779ImmmvEt+ywP8GjfpzgM57n1WZ26fBTVy+ZibiH\nhKCJuzU5vynu5M0eYCufRFQ7LG/xet1GvpBIch0U6ilZVnNDrsNUtQ03Hpen144g\nli5ldPh5Sm88ibDbi4yEIgti2JBKIuVE+iEdkIejF8DZps008TbLLoENM1VpHpUT\nCIJY657t+Xhz9GOhp1w3bVoKEoF/6psvc6IFHK8bUMq+4003VGDZe2BMlgazZPHc\n3o5CqNviajnRoo9QnLUH1qOljNMR+mkewNOkL2PRGvkCuHJrEk0qmUU7lX3iVN1J\nC7y1fax53ePXajLD+5/sQNeszVg1cIIUlXBy2Bx/F7UCAwEAAaNQME4wHQYDVR0O\nBBYEFAZ1+cZNMujQhCrtCRKfPm+NLDg0MB8GA1UdIwQYMBaAFAZ1+cZNMujQhCrt\nCRKfPm+NLDg0MAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQEFBQADggEBAAuyyYER\nT5+E+KuFbDL/COmxRWZhF/u7wIW9cC2o5LlKCUp5rRQfAVRNJlqasldM/G4Bg/uQ\n5khSYe2XNe2C3iajVkR2RlqXBvCdQuCtudhZVd4jSGWi/yI7Ub2/HyOTjZ5eG82o\nF6e3pNRCQwTw0y0orQmdh0s+UmHEVjIe8PfbdRymfeQO70EXTxncBJ5elZx8s0E9\nTPPdbl2knZmKJhwnZFKCaa4DmDA6CDa2GPz+2++DQl1rCIB3mIcxpg6wSdMA6C6l\nnBJBX5Wnxckldp8G0FNXTa1DYqjPZ5U84tkh46pFSLsbSse45xNhrMPeZDWlgHsp\nWfK01YbCWioNVGk=\n-----END CERTIFICATE-----"
  },
  "parameters": {},
  "result": {
    "action": "install",
    "message": "",
    "status": "success"
  },
  "revision": "01CP6XM0KVB9V1BQDZ9NK8VP29"
}
```

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
- Given an installation's name, return the contents of the outputs.

To satisfy these requirements, implementations of a CNAB package manager are expected to be able to store and retrieve state information. However, note that nothing in the CNAB specification tells _how or where_ this state information is to be stored. It is _not a requirement_ to store that state information inside of the invocation image. (In fact, this is discouraged.)

## Injecting Claim Data into an Invocation Image

Compliant CNAB implementations MUST conform to this section.

The claim is produced outside of the CNAB package. The following claim data is injected
into the invocation container at runtime:

- `$CNAB_INSTALLATION_NAME`: The value of `claim.name`.
- `$CNAB_BUNDLE_NAME`: The name of the present bundle.
- `$CNAB_ACTION`: The action to be performed (install, upgrade, ...)
- `$CNAB_REVISION`: The ULID for the present release revision. (On upgrade, this is the _new_ revision)

> Credential data, which is also injected into the invocation image, is _not_ managed by the claim system. Rules for injecting credentials are found in [the bundle runtime definition](103-bundle-runtime.md).

Parameters declared with an `env` key in the `destination` object MUST have their values injected as environment variables according to the name specified. Likewise, files MUST be injected if `path` is set on `destination`.

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
  "bundle": {
    "name": "technosophos.example_cnab",
    "uri": "hub.docker.com/technosophos/example_cnab",
    "version": "0.1.0"
  },
  "created": "TIMESTAMP",
  "modified": "TIMESTAMP",
  "name": "technosophos.my_first_install",
  "result": {
    "message": "yay!",    // From STDOUT (echo)
    "action": "install",  // Determined by the action
    "status": "success"   // if exit code == 0, success, else failure
  },
  "revision": "01CN530TF9Q095VTRYP1M8797C"
}
```

Tools that implement claims MAY then present `result` info to end users to show the result of running an invocation image.

## Credentials and Claims

Credential data MUST NOT be stored in claims. Credentials are identity-specific, while claims are identity-neutral.
