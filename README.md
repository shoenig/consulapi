consulapi
=========

Package `consulapi` provides a lightweight, easy-to-use Consul client that is
easy to write unit tests with.

[![Go Report Card](https://goreportcard.com/badge/gophers.dev/pkgs/consulapi)](https://goreportcard.com/report/gophers.dev/pkgs/consulapi)
[![Build Status](https://travis-ci.com/shoenig/consulapi.svg?branch=master)](https://travis-ci.com/shoenig/consulapi)
[![GoDoc](https://godoc.org/gophers.dev/pkgs/consulapi?status.svg)](https://godoc.org/gophers.dev/pkgs/consulapi)
[![NetflixOSS Lifecycle](https://img.shields.io/osslifecycle/shoenig/consulapi.svg)](OSSMETADATA)
[![GitHub](https://img.shields.io/github/license/shoenig/consulapi.svg)](LICENSE)

# Project Overview

Module `consulapi` is a consul client library for Go programs, focused on
the "90% use case". Although it is slightly feature limited compared to the
official `consul/api` library, it brings forward an easy to use API that most
consul users will appreciate.

The feature-oriented interfaces exposed by `consulapi` aim to be easily mockable,
making it easier to write unit tests that explore all possible outcomes of an
operation involving consul. Test those error conditions!

# Getting Started

The `consulapi` module can be installed by running
```bash
$ go get gophers.dev/pkgs/consulapi
```

#### Example Usage
Creating a client is very simple - just call `New` with the desired
`ClientOptions`.

```go
client := consulapi.New(consulapi.ClientOptions{
    Address:    "https://demo.consul.io",
    Logger:     loggy.New("elector-example"),
    // see client.go for full set of options
})

members, err := client.Members()
// etc ...
```

# Design
A few factors contribute to the simplicity of `consulapi`.

First, we export interfaces instead of concrete implementations.
This enabled both re-implementations if necessary, as well enables
the use of mocks in testing. A mock implementation for each

Second, the nature of the key-value store is significantly limited.
`consulapi` enforces a strongly opinionated design that all keys
and values must be strings, and that all keys may only be `/`
separated. This cuts down on a lot of type casting overhead.

Third, the source code itself is intended to be easy to read and
understand. It is centered around common http method calls, with
the intent of being a reduced reflection of the HTTP API.

# Contributing

The `gophers.dev/pkgs/consulapi` module is always improving with new features
and error corrections. For contributing bug fixes and new features please file
an issue.

# License

The `gophers.dev/pkgs/consulapi` module is open source under the [MIT](LICENSE) license.
