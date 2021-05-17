package group

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Middleware is just type alias of middleware function.
//
// Middleware は単にミドルウェア関数の型の別名.
type Middleware func(httprouter.Handle) httprouter.Handle

// RouteGroup is struct to retain group information.
//
// RouteGroup はグループ情報を保持する構造体.
type RouteGroup struct {
	children    []*RouteGroup
	handlers    map[string]httprouter.Handle
	middlewares []Middleware
	path        string
}

// New returns a new initialized RouteGroup.
// Define only base path first, other info define by method chaining.
//
// New は初期化された RouteGroup を返す.
// 初めはパスのみ定義し, 他の情報はメソッドチェーンで定義していく.
func New(path string) *RouteGroup {
	return &RouteGroup{
		children:    make([]*RouteGroup, 0),
		handlers:    make(map[string]httprouter.Handle),
		middlewares: make([]Middleware, 0),
		path:        path,
	}
}

// DELETE is a shortcut for (*RouteGroup).Handle(http.MethodDelete, handle).
func (r *RouteGroup) DELETE(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodDelete, handle)
}

// GET is a shortcut for (*RouteGroup).Handle(http.MethodGet, handle).
func (r *RouteGroup) GET(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodGet, handle)
}

// HEAD is a shortcut for (*RouteGroup).Handle(http.MethodHead, handle).
func (r *RouteGroup) HEAD(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodHead, handle)
}

// OPTIONS is a shortcut for (*RouteGroup).Handle(http.MethodOptions, handle).
func (r *RouteGroup) OPTIONS(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodOptions, handle)
}

// PATCH is a shortcut for (*RouteGroup).Handle(http.MethodPatch, handle).
func (r *RouteGroup) PATCH(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPatch, handle)
}

// POST is a shortcut for (*RouteGroup).Handle(http.MethodPost, handle).
func (r *RouteGroup) POST(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPost, handle)
}

// PUT is a shortcut for (*RouteGroup).Handle(http.MethodPut, handle).
func (r *RouteGroup) PUT(handle httprouter.Handle) *RouteGroup {
	return r.Handle(http.MethodPut, handle)
}

// Any relates its handler to all HTTP methods
// defined in `net/http` and returns itself.
// Handler will be ignored if HTTP method is already registered.
//
// Any はハンドラを `net/http` で定義された HTTP メソッド全てに紐付けて自身を返す.
// 登録済みの HTTP メソッドが指定された場合は無視される.
func (r *RouteGroup) Any(handler httprouter.Handle) *RouteGroup {
	return r.Match(
		handler,
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodOptions,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	)
}

// Children returns self includes specified groups.
//
// Children は指定されたグループを含む自身を返す.
func (r *RouteGroup) Children(children ...*RouteGroup) *RouteGroup {
	r.children = append(r.children, children...)

	return r
}

// Handle returns self includes specified pair of HTTP method and handler.
// Handler will be overwritten if a registered HTTP method is specified.
//
// Handle は指定された HTTP メソッドとハンドラのペアを含む自身を返す.
// 登録済みの HTTP メソッドが指定された場合, ハンドラは上書きされる.
func (r *RouteGroup) Handle(method string, handler httprouter.Handle) *RouteGroup {
	r.handlers[method] = handler

	return r
}

// Match relates the handler to all specified HTTP methods and returns itself.
// Handler will be ignored if HTTP method is already registered.
//
// Match はハンドラを指定された HTTP メソッド全てに紐付けて自身を返す.
// 登録済みの HTTP メソッドが指定された場合は無視される.
func (r *RouteGroup) Match(
	handler httprouter.Handle,
	methods ...string,
) *RouteGroup {
	for _, m := range methods {
		if _, exists := r.handlers[m]; !exists {
			r.handlers[m] = handler
		}
	}

	return r
}

// Middleware returns self includes specified middlewares.
//
// Middleware は指定されたミドルウェアを含む自身を返す.
func (r *RouteGroup) Middleware(middlewares ...Middleware) *RouteGroup {
	r.middlewares = append(r.middlewares, middlewares...)

	return r
}

// Routes returns a one-dimensional representation of
// the handlers registered in the parent and child hierarchies.
// The path and middleware are inherited,
// and the middleware is registered in the order of registration of
// the parent hierarchy, followed by the order of registration of
// the child hierarchies.
//
// Routes は自身と子階層に登録済みのハンドラを一次元化して返す.
// パスとミドルウェアは引き継がれ, ミドルウェアは親階層の登録順,
// 続いて子階層の登録順で実行されるように登録される.
func (r *RouteGroup) Routes() Routes {
	return r.routes("", nil)
}

func (r *RouteGroup) len() int {
	cnt := len(r.handlers)

	for _, rg := range r.children {
		cnt += rg.len()
	}

	return cnt
}

func (r *RouteGroup) routes(
	parentPath string,
	parentMiddlewares []Middleware,
) Routes {
	parentMiddlewares = append(parentMiddlewares, r.middlewares...)
	path := joinPath(parentPath, r.path)
	routes := make(Routes, 0, r.len())

	for m, h := range r.handlers {
		routes = append(routes, &Route{
			handler: middlewareWith(h, parentMiddlewares...),
			method:  m,
			path:    path,
		})
	}

	for _, rg := range r.children {
		routes = append(routes, rg.routes(path, parentMiddlewares)...)
	}

	return routes
}

func joinPath(ps ...string) string {
	var buf strings.Builder

	for _, p := range ps {
		if p := strings.Trim(p, "/"); p != "" {
			buf.WriteByte('/')
			buf.WriteString(p)
		}
	}

	if 0 < buf.Len() {
		return buf.String()
	}

	return "/"
}

func middlewareWith(h httprouter.Handle, ms ...Middleware) httprouter.Handle {
	for i := len(ms) - 1; 0 <= i; i-- {
		h = ms[i](h)
	}

	return h
}
