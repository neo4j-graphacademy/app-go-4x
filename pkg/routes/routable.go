package routes

import "net/http"

type Routable interface {
	Register(server *http.ServeMux)
}
