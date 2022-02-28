package routes

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"net/http"
)

type peopleRoutes struct {
	people services.PeopleService
}

func NewPeopleRoutes(people services.PeopleService) Routable {
	return &peopleRoutes{people: people}
}

func (m *peopleRoutes) AddRoutes(server *http.ServeMux) {
	server.HandleFunc("/api/people/",
		func(writer http.ResponseWriter, request *http.Request) {
			m.FindAllPeople(request, writer)
		})
}

func (m *peopleRoutes) FindAllPeople(request *http.Request, writer http.ResponseWriter) {
	page := paging.ParsePaging(request, paging.PersonSortableAttributes())
	people, err := m.people.FindAll(page)
	serializeJson(writer, people, err)
}
