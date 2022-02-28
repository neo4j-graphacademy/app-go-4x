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
}

func NewMovieRoutes(movies services.MovieService,
	ratings services.RatingService) Routable {
	return &movieRoutes{
		movies:  movies,
		ratings: ratings,
	}
}

func (m *movieRoutes) AddRoutes(server *http.ServeMux) {
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
				m.FindOneMovieById(path, writer)
			}
		})
}

func (m *movieRoutes) FindAllMovies(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	// TODO: extract userId
	movies, err := m.movies.FindAll("", page)
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindOneMovieById(id string, writer http.ResponseWriter) {
	// TODO: extract userId
	movies, err := m.movies.FindOneById(id, "")
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindAllMoviesBySimilarity(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	// TODO: extract userId
	movies, err := m.movies.FindAllBySimilarity(id, "", page)
	serializeJson(writer, movies, err)
}

func (m *movieRoutes) FindAllRatingsByMovieId(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.RatingSortableAttributes())
	movies, err := m.ratings.FindAllByMovieId(id, page)
	serializeJson(writer, movies, err)
}
