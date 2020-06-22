---
title: Standardization Process
weight: 901
---

# Appendix A: The Process of Developing the Specification

This process governs the motion of a specification from rough draft through full standardization. A document is termed a _deliverable_ as it passes through various phases from _pre-draft deliverable_ to _accepted deliverable_, at which point it becomes a final specification.

## Versioning

Specifications shall each track their own versions. A specification will call out a target version (such as 1.0), and then also include its phase (such as Pre-Draft, see below) to indicate progress towards the target version.

For example, one of the specifications in this repository is the CNAB Core specification. The CNAB Core specification should be formally referenced as _Cloud Native Application Bundle Core 1.0.0_. The phase of the process MAY be appended as a stability marker: _Cloud Native Application Bundle Core 1.0.0-FA_. It MAY be abbreviated to _CNAB1_ or _CNAB1-FA_.

Other specifications must follow the same nomenclature.

Inspired by SemVer 2, CNAB follows a rigid versioning scheme. Versions are presented in the form `X.Y.Z[-S]`, where `X` is the major version, `Y` is the minor version, and `Z` is the patch version. The optional `-S` is a draft stability marker.

- Major releases (`X`): Major releases contain breaking changes, including features, fixes, and reorganizations. Implementors should not assume that two major versions are compatible. For example, `1.9.9` is not to be considered compatible with `2.0.0`.
- Minor releases (`Y`): Minor releases contain features and fixes only. A feature MUST NOT remove or modify existing pieces of the spec (including schemata and file system layouts), but MAY add new things. A minor release SHOULD be backward compatible, though certain security concerns may override this requirement.
    - Proposed changes to the spec SHOULD be rejected if accepting them will lead to changes in the interpretation of a bundle. (e.g. repurposing an existing field in the `bundle.json` is not allowed by these compatibility rules)
    - A runtime SHOULD be able to run bundles with the same major version and different minor versions.
    - A runtime MAY behave differently when processing a package of a newer minor release version, provided that difference is due to new features of the specification.
    - Bundle authors SHOULD strive to make bundles compatible with all minor versions of the same major version, though they MAY take advantage of additive features.
- Patch releases (`Z`): Patch release contain fixes to the text of the specification. Patch releases MUST NOT change the behavior of the specification (except in cases where the specification was too vague and the patch clarifies).
    - Patch releases MUST be both forward and backward compatible to the minor version number
    - Patch releases MUST NOT make the schema harder to validate against (though they may relax the schema).
    - Our intention is that CNAB Runtimes SHOULD NOT have to specify behavior at the patch level, as all patch levels are compatible.

Stability markers provide a way to indicate that a bundle or runtime is experimenting with features or fixes. _If an object is tagged with a stability marker, it MUST be treated as incompatible with any other version number. E.g. `1.0.0-GA` MUST be considered incompatible with `1.0.0`. Production artifacts SHOULD NOT use stability markers.

