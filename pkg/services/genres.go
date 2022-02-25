package services

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type Genre = map[string]interface{}

type GenreService interface {
	FindAll() ([]Genre, error)
}

type neo4jGenreService struct {
	driver neo4j.Driver
}

func NewGenreService(driver neo4j.Driver) GenreService {
	return &neo4jGenreService{driver: driver}
}

func (gs *neo4jGenreService) FindAll() (genres []Genre, err error) {
	session := gs.driver.NewSession(neo4j.SessionConfig{})
	defer func() {
		err = ioutils.DeferredClose(session, err)
	}()

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
		records, err := result.Collect()
		if err != nil {
			return nil, err
		}
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
	genres = results.([]Genre)
	return genres, nil
}
