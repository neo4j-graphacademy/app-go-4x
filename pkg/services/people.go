package services

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Person = map[string]interface{}

type PeopleService interface {
	FindAll(page *paging.Paging) ([]Person, error)

	FindOneById(id string) (Person, error)

	FindAllBySimilarity(id string, page *paging.Paging) ([]Person, error)
}

type neo4jPeopleService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewPeopleService(loader *fixtures.FixtureLoader, driver neo4j.Driver) PeopleService {
	return &neo4jPeopleService{loader: loader, driver: driver}
}

// FindAll should return a paginated list of People (actors or directors),
// with an optional filter on the person's name based on the `q` parameter.
//
// Results should be ordered by the `sort` parameter and limited to the
// number passed as `limit`.  The `skip` variable should be used to skip a
// certain number of rows.
// tag::all[]
func (ps *neo4jPeopleService) FindAll(page *paging.Paging) (_ []Person, err error) {
	// TODO: Get a list of people from the database

	people, err := ps.loader.ReadArray("fixtures/people.json")
	if err != nil {
		return nil, err
	}
	return fixtures.Slice(people, page.Skip(), page.Limit()), nil
}

//end::all[]

// FindOneById finds a user by their ID.
// If no user is found, an error should be thrown.
// tag::findById[]
func (ps *neo4jPeopleService) FindOneById(id string) (_ Person, err error) {
	// TODO: Find a user by their ID

	return ps.loader.ReadObject("fixtures/pacino.json")
}

//end::findById[]

// FindAllBySimilarity gets a list of similar people to a Person, ordered by their similarity score
// in descending order.
// tag::getSimilarPeople[]
func (ps *neo4jPeopleService) FindAllBySimilarity(id string, page *paging.Paging) (_ []Person, err error) {
	// TODO: Get a list of similar people to the person by their id
	people, err := ps.loader.ReadArray("fixtures/people.json")
	if err != nil {
		return nil, err
	}
	return fixtures.Slice(people, page.Skip(), page.Limit()), nil
}

// end::getSimilarPeople[]
