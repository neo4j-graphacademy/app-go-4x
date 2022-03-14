package challenges_test

import (
	"fmt"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestHandleUniqueConstraints(outer *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(outer, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(outer, err)

	defer func() {
		assertNilError(outer, driver.Close())
	}()

	session := driver.NewSession(neo4j.SessionConfig{})

	// Check Constraint exists
	result, err := session.Run(`SHOW CONSTRAINTS
		YIELD entityType, labelsOrTypes, properties
		WHERE entityType = 'NODE' AND labelsOrTypes = ['User'] AND properties = ['email']
		RETURN count(*) AS count`, map[string]interface{}{})

	assertNilError(outer, err)

	first, err := result.Single()
	assertNilError(outer, err)

	assertEquals(outer, first.Values[0], int64(1))

	// Define variables
	email := "graphacademy@neo4j.com"
	password := "notletmein"
	name := "Graph Academy"

	// Delete any existing user
	session.Run("MATCH (u:User {email: $email}) DETACH DELETE u", map[string]interface{}{"email": email})

	// Create Service
	service := services.NewAuthService(driver, "secret", 10)

	// Create the user
	user, err := service.Save(email, password, name)

	assertNilError(outer, err)
	assertNotNil(outer, user)

	// Attempt to create the user again
	other, err := service.Save(email, password, name)
	fmt.Println(other)
	assertNil(outer, other)
	assertNotNil(outer, err)

	assertContains(outer, err.Error(), "already exists")
}
