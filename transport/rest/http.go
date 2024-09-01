package rest

import "net/http"

type Server struct {
	Router http.ServeMux
	Server http.Server
}
