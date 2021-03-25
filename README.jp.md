# httprouter-group

[English](./README.md)

## Status

[![Check codes](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/check_code.yml)

[![Create Release](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yml/badge.svg)](https://github.com/mythrnr/httprouter-group/actions/workflows/release.yml)

## Description

[`mythrnr/httprouter-group`](https://github.com/mythrnr/httprouter-group) は
[`julienschmidt/httprouter`](https://github.com/julienschmidt/httprouter) にルーティングをグループ化して定義する機能を追加する.

詳細は [GoDoc](https://pkg.go.dev/github.com/mythrnr/httprouter-group) を参照

## Feature

### [`RouteGroup`](https://github.com/mythrnr/httprouter-group/blob/master/route_group.go)

- `(*RouteGroup).Children(...*RouteGroup)` で登録した `RouteGroup`
  は親のパスとミドルウェアを引き継ぐことができる.
- `mythrnr/httprouter-group` は定義を簡素化するためのもので, 次の例の通り
  `julienschmidt/httprouter` のパフォーマンスには影響を与えない.

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
