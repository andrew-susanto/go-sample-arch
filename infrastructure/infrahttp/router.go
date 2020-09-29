package infrahttp

import (
	"net/http"
)

//go:generate mockgen -source=./router.go -destination=./router_mock.go -package=infrahttp

type Router interface {
	GetHandler() http.Handler
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request))
	GetRouteParams(r *http.Request) map[string]string
}

func NewRouter() Router {
	return NewRmux()
}
