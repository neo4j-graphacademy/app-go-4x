package services

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"

	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Rating = map[string]interface{}

type RatingService interface {
	FindAllByMovieId(id string, page *paging.Paging) ([]Rating, error)

	Save(rating int, movieId string, userId string) (Movie, error)
}

type neo4jRatingService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewRatingService(loader *fixtures.FixtureLoader, driver neo4j.Driver) RatingService {
	return &neo4jRatingService{loader: loader, driver: driver}
}

// FindAllByMovieId returns a paginated list of reviews for a Movie.
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
// tag::forMovie[]
func (rs *neo4jRatingService) FindAllByMovieId(movieId string, page *paging.Paging) (_ []Rating, err error) {
	return rs.loader.ReadArray("fixtures/ratings.json")
}

// end::forMovie[]

// Save adds a relationship between a User and Movie with a `rating` property.
// The `rating` parameter should be converted to a Neo4j Integer.
//
// If the User or Movie cannot be found, a NotFoundError should be thrown
// tag::add[]
func (rs *neo4jRatingService) Save(rating int, movieId string, userId string) (_ Movie, err error) {
	// TODO: Open a new session
	// TODO: Save the rating in the database
	// TODO: Return movie details and a rating

	return rs.loader.ReadObject("fixtures/goodfellas.json")
}

// end::add[]
