package services

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/fixtures"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
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
	// Open a new session
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Get a list of Genres from the database
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
		MATCH (g:Genre)
		WHERE g.name <> '(no genres listed)'
		CALL {
			WITH g
			MATCH (g)<-[:IN_GENRE]-(m:Movie)
			WHERE m.imdbRating IS NOT NULL
			AND m.poster IS NOT NULL
			RETURN m.poster AS poster
			ORDER BY m.imdbRating DESC LIMIT 1
		}
		RETURN g {
			.name,
			link: '/genres/'+ g.name,
			poster: poster,
			movies: size( (g)<-[:IN_GENRE]-() )
		} as genre
		ORDER BY g.name ASC`, nil)

		if err != nil {
			return nil, err
		}

		// Collect Results
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}

		// Get genres from results
		var results []map[string]interface{}
		for _, record := range records {
			genre, _ := record.Get("genre")
			results = append(results, genre.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]Genre), nil
}

// end::all[]

// FindOneByName should find a Genre node by its name and return a set of properties
// along with a `poster` image and `movies` count.
//
// If the genre is not found, an error should be thrown.
// tag::find[]
func (gs *neo4jGenreService) FindOneByName(name string) (_ Genre, err error) {
	// Open a new Session
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Get genre information from the database
	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
		MATCH (g:Genre {name: $name})<-[:IN_GENRE]-(m:Movie)
		WHERE m.imdbRating IS NOT NULL
		AND m.poster IS NOT NULL
		AND g.name <> '(no genres listed)'
		WITH g, m
		ORDER BY m.imdbRating DESC

		WITH g, head(collect(m)) AS movie

		RETURN g {
		  link: '/genres/'+ g.name,
		  .name,
		  movies: size((g)<-[:IN_GENRE]-()),
		  poster: movie.poster
		} AS genre`, map[string]interface{}{"name": name})

		if err != nil {
			return nil, err
		}

		// Attempt to get the first and only record
		records, err := result.Single()
		if err != nil {
			return nil, err
		}

		// Get genre information from the first record
		record, _ := records.Get("genre")
		return record, nil
	})

	// Return an error if the genre is not found
	if err != nil {
		return nil, err
	}

	return result.(Genre), nil
}

// end::find[]
