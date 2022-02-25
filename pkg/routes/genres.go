package routes

import (
	"encoding/json"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"strings"
)

type genreRoutes struct {
	genres services.GenreService
	movies services.MovieService
}

func NewGenreRoutes(genres services.GenreService, movies services.MovieService) Routable {
	return &genreRoutes{genres: genres, movies: movies}
}

func (g *genreRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/genres/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/genres/")
			switch {
			case path == "":
				g.FindAllGenres(writer)
			case strings.HasSuffix(path, "/movies"):
				genre := strings.TrimSuffix(path, "/movies")
				g.FindAllMoviesByGenre(genre, paging.ParsePaging(request), writer)
			default:
				g.FindOneGenreByName(path, writer)
			}
		})
}

func (g *genreRoutes) FindAllGenres(writer http.ResponseWriter) {
	genres, err := g.genres.FindAll()
	handleResponse(writer, genres, err)
}

func (g *genreRoutes) FindAllMoviesByGenre(genre string, page *paging.Paging, writer http.ResponseWriter) {
	// TODO: extract userId
	movies, err := g.movies.FindAllByGenre(genre, "", page)
	handleResponse(writer, movies, err)
}

func (g *genreRoutes) FindOneGenreByName(name string, writer http.ResponseWriter) {
	genre, err := g.genres.FindOneByName(name)
	handleResponse(writer, genre, err)
}

func handleResponse(writer http.ResponseWriter, result interface{}, err error) {
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	genreJson, err := json.Marshal(result)
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	_, _ = writer.Write(genreJson)
}
