package challenges_test

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestAuthentication(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	// Create Service
	service := services.NewAuthService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver, "secret", 10)

	email := "authenticated@neo4j.com"
	password := "AuthenticateM3!"
	name := "Authenticated User"

	// Delete any existing User
	session := driver.NewSession(neo4j.SessionConfig{})
	session.Run("MATCH (u:User {email: $email}) DETACH DELETE u", map[string]interface{}{"email": email})

	// Create User
	user, err := service.Save(email, password, name)

	assertNilError(t, err)
	assertEquals(t, email, user["email"])

	// Incorrect Username
	incorrectUsername, err := service.FindOneByEmailAndPassword("unknown", "password")
	assertTrue(t, incorrectUsername == nil)
	assertNotNil(t, err)

	// Incorrect Password
	incorrectPassword, err := service.FindOneByEmailAndPassword(email, "incorrectpassword")
	assertTrue(t, incorrectPassword == nil)
	assertNotNil(t, err)

	// Correct
	correct, err := service.FindOneByEmailAndPassword(email, password)

	assertNilError(t, err)
	assertEquals(t, correct["email"], email)
	assertEquals(t, correct["name"], name)
	assertNotNil(t, correct["token"])

	// GA: set a timestamp to verify that the tests have passed
	session.Run("MATCH (u:User {email: $email}) SET u.authenticatedAt = datetime()", map[string]interface{}{"email": email})

}
