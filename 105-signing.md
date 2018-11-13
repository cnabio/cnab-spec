# Signing, Attesting, Verifying, and Validating CNAB Bundles

Package signing is a method frequently used by package managers as a way of ensuring, to the user, that the package was _vetted by a trusted entity_ and _has since, not been altered_. Signing is a critical piece of infrastructure security.

This document outlines how CNAB bundles use a multi-layered fingerprinting and signing strategy to meet these criteria.

## Summary

- Every image and invocation image field in the `bundle.json` must have a `digest:` field that must contain a digest of the image.
- Digests are computed in accordance with the underlying image type (e.g. OCI bundles are validated by computing the top hash of a Merkle tree, VM images are computed by digest of the image)
- Signed bundles are clear-signed `bundle.json` files according to the Open PGP specification. When present, these are used in lieu of the unsigned `bundle.json` file.
- Authority is granted by the signed bundle, and integrity is granted via the image digests embedded in the bundle.json
- Attestations provide a mechanism for making additional guarantees about a bundle. Attesting a bundle may indicate that a release has been certified, or passed tests, or manual checked. It is a method to attach additional cryptographically based assurances to a bundle

## Image Integrity with Digests

CNAB correlates a number of images, of varying types, together in a single bundle. This section of the specification defines how image integrity is to be tested via digests and checksumming.

### Digests, OCI CAS, and Check Summing

A frequently used tool for validating that a file has not been changed between time T and time T+n is _checksumming_. In this strategy, the creator of the file runs a cryptographic hashing function (such as a SHA-512) on a file, and generates a _digest_ of the file's content. The digest can then be transmitted to a recipient separately from the file. The recipient can then re-run the same cryptographic hashing function on the file, and verify that the two functions are identical. There are more elaborate strategies for digesting complex objects, such as a [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree), which is what the OCI specification uses. In any event, the output _digest_ can be used to later verify the integrity of an object.

