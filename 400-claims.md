---
title: Claims System
weight: 400
---

# CNAB Claims 1.0 (CNAB-Claims1)
Draft, Feb. 2020

This specification describes the CNAB Claims system. This is not part of the [CNAB Core](100-CNAB.md) specification, but is a specification describing how records of CNAB installations may be formatted for storage. Implementations of CNAB Core can meet the CNAB Core specification without implementing this specification. Implementations that support claims MAY state that they comply with CNAB Claims 1.0-WD.

In CNAB, an _installation_ is a particular installed instance of a CNAB bundle. For example, if a bundle named `myApp` is installed in two places, we say there are _two installations of `myApp`_.

An installation MAY change over time, as a particular bundle is installed, then upgraded. In CNAB, each time a mutating (destructive) operation is performed (Such as `install`, `upgrade`, or custom operations that are not read-only), the claim gets a new _revision_. A revision is a unique identifier that identifies the combination of an _installation_ and a modification of that installation. For example, when a CNAB bundle is installed, the initial installation will have an initial revision ID. When that installation is upgraded, it will have a new revision ID. How revisions are used is outside of the scope of this document, but it is safe to assume that if a revision ID has changed, one or more artifacts owned by the installation has also been changed.

A _claim_ is a record of a CNAB _installation_. This specification defines the format of a claim, and describes certain necessary behaviors of a system that implements claims.

A _host environment_ is an environment, possibly shared between multiple CNAB runtimes, that provides persistence to a CNAB configuration. For example, a filesystem may contain a record of claims that two different CNAB clients share. Or a database may contain the environment that is shared by multiple tools at different locations in the network. Each of these is a CNAB host environment.

The word _claim_ was chosen because it represents the relationship between a certain CNAB host environment and the resources that were created by a CNAB runtime in that host environment. In this sense, an environment takes responsibility for those resources if and only if it can _claim_ them. A claim is an _external assertion of ownership_. That is, the claim itself is not "installed" into the host environment. It is stored separately, possibly in an entirely different location.

The claims system is designed to satisfy the requirements of the [Bundle Runtime specification](103-bundle-runtime.md) regarding the tracking `CNAB_REVISION` and `CNAB_LAST_REVISION`. It also provides a description of how the state of a bundle installation may be represented.

This specification uses the terminology defined in [the CNAB Core specification](100-CNAB.md).

## Managing State

Fundamentally, package managers provide a state management layer to keep records of what was installed. For example, [homebrew](http://homebrew.sh), a popular macOS package manager, stores records for all installed software in `/usr/local/Cellar`. Helm, the package manager for Kubernetes, stores state records in Kubernetes ConfigMaps located in the system namespace. The Debian Apt system stores state in `/var/run`. In all of these cases, the stored state allows the package managing system to be able to answer (quickly) the question of whether a given package is installed.

Example from Homebrew:

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

This is done so that implementors can standardize on a way of relating an installation claim (the record of an installation) to operations like `install`, `upgrade`, or `uninstall`. This, in turn, is necessary if CNAB bundles are expected to be executable by different implementations.

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
    "schemaVersion": "v1.0.0",
    "version": "0.1.0"
  },
  "bundleReference": "hub.example.com/my/bundle@sha256:eeeeeeeee...",
  "created": "2018-08-30T20:39:55.549002887-06:00",
  "custom": {},
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

The fields above are defined as follows:

