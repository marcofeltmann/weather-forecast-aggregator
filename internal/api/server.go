package api

import (
	"log/slog"
	"net/http"
)

/*
Inspiration for the server type is token from Mat Ryer's talk at GopherCon 2019:

"How I write HTTP Web Services after 8 years", watch it on YouTube:
https://www.youtube.com/watch?v=rWBSMsLG8po

But he's constantly updating it kinda each year like on the Grafana blog:
https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
*/
type Server struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

/*
NewServer returns a configured API server.

Currently the server itself only contains references, so it's safe to return a
value type.
*/
func NewServer(logger *slog.Logger) Server {
	/*
		As http.DefaultServeMux is a package-global reference type variable every
		third-party package might manipulate it. Even a different package in this
		codebase. Prometheus for instance registers it's /metrics handler there.

		So creating a new ServeMux type prevents unexpected side effects.
	*/
	mux := http.NewServeMux()
	addRoutes(mux, logger)
	return Server{
		logger: logger,
		mux:    mux,
	}
}

func (s Server) Handler() http.Handler {
	return s.mux
}
