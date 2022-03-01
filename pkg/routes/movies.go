package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"strings"
)

type movieRoutes struct {
	movies  services.MovieService
	ratings services.RatingService
	auth    services.AuthService
}

func NewMovieRoutes(movies services.MovieService,
	ratings services.RatingService,
	auth services.AuthService) Routable {
	return &movieRoutes{
		movies:  movies,
		ratings: ratings,
		auth:    auth,
	}
}

func (m *movieRoutes) Register(server *http.ServeMux) {
	server.HandleFunc("/api/movies/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/movies/")
			switch {
			case path == "":
				m.FindAllMovies(request, writer)
			case strings.HasSuffix(path, "/similar"):
				id := strings.TrimSuffix(path, "/similar")
				m.FindAllMoviesBySimilarity(id, request, writer)
			case strings.HasSuffix(path, "/ratings"):
				id := strings.TrimSuffix(path, "/ratings")
				m.FindAllRatingsByMovieId(id, request, writer)
			default:
				m.FindOneMovieById(path, request, writer)
			}
		})
}

func (m *movieRoutes) FindAllMovies(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	userId, err := extractUserId(request, m.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := m.movies.FindAll(userId, page)
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindOneMovieById(id string, request *http.Request, writer http.ResponseWriter) {
	userId, err := extractUserId(request, m.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := m.movies.FindOneById(id, userId)
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindAllMoviesBySimilarity(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	userId, err := extractUserId(request, m.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := m.movies.FindAllBySimilarity(id, userId, page)
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindAllRatingsByMovieId(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.RatingSortableAttributes())
	movies, err := m.ratings.FindAllByMovieId(id, page)
	serializeJson(writer, movies, err)
}
