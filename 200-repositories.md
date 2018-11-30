# CNAB Repositories

A CNAB repository is a location where packaged bundles can be stored and shared. They can be created by anyone to distribute their own bundles.

CNAB supports two HTTP based transfer protocols; A simple protocol which requires only a standard HTTP server on the server end of the connection, and a "smart" protocol with deeper integration with CNAB clients. This document describes both protocols.

As a design feature, smart clients can automatically upgrade simple protocol URLs to smart URLs.  This permits all users to have the same published URL, and the peers automatically select the most efficient transport available to them.

This document explains how to create and work with CNAB repositories.

## General Information

This section describes information that applies to both the "smart" and simple protocols.

### URL Format

URLs for CNAB repositories accessed by HTTP use the standard HTTP URL syntax documented by RFC 1738, so they are of the form:

    https://<host>:<port>/<path>?<query>

Within this documentation, the placeholder `$REPO_URL` will stand for the `https://` repository URL entered by the end-user.

It should also be noted that `<path>` can be many levels deep. `/foo` and `/foo/bar/car/star` are both valid paths.

### Authentication

Standard HTTP authentication is used if authentication is REQUIRED to access a repository, and MAY be configured and enforced by the HTTP server software. Clients SHOULD support Basic authentication as described by RFC 2617. Servers SHOULD support Basic authentication by relying upon the HTTP server placed in front of the CNAB server software.

Clients and servers MAY support other common forms of HTTP based authentication, such as Digest authentication or OAuth2.

### SSL

It is STRONGLY recommended that clients and servers support SSL, particularly to protect passwords during the authentication process.

### Session State

The CNAB transfer protocol is intended to be completely stateless from the perspective of the HTTP server side. All state MUST be retained and managed by the client. This permits simple round-robin load-balancing on the server side, without needing to worry about state management.

### General Request Processing

Except where noted, all standard HTTP behavior SHOULD be assumed by both client and server.  This includes (but is not necessarily limited to):

If there is no repository at `$REPO_URL`, or the resource pointed to by a location matching `$REPO_URL` does not exist, the server MUST NOT respond with `200 OK` response. A server SHOULD respond with `404 Not Found`, `410 Gone`, or any other suitable HTTP status code which does not imply the resource exists as requested.

If there is a repository at `$REPO_URL`, but access is not currently permitted, the server MUST respond with the `403 Forbidden` HTTP status code.

Servers SHOULD support both HTTP 1.0 and HTTP 1.1. Servers SHOULD support chunked encoding for both request and response bodies.

Clients SHOULD support both HTTP 1.0 and HTTP 1.1. Clients SHOULD support chunked encoding for both request and response bodies.

Servers MAY return ETag and/or Last-Modified headers.

Clients MAY revalidate cached entities by including If-Modified-Since and/or If-None-Match request headers.

Servers MAY return `304 Not Modified` if the relevant headers appear in the request and the entity has not changed.  Clients MUST treat `304 Not Modified` identical to `200 OK` by reusing the cached entity.

Clients MAY reuse a cached entity without revalidation if the Cache-Control and/or Expires header permits caching.  Clients and servers MUST follow RFC 2616 for cache controls.

## The Capabilities API

HTTP clients can determine a server's capabilities by making a `GET` request to the root URL, `/`, without any search/query parameters.

The server response SHOULD contain a header called `CNAB-Capabilities` with a whitespace-delimited list of server capabilities. These allow the server to declare what it can and cannot support to the client.

Clients SHOULD fall back to the simple protocol if the header is not present. When falling back to the simple protocol, clients SHOULD discard the response already in hand, even if the response code is not between 200-399. Clients MUST NOT continue if they do not support the simple protocol. HTTP servers that only support the simple protocol MAY return a `CNAB-Capabilities` header.

A server MUST support all the capabilities it lists in the returned header.

Example smart server reply:

```bash
200 OK
CNAB-Capabilities: "search simple-proto smart-proto upload-thin-bundle upload-thick-bundle auth-oauth2"
```

### search

The "search" capability came about as a way for clients to determine if the repository has a search API available at `/search`, while maintaining compatibility for simpler servers (e.g. a basic file server).

