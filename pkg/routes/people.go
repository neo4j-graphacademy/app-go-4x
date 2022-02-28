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
}

func NewPeopleRoutes(people services.PeopleService,
	movies services.MovieService) Routable {
	return &peopleRoutes{
		people: people,
		movies: movies,
	}
}

func (m *peopleRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/people/",
		func(writer http.ResponseWriter, request *http.Request) {
			path := strings.TrimPrefix(request.URL.Path, "/api/people/")
			switch {
			case path == "":
				m.FindAllPeople(request, writer)
			case strings.HasSuffix(path, "/similar"):
				id := strings.TrimSuffix(path, "/similar")
				m.FindAllPeopleBySimilarity(id, request, writer)
			case strings.HasSuffix(path, "/acted"):
				id := strings.TrimSuffix(path, "/acted")
				m.FindAllActedInMovies(id, request, writer)
			case strings.HasSuffix(path, "/directed"):
				id := strings.TrimSuffix(path, "/directed")
				m.FindAllDirectedMovies(id, request, writer)
			default:
				m.FindOnePersonById(path, writer)
			}
		})
}

func (m *peopleRoutes) FindAllPeople(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.PersonSortableAttributes())
	people, err := m.people.FindAll(page)
	serializeJson(writer, people, err)
}

func (m *peopleRoutes) FindOnePersonById(personId string, writer http.ResponseWriter) {
	person, err := m.people.FindOneById(personId)
	serializeJson(writer, person, err)
}

func (m *peopleRoutes) FindAllPeopleBySimilarity(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.PersonSortableAttributes())
	people, err := m.people.FindAllBySimilarity(id, page)
	serializeJson(writer, people, err)
}

func (m *peopleRoutes) FindAllActedInMovies(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	// TODO: retrieve userId
	movies, err := m.movies.FindAllByActorId(id, "", page)
	serializeJson(writer, movies, err)
}

func (m *peopleRoutes) FindAllDirectedMovies(id string, request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.MovieSortableAttributes())
	// TODO: retrieve userId
	movies, err := m.movies.FindAllByDirectorId(id, "", page)
	serializeJson(writer, movies, err)
}
