package neoflix;

import org.neo4j.driver.Driver;
import org.neo4j.driver.Transaction;
import org.neo4j.driver.Values;

import java.util.List;
import java.util.Map;

public class RatingService {
    private final Driver driver;

    /**
     * The constructor expects an instance of the Neo4j Driver, which will be
     * used to interact with Neo4j.
     */
    public RatingService(Driver driver) {
        this.driver = driver;
    }

    /**
     * Return a paginated list of reviews for a Movie.
     *
     * Results should be ordered by the `sort` parameter, and in the direction specified
     * in the `order` parameter.
     * Results should be limited to the number passed as `limit`.
     * The `skip` variable should be used to skip a certain number of rows.
     *
     * @param {string} id       The tmdbId for the movie
     * @param {string} sort  The field to order the results by
     * @param {string} order    The direction of the order (ASC/DESC)
     * @param {number} limit    The total number of records to return
     * @param {number} skip     The number of records to skip
     * @returns {Promise<Record<string, any>>}
     */
    // tag::forMovie[]
    public List<Map<String,Object>> forMovie(String id, NeoflixApp.Params params) {
        // Open a new database session
        try (var session = this.driver.session()) {

            // Get ratings for a Movie
            return session.readTransaction(tx -> {
                String query = String.format("""
                        MATCH (u:User)-[r:RATED]->(m:Movie {tmdbId: $id})
                        RETURN r {
                            .rating,
                            .timestamp,
                             user: u { .id, .name }
                        } AS review
                        ORDER BY r.`%s` %s
                        SKIP $skip
                        LIMIT $limit""", params.sort(NeoflixApp.Params.Sort.timestamp), params.order());
                var res = tx.run(query, Values.parameters("id", id, "limit", params.limit(), "skip", params.skip()));
                return res.list(row -> row.get("review").asMap());
            });
        }
    }
    // end::forMovie[]

}