package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
	"strings"
)

type peopleRoutes struct {
	people services.PeopleService
	movies services.MovieService
	auth   services.AuthService
}

func NewPeopleRoutes(people services.PeopleService,
	movies services.MovieService,
	auth services.AuthService) Routable {
	return &peopleRoutes{
		people: people,
		movies: movies,
		auth:   auth,
	}
}

func (p *peopleRoutes) Register(server *http.ServeMux) {
	server.HandleFunc("/api/people/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/people/")
			switch {
			case path == "":
				p.FindAllPeople(request, writer)
			case strings.HasSuffix(path, "/similar"):
				id := strings.TrimSuffix(path, "/similar")
				p.FindAllPeopleBySimilarity(id, request, writer)
			case strings.HasSuffix(path, "/acted"):
				id := strings.TrimSuffix(path, "/acted")
				p.FindAllActedInMovies(id, request, writer)
			case strings.HasSuffix(path, "/directed"):
				id := strings.TrimSuffix(path, "/directed")
				p.FindAllDirectedMovies(id, request, writer)
			default:
				p.FindOnePersonById(path, writer)
			}
		})
}

func (p *peopleRoutes) FindAllPeople(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.PersonSortableAttributes())
	people, err := p.people.FindAll(page)
	serializeJson(writer, people, err)
}

func (p *peopleRoutes) FindOnePersonById(personId string, writer http.ResponseWriter) {
	person, err := p.people.FindOneById(personId)
	serializeJson(writer, person, err)
}

func (p *peopleRoutes) FindAllPeopleBySimilarity(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.PersonSortableAttributes())
	people, err := p.people.FindAllBySimilarity(id, page)
	serializeJson(writer, people, err)
}

func (p *peopleRoutes) FindAllActedInMovies(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	userId, err := extractUserId(request, p.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := p.movies.FindAllByActorId(id, userId, page)
	serializeJson(writer, movies, err)
}

func (p *peopleRoutes) FindAllDirectedMovies(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	userId, err := extractUserId(request, p.auth)
	if err != nil {
		serializeError(writer, err)
		return
	}
	movies, err := p.movies.FindAllByDirectorId(id, userId, page)
	serializeJson(writer, movies, err)
}
