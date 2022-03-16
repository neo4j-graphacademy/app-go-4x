package challenges_test

import (
	"fmt"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestMovieList(outer *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(outer, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(outer, err)

	defer func() {
		assertNilError(outer, driver.Close())
	}()

	service := services.NewMovieService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)

	limit := 1

	output, err := service.FindAll("", paging.NewPaging("", "title", "ASC", 0, limit))
	assertNilError(outer, err)

	assertEquals(outer, len(output), limit)

	// Test Pagination
	next, err := service.FindAll("", paging.NewPaging("", "title", "ASC", 1, limit))

	assertNilError(outer, err)
	assertEquals(outer, len(output), limit)
	assertNotEquals(outer, next[0]["title"], output[0]["title"])

	// Test Ordering
	ordered, err := service.FindAll("", paging.NewPaging("", "imdbRating", "DESC", 0, limit))

	assertNilError(outer, err)
	assertEquals(outer, len(output), limit)
	assertNotEquals(outer, ordered[0]["title"], output[0]["title"])

	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("What is the title of the highest rated movie in the recommendations dataset?")
	fmt.Println("Copy and paste the following answer into the text box:")

	fmt.Println(ordered[0]["title"])
}
