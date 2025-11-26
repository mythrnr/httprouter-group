package group_test

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	group "github.com/mythrnr/httprouter-group"
)

//nolint:testableexamples
func Example() {
	// first, define routes, handlers, and middlewares.
	g := group.New("/").GET(
		func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
			func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
				func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
					w.Write([]byte("GET /users/:id\n"))
				},
			),
		),
		group.New("/users/:id").PUT(
			func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
				w.Write([]byte("PUT /users/:id\n"))
			},
		).DELETE(
			func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
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
	router.PanicHandler = func(_ http.ResponseWriter, _ *http.Request, rec any) {
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
	// for _, r := range g.Routes() {
	// 	router.Handle(r.Method(), r.Path(), r.Handler())
	// }

	// serve.
	// log.Fatal(http.ListenAndServe(":8080", router))
}
