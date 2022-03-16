package challenges_test

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestGenreDetails(t *testing.T) {
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
	assertNotNil(t, service)

	// Get Genre by Name
	name := "Action"

	genre, err := service.FindOneByName(name)

	assertNilError(t, err)
	assertNotNil(t, genre)

	assertEquals(t, name, genre["name"])

	fmt.Println()
	fmt.Println()

	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("How many movies are in the Action genre?")
	fmt.Println("Copy and paste the following answer into the text box:")

	fmt.Println()
	fmt.Println(genre["movies"])

	fmt.Println()

}
