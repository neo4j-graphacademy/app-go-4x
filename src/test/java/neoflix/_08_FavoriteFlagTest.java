// Task: Rewrite the AuthService to allow users to authenticate against the database
// Outcome: A user will be able to authenticate against their database record
package neoflix;

import neoflix.services.FavoriteService;
import neoflix.services.MovieService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import static neoflix.Params.Order.DESC;
import static neoflix.Params.Sort.imdbRating;
import static org.junit.jupiter.api.Assertions.*;

class _08_FavoriteFlagTest {
    private static Driver driver;

    private static final String userId = "fe770c6b-4034-4e07-8e40-2f39e7a6722c";
    private static final String email = "graphacademy.flag@neo4j.com";

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

    @BeforeEach
    void setUp() {
        try (var session = driver.session()) {
            session.writeTransaction(tx ->
                    tx.run("MATCH (u:User {userId: $userId})-[r:HAS_FAVORITE]->(m:Movie) DELETE r",
                            Values.parameters("userId", userId)));
        }
    }

    @Test
    void favoriteMovieReturnsFlaggedInMovieList() {
        MovieService movieService = new MovieService(driver);
        FavoriteService favoriteService = new FavoriteService(driver);

        // Get the most popular movie
        var topMovie = movieService.all(new Params(null, imdbRating, DESC, 1, 0), userId);

        // Add top movie to user favorites
        var topMovieId = topMovie.get(0).get("tmdbId").toString();
        var add = favoriteService.add(userId, topMovieId);
        assertEquals(topMovieId, add.get("tmdbId"));
        assertTrue((Boolean)add.get("favorite"), "top movie is favorite");

        var addCheck = favoriteService.all(userId, new Params(null, imdbRating, Params.Order.ASC, 999, 0));

        assertEquals(1, addCheck.size());
        var found = addCheck.stream().allMatch(movie -> movie.get("tmdbId").equals(topMovieId));
        assertTrue(found);

        var topTwo = movieService.all(new Params(null, imdbRating, DESC, 2, 0), userId);
        assertEquals(topMovieId, topTwo.get(0).get("tmdbId"));
        assertEquals(true, topTwo.get(0).get("favorite"));
        assertEquals(false, topTwo.get(1).get("favorite"));
    }
}
