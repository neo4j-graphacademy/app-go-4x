package services

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Movie = map[string]interface{}

type MovieService interface {
	FindAll(userId string, page *paging.Paging) ([]Movie, error)

	FindAllByGenre(genre, userId string, page *paging.Paging) ([]Movie, error)

	FindAllByActorId(actorId string, userId string, page *paging.Paging) ([]Movie, error)

	FindAllByDirectorId(actorId string, userId string, page *paging.Paging) ([]Movie, error)
}

type neo4jMovieService struct {
	driver neo4j.Driver
}

func NewMovieService(driver neo4j.Driver) MovieService {
	return &neo4jMovieService{driver: driver}
}

// FindAll should return a paginated list of movies ordered by the `sort`
// parameter and limited to the number passed as `limit`.  The `skip` variable should be
// used to skip a certain number of rows.
//
// If a userId value is supplied, a `favorite` boolean property should be returned to
// signify whether the user has aded the movie to their "My Favorites" list.
// tag::all[]
func (gs *neo4jMovieService) FindAll(userId string, page *paging.Paging) (movies []Movie, err error) {
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()
	// tag::allcypher[]
	// Execute a query in a new Read Transaction

	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get an array of IDs for the User's favorite movies
		favorites, err := getUserFavorites(tx, userId)
		if err != nil {
			return nil, err
		}
		// Retrieve a list of movies with the
		// favorite flag appended to the movie's properties
		sort := page.Sort()
		result, err := tx.Run(fmt.Sprintf(`
		MATCH (m:Movie)
		WHERE m.`+"`%[1]s`"+` IS NOT NULL
		RETURN m {
			.*,
			favorite: m.tmdbId IN $favorites
		} AS movie
		ORDER BY m.`+"`%[1]s`"+` %s
		SKIP $skip
		LIMIT $limit
		`, sort, page.Order()), map[string]interface{}{
			"favorites": favorites,
			"skip":      page.Skip(),
			"limit":     page.Limit(),
		})
		if err != nil {
			return nil, err
		}
		// tag::allmovies[]
		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			movie, _ := record.Get("movie")
			results = append(results, movie.(map[string]interface{}))
		}
		return results, nil
		// end::allmovies[]
	})
	// end::allcypher[]

	if err != nil {
		return nil, err
	}
	movies = results.([]Movie)
	// tag::return[]
	return movies, nil
	// end::return[]
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
func (gs *neo4jMovieService) FindAllByGenre(genre string, userId string, page *paging.Paging) (movies []Movie, err error) {
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		favorites, err := getUserFavorites(tx, userId)
		if err != nil {
			return nil, err
		}
		result, err := tx.Run(fmt.Sprintf(`
		MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})
		WHERE m.`+"`%[1]s`"+` IS NOT NULL
		RETURN m {
			.*,
			  favorite: m.tmdbId IN $favorites
		} AS movie
		ORDER BY m.`+"`%[1]s`"+` %s
		SKIP $skip
		LIMIT $limit`, page.Sort(), page.Order()), map[string]interface{}{
			"name":      genre,
			"favorites": favorites,
			"skip":      page.Skip(),
			"limit":     page.Limit(),
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			movie, _ := record.Get("movie")
			results = append(results, movie.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	movies = results.([]Movie)
	return movies, nil
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
func (gs *neo4jMovieService) FindAllByActorId(actorId string, userId string, page *paging.Paging) (movies []Movie, err error) {
	// Get Movies acted in by a Person
	// MATCH (:Person {tmdbId: $id})-[:ACTED_IN]->(m:Movie)

	// Open a new session
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Execute a query in a new Read Transaction
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get an array of IDs for the User's favorite movies
		favorites, err := getUserFavorites(tx, userId)
		if err != nil {
			return nil, err
		}

		// Retrieve a list of movies with the
		// favorite flag appended to the movie's properties
		result, err := tx.Run(fmt.Sprintf(`
			MATCH (:Person {tmdbId: $id})-[:ACTED_IN]->(m:Movie)
			WHERE m.`+"`%[1]s`"+` IS NOT NULL
			RETURN m {
			  .*,
			    favorite: m.tmdbId IN $favorites
			} AS movie
			ORDER BY m.`+"`%[1]s`"+` %s
			SKIP $skip
			LIMIT $limit`, page.Sort(), page.Order()), map[string]interface{}{
			"id":        actorId,
			"favorites": favorites,
			"skip":      page.Skip(),
			"limit":     page.Limit(),
		})
		if err != nil {
			return nil, err
		}

		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			movie, _ := record.Get("movie")
			results = append(results, movie.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	movies = results.([]Movie)
	return movies, nil
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
func (gs *neo4jMovieService) FindAllByDirectorId(actorId string, userId string, page *paging.Paging) (movies []Movie, err error) {
	// Get Movies directed by a Person
	// MATCH (:Person {tmdbId: $id})-[:DIRECTED]->(m:Movie)

	// Open a new session
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Execute a query in a new Read Transaction
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		// Get an array of IDs for the User's favorite movies
		favorites, err := getUserFavorites(tx, userId)
		if err != nil {
			return nil, err
		}

		// Retrieve a list of movies with the
		// favorite flag appended to the movie's properties
		result, err := tx.Run(fmt.Sprintf(`
			MATCH (:Person {tmdbId: $id})-[:DIRECTED]->(m:Movie)
			WHERE m.`+"`%[1]s`"+` IS NOT NULL
			RETURN m {
			  .*,
			    favorite: m.tmdbId IN $favorites
			} AS movie
			ORDER BY m.`+"`%[1]s`"+` %s
			SKIP $skip
			LIMIT $limit`, page.Sort(), page.Order()), map[string]interface{}{
			"id":        actorId,
			"favorites": favorites,
			"skip":      page.Skip(),
			"limit":     page.Limit(),
		})
		if err != nil {
			return nil, err
		}

		// Get a list of Movies from the Result
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
		var results []map[string]interface{}
		for _, record := range records {
			movie, _ := record.Get("movie")
			results = append(results, movie.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	movies = results.([]Movie)
	return movies, nil
}

// end::getForDirector[]

// getUserFavorites should return a list of tmdbId properties for the movies that
// the user has added to their 'My Favorites' list.
// tag::getUserFavorites[]
func getUserFavorites(tx neo4j.Transaction, userId string) ([]string, error) {
	if userId == "" {
		return nil, nil
	}
	results, err := tx.Run(`
	MATCH (u:User {userId: $userId})-[:HAS_FAVORITE]->(m)
	RETURN m.tmdbId AS id
`, map[string]interface{}{"userId": userId})
	if err != nil {
		return nil, err
	}
	var ids []string
	for results.Next() {
		record := results.Record()
		id, _ := record.Get("id")
		ids = append(ids, id.(string))
	}
	return ids, nil
}

// end::getUserFavorites[]
