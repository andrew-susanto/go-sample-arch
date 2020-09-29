package infrahttp

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Rmux struct {
	router *mux.Router
}

func NewRmux() Rmux {
	return Rmux{
		router: mux.NewRouter(),
	}
}

func (svc Rmux) GetHandler() http.Handler {
	return svc.router
}

func (svc Rmux) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) {
	svc.router.HandleFunc(path, f)
}

func (svc Rmux) GetRouteParams(r *http.Request) map[string]string {
	return mux.Vars(r)
}
