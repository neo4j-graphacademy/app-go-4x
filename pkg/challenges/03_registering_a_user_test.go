package challenges_test

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestRegisterUser(outer *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(outer, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(outer, err)

	defer func() {
		assertNilError(outer, driver.Close())
	}()

	// Create Service
	service := services.NewAuthService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver, "secret", 10)

	email := "graphacademy@neo4j.com"
	password := "notletmein"
	name := "Graph Academy"

	user, err := service.Save(email, password, name)

	assertNilError(outer, err)

	// Check return properties
	assertEquals(outer, email, user["email"])
	assertEquals(outer, name, user["name"])
	assertNil(outer, user["password"])

	// Check user in database
	session := driver.NewSession(neo4j.SessionConfig{})

	result, err := session.Run(
		"MATCH (u:User {email: $email}) RETURN u",
		map[string]interface{}{"email": email})

	assertNilError(outer, err)

	assertResultHasNextRecord(outer, result)

	node := result.Record().Values[0].(neo4j.Node)

	assertEquals(outer, email, node.Props["email"])
	assertEquals(outer, name, node.Props["name"])
	assertNotEquals(outer, password, node.Props["password"])
}
