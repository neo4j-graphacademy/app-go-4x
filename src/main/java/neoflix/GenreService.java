package neoflix;

import org.neo4j.driver.Driver;
import org.neo4j.driver.Values;

import java.util.List;

class GenreService {
    private final Driver driver;

    public GenreService(Driver driver) {
        this.driver = driver;
    }

    /**
     * This method should return a list of genres from the database with a
     * `name` property, `movies` which is the count of the incoming `IN_GENRE`
     * relationships and a `poster` property to be used as a background.
     *
     * [
     *   {
     *    name: 'Action',
     *    movies: 1545,
     *    poster: 'https://image.tmdb.org/t/p/w440_and_h660_face/qJ2tW6WMUDux911r6m7haRef0WH.jpg'
     *   }, ...
     *
     * ]
     *
     * @return List<Genre> genres
     */
    // tag::all[]
    public List<Genre> all() {
        // Open a new Session, close automatically at the end
        try (var session = driver.session()) {
            // Get a list of Genres from the database
            var query = """
                    MATCH (g:Genre)
                    CALL {
                      WITH g
                      MATCH (g)<-[:IN_GENRE]-(m:Movie)
                      WHERE m.imdbRating IS NOT NULL
                      AND m.poster IS NOT NULL
                      AND g.name <> '(no genres listed)'
                      RETURN m.poster AS poster
                      ORDER BY m.imdbRating DESC LIMIT 1
                    }
                    RETURN g {
                      .name,
                      link: '/genres/'+ g.name,
                      poster: poster,
                      movies: size( (g)<-[:IN_GENRE]-() )
                    } as genre
                    ORDER BY g.name ASC
                    """;
            var genres = session.readTransaction(
                    tx -> tx.run(query)
                            .list(row ->
                                row.get("genre")
                                .computeOrDefault(v ->
                                        new Genre(
                                                v.get("name").asString(),
                                                v.get("poster").asString(),
                                                v.get("movies").asLong(),
                                                v.get("link").asString()),
                                        null)
                            ));

                            // alternative .list(row ->row.get("genre").asMap()));

            // Return results
            return genres;
        }
    }
    // end::all[]

    /**
     * This method should find a Genre node by its name and return a set of properties
     * along with a `poster` image and `movies` count.
     *
     * If the genre is not found, a NotFoundError should be thrown.
     *
     * @param name                     The name of the genre
     * @return Genre  The genre information
     */
    // tag::find[]
    public Genre find(String name) {
        // Open a new Session, close automatically at the end
        try (var session = driver.session()) {
            // Get a list of Genres from the database
            var query = """
                    MATCH (g:Genre {name: $name})<-[:IN_GENRE]-(m:Movie)
                    WHERE m.imdbRating IS NOT NULL
                    AND m.poster IS NOT NULL
                    AND g.name <> '(no genres listed)'
                    WITH g, m
                    ORDER BY m.imdbRating DESC

                    WITH g, head(collect(m)) AS movie

                    RETURN g {
                      link: '/genres/'+ g.name,
                      .name,
                      movies: size((g)<-[:IN_GENRE]-()),
                      poster: movie.poster
                    } AS genre
                  """;
            var genre = session.readTransaction(
                    tx -> tx.run(query, Values.parameters("name", name))
                            // Throw a NoSuchRecordException if the genre is not found
                            .single().get("genre")
                                .computeOrDefault(v ->
                                                new Genre(
                                                        v.get("name").asString(),
                                                        v.get("poster").asString(),
                                                        v.get("movies").asLong(),
                                                        v.get("link").asString()
                                                        ),
                                        null)
                            );

            // alternative .list(row ->row.get("genre").asMap()));

            // Return results
            return genre;
        }
    }
    // end::find[]

    public record Genre(String name, String poster, long movies, String link) { }
}
