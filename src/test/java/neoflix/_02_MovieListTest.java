// mvn test -Dtest=neoflix._02_MovieListTest#applyOrderListAndSkip
// mvn test -Dtest=neoflix._02_MovieListTest#orderMoviesByRating
package neoflix;

import neoflix.services.MovieService;
import org.junit.jupiter.api.AfterAll;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;

import static neoflix.Params.Order.ASC;
import static neoflix.Params.Order.DESC;
import static neoflix.Params.Sort.imdbRating;
import static neoflix.Params.Sort.title;
import static org.junit.jupiter.api.Assertions.*;

class _02_MovieListTest {

    private static Driver driver;

    @BeforeAll
    static void initDriver() {
        AppUtils.loadProperties();
        driver = AppUtils.initDriver();
    }

    @AfterAll
    static void closeDriver() {
        driver.close();
    }

    @Test
    void applyOrderListAndSkip() {
        MovieService movieService = new MovieService(driver);
        var limit = 1;
        var output = movieService.all(new Params(null, title, ASC, limit, 0), null);
        assertNotNull(output);
        assertEquals(limit, output.size());
        assertNotNull(output.get(0));
        var firstTitle = output.get(0).get("title");
        assertNotNull(firstTitle);
        assertEquals("\"Great Performances\" Cats", firstTitle);

        var skip = 1;
        var next = movieService.all(new Params(null, Params.Sort.title, ASC, limit, skip), null);
        assertNotNull(next);
        assertEquals(limit, next.size());
        assertNotEquals(next.get(0).get("title"), firstTitle);
    }

    @Test
    void orderMoviesByRating() {
        var movieService = new MovieService(driver);
        var limit = 1;
        var output = movieService.all(new Params(null, imdbRating, DESC, limit, 0), null);
        assertNotNull(output);
        assertEquals(limit, output.size());
        assertNotNull(output.get(0));
        var firstTitle = output.get(0).get("title");
        assertNotNull(firstTitle);

        System.out.println("""
                
                Here is the answer to the quiz question on the lesson:
                What is the title of the highest rated movie in the recommendations dataset?
                Copy and paste the following answer into the text box:
                """);
        System.out.println(firstTitle);
    }
}