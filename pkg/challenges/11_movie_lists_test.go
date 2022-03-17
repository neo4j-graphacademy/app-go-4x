package challenges_test

import (
	"fmt"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func TestMoviePagination(t *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(t, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(t, err)

	defer func() {
		assertNilError(t, driver.Close())
	}()

	service := services.NewMovieService(
		&fixtures.FixtureLoader{Prefix: "../.."},
		driver)
	assertNotNil(t, service)

	tomHanks := "31"
	coppola := "1776"

	genre := "Comedy"
	movieLimit := 10

	// return a paginated list of movies by Genre
	firstByGenre, err := service.FindAllByGenre(genre, "", paging.NewPaging("", "title", "ASC", 0, movieLimit))

	assertNilError(t, err)
	assertNotNil(t, firstByGenre)
	assertEquals(t, movieLimit, len(firstByGenre))

	// Second Page
	secondByGenre, err := service.FindAllByGenre(genre, "", paging.NewPaging("", "title", "ASC", movieLimit, movieLimit))

	assertNilError(t, err)
	assertNotNil(t, secondByGenre)
	assertEquals(t, movieLimit, len(secondByGenre))
	assertNotEquals(t, firstByGenre[0]["title"], secondByGenre[0]["title"])

	// Reordered
	reorderedByGenre, err := service.FindAllByGenre(genre, "", paging.NewPaging("", "released", "ASC", movieLimit, movieLimit))

	assertNilError(t, err)
	assertEquals(t, movieLimit, len(reorderedByGenre))
	assertNotEquals(t, firstByGenre[0]["title"], reorderedByGenre[0]["title"])

	// return a paginated list of movies by Actor
	actorLimit := 2

	firstByActor, err := service.FindAllByActorId(tomHanks, "", paging.NewPaging("", "title", "ASC", 0, actorLimit))

	assertNilError(t, err)
	assertNotNil(t, firstByActor)
	assertEquals(t, actorLimit, len(firstByActor))

	secondByActor, err := service.FindAllByActorId(tomHanks, "", paging.NewPaging("", "title", "ASC", actorLimit, actorLimit))

	assertNotNil(t, secondByActor)
	assertEquals(t, actorLimit, len(firstByActor))
	assertNotEquals(t, firstByActor[0]["title"], secondByActor[0]["title"])

	// Reordered
	reorderedByActor, err := service.FindAllByActorId(tomHanks, "", paging.NewPaging("", "released", "ASC", 0, actorLimit))

	assertNilError(t, err)
	assertEquals(t, actorLimit, len(reorderedByActor))
	assertNotEquals(t, firstByActor[0]["title"], reorderedByActor[0]["title"])

	// return a paginated list of movies by Director
	directorLimit := 1

	firstByDirector, err := service.FindAllByDirectorId(tomHanks, "", paging.NewPaging("", "title", "ASC", 0, directorLimit))

	assertNilError(t, err)
	assertNotNil(t, firstByDirector)
	assertEquals(t, directorLimit, len(firstByDirector))

	secondByDirector, err := service.FindAllByDirectorId(tomHanks, "", paging.NewPaging("", "title", "ASC", directorLimit, directorLimit))

	assertNotNil(t, secondByDirector)
	assertEquals(t, directorLimit, len(firstByDirector))
	assertNotEquals(t, firstByDirector[0]["title"], secondByDirector[0]["title"])

	// Reordered
	reorderedByDirector, err := service.FindAllByDirectorId(tomHanks, "", paging.NewPaging("", "released", "ASC", 0, directorLimit))

	assertNilError(t, err)
	assertEquals(t, directorLimit, len(reorderedByDirector))
	assertNotEquals(t, firstByDirector[0]["title"], reorderedByDirector[0]["title"])

	// find films directed by Francis Ford Coppola
	copollaFilms, err := service.FindAllByDirectorId(coppola, "", paging.NewPaging("", "title", "ASC", 0, 100))

	assertEquals(t, 16, len(copollaFilms))

	fmt.Println()
	fmt.Println("Here is the answer to the quiz question on the lesson:")
	fmt.Println("How many films has Francis Ford Coppola directed?")
	fmt.Println("Copy and paste the following answer into the text box:")
	fmt.Println()

	fmt.Println(len(copollaFilms))

}
