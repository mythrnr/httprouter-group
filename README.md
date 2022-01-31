# httprouter-group

[日本語](./README.jp.md)

## Status

[![Check codes](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yaml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yaml)

[![Create Release](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yaml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yaml)

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
    "fmt"
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
    g := group.New("/").GET(
        func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
            w.Write([]byte("GET /\n"))
        },
    ).Middleware(
        func(h httprouter.Handle) httprouter.Handle {
            return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                w.Write([]byte("Middleware 1: before\n"))
                h(w, r, p)
                w.Write([]byte("Middleware 1: after\n"))
            }
        },
    ).Children(
        group.New("/users").GET(
            func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                w.Write([]byte("GET /users\n"))
            },
        ).Middleware(
            func(h httprouter.Handle) httprouter.Handle {
                return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    w.Write([]byte("Middleware 2: before\n"))
                    h(w, r, p)
                    w.Write([]byte("Middleware 2: after\n"))
                }
            },
        ).Children(
            group.New("/:id").GET(
                func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    w.Write([]byte("GET /users/:id\n"))
                },
            ),
        ),
        group.New("/users/:id").PUT(
            func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                w.Write([]byte("PUT /users/:id\n"))
            },
        ).DELETE(
            func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                w.Write([]byte("DELETE /users/:id\n"))
            },
        ).Middleware(
            func(h httprouter.Handle) httprouter.Handle {
                return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
                    w.Write([]byte("Middleware 3: before\n"))
                    h(w, r, p)
                    w.Write([]byte("Middleware 3: after\n"))
                }
            },
        ),
    )

    // next, set up and configure router.
    router := httprouter.New()
    router.PanicHandler = func(w http.ResponseWriter, r *http.Request, rec interface{}) {
        log.Fatal(rec)
    }

    // logging.
    //
    // GET     /
    // GET     /users
    // DELETE  /users/:id
    // GET     /users/:id
    // PUT     /users/:id
    fmt.Println(g.Routes().String())

    // finally, register routes to httprouter instance.
    for _, r := range g.Routes() {
        router.Handle(r.Method(), r.Path(), r.Handler())
    }

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
