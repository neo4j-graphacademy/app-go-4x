package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"sort"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestGenreList(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	service := services.NewGenreService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)

	// Should retrieve a list of genres
	output, err := service.FindAll()

	assertNilError(t, err)

	assertEquals(t, len(output), 19)
	assertEquals(t, "Action", output[0]["name"])
	assertEquals(t, "Western", output[18]["name"])

	// Get Genre with the most movies
	sort.Slice(output, func(i, j int) bool {
		return output[i]["movies"].(int64) > output[j]["movies"].(int64)
	})

	// Answer to question
	fmt.Println("")
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("Which genre has the highest movie count?")
	fmt.Println("Copy and paste the following answer into the text box:")
	fmt.Println("")
	fmt.Println(output[0]["name"])
	fmt.Println("")

}
