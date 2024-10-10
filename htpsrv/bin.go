package htpsrv

import (
	"fmt"
	"net/http"
	"net/url"
	"project/logger"
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

	s.Socket.Get("/*", processinGetQuery)
	s.Socket.Post("/*", processingPostQuery)

	err = http.ListenAndServe(":"+strconv.Itoa(port), s.Socket)
	return
}

func processinGetQuery(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		logger.Debug(fmt.Errorf("get query url parse error: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch u.Path {
	case "/test":
		w.Write([]byte("isOK"))
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func processingPostQuery(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		logger.Debug(fmt.Errorf("post query url parse error: %w", err).Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch u.Path {
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
