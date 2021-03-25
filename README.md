# httprouter-group

[日本語](./README.jp.md)

## Status

[![Check codes](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yml)

[![Create Release](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yml)

## Description

Package [`mythrnr/httprouter-group`](https://github.com/mythrnr/httprouter-group)
provides to define grouped routing with [`julienschmidt/httprouter`](https://github.com/julienschmidt/httprouter).

For more details, see [GoDoc](https://pkg.go.dev/github.com/mythrnr/httprouter-group).

## Feature

### [`RouteGroup`](https://github.com/mythrnr/httprouter-group/blob/master/route_group.go)

- The `RouteGroup` registered by `(*RouteGroup).Children(...*RouteGroup)`
  can take over parent path and middlewares.
- `mythrnr/httprouter-group` is intended to simplify definitions.
  This package does not affect the performance of `julienschmidt/httprouter`
  as shown in the following example.

```go
package main

import (
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"
    group "github.com/mythrnr/httprouter-group"
)

// This definition provides following routes.
//
// - `GET /` with middleware 1
// - `GET /users` with middleware 1, 2
// - `GET /users/:id` with middleware 1, 2
// - `PUT /users/:id` with middleware 1, 3
// - `DELETE /users/:id` with middleware 1, 3
//
func main() {
    // first, define routes, handlers, and middlewares.
    g := group.New("/").Handle(
        http.MethodGet,
        func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
            rw.Write([]byte("GET /\n"))
        },
    ).Middleware(
        func(h httprouter.Handle) httprouter.Handle {
            return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                rw.Write([]byte("Middleware 1: before\n"))
                h(rw, r, p)
                rw.Write([]byte("Middleware 1: after\n"))
            }
        },
    ).Children(
        group.New("/users").Handle(
            http.MethodGet,
            func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                rw.Write([]byte("GET /users\n"))
            },
        ).Middleware(
            func(h httprouter.Handle) httprouter.Handle {
                return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    rw.Write([]byte("Middleware 2: before\n"))
                    h(rw, r, p)
                    rw.Write([]byte("Middleware 2: after\n"))
                }
            },
        ).Children(
            group.New("/:id").Handle(
                http.MethodGet,
                func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    rw.Write([]byte("GET /users/:id\n"))
                },
            ),
        ),
        group.New("/users/:id").Handle(
            http.MethodPut,
            func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                rw.Write([]byte("PUT /users/:id\n"))
            },
        ).Handle(
            http.MethodDelete,
            func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                rw.Write([]byte("DELETE /users/:id\n"))
            },
        ).Middleware(
            func(h httprouter.Handle) httprouter.Handle {
                return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    rw.Write([]byte("Middleware 3: before\n"))
                    h(rw, r, p)
                    rw.Write([]byte("Middleware 3: after\n"))
                }
            },
        ),
    )

    // next, set up and configure router.
    router := httprouter.New()
    router.PanicHandler = func(rw http.ResponseWriter, r *http.Request, rec interface{}) {
        log.Fatal(rec)
    }

    // finally, register routes to httprouter instance.
    g.Register(router)

    // serve.
    log.Fatal(http.ListenAndServe(":8080", router))
}
```

## Requirements

Go 1.13 or above.

## Install

Get it with `go get`.

```bash
go get github.com/mythrnr/httprouter-group
```
