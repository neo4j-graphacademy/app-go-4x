package services

import (
	"math/rand"

	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Movie = map[string]interface{}

type MovieService interface {
	FindAll(userId string, page *paging.Paging) ([]Movie, error)

	FindAllByGenre(genre, userId string, page *paging.Paging) ([]Movie, error)

	FindAllByActorId(actorId string, userId string, page *paging.Paging) ([]Movie, error)

	FindAllByDirectorId(actorId string, userId string, page *paging.Paging) ([]Movie, error)

	FindOneById(id string, userId string) (Movie, error)

	FindAllBySimilarity(id string, userId string, page *paging.Paging) ([]Movie, error)
}

type neo4jMovieService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewMovieService(loader *fixtures.FixtureLoader, driver neo4j.Driver) MovieService {
	return &neo4jMovieService{loader: loader, driver: driver}
}

// FindAll should return a paginated list of movies ordered by the `sort`
// parameter and limited to the number passed as `limit`.  The `skip` variable should be
// used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
// tag::all[]
func (ms *neo4jMovieService) FindAll(userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Open an Session
	// TODO: Execute a query in a new Read Transaction
	// TODO: Get a list of Movies from the Result
	// TODO: Close the session

	popularMovies, err := ms.loader.ReadArray("fixtures/popular.json")
	if err != nil {
		return nil, err
	}

	return fixtures.Slice(popularMovies, page.Skip(), page.Limit()), err
}

// end::all[]

// FindAllByGenre should return a paginated list of movies that have a relationship to the
// supplied Genre.
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
//
// tag::getByGenre[]
func (ms *neo4jMovieService) FindAllByGenre(genre string, userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Get Movies in a Genre
	// MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})

	popularMovies, err := ms.loader.ReadArray("fixtures/popular.json")
	if err != nil {
		return nil, err
	}
	return fixtures.Slice(popularMovies, page.Skip(), page.Limit()), nil
}

// end::getByGenre[]

// FindAllByActorId should return a paginated list of movies that have an ACTED_IN relationship
// to a Person with the id supplied
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
// tag::getForActor[]
func (ms *neo4jMovieService) FindAllByActorId(actorId string, userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Get Movies acted in by a Person
	// MATCH (:Person {tmdbId: $id})-[:ACTED_IN]->(m:Movie)

	roles, err := ms.loader.ReadArray("fixtures/roles.json")
	if err != nil {
		return nil, err
	}
	return fixtures.Slice(roles, page.Skip(), page.Limit()), nil
}

// end::getForActor[]

// FindAllByDirectorId should return a paginated list of movies that have an DIRECTED relationship
// to a Person with the id supplied
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
// tag::getForDirector[]
func (ms *neo4jMovieService) FindAllByDirectorId(actorId string, userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Get Movies directed by a Person
	// MATCH (:Person {tmdbId: $id})-[:DIRECTED]->(m:Movie)

	popularMovies, err := ms.loader.ReadArray("fixtures/popular.json")
	if err != nil {
		return nil, err
	}
	return fixtures.Slice(popularMovies, page.Skip(), page.Limit()), nil
}

// end::getForDirector[]

// FindOneById finds a Movie node with the ID passed as the `id` parameter.
// Along with the returned payload, a list of actors, directors, and genres should
// be included.
// The number of incoming RATED relationships should also be returned as `ratingCount`
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
// tag::findById[]
func (ms *neo4jMovieService) FindOneById(id string, userId string) (_ Movie, err error) {
	// TODO: Find a movie by its ID
	// MATCH (m:Movie {tmdbId: $id})

	return ms.loader.ReadObject("fixtures/goodfellas.json")
}

// end::findById[]

// FindAllBySimilarity should return a paginated list of similar movies to the Movie with the
// id supplied.  This similarity is calculated by finding movies that have many first
// degree connections in common: Actors, Directors and Genres.
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has added the movie to their "My Favorites" list.
// tag::getSimilarMovies[]
func (ms *neo4jMovieService) FindAllBySimilarity(id string, userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Get similar movies based on genres or ratings
	popularMovies, err := ms.loader.ReadArray("fixtures/popular.json")
	if err != nil {
		return nil, err
	}
	results := fixtures.Slice(popularMovies, page.Skip(), page.Limit())
	for _, movie := range results {
		movie["score"] = rand.Intn(100)
	}
	return results, nil
}

// end::getSimilarMovies[]

// getUserFavorites should return a list of tmdbId properties for the movies that
// the user has added to their 'My Favorites' list.
// tag::getUserFavorites[]
func getUserFavorites(tx neo4j.Transaction, userId string) ([]string, error) {
	return nil, nil
}

// end::getUserFavorites[]
