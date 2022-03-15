package services

import (
	"fmt"

	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
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
	driver neo4j.Driver
}

func NewPeopleService(driver neo4j.Driver) PeopleService {
	return &neo4jPeopleService{driver: driver}
}

// FindAll should return a paginated list of People (actors or directors),
// with an optional filter on the person's name based on the `q` parameter.
//
// Results should be ordered by the `sort` parameter and limited to the
// number passed as `limit`.  The `skip` variable should be used to skip a
// certain number of rows.
// tag::all[]
func (n *neo4jPeopleService) FindAll(page *paging.Paging) (_ []Person, err error) {
	// TODO: Get a list of people from the database

	// people, err := fixtures.ReadArray("fixtures/people.json")
	// if err != nil {
	// 	return nil, err
	// }
	// return fixtures.Slice(people, page.Skip(), page.Limit()), nil

	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(fmt.Sprintf(`
			MATCH (p:Person)
			WHERE $q IS NULL OR toLower(p.name) CONTAINS toLower($q)
			RETURN p { .* } AS person
			ORDER BY p.`+"`%s`"+` %s
			SKIP $skip
			LIMIT $limit`, page.Sort(), page.Order()),
			map[string]interface{}{
				"q":     page.Query(),
				"skip":  page.Skip(),
				"limit": page.Limit(),
			})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			person, _ := record.Get("person")
			results = append(results, person.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]Person), nil
}

//end::all[]

// FindOneById finds a user by their ID.
// If no user is found, an error should be thrown.
// tag::findById[]
func (n *neo4jPeopleService) FindOneById(id string) (_ Person, err error) {
	// TODO: Find a user by their ID

	// return fixtures.ReadObject("fixtures/pacino.json")

	// Open a new database session
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Get a person from the database
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (p:Person {tmdbId: $id})
				RETURN p {
					.*,
					actedCount: size((p)-[:ACTED_IN]->()),
					directedCount: size((p)-[:DIRECTED]->())
				} AS person`,
			map[string]interface{}{"id": id})
		if err != nil {
			return nil, err
		}
		record, err := result.Single()
		if err != nil {
			return nil, err
		}
		person, _ := record.Get("person")
		return person.(map[string]interface{}), nil
	})

	if err != nil {
		return nil, err
	}
	return result.(Person), nil
}

//end::findById[]

// FindAllBySimilarity gets a list of similar people to a Person, ordered by their similarity score
// in descending order.
// tag::getSimilarPeople[]
func (n *neo4jPeopleService) FindAllBySimilarity(id string, page *paging.Paging) (_ []Person, err error) {
	// TODO: Get a list of similar people to the person by their id
	// people, err := fixtures.ReadArray("fixtures/people.json")
	// if err != nil {
	// 	return nil, err
	// }
	// return fixtures.Slice(people, page.Skip(), page.Limit()), nil

	// Open a new database session
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Get a list of similar people to the person by their id
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
				MATCH (:Person {tmdbId: $id})-[:ACTED_IN|DIRECTED]->(m)<-[r:ACTED_IN|DIRECTED]-(p)
				RETURN p {
					.*,
					actedCount: size((p)-[:ACTED_IN]->()),
					directedCount: size((p)-[:DIRECTED]->()),
					inCommon: collect(m {.tmdbId, .title, type: type(r)})
				} AS person
				ORDER BY size(person.inCommon) DESC
				SKIP $skip
				LIMIT $limit`,
			map[string]interface{}{
				"id":    id,
				"skip":  page.Skip(),
				"limit": page.Limit(),
			})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			person, _ := record.Get("person")
			results = append(results, person.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]Person), nil
}

// end::getSimilarPeople[]
