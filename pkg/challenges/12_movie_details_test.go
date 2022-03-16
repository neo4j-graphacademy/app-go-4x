package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestMovieDetails(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	// get a movie by tmdbId
	lockStock := "100"

	service := services.NewMovieService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)
	assertNotNil(t, service)

	movieById, err := service.FindOneById(lockStock, "")

	assertNilError(t, err)
	assertEquals(t, movieById["tmdbId"], lockStock)
	assertEquals(t, movieById["title"], "Lock, Stock & Two Smoking Barrels")

	// get similar movies ordered by similarity score
	limit := 1

	output, err := service.FindAllBySimilarity(lockStock, "", paging.NewPaging("", "title", "ASC", 0, limit))

	assertNilError(t, err)

	paginated, err := service.FindAllBySimilarity(lockStock, "", paging.NewPaging("", "title", "ASC", 1, limit))

	assertNilError(t, err)
	assertNotNil(t, output)
	assertEquals(t, limit, len(output))
	assertEquals(t, limit, len(paginated))
	assertNotEquals(t, paginated[0]["tmdbId"], output[0]["tmdbId"])

	fmt.Println()
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("What is the title of the most similar movie to Lock, Stock & Two Smoking Barrels?")
	fmt.Println("Copy and paste the following answer into the text box:")
	fmt.Println()

	fmt.Println(output[0]["title"])

}
