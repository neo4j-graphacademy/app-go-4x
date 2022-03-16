package services

import (
	"fmt"

	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Genre = map[string]interface{}

type GenreService interface {
	FindAll() ([]Genre, error)

	FindOneByName(name string) (Genre, error)
}

type neo4jGenreService struct {
	loader *fixtures.FixtureLoader
	driver neo4j.Driver
}

func NewGenreService(loader *fixtures.FixtureLoader, driver neo4j.Driver) GenreService {
	return &neo4jGenreService{loader: loader, driver: driver}
}

// FindAll should return a list of genres from the database with a
// `name` property, `movies` which is the count of the incoming `IN_GENRE`
// relationships and a `poster` property to be used as a background.
//
// [
//   {
//    name: 'Action',
//    movies: 1545,
//    poster: 'https://image.tmdb.org/t/p/w440_and_h660_face/qJ2tW6WMUDux911r6m7haRef0WH.jpg'
//   }, ...
//
// ]
//
// tag::all[]
func (gs *neo4jGenreService) FindAll() (_ []Genre, err error) {
	// TODO: Open a new session
	// TODO: Get a list of Genres from the database

	return gs.loader.ReadArray("fixtures/genres.json")
}

// end::all[]

// FindOneByName should find a Genre node by its name and return a set of properties
// along with a `poster` image and `movies` count.
//
// If the genre is not found, an error should be thrown.
// tag::find[]
func (gs *neo4jGenreService) FindOneByName(name string) (_ Genre, err error) {
	// TODO: Open a new session
	// TODO: Get Genre information from the database
	// TODO: Return an error if the genre is not found

	genres, err := gs.loader.ReadArray("fixtures/genres.json")
	if err != nil {
		return nil, err
	}
	for _, genre := range genres {
		if genre["name"] == name {
			return genre, nil
		}
	}
	return nil, fmt.Errorf("genre %q not found", name)
}

// end::find[]
