package neoflix;

import static spark.Spark.*;

import com.google.gson.Gson;
import neoflix.routes.*;
import neoflix.services.AuthService;
import org.neo4j.driver.*;
import spark.Request;

import java.util.Map;
import java.util.Properties;

public class NeoflixApp {

    public static void main(String[] args) throws Exception {
        Properties props = new Properties();
        props.load(NeoflixApp.class.getResourceAsStream("/application.properties"));

        String jwtSecret = props.getProperty("JWT_SECRET");
        int port = Integer.parseInt(props.getProperty("APP_PORT", "3000"));
        port(port);
        staticFiles.location("/public");

        AuthToken auth = AuthTokens.basic(props.getProperty("NEO4J_USERNAME"), props.getProperty("NEO4J_PASSWORD"));
        Driver driver = GraphDatabase.driver(props.getProperty("NEO4J_URI"), auth);
        Gson gson = GsonUtils.gson();

        before((req, res) -> {
            String token = req.headers("Authorization");
            String bearer = "Bearer ";
            if (token != null && !token.isBlank() && token.startsWith(bearer)) {
                // verify token
                token = token.substring(bearer.length());
                String userId = AuthUtils.verify(token, jwtSecret);
                req.attribute("user", userId);
            }
        });
        path("/api", () -> {
            path("/movies", new MovieRoutes(driver,gson));
            path("/genres", new GenreRoutes(driver,gson));
            path("/auth", new AuthRoutes(driver,gson, jwtSecret));
            path("/account", new AccountRoutes(driver,gson));
            path("/people", new PeopleRoutes(driver,gson));
        });
        System.out.printf("Started server at port %d%n",port);
    }

    public static String getUserId(Request req) {
        Object user = req.attribute("user");
        if (user == null) return null;
        return user.toString();
    }

}