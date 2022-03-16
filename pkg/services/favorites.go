package services

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FavoriteService interface {
	Save(userId, movieId string) (Movie, error)

	FindAllByUserId(userId string, page *paging.Paging) ([]Movie, error)

	Delete(userId, movieId string) (Movie, error)
}

type neo4jFavoriteService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewFavoriteService(loader *fixtures.FixtureLoader, driver neo4j.Driver) FavoriteService {
	return &neo4jFavoriteService{loader: loader, driver: driver}
}

// Save should create a `:HAS_FAVORITE` relationship between
// the User and Movie ID nodes provided.
//
// If either the user or movie cannot be found, a `NotFoundError` should be thrown.
// tag::add[]
func (fs *neo4jFavoriteService) Save(userId, movieId string) (_ Movie, err error) {
	// TODO: Open a new Session
	// TODO: Create HAS_FAVORITE relationship within a Write Transaction
	// TODO: Close the session
	// TODO: Return movie details and `favorite` property

	result, err := fs.loader.ReadObject("fixtures/goodfellas.json")
	if err != nil {
		return nil, err
	}
	result["favorite"] = true
	return result, nil
}

// end::add[]

// FindAllByUserId should retrieve a list of movies that have an incoming :HAS_FAVORITE
// relationship from a User node with the supplied `userId`.
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
// tag::all[]
func (fs *neo4jFavoriteService) FindAllByUserId(userId string, page *paging.Paging) (_ []Movie, err error) {
	// TODO: Open a new session
	// TODO: Retrieve a list of movies favorited by the user
	// TODO: Close session

	return fs.loader.ReadArray("fixtures/popular.json")
}

// end::all[]

// Delete should remove the `:HAS_FAVORITE` relationship between
// the User and Movie ID nodes provided.
// If either the user, movie or the relationship between them cannot be found,
// a `NotFoundError` should be thrown.
// tag::remove[]
func (fs *neo4jFavoriteService) Delete(userId, movieId string) (_ Movie, err error) {
	// TODO: Open a new Session
	// TODO: Delete the HAS_FAVORITE relationship within a Write Transaction
	// TODO: Close the session
	// TODO: Return movie details and `favorite` property

	result, err := fs.loader.ReadObject("fixtures/goodfellas.json")
	if err != nil {
		return nil, err
	}
	result["favorite"] = false
	return result, nil
}

// end::remove[]
