package neoflix.services;

import neoflix.Params;
import org.neo4j.driver.Driver;
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
    public List<Map<String,Object>> forMovie(String id, Params params) {
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
                        LIMIT $limit""", params.sort(Params.Sort.timestamp), params.order());
                var res = tx.run(query, Values.parameters("id", id, "limit", params.limit(), "skip", params.skip()));
                return res.list(row -> row.get("review").asMap());
            });
        }
    }
    // end::forMovie[]


    /**
     * Add a relationship between a User and Movie with a `rating` property.
     * The `rating` parameter should be converted to a Neo4j Integer.
     *
     * If the User or Movie cannot be found, a NotFoundError should be thrown
     *
     * @param {string} userId   the userId for the user
     * @param {string} movieId  The tmdbId for the Movie
     * @param {number} rating   An integer representing the rating from 1-5
     * @returns {Promise<Record<string, any>>}  A movie object with a rating property appended
     */
    // tag::add[]
    public Map<String,Object> add(String userId, String movieId, int rating) {
        // tag::write[]
        // Save the rating in the database

        // Open a new session
        try (var session = this.driver.session()) {

            // Run the cypher query
            var movies = session.writeTransaction(tx -> {
                String query = """
                        MATCH (u:User {userId: $userId})
                        MATCH (m:Movie {tmdbId: $movieId})

                        MERGE (u)-[r:RATED]->(m)
                        SET r.rating = $rating, r.timestamp = timestamp()
                                
                        RETURN m { .*, rating: r.rating } AS movie
                        """;
                var res = tx.run(query, Values.parameters("userId", userId, "movieId", movieId, "rating", rating));
                return res.list(row -> row.get("movie").asMap()).stream();
            });
            // end::write[]

            // tag::throw[]
            var movie = movies.findFirst().orElseThrow(() -> new RuntimeException(String.format("Could not create rating for Movie %s by User %s", movieId, userId)));
            // end::throw[]

            // tag::addreturn[]
            // Return movie details and a rating
            return movie;
            // end::addreturn[]
        }
    }
    // end::add[]
}