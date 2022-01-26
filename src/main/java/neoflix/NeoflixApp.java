package neoflix;

import static spark.Spark.*;

import com.google.gson.Gson;
import neoflix.routes.AccountRoutes;
import neoflix.routes.GenreRoutes;
import neoflix.routes.MovieRoutes;
import neoflix.routes.PeopleRoutes;
import org.neo4j.driver.*;
import spark.Request;

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

    public static String getUserId(Request req) {
        Object user = req.attribute("user");
        if (!(user instanceof Map)) return null;
        return (String) ((Map<String, Object>) user).get("userId"); // todo
    }

}