package neoflix;

import static neoflix.NeoflixApp.Params.Sort.*;
import static spark.Spark.*;

import com.google.gson.Gson;
import org.neo4j.driver.*;
import spark.Request;
import spark.RouteGroup;

import java.util.EnumSet;
import java.util.Map;
import java.util.Properties;

public class NeoflixApp {

    public static void main(String[] args) throws Exception {
        Properties props = new Properties();
        props.load(NeoflixApp.class.getResourceAsStream("/application.properties"));

        int port = Integer.parseInt(props.getProperty("APP_PORT", "3000"));
        port(port);
        staticFiles.location("/public");

        AuthToken auth = AuthTokens.basic(props.getProperty("NEO4J_USERNAME"), props.getProperty("NEO4J_PASSWORD"));
        Driver driver = GraphDatabase.driver(props.getProperty("NEO4J_URI"), auth);
        Gson gson = GsonUtils.gson();

        path("/api", () -> {
            path("/movies", new MovieRoutes(driver,gson));
            path("/genres", new GenreRoutes(driver,gson));
            // path("/auth", new AuthRoutes(driver,gson));
            path("/account", new AccountRoutes(driver,gson));
            path("/people", new PeopleRoutes(driver,gson));
        });
        System.out.printf("Started server at port %d%n",port);
    }

    private static class GenreRoutes implements RouteGroup {
        private final Gson gson;
        private final GenreService genreService;
        private final MovieService movieService;

        public GenreRoutes(Driver driver, Gson gson) {
            genreService = new GenreService(driver);
            movieService = new MovieService(driver);
            this.gson = gson;
        }

        @Override
        public void addRoutes() {
            /*
             * @GET /genres/
             *
             * This route should retrieve a full list of Genres from the
             * database along with a poster and movie count.
             */
            get("", (req, res) -> genreService.all(), gson::toJson);

            /*
             * @GET /genres/:name
             *
             * This route should return information on a genre with a name
             * that matches the :name URL parameter.  If the genre is not found,
             * a 404 should be thrown.
             */
            get("/:name", (req, res) -> genreService.find(req.params(":name")), gson::toJson);

            /**
             * @GET /genres/:name/movies
             *
             * This route should return a paginated list of movies that are listed in
             * the genre whose name matches the :name URL parameter.
             */
            get("/:name/movies", (req, res) -> {
                String userId = getUserId(req); // TODO
                return movieService.byGenre(req.params(":name"), Params.parse(req, Params.MOVIE_SORT), userId);
            }, gson::toJson);
        }

    }
    private static class PeopleRoutes implements RouteGroup {
        private final Gson gson;
        private final PeopleService peopleService;
        private final MovieService movieService;

        public PeopleRoutes(Driver driver, Gson gson) {
            this.gson = gson;
            peopleService = new PeopleService(driver);
            movieService = new MovieService(driver);
        }

        @Override
        public void addRoutes() {
            /*
             * @GET /people/
             *
             * This route should return a paginated list of People from the database
             */
            get("", (req, res) -> peopleService.all(Params.parse(req, Params.PEOPLE_SORT)), gson::toJson);

            /*
             * @GET /people/:id
             *
             * This route should the properties of a Person based on their tmdbId
             */
            get("/:id", (req, res) -> peopleService.findById(req.params(":id")), gson::toJson);

            /*
             * @GET /people/:id/similar
             *
             * This route should return a paginated list of similar people to the person
             * with the :id supplied in the route params.
             */
            get("/:id/similar", (req, res) -> peopleService.getSimilarPeople(req.params(":id"), Params.parse(req, Params.PEOPLE_SORT)), gson::toJson);

            /*
             * @GET /people/:id/acted
             *
             * This route should return a paginated list of movies that the person
             * with the :id has acted in.
             */
            get("/:id/acted", (req, res) -> {
                String userId = getUserId(req); // TODO
                return movieService.getForActor(req.params(":id"), Params.parse(req, Params.MOVIE_SORT), userId);
            }, gson::toJson);

            /*
             * @GET /people/:id/directed
             *
             * This route should return a paginated list of movies that the person
             * with the :id has acted in.
             */
            get("/:id/directed", (req, res) -> {
                String userId = getUserId(req); // TODO
                return movieService.getForDirector(req.params(":id"), Params.parse(req, Params.MOVIE_SORT), userId);
            }, gson::toJson);
        }

    }
    private static class AccountRoutes implements RouteGroup {
        private final Gson gson;
        private final FavoriteService favoriteService;
        private final RatingService ratingService;

        public AccountRoutes(Driver driver, Gson gson) {
            this.gson = gson;
            favoriteService = new FavoriteService(driver);
            ratingService = new RatingService(driver);
        }

