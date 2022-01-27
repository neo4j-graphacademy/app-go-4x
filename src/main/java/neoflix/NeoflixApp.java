package neoflix;

import static spark.Spark.*;

import com.google.gson.Gson;
import neoflix.routes.*;
import org.neo4j.driver.*;
import spark.Request;

import java.io.IOException;
import java.util.Properties;

public class NeoflixApp {

    static final Properties PROPS = new Properties() {{
        try {
            load(NeoflixApp.class.getResourceAsStream("/application.properties"));
        } catch (IOException e) {
            throw new RuntimeException("Error loading application.properties",e);
        }
    }};

    public static void main(String[] args) throws Exception {
        int port = Integer.parseInt(PROPS.getProperty("APP_PORT", "3000"));
        port(port);
        Driver driver = initDriver();
        Gson gson = GsonUtils.gson();

        staticFiles.location("/public");
        String jwtSecret = PROPS.getProperty("JWT_SECRET");
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

    static Driver initDriver() {
        AuthToken auth = AuthTokens.basic(PROPS.getProperty("NEO4J_USERNAME"), PROPS.getProperty("NEO4J_PASSWORD"));
        Driver driver = GraphDatabase.driver(PROPS.getProperty("NEO4J_URI"), auth);
        driver.verifyConnectivity();
        return driver;
    }

    public static String getUserId(Request req) {
        Object user = req.attribute("user");
        if (user == null) return null;
        return user.toString();
    }

}