When enabled, this capability means that the server has a search API available.

HTTP clients that support this capability MAY search for bundles by making a request to `/search`.

NOTE: HTTP clients that support only the simple protocol MAY still use this capability, provided that they understand the search API demonstrated below.

### simple-proto

The simple protocol is a vastly simplified API compared to the "smart" protocol. This protocol is called "simple" because it requires no CNAB-specific code on the server side during the transport process; the fetch process is a series of HTTP GET requests, where the client can assume the layout of the CNAB repository on the server.

When enabled, this capability means that the server can handle requests using the simple protocol.

### smart-proto

When enabled, this capability means that the server can handle requests using the smart protocol. This capability is REQUIRED to handle uploads.

### upload-thin-bundle

A thin bundle is one which reference container images not contained within the bundle (but are known to exist at the receiving end).

When enabled, this capability means that the server can receive and host thin bundles. Supporting this feature compared to "upload-thick-bundle" can reduce the network traffic significantly.

This feature MUST NOT be enabled if "smart-proto" is disabled, as the simple protocol does not define a way to handle uploads.

HTTP servers that support the "smart" protocol MAY support either the "upload-thin-bundle" or "upload-thick-bundle" capability, or both.

If both the "upload-thin-bundle" and "upload-thick-bundle" capabilities are not present, the server is considered to be in read-only mode.

### upload-thick-bundle

A thick bundle is one which reference container images are contained within the bundle.

When enabled, this capability means that the server can receive and host thick bundles. Supporting this feature increases the network traffic significantly.

This feature MUST NOT be enabled if "smart-proto" is disabled, as the simple protocol does not define a way to handle uploads.

HTTP servers that support the "smart" protocol MAY support either the "upload-thin-bundle" or "upload-thick-bundle" capability, or both.

### auth-oauth2, auth-basic, auth-digest

The "auth-" capabilities came about as a way for clients to determine the authentication strategy to be used against the server for authentication/authorization.

If no "auth-" capability is present, the server supports no auth strategy.

## The Simple Protocol

The simple protocol is a vastly simplified API compared to the "smart" protocol. This protocol is called simple because it requires no CNAB-specific code on the server side during the transport process; the fetch process is a series of HTTP GET requests, where the client can assume the layout of the CNAB repository on the server.

Let’s follow the process for fetching a CNAB bundle. In this example, we'll attempt to fetch version 2.0.0 of the `helloworld` bundle.

### Fetching a Bundle

Simple HTTP clients MUST first make a `GET` request to `$REPO_URL/info/refs`, without any search/query parameters.

The Content-Type of the returned entity SHOULD be `application/json`, but MAY be any content type. Clients MUST attempt to validate the content against the returned Content-Type.

When examining the response, clients SHOULD only examine the HTTP status code. The only valid response codes are `200 OK` and the redirection status codes 300-399; anything else should be considered an error.

The returned content is a list of remote references in the repository as well as their content digests. References can be anything (for example, a particular version of a bundle) and are stored in the repository's "blob store".

If the Content-Type of the returned entity is `application/json`, clients SHOULD to attempt to verify the bundle's signature. When the signature fails, clients SHOULD NOT continue unless they intentionally choose to ignore the signature failure.

Example simple server reply:

```bash
200 OK
Content-Type: application/json

{
    "refs/tags": {
        "1.0.0": "ca82a6dff817ec66f44342007202690a93763949",
        "2.0.0": "6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6"
    }
}
```

Now you have a list of the remote references and their respective digests. As you can see, versioned bundles are always stored under the `refs/tags` prefix.

At this point, the client is ready to start the walking process. Because the starting point is the `6c3c624` object we saw in the `info/refs` file, the client starts by fetching the digest.

Simple HTTP clients MUST make a `GET` request to `$REPO_URL/blobs/6c3c624b58dbbcd3c0dd82b4c53f04194d1247c6` without any search/query parameters.

The Content-Type of the returned entity SHOULD be `multipart/signed; value="application/json"`, but MAY be any content type. Clients MUST attempt to validate the content against the returned Content-Type.

Cache-Control headers MAY be returned to disable caching of the returned entity.

