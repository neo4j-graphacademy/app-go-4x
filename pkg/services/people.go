package services

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Person = map[string]interface{}

type PeopleService interface {
	FindAll(page *paging.Paging) ([]Movie, error)
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
func (n *neo4jPeopleService) FindAll(page *paging.Paging) (persons []Person, err error) {
	session := n.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(fmt.Sprintf(`
			MATCH (p:Person)
			WHERE $q IS null OR p.name CONTAINS $q
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
	persons = results.([]Person)
	return persons, nil
}

//end::all[]
