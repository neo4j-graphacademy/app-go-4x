package neoflix.services;

import neoflix.NeoflixApp;
import neoflix.Params;
import org.neo4j.driver.Driver;
import org.neo4j.driver.Transaction;
import org.neo4j.driver.Values;

import java.util.List;
import java.util.Map;

public class MovieService {
    private final Driver driver;

    /**
     * The constructor expects an instance of the Neo4j Driver, which will be
     * used to interact with Neo4j.
     */
    public MovieService(Driver driver) {
        this.driver = driver;
    }

    /**
     * This method should return a paginated list of movies ordered by the `sort`
     * parameter and limited to the number passed as `limit`.  The `skip` variable should be
     * used to skip a certain number of rows.
     *
     * If a userId value is supplied, a `favorite` boolean property should be returned to
     * signify whether the user has aded the movie to their "My Favorites" list.
     *
     * @param params query params (query, sort, order, limit, skip)
     * @param userId
     * @returns {Promise<Record<string, any>[]>}
     */
    // tag::all[]
    public List<Map<String,Object>> all(Params params, String userId) {
        // Open a new session
        try (var session = this.driver.session()) {
            // tag::allcypher[]
            // Execute a query in a new Read Transaction
            var movies = session.readTransaction(tx -> {
                // Get an array of IDs for the User's favorite movies
                var favorites =  getUserFavorites(tx, userId);

                // Retrieve a list of movies with the
                // favorite flag appened to the movie's properties
                Params.Sort sort = params.sort(Params.Sort.title);
                String query = String.format("""
                    MATCH (m:Movie)
                    WHERE m.`%s` IS NOT NULL
                    RETURN m {
                      .*,
                      favorite: m.tmdbId IN $favorites
                    } AS movie
                    ORDER BY m.`%s` %s
                    SKIP $skip
                    LIMIT $limit
                    """, sort, sort, params.order());
                var res= tx.run(query, Values.parameters( "skip", params.skip(), "limit", params.limit(), "favorites",favorites));
                // tag::allmovies[]
                // Get a list of Movies from the Result
                return res.list(row -> row.get("movie").asMap());
                // end::allmovies[]
            });
            // end::allcypher[]

            // tag::return[]
            return movies;
            // end::return[]
        }
    }
    // end::all[]


    /**
     * @public
     * This method find a Movie node with the ID passed as the `id` parameter.
     * Along with the returned payload, a list of actors, directors, and genres should
     * be included.
     * The number of incoming RATED relationships should also be returned as `ratingCount`
     *
     * If a userId value is suppled, a `favorite` boolean property should be returned to
     * signify whether the user has aded the movie to their "My Favorites" list.
     *
     * @param {string} id
     * @returns {Promise<Record<string, any>>}
     */
    // tag::findById[]
    public Map<String,Object> findById(String id, String userId) {
        // Find a movie by its ID
        // MATCH (m:Movie {tmdbId: $id})

        // Open a new database session
        try (var session = this.driver.session()) {
            // Find a movie by its ID
            return session.readTransaction(tx -> {
                var favorites = getUserFavorites(tx, userId);

                String query = """
                                MATCH (m:Movie {tmdbId: $id})
                                RETURN m {
                                  .*,
                                    actors: [ (a)-[r:ACTED_IN]->(m) | a { .*, role: r.role } ],
                                    directors: [ (d)-[:DIRECTED]->(m) | d { .* } ],
                                    genres: [ (m)-[:IN_GENRE]->(g) | g { .name }],
                                    ratingCount: size((m)<-[:RATED]-()),
                                    favorite: m.tmdbId IN $favorites
                                } AS movie
                                LIMIT 1
                        """;
                var res = tx.run(query, Values.parameters("id", id, "favorites", favorites));
                return res.single().get("movie").asMap();
            });

        }
    }
    // end::findById[]

    /**
     * This method should return a paginated list of similar movies to the Movie with the
     * id supplied.  This similarity is calculated by finding movies that have many first
     * degree connections in common: Actors, Directors and Genres.
     *
     * Results should be ordered by the `sort` parameter, and in the direction specified
     * in the `order` parameter.
     * Results should be limited to the number passed as `limit`.
     * The `skip` variable should be used to skip a certain number of rows.
     *
     * If a userId value is suppled, a `favorite` boolean property should be returned to
     * signify whether the user has aded the movie to their "My Favorites" list.
     *
     * @param id
     * @param params Query parameters for pagination and ordering
     * @param userId
     * @returns List<Movie> similar movies
     */
    // tag::getSimilarMovies[]
    public List<Map<String,Object>> getSimilarMovies(String id, Params params, String userId) {
        // Get similar movies based on genres or ratings
        // MATCH (:Movie {tmdbId: $id})-[:IN_GENRE|ACTED_IN|DIRECTED]->()<-[:IN_GENRE|ACTED_IN|DIRECTED]-(m)

        // Open an Session
        try (var session = this.driver.session()) {

            // Get similar movies based on genres or ratings
            var movies = session.readTransaction(tx -> {
                var favorites = getUserFavorites(tx, userId);
                String query = """            
                          MATCH (:Movie {tmdbId: $id})-[:IN_GENRE|ACTED_IN|DIRECTED]->()<-[:IN_GENRE|ACTED_IN|DIRECTED]-(m)
                          WHERE m.imdbRating IS NOT NULL
                            
                          WITH m, count(*) AS inCommon
                          WITH m, inCommon, m.imdbRating * inCommon AS score
                          ORDER BY score DESC
                            
                          SKIP $skip
                          LIMIT $limit
                            
                          RETURN m {
                            .*,
                              score: score,
                              favorite: m.tmdbId IN $favorites
                          } AS movie
                        """;
                var res = tx.run(query, Values.parameters("id", id, "skip", params.skip(), "limit", params.limit(), "favorites", favorites));
                // Get a list of Movies from the Result
                return res.list(row -> row.get("movie").asMap());
            });

            return movies;
        }
    }
    // end::getSimilarMovies[]


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
     * @param params Query parameters for pagination and ordering
     * @param userId
     * @return List<Movie> movies for that genre
     */
    // tag::getByGenre[]
    public List<Map<String,Object>> byGenre(String name, Params params, String userId) {
        // Get Movies in a Genre
        // MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})

        // Open a new session and close at the end
        try (var session = driver.session()) {
            // Execute a query in a new Read Transaction
            return session.readTransaction((tx) -> {
                // Get an array of IDs for the User's favorite movies
                var favorites = getUserFavorites(tx, userId);

                // Retrieve a list of movies with the
                // favorite flag append to the movie's properties
                var result = tx.run(
                        String.format("""
                                  MATCH (m:Movie)-[:IN_GENRE]->(:Genre {name: $name})
                                  WHERE m.`%s` IS NOT NULL
                                  RETURN m {
                                    .*,
                                      favorite: m.tmdbId IN $favorites
                                  } AS movie
                                  ORDER BY m.`%s` %s
                                  SKIP $skip
                                  LIMIT $limit
                                """, params.sort(), params.sort(), params.order()),
                        Values.parameters("skip", params.skip(), "limit", params.limit(),
                                "favorites", favorites, "name", name));
                var movies = result.list(row -> row.get("movie").asMap());
                return movies;
            });
        }
    }
    // end::getByGenre[]

