package httphandler

import (
	"io"
	"net/http"

	"github.com/joez-tkpd/go-sample-arch/infrastructure/infrahttp"
)

type Handler struct {
	router infrahttp.Router
	user   UserUsecase
	// middleware
}

func NewHandler(user UserUsecase) Handler {
	return Handler{user: user}
}

func (h Handler) Register(router infrahttp.Router) {
	h.router = router

	router.HandleFunc("/", h.YourHandler)
	router.HandleFunc("/health", h.HealthCheckHandler)
	router.HandleFunc("/user/{id}", h.GetUserHandler)
}

func (h Handler) YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func (h Handler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}
