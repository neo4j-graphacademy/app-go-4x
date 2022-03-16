package challenges_test

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestHandleUniqueConstraints(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	session := driver.NewSession(neo4j.SessionConfig{})

	// Check Constraint exists
	result, err := session.Run(`SHOW CONSTRAINTS
		YIELD entityType, labelsOrTypes, properties
		WHERE entityType = 'NODE' AND labelsOrTypes = ['User'] AND properties = ['email']
		RETURN count(*) AS count`, map[string]interface{}{})

	assertNilError(t, err)

	first, err := result.Single()
	assertNilError(t, err)

	assertEquals(t, first.Values[0], int64(1))

	// Define variables
	email := "graphacademy@neo4j.com"
	password := "notletmein"
	name := "Graph Academy"

	// Delete any existing user
	session.Run("MATCH (u:User {email: $email}) DETACH DELETE u", map[string]interface{}{"email": email})

	// Create Service
	service := services.NewAuthService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver, "secret", 10)

	// Create the user
	user, err := service.Save(email, password, name)

	assertNilError(t, err)
	assertFalse(t, user == nil)

	// Attempt to create the user again
	other, err := service.Save(email, password, name)
	assertTrue(t, other == nil)
	assertNotNil(t, err)

	assertContains(t, err.Error(), "already exists")
}
