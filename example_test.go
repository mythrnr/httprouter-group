package group_test

import (
	"log"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	group "github.com/mythrnr/httprouter-group"
)

func Example() {
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

	// GET     /
	// GET     /users
	// DELETE  /users/:id
	// GET     /users/:id
	// PUT     /users/:id
	log.Print(strings.Join(g.List(), "\n"))

	// finally, register routes to httprouter instance.
	// g.Register(router)

	// serve.
	// log.Fatal(http.ListenAndServe(":8080", router))
}
