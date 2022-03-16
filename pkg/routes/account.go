package routes

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

type accountRoutes struct {
	ratings   services.RatingService
	auth      services.AuthService
	favorites services.FavoriteService
}

func NewAccountRoutes(ratings services.RatingService,
	auth services.AuthService,
	favorites services.FavoriteService) Routable {
	return &accountRoutes{
		ratings:   ratings,
		auth:      auth,
		favorites: favorites,
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
			case strings.HasPrefix(path, "favorites/"):
				movieId := strings.TrimPrefix(path, "favorites/")
				switch request.Method {
				case "POST":
					a.SaveFavorite(movieId, request, writer)
				case "DELETE":
					a.DeleteFavorite(movieId, request, writer)
				}

			case path == "favorites":
				page := paging.ParsePaging(request, paging.MovieSortableAttributes())
				a.FindAllFavorites(page, request, writer)
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
	serializeJson(writer, movie, err)
}

func (a *accountRoutes) SaveFavorite(movieId string, request *http.Request, writer http.ResponseWriter) {
	userId, err := extractUserId(request, a.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movie, err := a.favorites.Save(userId, movieId)
	serializeJson(writer, movie, err)
}

func (a *accountRoutes) FindAllFavorites(page *paging.Paging, request *http.Request, writer http.ResponseWriter) {
	userId, err := extractUserId(request, a.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := a.favorites.FindAllByUserId(userId, page)
	serializeJson(writer, movies, err)
}

func (a *accountRoutes) DeleteFavorite(movieId string, request *http.Request, writer http.ResponseWriter) {
	userId, err := extractUserId(request, a.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movie, err := a.favorites.Delete(userId, movieId)
	serializeJson(writer, movie, err)
}

func extractUserId(request *http.Request, auth services.AuthService) (string, error) {
	bearer := strings.TrimPrefix(request.Header.Get("Authorization"), "Bearer ")
	// FIXME remove once frontend bug fixed
	if bearer == "undefined" {
		bearer = ""
	}
	return auth.ExtractUserId(bearer)
}

// FIXME remove once frontend bug fixed - rating should always be a number
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
