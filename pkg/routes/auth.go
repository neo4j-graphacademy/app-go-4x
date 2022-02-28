package routes

import (
	"encoding/json"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"io/ioutil"
	"net/http"
	"strings"
)

type authRoutes struct {
	auth services.AuthService
}

func NewAuthRoutes(auth services.AuthService) Routable {
	return &authRoutes{auth: auth}
}

func (a *authRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/auth/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			switch {
			case strings.HasSuffix(path, "/register"):
				a.Register(request, writer)
			case strings.HasSuffix(path, "/login"):
				a.Login(request, writer)
			}
		})
}

func (a *authRoutes) Register(request *http.Request, writer http.ResponseWriter) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		serializeError(writer, err)
		return
	}
	var userData map[string]interface{}
	if err = json.Unmarshal(body, &userData); err != nil {
		serializeError(writer, err)
		return
	}
	user, err := a.auth.Register(
		userData["email"].(string),
		userData["password"].(string),
		userData["name"].(string),
	)
	serializeJson(writer, user, err)
}

func (a *authRoutes) Login(request *http.Request, writer http.ResponseWriter) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		serializeError(writer, err)
		return
	}
	var userData map[string]interface{}
	if err = json.Unmarshal(body, &userData); err != nil {
		serializeError(writer, err)
		return
	}
	user, err := a.auth.LogIn(
		userData["email"].(string),
		userData["password"].(string),
	)
	serializeJson(writer, user, err)
}
