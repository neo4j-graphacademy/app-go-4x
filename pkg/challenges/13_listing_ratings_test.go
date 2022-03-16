package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestListingRatings(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	// retrieve a list of ratings from the database
	pulpFiction := "680"
	limit := 10

	service := services.NewRatingService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)
	assertNotNil(t, service)

	first, err := service.FindAllByMovieId(pulpFiction, paging.NewPaging("", "timestamp", "ASC", 0, limit))

	assertNilError(t, err)
	assertNotNil(t, first)
	assertEquals(t, limit, len(first))

	paginated, err := service.FindAllByMovieId(pulpFiction, paging.NewPaging("", "timestamp", "ASC", limit, limit))

	assertNilError(t, err)
	assertNotNil(t, paginated)
	assertEquals(t, limit, len(paginated))

	assertNotEquals(t, first[0]["rating"], paginated[0]["rating"])

	// apply an ordering and pagination to the query
	latest, err := service.FindAllByMovieId(pulpFiction, paging.NewPaging("", "timestamp", "DESC", 0, limit))

	assertNotEquals(t, latest[0]["rating"], first[0]["rating"])

	fmt.Println()
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("What is the name of the first person to rate the movie Pulp Fiction?")
	fmt.Println("Copy and paste the following answer into the text box:")
	fmt.Println()

	firstReview := first[0]
	firstUser := firstReview["user"].(services.Rating)

	fmt.Println(firstUser["name"])

}
