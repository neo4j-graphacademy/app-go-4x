package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"strings"
)

type genreRoutes struct {
	genres services.GenreService
	movies services.MovieService
	auth   services.AuthService
}

func NewGenreRoutes(genres services.GenreService,
	movies services.MovieService,
	auth services.AuthService) Routable {

	return &genreRoutes{
		genres: genres,
		movies: movies,
		auth:   auth,
	}
}

func (g *genreRoutes) Register(server *http.ServeMux) {
	server.HandleFunc("/api/genres/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/genres/")
			switch {
			case path == "":
				g.FindAllGenres(writer)
			case strings.HasSuffix(path, "/movies"):
				genre := strings.TrimSuffix(path, "/movies")
				pagingParams := paging.ParsePaging(request, paging.MovieSortableAttributes())
				g.FindAllMoviesByGenre(genre, pagingParams, request, writer)
			default:
				g.FindOneGenreByName(path, writer)
			}
		})
}

func (g *genreRoutes) FindAllGenres(writer http.ResponseWriter) {
	genres, err := g.genres.FindAll()
	serializeJson(writer, genres, err)
}

func (g *genreRoutes) FindAllMoviesByGenre(genre string,
	page *paging.Paging,
	request *http.Request,
	writer http.ResponseWriter) {

	userId, err := extractUserId(request, g.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := g.movies.FindAllByGenre(genre, userId, page)
	serializeJson(writer, movies, err)
}

func (g *genreRoutes) FindOneGenreByName(name string, writer http.ResponseWriter) {
	genre, err := g.genres.FindOneByName(name)
	serializeJson(writer, genre, err)
}
