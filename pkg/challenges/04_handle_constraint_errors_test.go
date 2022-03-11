package challenges_test

import (
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
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

	// Create Service
	service := services.NewAuthService(driver, "secret", 10)

	email := "graphacademy@neo4j.com"
	password := "notletmein"
	name := "Graph Academy"

	user, err := service.Save(email, password, name)

	assertNilError(outer, err)
	assertNotNil(outer, user)

	// Attempt to create the user again
	other, err := service.Save(email, password, name)

	assertNil(outer, other)
	assertNotNil(outer, err)

	assertContains(outer, err.Error(), "already exists")
}
