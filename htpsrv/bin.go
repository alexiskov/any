package htpsrv

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type (
	Server struct {
		Socket *chi.Mux
	}
)

func Start(port int) (s Server, err error) {
	s.Socket = chi.NewRouter()

	s.Socket.Get("/", processinGetQuery)
	s.Socket.Post("/", processingPostQuery)

	err = http.ListenAndServe(":"+strconv.Itoa(port), s.Socket)
	return
}

func processinGetQuery(w http.ResponseWriter, r *http.Request) {

}

func processingPostQuery(w http.ResponseWriter, r *http.Request) {

}
