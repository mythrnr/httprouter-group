package group

import (
	"fmt"
	"sort"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Route is struct to retain routing information.
// Route does not keep the information of
// the parent hierarchy (path, middleware).
//
// Route はルーティングの情報を保持する構造体.
// 親階層の情報（パス, ミドルウェア）は保持しない.
type Route struct {
	handler httprouter.Handle
	method  string
	path    string
}

// Handler just returns handler.
//
// Handler は単にハンドラを返す.
func (r *Route) Handler() httprouter.Handle { return r.handler }

// Method just returns HTTP method.
//
// Method は単に HTTP メソッドを返す.
func (r *Route) Method() string { return r.method }

// Path just returns path.
//
// Path は単にパスを返す.
func (r *Route) Path() string { return r.path }

// Routes is aggregation of `Route`.
//
// Routes は `Route` の集合体.
type Routes []*Route

// String sorts the registered HTTP method and URI pairs by string
// and returns them as a newline-separated string.
// This includes the information of children.
//
// String は登録済みの HTTP メソッドと URI のペアを
// 文字列でソートし, 改行区切りの文字列で返す.
// 自身を起点とした子階層の情報も取得して返す.
//
// Example.
//
//	GET     /
//	GET     /users
//	DELETE  /users/:id
//	GET     /users/:id
//	PUT     /users/:id
func (r Routes) String() string {
	// 8 = utf8.RuneCountInString(http.MethodOptions) + 1.
	const format = "%-8s%s"

	l := make([][2]string, 0, len(r))
	ll := make([]string, 0, len(l))

	for i := range r {
		l = append(l, [2]string{r[i].method, r[i].path})
	}

	sort.Sort(_sort(l))

	for i := range l {
		ll = append(ll, fmt.Sprintf(format, l[i][0], l[i][1]))
	}

	return strings.Join(ll, "\n")
}

// _sort is defined for sorting when stringing `Routes`.
//
// _sort は `Routes` の文字列化のときにソートするための定義.
type _sort [][2]string

// Len is an implementation of `sort.Interface`.
//
// Len は `sort.Interface` の実装.
func (s _sort) Len() int { return len(s) }

// Swap is an implementation of `sort.Interface`.
//
// Swap は `sort.Interface` の実装.
func (s _sort) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is an implementation of `sort.Interface`.
// Sort by path, HTTP method priority.
//
// Less は `sort.Interface` の実装.
// パス, HTTP メソッドの優先順でソートする.
func (s _sort) Less(i, j int) bool {
	if s[i][1] == s[j][1] {
		return s[i][0] < s[j][0]
	}

	return s[i][1] < s[j][1]
}
