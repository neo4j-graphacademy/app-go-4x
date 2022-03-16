package challenges_test

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestFavoritesFlag(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	fixtureLoader := &fixtures.FixtureLoader{Prefix: "../.."}
	// Create Services
	favoriteService := services.NewFavoriteService(
		fixtureLoader,
		driver)
	movieService := services.NewMovieService(
		fixtureLoader,
		driver)

	assertNotNil(t, favoriteService)
	assertNotNil(t, movieService)

	userId := "fe770c6b-4034-4e07-8e40-2f39e7a6722c"
	email := "graphacademy.flag@neo4j.com"

	// Create the User
	session := driver.NewSession(neo4j.SessionConfig{})
	session.Run(`
		MERGE (u:User {userId: $userId}) SET u.email = $email
		FOREACH (r IN [(u)-[r:HAS_FAVORITE]->() | r ] | DELETE r)
	`, map[string]interface{}{"userId": userId, "email": email})

	// Get the most popular movie
	firstCall, err := movieService.FindAll(userId, paging.NewPaging("", "imdbRating", "DESC", 0, 1))

	assertNilError(t, err)
	assertNotNil(t, firstCall)

	movieId := firstCall[0]["tmdbId"].(string)
	assertEquals(t, false, firstCall[0]["favorite"])

	// Add it to user favorites
	favorite, err := favoriteService.Save(userId, movieId)

	assertNilError(t, err)

	assertEquals(t, movieId, favorite["tmdbId"])
	assertEquals(t, true, favorite["favorite"])

	// Get most popular movie again
	secondCall, err := movieService.FindAll(userId, paging.NewPaging("", "imdbRating", "DESC", 0, 1))

	assertNilError(t, err)
	assertNotNil(t, secondCall)

	assertEquals(t, movieId, secondCall[0]["tmdbId"])
	assertEquals(t, true, secondCall[0]["favorite"])
}
