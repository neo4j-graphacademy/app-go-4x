package neoflix;

import static spark.Spark.*;

import com.google.gson.Gson;
import org.neo4j.driver.*;
import spark.Request;
import spark.Response;
import spark.Route;

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
        Gson gson = new Gson();

        get("/hello", new HelloRoute(driver));
        get("/api/genres", new GenresRoute(driver),gson::toJson);
        System.out.printf("Started server at port %d%n",port);

    }

    private static class HelloRoute implements Route {
        private final Driver driver;

        public HelloRoute(Driver driver) {
            this.driver = driver;
        }

        @Override
        public Object handle(Request req, Response res) throws Exception {
            try (var session = driver.session()) {
                var count = session.readTransaction(tx -> tx.run("MATCH (n) RETURN count(*) as c").single().get("c").asLong());
                return "The graph has " + count + " nodes.";
            }

        }
    }
    private static class GenresRoute implements Route {
        private final GenreService service;

        public GenresRoute(Driver driver) {
            service = new GenreService(driver);
        }

        @Override
        public Object handle(Request req, Response res) throws Exception {
            return service.all();
        }
    }
}