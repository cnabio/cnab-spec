---
title: Dependencies
weight: 500
---

# CNAB Dependencies 1.0
Working Draft, April. 2019

This specification describes how bundles can define dependencies on other
bundles.

This specification, the CNAB Dependencies specification, is not part of the CNAB Core specification. An implementation of CNAB Core MAY NOT implement this specification, yet still claim compliance with CNAB Core. A CNAB implementation MUST NOT claim to be CNAB Dependencies-compliant unless it meets this specification.

# Overview

One way for a bundle author to manage the complexity of a large bundle and reuse
common logic is to create bundles for discrete components, such as managing
a database as a service on a cloud platform, and then use that bundle as a
dependency.

This specification defines [dependencies metadata](#dependencies-metadata) in the
bundle.json for specifying dependencies but does not dictate that the metadata is
specifically used at a particular time, or how the dependency graph is resolved.
Each tool may choose to use the information differently, for example [porter](https://porter.sh) resolves dependencies when the bundle is built, but another tool could
use this information at runtime.

There are two cases for how a bundle may need to depend upon another bundle:

1. [Depend on a named bundle](#depend-on-a-named-bundle)
1. [Depend on a bundle that provides a capability](#depend-on-a-bundle-that-provides-a-capability)

## Depend on a named bundle

The bundle depends on a specific named bundle that is known in advance. It is 
stored in the custom extensions section of the bundle.

```json
{
  "custom": {
    "dependencies": {
      "requires": {
        "storage": {
          "bundle": "somecloud/blob-storage"
        },
        "mysql": {
          "bundle": "somecloud/mysql",
          "version": {
            "prereleases": true,
            "ranges": ["5.7.x"]
          }
        }
    },
  },
  "name": "wordpress"
}
```

## Depend on a bundle that provides a capability

This section is a placeholder and will be completed in a follow-up pull request.

## Dependencies Metadata

This specification introduces a `dependencies` object in the bundle.json
that defines metadata necessary to specify a dependency.

The entry `dependencies` in the custom extension map, `custom`, is reserved and
MUST only be used for this CNAB Dependencies Specification.

### requires

The `requires` map defines the criteria for the dependent bundle. The key for
each dependency is a way for the bundle to reference the dependency.

* `bundle`: A reference to a bundle in the format REGISTRY/NAME.
* `version`: A set of criteria applied to the bundle version when selecting an
    acceptable version of the bundle.
    * `ranges`: An array of allowed version ranges.

        Versions are specified using semver, with or without a leading "v".
        An `x` in the place of the minor or patch place can be used to specify
        a wildcard. Ranges can be specified by separating the two versions with
        a dash, the dash must be surrounded by spaces in order to disambiguate
        from prerelease tags.

        Below are some example ranges:
        * `1.2.3` - Requires an exact version.
        * `1.2.x` - Restricts the major.minor to `1.2` and allows any patch.
        * `1.x` - Restricts the major to `1` and allows any minor.
        * `1.5.x - 3` - Allows a range of versions from `1.5` up to but not including 4.
    * `prereleases`: A boolean field, defaults to `false`, which specifies if
        prerelease versions of the bundle are allowed.