    /**
     * This method should return a paginated list of movies that have an ACTED_IN relationship
     * to a Person with the id supplied
     *
     * Results should be ordered by the `sort` parameter, and in the direction specified
     * in the `order` parameter.
     * Results should be limited to the number passed as `limit`.
     * The `skip` variable should be used to skip a certain number of rows.
     *
     * If a userId value is suppled, a `favorite` boolean property should be returned to
     * signify whether the user has aded the movie to their "My Favorites" list.
     *
     * @param actorId actor id
     * @param params query params with sorting and pagination
     * @param userId
     * @return List<Movie>
     */
    // tag::getForActor[]
    public List<Map<String,Object>> getForActor(String actorId, Params params,String userId) {
        // Get Movies acted in by a Person
        // MATCH (:Person {tmdbId: $id})-[:ACTED_IN]->(m:Movie)

        // Open a new session
        try (var session = this.driver.session()) {

            // Execute a query in a new Read Transaction
            var movies = session.readTransaction(tx -> {
                // Get an array of IDs for the User's favorite movies
                var favorites = getUserFavorites(tx, userId);
                var sort = params.sort(Params.Sort.title);

                // Retrieve a list of movies with the
                // favorite flag appended to the movie's properties
                String query = String.format("""
                          MATCH (:Person {tmdbId: $id})-[:ACTED_IN]->(m:Movie)
                          WHERE m.`%s` IS NOT NULL
                          RETURN m {
                            .*,
                              favorite: m.tmdbId IN $favorites
                          } AS movie
                          ORDER BY m.`%s` %s
                          SKIP $skip
                          LIMIT $limit
                        """, sort, sort, params.order());
                var res = tx.run(query, Values.parameters("skip", params.skip(), "limit", params.limit(), "favorites", favorites, "id", actorId));
                // Get a list of Movies from the Result
                return res.list(row -> row.get("movie").asMap());
            });
            return movies;
        }
    }
    // end::getForActor[]

    /**
     * This method should return a paginated list of movies that have an DIRECTED relationship
     * to a Person with the id supplied
     *
     * Results should be ordered by the `sort` parameter, and in the direction specified
     * in the `order` parameter.
     * Results should be limited to the number passed as `limit`.
     * The `skip` variable should be used to skip a certain number of rows.
     *
     * If a userId value is suppled, a `favorite` boolean property should be returned to
     * signify whether the user has aded the movie to their "My Favorites" list.
     *
     * @param directorId director id
     * @param params query params with sorting and pagination
     * @param userId
     * @return List<Movie>
     */
    // tag::getForActor[]
    public List<Map<String,Object>> getForDirector(String directorId, Params params,String userId) {
        // Get Movies acted in by a Person
        // MATCH (:Person {tmdbId: $id})-[:DIRECTED]->(m:Movie)

        // Open a new session
        try (var session = this.driver.session()) {

            // Execute a query in a new Read Transaction
            var movies = session.readTransaction(tx -> {
                // Get an array of IDs for the User's favorite movies
                var favorites = getUserFavorites(tx, userId);
                var sort = params.sort(Params.Sort.title);

                // Retrieve a list of movies with the
                // favorite flag appended to the movie's properties
                String query = String.format("""
                          MATCH (:Person {tmdbId: $id})-[:DIRECTED]->(m:Movie)
                          WHERE m.`%s` IS NOT NULL
                          RETURN m {
                            .*,
                              favorite: m.tmdbId IN $favorites
                          } AS movie
                          ORDER BY m.`%s` %s
                          SKIP $skip
                          LIMIT $limit
                        """, sort, sort, params.order());
                var res = tx.run(query, Values.parameters("skip", params.skip(), "limit", params.limit(), "favorites", favorites, "id", directorId));
                // Get a list of Movies from the Result
                return res.list(row -> row.get("movie").asMap());
            });
            return movies;
        }
    }
    // end::getForActor[]


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
