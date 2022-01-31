// Task: Rewrite the AuthService to allow users to authenticate against the database
// Outcome: A user will be able to authenticate against their database record
package neoflix;

import neoflix.services.FavoriteService;
import neoflix.services.RatingService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import static org.junit.jupiter.api.Assertions.*;

class _07_FavoritesListTest {
    private static Driver driver;

    private static String toyStory = "862";
    private static String goodfellas = "769";
    private static String userId = "9f965bf6-7e32-4afb-893f-756f502b2c2a";
    private static String email = "graphacademy.favorite@neo4j.com";

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
    void notFoundIfMovieOrUserNotExist() {
        FavoriteService favoriteService = new FavoriteService(driver);

        try {
            favoriteService.add("unknown", "x999");
            fail("Adding favorite with unknown userId or movieId should fail");
        } catch (Exception e) {
            assertEquals(e.getMessage(), "Couldn't create a favorite relationship for User unknown and Movie x999");
        }
    }

    @Test
    void saveMovieToUserFavorites() {
        FavoriteService favoriteService = new FavoriteService(driver);

        var output = favoriteService.add(userId, toyStory);

        assertNotNull(output);
        assertEquals(output.get("tmdbId"), toyStory);
        assertEquals(output.get("favorite"), true);

        var favorites = favoriteService.all(userId, new Params(null, Params.Sort.title, Params.Order.DESC, 10, 0));

        var movieFavorite = favorites.stream().filter(movie -> movie.get("tmdbId").equals(toyStory));
        assertNotNull(movieFavorite);
        assertEquals(movieFavorite.findAny().get().get("tmdbId"), toyStory);
    }

    @Test
    void addAndRemoveMovieFromFavorites() {
        FavoriteService favoriteService = new FavoriteService(driver);

        var add = favoriteService.add(userId, goodfellas);
        assertEquals(add.get("tmdbId"), goodfellas);
        assertEquals(add.get("favorite"), true);

        var addCheck = favoriteService.all(userId, new Params(null, Params.Sort.title, Params.Order.DESC, 10, 0));
        var found = addCheck.stream().filter(movie -> movie.get("tmdbId").equals(goodfellas));
        assertNotNull(found);
        assertEquals(found.findAny().get().get("tmdbId"), goodfellas);

        var remove = favoriteService.remove(userId, goodfellas);
        assertEquals(remove.get("tmdbId"), goodfellas);
        assertEquals(remove.get("favorite"), false);

        var removeCheck = favoriteService.all(userId, new Params(null, Params.Sort.title, Params.Order.DESC, 10, 0));
        var notFound = removeCheck.stream().filter(movie -> movie.get("tmdbId").equals(goodfellas));
        assertTrue(notFound.findAny().isEmpty());
    }
}
