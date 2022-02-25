package services

import (
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes/paging"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Movie = map[string]interface{}

type MovieService interface {
	FindAllByGenre(genre, userId string, page *paging.Paging) ([]Movie, error)
}

type neo4jMovieService struct {
	driver neo4j.Driver
}

func NewMovieService(driver neo4j.Driver) MovieService {
	return &neo4jMovieService{driver: driver}
}

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
		WHERE m.%[1]s IS NOT NULL
		RETURN m {
			.*,
			  favorite: m.tmdbId IN $favorites
		} AS movie
		ORDER BY m.%[1]s %s
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
