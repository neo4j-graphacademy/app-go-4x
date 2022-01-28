package neoflix;

import neoflix.services.AuthService;
import neoflix.services.MovieService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import static neoflix.Params.Order.ASC;
import static neoflix.Params.Order.DESC;
import static neoflix.Params.Sort.imdbRating;
import static neoflix.Params.Sort.title;
import static org.junit.jupiter.api.Assertions.*;

class _03_RegisterUserTest {

    private static Driver driver;
    private static String jwtSecret;

    private static String email = "graphacademy@neo4j.com";
    private static String password = "letmein";
    private static String name = "Graph Academy";

    @BeforeAll
    static void initDriverAuth() {
        AppUtils.loadProperties();
        driver = AppUtils.initDriver();
        jwtSecret = AppUtils.getJwtSecret();

        driver.session().writeTransaction(tx -> tx.run("MATCH (u:User {email: $email}) DETACH DELETE u", Values.parameters("email", email)));
    }

    @AfterAll
    static void closeDriver() {
        driver.close();
    }

    @Test
    void registerUser() {
        AuthService authService = new AuthService(driver, jwtSecret);
        var limit = 1;
        var output = authService.register(email, password, name);
        assertNotNull(output);
        assertEquals(1, output.size());

        var userEmail = output.get("email");
        assertNotNull(email);
        assertEquals(email, userEmail);

        var userPassword = output.get("password");
        assertNotNull(password);

        var userName = output.get("name");
        assertNotNull(name);
        assertEquals(name, name);

        var token = output.get("token");
        assertNotNull(name);
    }
}