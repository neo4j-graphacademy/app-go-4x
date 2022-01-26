package neoflix.routes;

import com.google.gson.Gson;
import neoflix.NeoflixApp;
import neoflix.services.AuthService;
import org.neo4j.driver.Driver;
import spark.RouteGroup;

import static spark.Spark.*;

public class AuthRoutes implements RouteGroup {
    private final Gson gson;
    private final AuthService authService;

    public AuthRoutes(Driver driver, Gson gson, String jwtSecret) {
        this.gson = gson;
        authService = new AuthService(driver, jwtSecret);
    }

    static class UserData { String email, name, password; };

    @Override
    public void addRoutes() {
        /*
         * @POST /auth/login
         *
         * TODO This route invokes the `Neo4jStrategy` in `src/passport/neo4j.strategy.js`,
         * which, when implemented, attempts to authenticate the user against the
         * Neo4j database.
         *
         * TODO The req.user object assigned by the strategy should include a `token` property,
         * which holds the JWT token.  This token is then used in the `JwtStrategy` from
         * `src/passport/jwt.strategy.js` to authenticate the request.
         */
        // tag::login[]
        post("/login", (req, res) -> {
            UserData userData = gson.fromJson(req.body(), UserData.class);
            var user = authService.authenticate(userData.email, userData.password);
            req.attribute("user", user.get("userId"));
            return user;
        }, gson::toJson);
        // end::login[]

        /*
         * @POST /auth/register
         *
         * This route should use the AuthService to create a new User node
         * in the database with an encrypted password before returning a User record which
         * includes a `token` property.  This token is then used in the `JwtStrategy` from
         * `src/passport/jwt.strategy.js` to authenticate the request.
         */
        // tag::register[]
        post("/register", (req, res) -> {
            String userId = NeoflixApp.getUserId(req);
            UserData userData = gson.fromJson(req.body(), UserData.class);

            return authService.register(userData.email, userData.password, userData.name);
        }, gson::toJson);
        // end::register[]
    }
}
