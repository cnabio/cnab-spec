# The Base Bundle Pattern

This section of the specification is non-normative. It describes a pattern for using the contents of one bundle in another bundle as a form of extension or inheritance. While this is non-normative to the specification, facilitating this pattern was a goal of CNAB's design.

The _base bundle pattern_ is a pattern for sharing common bundle tooling across multiple bundles. Succinctly expressed, a _base bundle_ is a bundle whose contents are inherited by another bundle. The most frequent way in which is occurs is when one bundle uses a `Dockerfile` to import another bundle's invocation image in its `FROM` line.

A bundle that uses a base bundle as a source is referred to here as an _extending bundle_ or _extension bundle_, drawing from the object oriented terminology of base classes and extension classes.

For example, a base CNAB image MAY provide something like this:

```Dockerfile
FROM ubuntu:latest

COPY ./some-chart /cnab/app/charts/some-chart
```

If the above bundle is built as `base-bundle:latest`, then it MAY be referenced by other images.

The `Dockerfile` for an extending bundle can import `base-bundle:latest` as a starting point:

```Dockerfile
FROM base-bundle:latest

RUN helm inspect values /cnab/app/charts/some-chart > ./myvals.yaml &&  sed ...
```

The example above is simply intended to show how by reserving the `/cnab` directory, we can make images extensible, while not worrying about non-CNAB images putting data in places CNAB treats as special.

The mechanisms for taking a base bundle and creating an extension bundle are not formalized in the definition, as tooling MAY implement this pattern as fit for the tooling's domain. However, there are particular points of consideration that implementations MAY wish to address:

- Run tools make good candidates for reuse, and MAY then be well suited for base bundles
- A base bundle MAY define parameters and credentials in its bundle definition. In such cases, implementations of the base bundle pattern MAY wish to bubble this configuration up to the extension bundle via tooling.
- The `/cnab/build` directory MAY be used to store tooling in a base bundle that extension bundles can use during bundle assembly.
- Docker and other OCI tooling MAY provide mechanisms for creating a bundle by composing via more than one base image. The present specification does not prohibit composition, and this section of the specification is non-normative.

Next section: [Process](901-process.md)