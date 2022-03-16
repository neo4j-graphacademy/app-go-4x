package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestPersonList(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	service := services.NewPeopleService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)
	assertNotNil(t, service)

	// retrieve a paginated list people from the database
	limit := 10

	output, err := service.FindAll(paging.NewPaging("", "name", "asc", 0, limit))

	assertNilError(t, err)
	assertNotNil(t, output)
	assertEquals(t, limit, len(output))

	paginated, err := service.FindAll(paging.NewPaging("", "name", "asc", limit, limit))

	assertNilError(t, err)
	assertNotNil(t, paginated)
	assertEquals(t, limit, len(paginated))

	assertNotEquals(t, output[0]["name"], paginated[0]["name"])

	// apply a filter, ordering and pagination to the query
	q := "A"

	filteredFirst, err := service.FindAll(paging.NewPaging(q, "name", "asc", 0, 1))

	assertNilError(t, err)
	assertNotNil(t, filteredFirst)
	assertEquals(t, 1, len(filteredFirst))

	filteredLast, err := service.FindAll(paging.NewPaging(q, "name", "desc", 0, 1))

	assertNilError(t, err)
	assertNotNil(t, filteredLast)
	assertEquals(t, 1, len(filteredLast))

	assertNotEquals(t, filteredLast[0]["name"], filteredFirst[0]["name"])

	// Quiz answer
	fmt.Println()
	fmt.Println()
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("What is the name of the first person in the database in alphabetical order?")
	fmt.Println("Copy and paste the following answer into the text box:")

	fmt.Println(output[0]["name"])
}
