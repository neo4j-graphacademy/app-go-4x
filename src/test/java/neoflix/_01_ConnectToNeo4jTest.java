package neoflix;

import org.junit.jupiter.api.Test;
import org.neo4j.driver.Driver;

import static neoflix.NeoflixApp.PROPS;
import static org.junit.jupiter.api.Assertions.*;

class _01_ConnectToNeo4jTest {

    @Test
    void createDriverAndConnectToServer() {
        assertNotNull(PROPS.getProperty("NEO4J_URI"), "neo4j uri defined");
        assertNotNull(PROPS.getProperty("NEO4J_USERNAME"), "username defined");
        assertNotNull(PROPS.getProperty("NEO4J_PASSWORD"), "password defined");

        Driver driver = NeoflixApp.initDriver();
        assertNotNull(driver, "driver instantiated");
        assertDoesNotThrow(driver::verifyConnectivity,"unable to verify connectivity");
    }
}