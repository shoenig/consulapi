consulapi
=========

### About
[consulapi](https://github.com/shoenig/consulapi) is a consul
client library for Go programs, targeted at the "99% use case".
While an official Go client provided by Hashicorp exists and
exposes the complete functionality of consul, it is sometimes
difficult to use and is always extremely painful to work with
in tests.

This consul client library for Go aims to be easily
mockable, and provides interfaces that are very easy to understand.

### Install
Like any go library, just use `go get` to install. If the Go team
ever officially blesses a package manager, we'll switch to that.

`go get github.com/shoenig/consulapi`

### Usage
Creating a client is very simple - just call `New` with the desired
`ClientOptions`.

```go
options := consulapi.ClientOptions{
    Address: "http://localhost:8500", // default
    HTTPTimeout: 10 * time.Seconds, // default
    SkipTLSVerification: false, // http used by default
 }

client := consulapi.New(options)
// client implements the consulapi.Client interface

members, err := client.Members()
// etc ...
```