        @Override
        public void addRoutes() {
            /*
             * @GET /account/
             *
             * This route simply returns the claims made in the JWT token
             */
            get("", (req, res) -> req.attribute("user"), gson::toJson);

            /*
             * @GET /account/favorites/
             *
             * This route should return a list of movies that a user has added to their
             * Favorites link by clicking the Bookmark icon on a Movie card.
             */
            // tag::list[]
            get("/favorites", (req, res) -> {
                String userId = getUserId(req);
                return favoriteService.all(userId, Params.parse(req, Params.MOVIE_SORT));
            }, gson::toJson);
            // end::list[]

            /*
             * @POST /account/favorites/:id
             *
             * This route should create a `:HAS_FAVORITE` relationship between the current user
             * and the movie with the :id parameter.
             */
            // tag::add[]
            post("/favorites/:id", (req, res) -> {
                String userId = getUserId(req);
                return favoriteService.add(req.params(":id"), userId);
            }, gson::toJson);
            // end::add[]

            /*
             * @DELETE /account/favorites/:id
             *
             * This route should remove the `:HAS_FAVORITE` relationship between the current user
             * and the movie with the :id parameter.
             */
            // tag::delete[]
            delete("/favorites/:id", (req, res) -> {
                String userId = getUserId(req); // TODO
                return favoriteService.remove(req.params(":id"), userId);
            }, gson::toJson);
            // end::delete[]

            /*
             * @POST /account/ratings/:id
             *
             * This route should create a `:RATING` relationship between the current user
             * and the movie with the :id parameter.  The rating value will be posted as part
             * of the post body.
             */
            // tag::rating[]
            get("/ratings/:id", (req, res) -> {
                String userId = getUserId(req); // TODO
                int rating = Integer.parseInt(req.body());
                return ratingService.add(userId, req.params(":id"), rating);
            }, gson::toJson);
            // end::rating[]
        }

    }

    private static class MovieRoutes implements RouteGroup {
        private final Gson gson;
        private final MovieService movieService;
        private final RatingService ratingService;

        public MovieRoutes(Driver driver, Gson gson) {
            this.gson = gson;
            movieService = new MovieService(driver);
            ratingService = new RatingService(driver);
        }

        @Override
        public void addRoutes() {
            /*
             * @GET /movies
             *
             * This route should return a paginated list of movies, sorted by the
             * `sort` query parameter,
             */
            // tag::list[]
            get("", (req, res) -> {
                String userId = getUserId(req); // TODO
                return movieService.all(Params.parse(req, Params.MOVIE_SORT), userId);
            }, gson::toJson);
            // end::list[]

            /*
             * @GET /movies/:id
             *
             * This route should find a movie by its tmdbId and return its properties.
             */
            // tag::get[]
            get("/:id", (req, res) -> {
                String userId = getUserId(req); // TODO
                Map<String, Object> movie = movieService.findById(req.params(":id"), userId);
                return movie;
            }, gson::toJson);

            /*
             * @GET /movies/:id/ratings
             *
             *
             * This route should return a paginated list of ratings for a movie, ordered by either
             * the rating itself or when the review was created.
             */
            // tag::ratings[]
            get("/:id/ratings", (req, res) -> ratingService.forMovie(req.params(":id"), Params.parse(req, Params.RATING_SORT)), gson::toJson);
            // end::ratings[]

            /*
             * @GET /movies/:id/similar
             *
             * This route should return a paginated list of similar movies, ordered by the
             * similarity score in descending order.
             */
            // tag::similar[]
            get("/:id/similar", (req, res) -> {
                String userId = getUserId(req); // TODO
                return movieService.getSimilarMovies(req.params(":id"), Params.parse(req, Params.MOVIE_SORT), userId);
            }, gson::toJson);
            // end::similar[]
        }

    }

    private static String getUserId(Request req) {
        Object user = req.attribute("user");
        if (!(user instanceof Map)) return null;
        return (String) ((Map<String, Object>) user).get("userId"); // todo
    }

    record Params(String query, Sort sort, Order order, int limit, int skip) {
        public Sort sort(Sort defaultSort) {
            return sort == null ? defaultSort : sort;
        }
        enum Order { ASC, DESC;
            static Order of(String value) {
                if (value == null || value.isBlank() || !"DESC".equalsIgnoreCase(value)) return ASC;
                return DESC;
            }
        };
        public enum Sort { /* Movie */ title, released, imdbRating, score,
            /* Person */ name, born, movieCount,
            /* */ rating, timestamp;
            static Sort of(String name) {
                if (name == null || name.isBlank()) return null;
                return valueOf(name);
            }
        }
        static final EnumSet<Sort> MOVIE_SORT = EnumSet.of(title, released, imdbRating, score);
        static final EnumSet<Sort> PEOPLE_SORT = EnumSet.of(name,born,movieCount);
        static final EnumSet<Sort> RATING_SORT = EnumSet.of(rating,timestamp);

        static Params parse(Request req, EnumSet<Sort> validSort) {
            String q = req.queryParamsSafe("q");
            Sort sort = Sort.of(req.queryParamsSafe("sort"));
            Order order = Order.of(req.queryParamsSafe("order"));
            int limit = Integer.parseInt(req.queryParamOrDefault("limit", "6"));
            int skip = Integer.parseInt(req.queryParamOrDefault("skip", "0"));
            // Only accept valid sort fields
            if (!validSort.contains(sort)) {
                sort = null;
            }
            return new Params(q, sort, order, limit, skip);
        }
    }
}