The OCI specification contains a [standard for representing digests](https://github.com/opencontainers/image-spec/blob/master/descriptor.md#digests). In its simplest case, it looks like this:

```text
ALGO:DIGEST
```

Where `ALGO` is the name of the cryptographic hashing function (`sha512`, `md5`, `blake2`...) plus some optional metadata, and DIGEST is the ASCII representation of the hash (typically as a hexadecimal number).

> Note: The OCI specification only allows `sha256` and `sha512`. This is not a restriction we make here.

For example:

```text
sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b
```

### Digesting Objects in the `bundle.json`

CNAB is composed of a `bundle.json` and a number of supporting images. Those images are referenced by the `bundle.json`. Thus, digesting those artifacts and including their digest in the `bundle.json` provides a convenient way to store (and locate) digests.

To that end, anything that shows up in the `invocationImages` or `images` section of the `bundle.json` _must_ have a `digest`:

```json
{
    "name": "helloworld",
    "version": "0.1.2",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:0.1.0",
            "digest": "sha256:6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b"
        }
    ],
    "images": [
        {
            "name": "image1",
            "digest": "sha256:aaaa624b58dbbcd3c0dd82b4c53f04194d1247c6eebdaab7c610cf7d66709b3b",
            "uri": "urn:image1uri",
            "refs": [
                {
                    "path": "image1path",
                    "field": "image.1.field"
                }
            ],
        }
    ]
}
```

Objects must contain a `digest` field even if the digest is present in another field. This is done to provide _access uniformity_.

> OCI images, for example, may embed a digest in the image's _version_ field. According to this specification, while this is allowed, it does not remove the requirement that the `digest` field be present and filled.

Different formats (viz. OCI) provide definitions for validating a digest. When possible, images should be validated using these definitions, according to their `imageType`. If a particular image type does not already define what it means to have a digest verified, the default method is to retrieve the object as-is, and checksum it in the format in which it was delivered when accessed.

Drivers may choose to accept the digesting by another trusted agent in lieu of performing the digest algorithm themselves. For example, if a driver requests that a remote agent install an image on its behalf, it may trust that the image digest given by that remote agent is indeed the digest of the object in question. And it may then compare that digest to the `bundle.json`'s digest. In such cases, a driver _should_ ensure that the channel between the driver itself and the trusted remote agent is itself secured (for example, via TLS). Failure to do so will invalidate the integrity of the check.

## Signing the `bundle.json`

The `bundle.json` file will contain the digests of all executable objects. That is, everything in the `invocationImages` and `images` sections will have digests that will make it possible to ensure that their content has not been tampered with.

Consequently, the `bundle.json` acts as an authoritative resource for image integrity. To act as an authoritative source, however, it must provide an additional assertion: The `bundle.json` must assert the intention of the bundle creator, in marking this as a _verified bundle_.

This is accomplished by _signing the bundle_.

The signature method used by CNAB is defined by the [Open PGP standard](https://tools.ietf.org/html/rfc4880)'s digital signatures specification. In short, a _packaging authority_ (the individual responsible for packaging or guaranteeing the package), signs the bundle with a _private key_. The packaging authority distributes the accompanying public key via other channels (not specified herein, but including trusted HTTP servers, Keybase, etc.)

The _package recipient_ (the consumer of the package) may then retrieve the public keys. Upon fetching a signed bundle, the package recipient may then _verify_ the signature on the bundle by testing it against the public key.

An Open PGP signature follows [the format in Section 7 of the specification](https://tools.ietf.org/html/rfc4880#section-7):

```text
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

   <BODY>

-----BEGIN PGP SIGNATURE-----
Comment: <GENERATOR>

<SIGNATURE>
-----END PGP SIGNATURE-----
```

In the above, `<BODY>` is the entire contents of the `bundle.json`, `<GENERATOR>` is the optional name of the program that generated the signature, and `<SIGNATURE>` is the signature itself.

For example, here is a `bundle.json`:

```json
{
    "name": "foo",
    "version": "1.0",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:0.1.2",
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685"
        }
    ],
    "images": [],
    "parameters": {},
    "credentials": {}
}
```

This is signed using the technique called _clear signing_ (OpenPGP, Section 7), which preserves the input along with the cryptographic signature:

```text
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

{
    "name": "foo",
    "version": "1.0",
    "invocationImages": [
        {
            "imageType": "docker",
            "image": "technosophos/helloworld:0.1.2",
            "digest": "sha256:aca460afa270d4c527981ef9ca4989346c56cf9b20217dcea37df1ece8120685"
        }
    ],
    "images": [],
    "parameters": {},
    "credentials": {}
}
-----BEGIN PGP SIGNATURE-----
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCgAdFiEE+yumTPTtBoSSRd/j3NX15e8yw0UFAluEHv0ACgkQ3NX15e8y
w0UyKQf/Tb/mURLiHWmw68q7cjAHx7wVgjClo34oB07uY1RCvjMiK4sXaoKC+0KO
pQOC/15HY9f2aazPHig//aqNpFyyHcpX9XjvH51CbXiNcFvIv7JgmFwr7WIUY7cS
FsaFSCS53Z5HqCQ/SYB9OU4R+uwBW/gKmP7vBGieNkEhqHQklQG9vc9zUQjuTlYp
KIW9cGd0rKWzs8wwiW9FytIM43x54sHmtXRnWxu6zNReXr32u6ZFPrfVA0yoAJQ4
7iDhcM/VqL4j0xxfFmZuqkRCtsbgD6iUmL8VzINODGsF4lHFQrl2sKXAqMoIXyCw
ANjudClHNUNQFojriAX8YAO4V2yGVg==
=OoBW
-----END PGP SIGNATURE-----
```

Note that because we have _clear signed_ the `bundle.json`, there is no longer any need to transmit the `bundle.json` separately from the signed bundle. In fact, due to encoding differences, it is _preferable to use the signed bundle in lieu of the unsigned bundle_. 

## Attestations

The purpose of a clear-signed signature is to assert that a particular artifact (the bundle) has been generated by an authority, where authority is marked by a cryptographic signature.

Attestations provide a way to add multiple assurances to a bundle.

An attestation is used as a placeholder for the statement "The signing party has certified (attested) that condition X has been met". For example, when certification steps need to occur for a given bundle, one may use attestations to prove that the certification has been performed. In this case, the certifying party performs the process, and then (upon the bundle's passing), the certifying party adds a _detached signature_ to the set of signatures associated with the bundle.

In a more complex case, a signed package may be certified by one party for use in one way (we'll call this "certification A"). A different party may certify the bundle for a different case ("certification B"). Note that in this case, the certifications are _independent_, and are presumably done with separate justifications. ("certificate A ensures this bundle is suitable for internal use", "certificate B ensures this bundle is suitable for use by our partners".) Because of this feature, attestations _are not chainlike_. Each individual attestation must be verifiable without reference to any other attestations (including the original clear-signed signature).

Detached signatures are described in [section 11.4](https://tools.ietf.org/html/rfc4880#section-11.4) of the OpenPGP specification. Attestations are to be performed by extracting the `bundle.json` from a signed bundle, and then signing that same text object. A verification of a detached signature should use the `bundle.json` text as a basis for its verification. A bundle is considered _attested_ (or _attestation verified_) when a the bundle verification passes for the expected key used in that attestation.

They key used for an attestation must be used _only for that attestation_. Implementations _may_ use a subkey (of a master key) for specific attestations, while preserving other subkeys to perform other attestations or signings.

Implementations of CNAB _may_ support creating and verifying attestations. Implementations of CNAB _may_ support favoring an attestation with equal weight as the original signature. Attestations _may_ be stored with the signed bundle, though there is no requirement that attestations be stored in a specific place.

Attestations _may_ use the `COMMENT:` field of a detached signature to indicate, in a human-friendly way, what the attestation is for. However, agents _must not_ consider this information definitive. Comment fields are not calculated into the signature and can be easily modified. Instead, attestation _must_ be based solely on the key.

```text
-----BEGIN PGP SIGNATURE-----
Comment: Attestation - Certified for Internal Use

iQEzBAABCgAdFiEE+yumTPTtBoSSRd/j3NX15e8yw0UFAluayWsACgkQ3NX15e8y
w0VvVQgA1FtF03jqQgiAkxd707ELtmrKX5dcfIYtEr3o5fBGtckUebV5RYFwfQqZ
fYoTVEiAzgtR6ceXQB+SjCj8KD5uhf2nzX5eKIAmhyCLKibVBVCaTlTsKzNR/Xe4
fPWp/nSlNo6Xc2kwx6RRPPMpYk/7WhXm7iIl7MmHmveHmTM1oTdrzhf/y1ZTc0Vu
qdBSRvsDJMnaf+iB2g9r113ee12UBta9pbLIXjlWpv4PknL7QNsp2B0KeExQXgvZ
2KKrz+ndWr3I5aONa6Zr9hh3NdZc/oa1peqJaCJtsrLj08/+WiwdWTWG3/8k+toW
UkmNdIkOChvHv42XkWnF1t1Hyi51ig==
=sQ9B
-----END PGP SIGNATURE-----
```

## Key Revocation and Expiration

When public keys are expired or revoked, bundles signed with those keys become invalid.
They must be re-signed with a valid key.

CNAB verification tools _should_ handle the key revocation case.

## Bundle Retractions

Cases may arise where a particular version (or versions) of a bundle should no longer be used. For example, if a version of a bundle is discovered to be insecure in significant ways, bundle authors may wish to _explicitly mark_ that bundle as insecure. This process must be done in a way that retains the integrity of the bundle.

> This definition does not preclude the mere deletion of problematic bundles. Operators of a bundle repository, for instance, _may_ opt to merely remove insecure bundles from their servers rather than mark them and leave them. However, there are cases where historic (while insecure) packages may be retained and still made available for installation.

Reusing a release version to replace an insecure release with a secure one is _expressly prohibited_. For example, if release 2.3.1 of a bundle is deemed insecure, operators _must not_ re-release a modified bundle as 2.3.1. The fixed version _must_ modify a semantic component of the version number. For example, `2.3.2`, `2.4.0`, `3.0.0`, and even `2.3.2-alpha.1` are all acceptable increments. `2.3.1` and `2.3.1+1` are examples of forbidden version increments. Likewise, release `2.3.1` _must not_ be renamed by semantic component. (e.g. 2.3.1-insecure is illegal, while 2.3.1+insecure is legal). For clarification on this policy, see the [SemVer 2 specification](https://semver.org).

Instead, the prefered pattern is to retain the insecure release at its given release number, but issue a _retraction_.

A *retraction* is a cryptographically signed indicator that a bundle _may_ be installed.

```
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

[
    {
    "bundle": "petstore",
    "version": "2.3.1",
    "signature": "wsBcBAEBCAAQBQJbvomsCRD4pCbFUsuABgAAC/IIAI3LD89Fn9aJu/+eNsJnTyJ17T9KQFkekAe681eMkVMUY1NDjYfcQjaw0BZqSxOrs7Tunjwxxxm4pG1ua3sDp99aNiB2tJN6AOKWXfs6zg3d8igskANv1ArmKqEiUyL69O8eBO0fz2dfUw67JazWu6HE+MYpurRph8w5Sz9Ay3STntsFngGEgB87P/UMFFioY1KebJpBNMhuGa6SrT8kxNifERQachtjnsZiPQddPo2AJYFuN4XxbHpRvi+N8F8T2gQIjP9Ux7muegUI3qU9q9PUVaefYa8rHJpw3VIt+1qf0RoiW53zJD+dYhSwTH4MBeagyDOjmQiLbXRI4Ofbc1s=\n=JinU",
    "reason": "A remove execution vulnerability was found in the sanitization function emptyCatBox()."
    }
]
-----BEGIN PGP SIGNATURE-----
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCgAdFiEE+yumTPTtBoSSRd/j3NX15e8yw0UFAlvBKt8ACgkQ3NX15e8y
w0WlpQgAnf8UGUuW8c63M3HN/oYGC+glOww5xqPPkJVj5k/BWqpYnVYcUzLlHG4x
wGTO1qb/RewctKJweU4VC3ACiIKKeQzxLiUuFvkhhIWpoq8qY0xDLlwM7Ccc/7Vc
q5CwbQZ6apwoouZH/Yw0e/LedRspifQ+qxt0lyTnZQXV51o2/ubXchbgdwP2fzQs
m4NWTqL6US0C+SEVysesdPCydUIPC/oj8zo+M/aN4MTFmqHlMhvEMPeAa0O+wbIt
FBwmPYfhWfIogk9b63htdJMJxmCvkQWqGKHm8Y8IwFCkoTunvZEygidv5VkdJa9t
IMrgnkGly/iibIz2oagOfJhtAnthnQ==
=6DUE
-----END PGP SIGNATURE-----
```

A retraction is a clear signed JSON array.

- The top field is an array of individual retractions
- A retraction object has the following fields:
    - bundle: The name of the bundle being retracted
    - version: The version being retracted
    - signature: The signature of the specific bundle being retracted
    - reason: A human-oriented text description of why the bundle was retracted.

The `signature` field is optional, but provides a content-specific test on the content retracted. It is only applicable to cases where a specific version is being retracted.

The `version` field is optional. If omitted, the entire Bundle is considerered retracted. When `version` is omitted, `signature` _must_ be omitted.

To specify a range of versions, a _SemVer range_ may be provided in the `version` field. In this case, a `signature` _must_ be omitted.

The `reason` field is optional, though a retraction _may_ have one. This may be used by a user agent to explain the reason for the retraction.

The following examples shows all three methods of specifying a retraction:

```
-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

[
    {
        "bundle": "petstore",
        "version": "2.3.1",
        "signature": "wsBcBAEBCAAQBQJbvomsCRD4pCbFUsuABgAAC/IIAI3LD89Fn9aJu/+eNsJnTyJ17T9KQFkekAe681eMkVMUY1NDjYfcQjaw0BZqSxOrs7Tunjwxxxm4pG1ua3sDp99aNiB2tJN6AOKWXfs6zg3d8igskANv1ArmKqEiUyL69O8eBO0fz2dfUw67JazWu6HE+MYpurRph8w5Sz9Ay3STntsFngGEgB87P/UMFFioY1KebJpBNMhuGa6SrT8kxNifERQachtjnsZiPQddPo2AJYFuN4XxbHpRvi+N8F8T2gQIjP9Ux7muegUI3qU9q9PUVaefYa8rHJpw3VIt+1qf0RoiW53zJD+dYhSwTH4MBeagyDOjmQiLbXRI4Ofbc1s=\n=JinU"
    },
    {
        "bundle": "helloseattle",
        "reason": "the helloseattle tool has been removed do to licensing restrictions"
    },
    {
        "bundle": "fireflies",
        "version": "<=2.3.1"
    }
]
-----BEGIN PGP SIGNATURE-----
Comment: GPGTools - https://gpgtools.org

iQEzBAEBCgAdFiEE+yumTPTtBoSSRd/j3NX15e8yw0UFAlvBLI8ACgkQ3NX15e8y
w0VqqQf/fm3u48QLpa576pJISh/whopjjnnQm1vuCHfatyumC67W802ZEMCNu+pQ
5G6Ffsli5fifG7E8vsmFS9gxOTAYcsnrFvwDTD35zdPpctlFd+mUGSumgtXEHAWm
O5xDIlri6xntoJI/4MBYStzdCg0/Sj+qQRn/w8fGJyrViCczjGiZKwWldKNkKDRn
svIlu4itPHrwBnLOBMv1NulAkClGxDzD1VgTnHqgY8rPV+8dWj1VZOkhDZvDB7yR
iIO3klY4RNsy5EfKpWkrhyiY5LW680qZHDbC+Lvtv0J9HW4K0lFqcrzTRJEKiVDK
h3sGAYdx5fA5PfmweTCvc34qUvPVnw==
=dEXV
-----END PGP SIGNATURE-----
```

- Only one specific version of `petstore` is retracted
- All versions of `helloseattle` are retracted
- Versions of `fireflies` that are less than or equal to `2.3.1` are retracted

When an agent has access to a retractions list or lists, it _may_ evaluate the retractions for each request that would require loading the bundle. An agent _must not_ install or upgrade a bundle it knows to be retracted without the express consent of the user.

Next section: [declarative images](106-declarative-images.md)
