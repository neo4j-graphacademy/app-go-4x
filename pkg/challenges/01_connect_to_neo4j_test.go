package challenges_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
)

func TestNeo4jConnection(outer *testing.T) {
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(outer, err)

	driver, err := config.NewDriver(settings)
	assertNilError(outer, err)
	defer func() {
		assertNilError(outer, driver.Close())
	}()

	outer.Run("Should create a driver instance and connect to server", func(t *testing.T) {
		assertStringNotEmpty(t, settings.Uri)
		assertStringNotEmpty(t, settings.Username)
		assertStringNotEmpty(t, settings.Password)
	})

	outer.Run("Driver has been instantiated", func(t *testing.T) {
		assertNotNil(t, driver)

		configuredUri := driver.Target()
		if !strings.HasPrefix(configuredUri.Scheme, "neo4j") &&
			!strings.HasPrefix(configuredUri.Scheme, "bolt") {
			t.Fatalf("expected URI %s to start with bolt or neo4j scheme",
				configuredUri.String())
		}

		fmt.Println("There are two tests in this suite")
	})
}