- `bundle` (REQUIRED): The bundle, as defined in [the Bundle Definition](101-bundle-json.md).
- `created` (OPTIONAL): A timestamp indicating when this release claim was first created. This MUST not be changed after initial creation.
- `bundleReference` (OPTIONAL): A canonical reference to the bundle used in the last action. This bundle reference SHOULD be digested to identify a specific version of the referenced bundle.
- `custom` (OPTIONAL): A section for custom extension data applicable to a given runtime.
- `modified` (OPTIONAL): A timestamp indicating the last time this release claim was modified. Executing the installation action MUST set this to a timestamp that matches `created`. Executing any modifying action (`upgrade`, `uninstall`, or a custom action with the property `"modifies": true`) MUST set this field to the time of the operation.
- `installation` (REQUIRED): The name of the _installation_. This can be automatically generated, though humans may need to interact with it. It MUST be unique within the installation environment, though that constraint MUST be imposed externally. Elsewhere, this field is referenced as the _installation name_. The format of this field must follow the same format used for the `installation` field in the [bundle.json file specification](101-bundle-json.md#the-bundlejson-file).
- `outputs` (OPTIONAL): Key/value pairs that were created by the operation. These are stored so that the user can access them after the operation completes. Some implementations MAY choose not to store these for security or portability reasons. See the [Outputs](#outputs) section for details. If this field is not present, it may be assumed that no outputs were generated as a result of the operation.
- `parameters` (OPTIONAL): Key/value pairs that were passed in during the operation. These are stored so that the operation can be re-run. Some implementations MAY choose not to store these for security or portability reasons. However, there are some caveats. See the [Parameters](#parameters) section below for details.
- `result` (REQUIRED): The outcome of the bundle's last action (e.g. if action is install, this indicates the outcome of the installation.). It is an object with the following fields:
  - `action` (REQUIRED): Indicates the action that the current bundle is in. This may be any of the built-in actions (`install`, `upgrade`, `uninstall`) as well as any custom actions as defined in the bundle descriptor. The special name `unknown` MAY be used in the case where the CNAB Runtime cannot determine the action name of a claim.
  - `message` (OPTIONAL): A human-readable string that communicates the outcome. Error messages MAY be included in `failure` conditions.
  - `status` (REQUIRED): Indicates the status of the last phase transition. Valid statuses are:
    - `failure`: failed before completion
    - `underway`: in progress. This should only be used if the invocation container MUST exit before it can determine whether all operations are complete. Note that `underway` is a _long term status_ that indicates that the installation's final state cannot be determined by the system. For this reason, it should be avoided. When used, `underway` should be considered a temporary status, and the runtime SHOULD work to resolve this to either `failure` or `success`.
    - `unknown`: the state is unknown. This is an error condition.
    - `success`: completed successfully
- `revision` (REQUIRED): An [ULID](https://github.com/ulid/spec) that MUST change each time the claim is modified. It MUST NOT change when a [non-modifying operation](https://github.com/cnabio/cnab-spec/blob/master/101-bundle-json.md#custom-actions) is performed on an installation.

> Note that credential data is _never_ stored in a claim. For this reason, a claim is not considered _trivially repeatable_. Credentials MUST be supplied on each action. Implementations of the claims specification are expected to retrieve the credential requirements from the `bundle` field.

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
    "credentials":{ 

    },
    "definitions":{ 
        "x509Certificate":{ 
          "contentEncoding":"base64",
          "contentMediaType":"application/x-x509-user-cert",
          "type":"string",
          "writeOnly":true
        }
    },
    "images":{ 

    },
    "invocationImages":[ 
        { 
          "image":"technosophos/demo2:0.2.0",
          "imageType":"docker"
        }
    ],
    "name":"technosophos.hellohelm",
    "outputs":{ 
        "clientCert":{ 
          "definition":"x509Certificate",
          "path":"/cnab/app/outputs/clientCert"
        }
    },
    "parameters":{ 

    },
    "schemaVersion":"v1.0.0",
    "version":"0.1.0"
  },
  "created": "2018-08-30T20:39:55.549002887-06:00",
  "installation": "technosophos.hellohelm",
  "modified": "2018-08-30T20:39:59.611068556-06:00",
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

### How the Claim is Used

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

### Environment Variables

The claim is produced outside of the CNAB package. The following claim data is injected
into the invocation container at runtime:

- `$CNAB_INSTALLATION_NAME`: The value of `claim.installation`.
- `$CNAB_BUNDLE_NAME`: The name of the present bundle.
- `$CNAB_ACTION`: The action to be performed (install, upgrade, ...)
- `$CNAB_REVISION`: The ULID for the present revision. (On upgrade, this is the _new_ revision)
- `$CNAB_CLAIMS_VERSION`: The version of this specification (currently `CNAB-Claims-1.0-WD`)

> Credential data, which is also injected into the invocation image, is _not_ managed by the claim system. Rules for injecting credentials are found in [the bundle runtime definition](103-bundle-runtime.md).

Parameters declared with an `env` key in the `destination` object MUST have their values injected as environment variables according to the name specified. Likewise, files MUST be injected if `path` is set on `destination`.

### Files

The invocation image may benefit from accessing the claim. Consequently, a claim MUST be attached to the invocation image when the invocation image is started.

The claim MUST be mounted at the path `/cnab/claim.json` inside of the bundle. The version of this claim that is to be mounted is the _version prior to this operation beginning_. In other words, when a bundle is installed, it creates the original installation claim. On the first upgrade, the claim describing the _installation_ is located at `/cnab/claim.json`. This allows the invocation image to inspect the former state and compare it to the desired new state.

Note: Systems may be compliant with the core specification but not support the claims specification. If `$CNAB_CLAIMS_VERSION` is not present, a runtime SHOULD assume that the claims specification is not implemented. Bundle authors may therefore have to take care when relying upon `/cnab/claim.json`, accommodating the case where the runtime does not support claims.

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
  "installation": "technosophos.my_first_install",
  "modified": "TIMESTAMP",
  "result": {
    "message": "yay!",    // From STDOUT (echo)
    "action": "install",  // Determined by the action
    "status": "success"   // success (exit == 0), failure (exit > 0), or underway (connection terminated before exit code was received)
  },
  "revision": "01CN530TF9Q095VTRYP1M8797C"
}
```

Tools that implement claims MAY then present `result` info to end users to show the result of running an invocation image.

## Credentials, Parameters, and Claims

The claims specification makes two assumptions about security:

- Claims are safe for writing sensitive information
- Claims are not a means for proxying identity

Parameters may contain sensitive information, such as database passwords. And this specification assumes that storing such information (un-redacted) can be done securely. However, the precise security mechanisms applied at the storage layer are beyond the scope of this specification.

Parameter values MUST be stored in claims. Implementations of the claims spec MUST NOT provide methods to block certain parameters from having their values written to the claim.

Credentials, on the other hand, are proxies for user identity, including potentially both authentication information and authorization information. Claims are not intended to provide a proxy for identity. In other words, a claim record should allow two users with identical permissions the ability to install the same configuration, it should not provide a way for one user to _gain_ the permissions of another user.

Credential data MUST NOT be stored in claims. Credentials are identity-specific, while claims are identity-neutral.
