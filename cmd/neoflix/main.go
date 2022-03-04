package main

import (
	"fmt"
	"net/http"

	config "github.com/neo4j-graphacademy/neoflix/pkg/config"

	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"github.com/neo4j-graphacademy/neoflix/pkg/routes"
	"github.com/neo4j-graphacademy/neoflix/pkg/services"
)

func main() {
	settings, err := config.ReadConfig("config.json")
	ioutils.PanicOnError(err)
	// tag::useDriver[]
	// tag::driver[]
	driver, err := config.NewDriver(settings)
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
		services.NewAuthService(driver, settings.JwtSecret, settings.SaltRounds),
		services.NewFavoriteService(driver))
	// end::useDriver[]

	server := newHttpServer()
	for _, route := range allRoutes {
		route.Register(server)
	}

	fmt.Printf("Server started on port %d\n", settings.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", settings.Port), server); err != nil {
		ioutils.PanicOnError(err)
	}
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
	authService services.AuthService,
	favoriteService services.FavoriteService) []routes.Routable {

	return []routes.Routable{
		routes.NewGenreRoutes(genreService, movieService, authService),
		routes.NewMovieRoutes(movieService, ratingService, authService),
		routes.NewPeopleRoutes(peopleService, movieService, authService),
		routes.NewAuthRoutes(authService),
		routes.NewAccountRoutes(ratingService, authService, favoriteService),
	}
}
