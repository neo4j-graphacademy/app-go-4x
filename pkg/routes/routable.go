package routes

import "net/http"

type Routable interface {
	AddRoutes(server *http.ServeMux)
}
