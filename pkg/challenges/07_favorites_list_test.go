package challenges_test

import (
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestMyFavoritesList(t *testing.T) {
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
	service := services.NewFavoriteService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)

	assertNotNil(t, service)

	toyStory := "862"
	goodfellas := "769"
	userId := "9f965bf6-7e32-4afb-893f-756f502b2c2a"
	email := "graphacademy.favorite@neo4j.com"

	// Create the User
	session := driver.NewSession(neo4j.SessionConfig{})
	session.Run(`
		MERGE (u:User {userId: $userId}) SET u.email = $email
		FOREACH (r IN [(u)-[r:HAS_FAVORITE]->() | r ] | DELETE r)
	`, map[string]interface{}{"userId": userId, "email": email})

	// Should throw an error if user or movie do not exist
	unknown, err := service.Save("unknown", "x999")
	assertFalse(t, unknown != nil)
	assertNotNil(t, err)

	unknownMovie, err := service.Save(userId, "x999")
	assertFalse(t, unknownMovie != nil)
	assertNotNil(t, err)

	unknownUser, err := service.Save("unknown", toyStory)
	assertFalse(t, unknownUser != nil)
	assertNotNil(t, err)

	// Add to list
	saved, err := service.Save(userId, toyStory)
	assertNilError(t, err)
	assertEquals(t, toyStory, saved["tmdbId"])
	assertEquals(t, true, saved["favorite"])

	all, err := service.FindAllByUserId(userId, paging.NewPaging("", "createdAt", "desc", 0, 1))

	assertNilError(t, err)
	assertEquals(t, len(all), 1)
	assertEquals(t, all[0]["tmdbId"], toyStory)

	remove, err := service.Delete(userId, toyStory)
	assertNilError(t, err)
	assertEquals(t, toyStory, remove["tmdbId"])
	assertEquals(t, false, remove["favorite"])

	// Add & Remove from list
	add, err := service.Save(userId, goodfellas)

	assertNilError(t, err)
	assertEquals(t, goodfellas, add["tmdbId"])
	assertEquals(t, true, add["favorite"])

	removeGoodfellas, err := service.Delete(userId, goodfellas)
	assertNilError(t, err)
	assertEquals(t, goodfellas, removeGoodfellas["tmdbId"])
	assertEquals(t, false, removeGoodfellas["favorite"])

	// Re-add the Toy Story Favorite for test
	readd, err := service.Save(userId, toyStory)
	assertNilError(t, err)
	assertNotNil(t, readd)
}
