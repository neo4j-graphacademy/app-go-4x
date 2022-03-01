package routes

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type accountRoutes struct {
	ratings services.RatingService
	auth    services.AuthService
}

func NewAccountRoutes(ratings services.RatingService,
	auth services.AuthService) Routable {
	return &accountRoutes{
		ratings: ratings,
		auth:    auth,
	}
}

func (a *accountRoutes) Register(server *http.ServeMux) {
	server.HandleFunc("/api/account/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/account/")
			switch {
			case strings.HasPrefix(path, "ratings/"):
				movieId := strings.TrimPrefix(path, "ratings/")
				a.SaveRating(movieId, request, writer)
			}
		})
}

func (a *accountRoutes) SaveRating(movieId string, request *http.Request, writer http.ResponseWriter) {
	ratingData, err := ioutils.ReadJson(request.Body)
	if err != nil {
		serializeError(writer, err)
		return
	}
	userId, err := extractUserId(request, a.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	rating, err := parseIntRating(ratingData["rating"])
	if err != nil {
		writer.WriteHeader(400)
		writer.Header().Set("Content-Type", "text/plain")
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	movie, err := a.ratings.Save(rating, movieId, userId)
	if err != nil {
		serializeError(writer, err)
		return
	}
	serializeJson(writer, movie, err)
}

func extractUserId(request *http.Request, auth services.AuthService) (string, error) {
	bearer := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	return auth.ExtractUserId(bearer)
}

func parseIntRating(rating interface{}) (int, error) {
	if ratingStr, ok := rating.(string); ok {
		return strconv.Atoi(ratingStr)
	}
	if ratingNumber, ok := rating.(float64); ok {
		return int(ratingNumber), nil
	}
	return -1, fmt.Errorf(
		"unsupported rating type: %s, cannot parse",
		reflect.TypeOf(rating),
	)
}
