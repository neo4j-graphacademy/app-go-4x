package neoflix;

import org.neo4j.driver.Driver;
import org.neo4j.driver.Value;

import java.util.List;
import java.util.Map;

class GenreService {
    private final Driver driver;

    public GenreService(Driver driver) {
        this.driver = driver;
    }

    // tag::all[]
    public List<Genre> all() throws Exception {
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
                      .*,
                      poster: poster
                    } as genre
                    ORDER BY g.name ASC
                    """;
            var genres = session.readTransaction(
                    tx -> tx.run(query)
                            .list(row ->
                                row.get("genre")
                                .computeOrDefault(v ->
                                        new Genre(v.get("name").asString(),v.get("poster").asString()),
                                        null)
                            ));

                            // alternative .list(row ->row.get("genre").asMap()));

            // Return results
            return genres;
        }
    }
    // end::all[]

    record Genre(String name, String poster) {}
}