A small number of stability markers are allowed, as determined by [the foundation's governance documents](https://github.com/cnabio/community/blob/master/governance.md):

- `PD`: Pre-draft indicates that the version of the resource is an unstable testing version.
- `DRAFT`: Draft indicates that the version is an unstable in-development version
- `GA`: Working Group Approval indicates that this version has been accepted by the maintainers, but not accepted by the Executive Directors. In practice, this marker MUST only be used internally for testing, and never for production.
- `AD`: Final approval by the Executive Directors.

The tags `GA` and `AD` should never be used in a SemVer stability marker. `GA` is not a milestone that warrants publication, and `AD` is synonymous with the final release.

The development process defines `Publication` and `Submission to Other Standards Bodies` as process steps. However, neither of these reflects a change in the specification. These are people processes.

The stability markers `ALPHA`, `BETA`, and so on are _disallowed_ under this specification, and MUST NOT be used to express CNAB versions.

Finally, certain small errata may be fixed on an existing release without incrementing the release version. The following changes are allowed as errata fixes:

- Correcting spelling or typographical errors, where changing these does not alter the meaning of the specification.
- Correcting minor grammatical mistakes.
- Adding a revised link when a broken link appears. This should be done by appending the text `(Updated link: http://example.com...)`. The text may be corrected fully during the next version change.
- In extenuating circumstances, the Executive Directors may approve retroactively editing text to meet legal requirements. In such cases, the directors will not approve changes that break the specification. Under such circumstances, the directors may issue a _retraction of a specification_ (removing a published specification), and publish a new specification version that meets the legal requirements. For example, an intellectual property infringement may only be correctable by a retraction.

## Git Release Flow

This section deals with the practical considerations of versioning in Git, this repo's version control system.

### Patch releases

When a patch release of a specification is required, the working group must approve the scope of commits proposed for inclusion. The patch commit(s) should be merged to the `master` branch when ready. Next, a new branch should be created for the designated patch. For example, if the previous most recent branch name of the specification is `cnab-core-1.0.0-ga`, the new branch would be created from `cnab-core-1.0.0-ga` and named `cnab-core-1.0.1-ga`. The patch commit(s) should then be cherry-picked into this new branch.

When the final release is approved, a Git tag should also be pushed, which triggers schema artifact publishing. Extending the example above, a `cnab-core-1.0.1` tag should be created from the `cnab-core-1.0.1-ga` branch and pushed to origin. We drop the `-ga` suffix as branches and tags may not have the same name in Git.

### Minor releases

When a minor release of a specification is required, the working group must approve the scope of commits proposed for inclusion. Likely this will be the `master` branch once the approved commit(s) are merged. Next, a new branch should be created from `master` and named `cnab-core-1.1.0-ga` (here assuming that the version immediately prior was `cnab-core-1.0.0`).

When the final release is approved, a Git tag should also be pushed, which triggers schema artifact publishing. Extending the example above, a `cnab-core-1.1.0` tag should be created from the `cnab-core-1.1.0-ga` branch and pushed to origin. We drop the `-ga` suffix as branches and tags may not have the same name in Git.

### Major releases

When a major release of a specification is required, the working group must approve the scope of commits proposed for inclusion. Likely this will be the `master` branch once the approved commit(s) are merged. Next, a new branch should be created from `master` and named `cnab-core-2.0.0-ga` (here assuming that the version immediately prior was `cnab-core-1.0.0`).

When the final release is approved, a Git tag should also be pushed, which triggers schema artifact publishing. Extending the example above, a `cnab-core-2.0.0` tag should be created from the `cnab-core-2.0.0-ga` branch and pushed to origin. We drop the `-ga` suffix as branches and tags may not have the same name in Git.

### Ad hoc schema releases

In addition to the scenarios above, if schemas at a certain commit need to be preserved in the form of artifacts and published, ad hoc versioning (i.e. not tied to a release branch) is permitted via the following Git tag flow. This is intended for specifications in `DRAFT` state which are perhaps under heavy development.

To enable implementations to pin at a certain version prior to an official release, we can issue a Git tag and CI will handle publishing the schemas. For example, if the Claims specification is still in the `DRAFT` state but its schemas at a particular commit are needed for implementation verification, we can push an appropriate tag to origin. The tag form is: `cnab-claim-1.0.0-DRAFT+abc1234`, where `cnab-claim-1.0.0-DRAFT` is the current working version and `abc1234` is the short SHA of the commit the tag will be created from.

### Non-normative exclusions

For non-normative documents (any in the `8xx` section), no formal release processes apply.  As long as the associated working group approves of additions or modifications to documents in this section (say, via a pull request), no further action or process is needed.

## Development Process

The specification will proceed through the following phases:

- *Pre-Draft (PD):* Any working group participant or contributor may submit a proposed initial draft document as a candidate for this working group. The Working Group Chair may designate such a submission as a Pre-Draft Document.
- *Draft (Draft):* A Pre-Draft may be approved by the working group, in which case it becomes an official draft under the auspices of the CNAB specification working group. The working group will continue to revise the draft until it is in a state the group sees as fit for standardization.
- *Working Group Approval (GA):* When the working group believes a draft fit for standardization, the group formally approves the draft and submits it for final approval.
- *Final Approval (AD):* The Executive Director (or a named proxy) may grant Final Approval to a draft with Working Group Approval. At this point, the work is now designated an Approved Deliverable and is no longer a draft.
- *Publication:* When Final Approval is granted, the Approved Deliverable is made publicly available under the terms of the Open Web Foundation 1.0 license.
- *Submission to Other Standards Bodies:* With the approval of the Working Group, the Executive Director may submit an Approved Deliverable to another standards body in collaboration with the JDF.

Documents that have reached AD are considered complete. Errata may be captured in a separate section of the document, but the document itself is not changed except to correct typographical and formatting errors where necessary.

When the content of a document needs changes that cannot be captured as errata, a new _version_ of the specification must be created, and must proceed through the stages outlined above.

## Changes

An earlier "provisional" process was outlined here, based on W3's model. That provisional process has now been replaced with the process described herein.

Anything previously marked "Working Draft" is now considered to be a Draft, as they have been accepted for work by the working group.

The [CNAB Foundation's governance documents](https://github.com/cnabio/community/blob/master/governance.md) cover the acceptance process in more detail.