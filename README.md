# Simple Oauth 2.0

A Golang implementation for Oauth 2.0 + OIDC mechanism with clean architechture approach. 

## Oauth 2.0

From the book OAuth 2.0 in Action

OAuth 2.0 is a delegation protocol, a means of letting someone who controls a resource allow a software application to access that resource on their behalf without impersonating them. 

The application requests authorization from the owner of the resource and **receives tokens that it can use to access the resource**.

We can think of the OAuth token as a “valet key” for the web. Simple valet keys **limit the valet to accessing** the ignition and doors but not the trunk or glove box.

OAuth tokens can limit the client’s access to only the actions that the resource owner has delegated.

From [RFC 6749](https://tools.ietf.org/html/rfc6749)

The OAuth 2.0 authorization framework enables a third-party application to obtain limited access to an HTTP service, either on behalf of a resource owner by orchestrating an approval interaction between the resource owner and the HTTP service, or by allowing the third-party application to obtain access on its own behalf.

### Component Oauth 2.0

- The resource owner
- The protected resource
- The client
- The authorization server

## Starting the app

```bash
make all
```

### Using docker

```bash
make build
docker run --rm --name simple-oauth2 -p 8080-8083:8080-8083 zeihanaulia/simple-oauth2
```

## TODO

1. [Build template for client, authorization and protected service](https://github.com/zeihanaulia/simple-oauth2/pull/1)
2. [Implementation authorization code grant type code](https://github.com/zeihanaulia/simple-oauth2/pull/2)
3. [Implementation authorization code grant type implicit](https://github.com/zeihanaulia/simple-oauth2/pull/8)
4. [Implementation authorization code grant type credentials](https://github.com/zeihanaulia/simple-oauth2/pull/9)

## Reference

- [OAuth 2.0 in Action](https://learning.oreilly.com/library/view/oauth-2-in/9781617293276/)
- [OAuth 2.0 and OpenID Connect (in plain English)](https://www.youtube.com/watch?v=996OiexHze0)