package services

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Rating = map[string]interface{}

type RatingService interface {
	FindAllByMovieId(id string, page *paging.Paging) ([]Rating, error)

	Save(rating int, movieId string, userId string) (Movie, error)
}

type neo4jRatingService struct {
	driver neo4j.Driver
}

func NewRatingService(driver neo4j.Driver) RatingService {
	return &neo4jRatingService{driver: driver}
}

// FindAllByMovieId returns a paginated list of reviews for a Movie.
//
// Results should be ordered by the `sort` parameter, and in the direction specified
// in the `order` parameter.
// Results should be limited to the number passed as `limit`.
// The `skip` variable should be used to skip a certain number of rows.
// tag::forMovie[]
func (rs *neo4jRatingService) FindAllByMovieId(movieId string, page *paging.Paging) (_ []Rating, err error) {
	// Open a new database session
	session := rs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Get ratings for a Movie
	results, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(fmt.Sprintf(`
			MATCH (u:User)-[r:RATED]->(m:Movie {tmdbId: $id})
			RETURN r {
				.rating,
				.timestamp,
			     user: u { .id, .name }
			} AS review
			ORDER BY r.`+"`%s`"+` %s
			SKIP $skip
			LIMIT $limit`, page.Sort(), page.Order()),
			map[string]interface{}{
				"id":    movieId,
				"skip":  page.Skip(),
				"limit": page.Limit(),
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
			review, _ := record.Get("review")
			results = append(results, review.(map[string]interface{}))
		}
		return results, nil
	})

	if err != nil {
		return nil, err
	}
	return results.([]Rating), nil
}

// end::forMovie[]

// Save adds a relationship between a User and Movie with a `rating` property.
// The `rating` parameter should be converted to a Neo4j Integer.
//
// If the User or Movie cannot be found, a NotFoundError should be thrown
// tag::add[]
func (rs *neo4jRatingService) Save(rating int, movieId string, userId string) (_ Movie, err error) {
	// tag::write[]
	// Save the rating in the database

	// Open a new session
	session := rs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

	// Run the cypher query
	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		result, err := tx.Run(`
			MATCH (u:User {userId: $userId})
			MATCH (m:Movie {tmdbId: $movieId})

			MERGE (u)-[r:RATED]->(m)
			SET r.rating = $rating, r.timestamp = timestamp()

			RETURN m { .*, rating: r.rating } AS movie
`, map[string]interface{}{
			"userId":  userId,
			"movieId": movieId,
			"rating":  rating,
		})
		if err != nil {
			return nil, err
		}

		record, err := result.Single()
		if err != nil {
			return nil, err
		}

		movie, _ := record.Get("movie")
		return movie.(map[string]interface{}), nil
	})
	// end::write[]

	// tag::throw[]
	if err != nil {
		return nil, err
	}
	// end::throw[]

	// tag::addreturn[]
	// Return movie details and a rating
	return result.(Movie), nil
	// end::addreturn[]
}

// end::add[]