When examining the response, clients SHOULD only examine the HTTP status code. The only valid response codes are `200 OK` and the redirection status codes 300-399; anything else should be considered an error.

The returned content is a clear-signed bundle as described in [The bundle.json File][bundle.json].

If the Content-Type of the returned entity is `multipart/signed; value="application/json"`, clients SHOULD to attempt to verify the bundle's signature. When the signature fails, clients SHOULD NOT continue unless they intentionally choose to ignore the signature failure.

Example simple server reply:

```bash
200 OK
Content-Type: multipart/signed; value="application/json"

-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA256

{
    "schemaVersion": "v1",
    "name": "technosophos.helloworld",
    "version": "2.0.0",
    "description": "An example 'thin' helloworld Cloud-Native Application Bundle",
    ...
}
-----BEGIN PGP SIGNATURE-----
Comment: helloworld v0.1.2

wsBcBAEBCAAQBQJbvomsCRD4pCbFUsuABgAAC/IIAI3LD89Fn9aJu/+eNsJnTyJ1
7T9KQFkekAe681eMkVMUY1NDjYfcQjaw0BZqSxOrs7Tunjwxxxm4pG1ua3sDp99a
NiB2tJN6AOKWXfs6zg3d8igskANv1ArmKqEiUyL69O8eBO0fz2dfUw67JazWu6HE
+MYpurRph8w5Sz9Ay3STntsFngGEgB87P/UMFFioY1KebJpBNMhuGa6SrT8kxNif
ERQachtjnsZiPQddPo2AJYFuN4XxbHpRvi+N8F8T2gQIjP9Ux7muegUI3qU9q9PU
VaefYa8rHJpw3VIt+1qf0RoiW53zJD+dYhSwTH4MBeagyDOjmQiLbXRI4Ofbc1s=
=JinU
-----END PGP SIGNATURE-----
```

## The "smart" Protocol

The simple protocol is, well, simple. It's a a bit inefficient and it can’t handle writing of data from the client to the server. The smart protocol is an alternative method of transferring data, but it requires a remote end that is intelligent about CNAB – it can read local data, figure out what the client has and needs, and generate a custom responses for it.

HTTP clients that support the "smart" protocol (or both the "smart" and simple protocols) discover bundles by making a parameterized request for the bundle.

The request MUST contain exactly one query parameter, `service=$servicename`, where `$servicename` MUST be the service name the client wishes to contact to complete the operation. The request MUST NOT contain additional query parameters.

Service names are defined as a particular action (e.g. "upload-bundle") that both the client and the server understand for a particular operation, such as uploading bundles, logging in, fetching bundles, etc.

Because the communication protocol between the client and the server is left up entirely to the implementation of the requested service name, no request type (such as GET/POST/PUT/PATCH) is REQUIRED, but for demonstration purposes we are demonstrating the protocol using a GET request.

```bash
GET $REPO_URL?service=upload-bundle HTTP/1.1
```

Example smart server reply:

```bash
200 OK
Content-Type: application/x-upload-bundle-advertisement
```

If the server does not recognize the requested service name, or the requested service name has been disabled by the server administrator, the server MUST respond with the `403 Forbidden` HTTP status code.

Otherwise, smart servers MUST respond with the smart server reply format for the requested service name.

Cache-Control headers SHOULD be used to disable caching of the returned entity.

The Content-Type MUST be `application/x-$servicename-advertisement`.

Clients SHOULD fall back to the simple protocol if another content type is returned. When falling back to the simple protocol, clients MUST discard the response already in hand and make an additional request to `$REPO_URL/info/refs` following the documentation for fetching a bundle using the simple protocol. Clients MUST NOT continue if they do not support the simple protocol.

Further content negotiation and the communication protocol between the client and the server is left up entirely to the custom reply format for the requested service name.

## The Search API

HTTP clients that support the "search" capability MAY search for repositories by making a request to `/search`. The `/search` API is a global resource, used to search for CNAB repositories.

HTTP servers that only support the simple protocol DO NOT need to implement a search API, but are welcome to do so if they desire. If they choose to do so, these servers MUST support the capabilities API.

