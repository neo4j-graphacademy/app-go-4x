package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestPersonProfile(t *testing.T) {
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

	coppola := "1776"

	// find a person by their ID

	output, err := service.FindOneById(coppola)

	assertNilError(t, err)
	assertNotNil(t, output)
	assertEquals(t, coppola, output["tmdbId"])
	assertEquals(t, "Francis Ford Coppola", output["name"])
	assertEquals(t, int64(16), output["directedCount"].(int64))
	assertEquals(t, int64(2), output["actedCount"].(int64))

	// return a paginated list of similar people to a person by their ID

	limit := 2

	first, err := service.FindAllBySimilarity(coppola, paging.NewPaging("", "", "", 0, limit))

	assertNilError(t, err)
	assertNotNil(t, first)
	assertEquals(t, limit, len(first))

	second, err := service.FindAllBySimilarity(coppola, paging.NewPaging("", "", "", limit, limit))

	assertNilError(t, err)
	assertNotNil(t, second)
	assertEquals(t, limit, len(second))
	assertNotEquals(t, first[0]["name"], second[0]["name"])

	fmt.Println()
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("According to our algorithm, who is the most similar person to Francis Ford Coppola?")
	fmt.Println("Copy and paste the following answer into the text box:")
	fmt.Println()

	fmt.Println(first[0]["name"])

}
