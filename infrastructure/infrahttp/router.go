package infrahttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

//go:generate mockgen -source=./router.go -destination=./router_mock.go -package=infrahttp

type router struct {
	router *mux.Router
}

type Router interface {
	GetHandler() http.Handler
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request))
	GetRouteParams(r *http.Request) map[string]string
}

func NewRouter() Router {
	return router{
		router: mux.NewRouter(),
	}
}

func (svc router) GetHandler() http.Handler {
	return svc.router
}

func (svc router) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	svc.router.HandleFunc(path, f)
}

func (svc router) GetRouteParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}
