package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"strings"
)

type authRoutes struct {
	auth services.AuthService
}

func NewAuthRoutes(auth services.AuthService) Routable {
	return &authRoutes{auth: auth}
}

func (a *authRoutes) Register(server *http.ServeMux) {
	server.HandleFunc("/api/auth/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			switch {
			case strings.HasSuffix(path, "/register"):
				a.Save(request, writer)
			case strings.HasSuffix(path, "/login"):
				a.Login(request, writer)
			}
		})
}

func (a *authRoutes) Save(request *http.Request, writer http.ResponseWriter) {
	userData, err := ioutils.ReadJson(request.Body)
	if err != nil {
		serializeError(writer, err)
		return
	}
	user, err := a.auth.Save(
		userData["email"].(string),
		userData["password"].(string),
		userData["name"].(string),
	)
	serializeJson(writer, user, err)
}

func (a *authRoutes) Login(request *http.Request, writer http.ResponseWriter) {
	userData, err := ioutils.ReadJson(request.Body)
	if err != nil {
		serializeError(writer, err)
		return
	}
	user, err := a.auth.FindOneByEmailAndPassword(
		userData["email"].(string),
		userData["password"].(string),
	)
	serializeJson(writer, user, err)
}
