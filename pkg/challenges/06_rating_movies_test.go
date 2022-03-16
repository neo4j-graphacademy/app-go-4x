package challenges_test

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestRatingMovies(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	// Create Services
	service := services.NewRatingService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)

	movieId := "769"
	userId := "1185150b-9e81-46a2-a1d3-eb649544b9c4"
	email := "graphacademy.reviewer@neo4j.com"
	rating := 5

	// Create the User
	session := driver.NewSession(neo4j.SessionConfig{})
	session.Run("MERGE (u:User {userId: $userId}) SET u.email = $email", map[string]interface{}{"userId": userId, "email": email})

	// Create the rating
	output, err := service.Save(rating, movieId, userId)

	assertNilError(t, err)
	assertEquals(t, movieId, output["tmdbId"])
	assertEquals(t, int64(rating), output["rating"])
}
