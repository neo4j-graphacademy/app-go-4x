package routes

import (
	"encoding/json"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
)

type genreRoutes struct {
	service services.GenreService
}

func NewGenreRoutes(service services.GenreService) Routable {
	return &genreRoutes{service: service}
}

func (g *genreRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/genres", func(writer http.ResponseWriter, _ *http.Request) {
		g.FindAll(writer)
	})
}

func (g *genreRoutes) FindAll(writer http.ResponseWriter) {
	genres, err := g.service.FindAll()
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	genreJson, err := json.Marshal(genres)
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	_, _ = writer.Write(genreJson)
}
