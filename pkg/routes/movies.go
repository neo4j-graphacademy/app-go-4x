package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
)

type movieRoutes struct {
	movies services.MovieService
}

func NewMovieRoutes(movies services.MovieService) Routable {
	return &movieRoutes{movies: movies}
}

func (m *movieRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/movies/",
		func(writer http.ResponseWriter, request *http.Request) {
			m.FindAllMovies(request, writer)
		})
}

func (m *movieRoutes) FindAllMovies(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	// TODO: extract userId
	movies, err := m.movies.FindAll("", page)
	serializeJson(writer, movies, err)
}
