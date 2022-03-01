package main

import (
	"encoding/json"
	"fmt"
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Uri      string `json:"NEO4J_URI"`
	Username string `json:"NEO4J_USERNAME"`
	Password string `json:"NEO4J_PASSWORD"`

	Port       int    `json:"APP_PORT"`
	JwtSecret  string `json:"JWT_SECRET"`
	SaltRounds int    `json:"SALT_ROUNDS"`
}

func main() {
	config, err := readConfig()
	ioutils.PanicOnError(err)
	// tag::driver[]
	driver, err := neo4j.NewDriver(
		config.Uri,
		neo4j.BasicAuth(config.Username, config.Password, ""),
	)
	// end::driver[]
	ioutils.PanicOnError(err)
	defer func() {
		ioutils.PanicOnError(driver.Close())
	}()

	allRoutes := allRoutes(
		services.NewMovieService(driver),
		services.NewGenreService(driver),
		services.NewRatingService(driver),
		services.NewPeopleService(driver),
		services.NewAuthService(driver, config.JwtSecret, config.SaltRounds))

	server := newHttpServer()
	for _, route := range allRoutes {
		route.Register(server)
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Port), server); err != nil {
		ioutils.PanicOnError(err)
	}
}

func readConfig() (*Config, error) {
	file, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func newHttpServer() *http.ServeMux {
	server := http.NewServeMux()
	server.Handle("/", http.FileServer(http.Dir("public")))
	return server
}

func allRoutes(
	movieService services.MovieService,
	genreService services.GenreService,
	ratingService services.RatingService,
	peopleService services.PeopleService,
	authService services.AuthService) []routes.Routable {

	return []routes.Routable{
		routes.NewGenreRoutes(genreService, movieService, authService),
		routes.NewMovieRoutes(movieService, ratingService, authService),
		routes.NewPeopleRoutes(peopleService, movieService, authService),
		routes.NewAuthRoutes(authService),
		routes.NewAccountRoutes(ratingService, authService),
	}
}
