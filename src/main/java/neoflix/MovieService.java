package neoflix;

import org.neo4j.driver.Driver;
import org.neo4j.driver.Transaction;
import org.neo4j.driver.Values;

import java.util.List;
import java.util.Map;

public class MovieService {
    private final Driver driver;

    public MovieService(Driver driver) {
        this.driver = driver;
    }

    /**
     * This method should return a paginated list of movies that have a relationship to the
     * supplied Genre.
     *
     * Results should be ordered by the `sort` parameter, and in the direction specified
     * in the `order` parameter.
     * Results should be limited to the number passed as `limit`.
     * The `skip` variable should be used to skip a certain number of rows.
     *
     * If a userId value is supplied, a `favorite` boolean property should be returned to
     * signify whether the user has added the movie to their "My Favorites" list.
     *
     * @param name
     * @param query
     * @param userId
     * @return List<Movie> movies for that genre
     */
    // tag::getByGenre[]
    public List<Map<String,Object>> byGenre(String name, NeoflixApp.Query query, String userId) {
        // Get Movies in a Genre
        // MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})

        // Open a new session and close at the end
        try (var session = driver.session()) {
            // Execute a query in a new Read Transaction
            return session.readTransaction((tx) -> {
                // Get an array of IDs for the User's favorite movies
                var favorites = getUserFavorites(tx, userId);

                // Retrieve a list of movies with the
                // favorite flag appened to the movie's properties
                var result = tx.run(
                                  String.format(
                                  """
                                  MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})
                                  WHERE m.`%s` IS NOT NULL
                                  RETURN m {
                                    .*,
                                      favorite: m.tmdbId IN $favorites
                                  } AS movie
                                  ORDER BY m.`%s` %s
                                  SKIP $skip
                                  LIMIT $limit
                                """,query.sort(),query.sort(),query.order()),
                        Values.parameters("skip", query.skip(), "limit", query.limit(),
                                "favorites", favorites, "name", name));
                var movies = result.list(row -> row.get("movie").asMap());
                return movies;
            });
        }
    }
    // end::getByGenre[]

    /**
     * This function should return a list of tmdbId properties for the movies that
     * the user has added to their 'My Favorites' list.
     *
     * @param tx The open transaction
     * @param userId The ID of the current user
     * @return List<String> movieIds of favorite movies
     */
    // tag::getUserFavorites[]
    private List<String> getUserFavorites(Transaction tx, String userId) {
        // If userId is not defined, return an empty list
        if (userId == null) return List.of();
        var favoriteResult =  tx.run("""
                    MATCH (u:User {userId: $userId})-[:HAS_FAVORITE]->(m)
                    RETURN m.tmdbId AS id
                """, Values.parameters("userId",userId));
            // Extract the `id` value returned by the cypher query
            return favoriteResult.list(row -> row.get("id").asString());
    }
    // end::getUserFavorites[]

    record Movie() {} // todo
}
