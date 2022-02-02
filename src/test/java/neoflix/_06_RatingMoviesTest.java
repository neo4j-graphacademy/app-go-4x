// Task: Rewrite the AuthService to allow users to authenticate against the database
// Outcome: A user will be able to authenticate against their database record
package neoflix;

import neoflix.services.RatingService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import java.util.UUID;

import static org.junit.jupiter.api.Assertions.*;

class _06_RatingMoviesTest {
    private static Driver driver;

    private static String email = "graphacademy.reviewer@neo4j.com";
    private static String movieId = "769";
    private static String userId = "1185150b-9e81-46a2-a1d3-eb649544b9c4";
    private static int rating = 5;

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
    void writeMovieRatingAsInt() {
        RatingService ratingService = new RatingService(driver);

        var output = ratingService.add(userId, movieId, rating);

        assertNotNull(output);
        assertEquals(output.get("tmdbId"), movieId);
        assertEquals(Integer.parseInt(output.get("rating").toString()), rating);
    }
}