Using the previous URL example, clients send requests to the following endpoint:

    https://<host>:<port>/search?q=<searchpart>

The following query parameters are supported:

- `q=$query`, where `$query` MUST be a string of keywords the client wishes to use to search across the repository
- `service=$servicename`, where `$servicename` MUST be the service name the client wishes to contact to complete the operation.

Servers that support the search capability MUST support at least the `q=$query` parameter, but in the case of only supporting the simple protocol, DOES NOT need to support the `service=$servicename` parameter.

```bash
GET https://<host>:<port>/search?q=helloworld HTTP/1.1
```

Example smart server reply:

```bash
200 OK
Content-Type: application/json
Link: <https://<host>:<port>/search?q=helloworld&page=2>; rel="next", <https://<host>:<port>/search?q=helloworld&page=50>; rel="last"

{
    "apiVersion": "v1",
    "bundles": {
        "bacongobbler.helloworld": "https://<host>:<port>/v2/bacongobbler/helloworld",
        "radu-matei.helloworld": "https://<host>:<port>/v2/radu-matei/helloworld",
    }
}
```

The Content-Type of the returned entity SHOULD be `application/json`, but MAY be any content type. Clients MUST attempt to validate the content against the returned Content-Type.

The basic response is paginated based on the number of repositories, however the default number of pages and the number of entries per page are left up to the server.

In the above example, there are 2 entries; one for each repository and its respective URL.

The Link header includes pagination information. It's important for clients to form calls using Link header values instead of constructing your own URLs.

The possible `rel` values are:

| Name  | Description                                                   |
|-------|---------------------------------------------------------------|
| next  | The link relation for the immediate next page of results.     |
| last  | The link relation for the last page of results.               |
| first | The link relation for the first page of results.              |
| prev  | The link relation for the immediate previous page of results. |

Clients can traverse through the paginated response by adding another query parameter, `page=$pagenumber`, where `$pagenumber` MUST be an integer between 1 and the value in "pages".

If no `page` query parameter is set, the response MUST be the first page.

Pages are one-indexed, such that `page=1` is the FIRST page. `page=0` is NOT a valid page number.

HTTP clients MAY also add a third query parameter, `per_page=$numentries`, where `$numentries` is the number of entries the client wishes to view in a single page. Servers MAY choose to ignore this query parameter, and clients should be prepared for that.

If the client includes the `service=$servicename` parameter and the server supports the smart protocol, the Content-Type MUST be `application/x-$servicename-searchresults`. Smart clients MAY fall back to the standard search protocol if another Content-Type is returned. When falling back to the dumb protocol, clients SHOULD re-use the response already in hand. Clients MUST NOT continue if they do not support the simple protocol.

If it's a smart response, content negotiation and the communication protocol between the client and the server is left up entirely to the custom reply format for the requested service name, so the response may differ entirely from what's shown here.

## Motivation

CNAB repositories are a centralized location where packaged bundles can be stored and shared. They can be created by anyone to distribute their own bundles, and users can use these repositories to share, collaborate and consume bundles created by the community. It makes searching, fetching, and sharing bundles easier, secure, and manageable for both the producer and the consumer of these bundles.

## Rationale

In early versions of reference implementations of the CNAB spec, we experimented with repositories being hosted from git repositories. We swapped this out with the HTTP-based approach after feedback for a couple reasons:

- Scalability: users can leverage their existing "object storage" platforms for hosting bundles (e.g. Azure Artifacts, Google Cloud Storage, AWS S3)
- Ecosystem compliance: CNAB repositories align more closely with similar distribution models like [Docker's distribution project](https://github.com/docker/distribution)

## Reference Implementations

CNAB clients should be able to handle the entire repository lifecycle, from logging in, searching, publishing, and fetching bundles.

As a reference implementation to the CNAB spec, Duffle implements the following commands to handle a CNAB repository's lifecycle:

- `duffle login` logs in to a bundle repository
- `duffle logout` logs out from a bundle repository
- `duffle pull` pulls bundles from a bundle repository
- `duffle push` pushes bundles to a bundle repository
- `duffle search` searches across logged in bundle repositories for bundles

[bundle.json]: 101-bundle-json.md
