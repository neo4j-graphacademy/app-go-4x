package neoflix;

import static neoflix.NeoflixApp.Query.Sort.*;
import static spark.Spark.*;

import com.google.gson.Gson;
import org.neo4j.driver.*;
import spark.Request;
import spark.RouteGroup;

import java.util.EnumSet;
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
            path("/genres", new GenreRoutes(driver,gson));
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
                var userId = req.headers("user"); // TODO
                return movieService.byGenre(req.params(":name"), Query.parse(req, Query.MOVIE_SORT), userId);
            }, gson::toJson);
        }

    }
    private static class PeopleRoutes implements RouteGroup {
        private final Gson gson;
        private final PeopleService peopleService;

        public PeopleRoutes(Driver driver, Gson gson) {
            this.gson = gson;
            peopleService = new PeopleService(driver);
        }

        @Override
        public void addRoutes() {
            /*
             * @GET /people/
             *
             * This route should return a paginated list of People from the database
             */
            get("", (req, res) -> peopleService.all(Query.parse(req,Query.PEOPLE_SORT)), gson::toJson);

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
            get("/:id/similar", (req, res) -> peopleService.getSimilarPeople(req.params(":id"), Query.parse(req, Query.PEOPLE_SORT)), gson::toJson);
        }

    }

    record Query(String query, Sort sort, Order order, int limit, int skip) {
        public Sort sort(Sort defaultSort) {
            return sort == null ? defaultSort : sort;
        }
        enum Order { ASC, DESC;
            static Order of(String value) {
                if (value == null || value.isBlank() || !"DESC".equalsIgnoreCase(value)) return ASC;
                return DESC;
            }
        };
        static final EnumSet<Order> ORDERS = EnumSet.allOf(Order.class);
        public enum Sort { /* Movie */ title, released, imdbRating,
            /* Person */ name, born, movieCount,
            /* */ rating, timestamp;
            static Sort of(String name) {
                if (name == null || name.isBlank()) return null;
                return valueOf(name);
            }
        }
        static final EnumSet<Sort> MOVIE_SORT = EnumSet.of(title, released, imdbRating);
        static final EnumSet<Sort> PEOPLE_SORT = EnumSet.of(name,born,movieCount);
        static final EnumSet<Sort> RATING_SORT = EnumSet.of(rating,timestamp);

        static Query parse(Request req, EnumSet<Sort> validSort) {
            String q = req.queryParamsSafe("q");
            Sort sort = Sort.of(req.queryParamsSafe("sort"));
            Order order = Order.of(req.queryParamsSafe("order"));
            int limit = Integer.parseInt(req.queryParamOrDefault("limit", "6"));
            int skip = Integer.parseInt(req.queryParamOrDefault("skip", "0"));
            // Only accept valid sort fields
            if (!validSort.contains(sort)) {
                sort = null;
            }
            return new Query(q, sort, order, limit, skip);
        }
    }
}