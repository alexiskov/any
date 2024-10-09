package htpsrv

import "net/http"

type (
	Server struct {
		Socket *http.Server
	}
)
