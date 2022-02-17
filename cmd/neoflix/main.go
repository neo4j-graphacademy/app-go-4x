package main

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Uri      string `json:"NEO4J_URI"`
	Username string `json:"NEO4J_USERNAME"`
	Password string `json:"NEO4J_PASSWORD"`

	Port      int    `json:"APP_PORT"`
	JwtSecret string `json:"JWT_SECRET"`
	SaltRound int    `json:"SALT_ROUNDS"`
}

func main() {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = json.Unmarshal(file, &config)
	if err != nil {
		panic(err)
	}

	// tag::driver[]
	driver, err := neo4j.NewDriver(config.Uri, neo4j.BasicAuth(config.Username,
		config.Password, ""))
	// end::driver[]
	if err != nil {
		panic(err)
	}

	server := http.NewServeMux()
	server.Handle("/", http.FileServer(http.Dir("public")))
	server.HandleFunc("/api/genres", func(writer http.ResponseWriter,
		request *http.Request) {
		session := driver.NewSession(neo4j.SessionConfig{})
		defer session.Close()

		var query = `
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
		ORDER BY g.name ASC`
		genres, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
			result, err := tx.Run(query, nil)
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
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		genreJson, err := json.Marshal(genres)
		if err != nil {
			writer.WriteHeader(500)
			writer.Write([]byte(err.Error()))
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(200)
		writer.Write(genreJson)
	})
	//server.HandleFunc()
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port),
		server); err != nil {
		panic(err)
	}
}
