---
title: Appendix A: Preliminary Release Process
weight: 901
---

# Appendix A: Preliminary Release Process

While CNAB has not officially been assigned to a standards body, it is good practice to put in place a provisional process for indicating the status, stability, and reliability of the specification.

## Versioning

Specifications shall each track their own versions. A specification will call out a target version (such as 1.0), and then also include its phase (such as Working Draft, see below) to indicate progress towards the target version.

The CNAB Core specification shall therefore be formally referenced as _Cloud Native Application Bundle Core 1.0.0_. The phase of the process MAY be appended as a stability marker: _Cloud Native Application Bundle Core 1.0.0-WD_. It MAY be abbreviated to _CNAB1_.

## Process

The provisional process is outlined here, and is derived from the [W3C Process](https://www.w3.org/2017/Process-20170301/)

The specification will proceed through the following phases:

- *Working Draft (WD)*: During this phase, the specification is open for review. During this phase, major (breaking) changes may be made.
- *Candidate Recommendation (CR)*: This phase indicates that a final review is immanent. During this phase, minor (non-breaking) changes may be made, and major changes will be considered if they present an immanent problem with the specification.
- *Recommendation (REC)*: This phase indicates that the specification is stable, and the present version is no longer open for changes to the conformance of the specification
  - W3C [class 1 and 2 changes](https://www.w3.org/2017/Process-20170301/#correction-classes) are allowed on a recommendation.
  - W3C class 3 and 4 changes can be addressed in a _minor version update_.

And [errata section](https://www.w3.org/2017/Process-20170301/#rec-modify) may be added to a specification as an appendix, as a location for tracking clarifications and changes.
