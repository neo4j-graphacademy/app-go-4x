// Task: Rewrite the AuthService to allow users to authenticate against the database
// Outcome: A user will be able to authenticate against their database record
package neoflix;

import neoflix.services.MovieService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import static neoflix.Params.Sort.released;
import static neoflix.Params.Sort.title;
import static org.junit.jupiter.api.Assertions.*;

class _11_MovieListTest {
    private static Driver driver;

    private static String userId = "fe770c6b-4034-4e07-8e40-2f39e7a6722c";
    private static String email = "graphacademy.movielists@neo4j.com";
    private static String tomHanks = "31";
    private static String coppola = "1776";

    @BeforeAll
    static void initDriver() {
        AppUtils.loadProperties();
        driver = AppUtils.initDriver();

        driver.session().writeTransaction(tx -> tx.run("""
                MERGE (u:User {userId: $userId}) SET u.email = $email
                """, Values.parameters("userId", userId, "email", email)));
    }

    @AfterAll
    static void closeDriver() {
        driver.close();
    }

    @Test
    void getPaginatedMoviesByGenre() {
        MovieService movieService = new MovieService(driver);

        var genreName = "Comedy";
        var limit = 10;

        var output = movieService.byGenre(genreName, new Params(null, title, Params.Order.ASC, limit, 0), userId);
        assertNotNull(output);
        assertEquals(limit, output.size());

        var secondOutput = movieService.byGenre(genreName, new Params(null, title, Params.Order.ASC, limit, limit), userId);
        assertNotNull(secondOutput);
        assertEquals(limit, secondOutput.size());

        assertNotEquals(output.get(0).get("title"), secondOutput.get(0).get("title"));

        var reordered = movieService.byGenre(genreName, new Params(null, released, Params.Order.ASC, limit, limit), userId);
        assertNotEquals(output.get(0).get("title"), reordered.get(0).get("title"));
    }
}
