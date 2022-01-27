package neoflix;

import org.neo4j.driver.AuthToken;
import org.neo4j.driver.AuthTokens;
import org.neo4j.driver.Driver;
import org.neo4j.driver.GraphDatabase;
import spark.Request;

import java.io.IOException;

public class AppUtils {
    static void loadProperties() {
        try {
            System.getProperties().load(ClassLoader.getSystemResourceAsStream("application.properties"));
        } catch (IOException e) {
            throw new RuntimeException("Error loading application.properties", e);
        }
    }

    public static String getUserId(Request req) {
        Object user = req.attribute("user");
        if (user == null) return null;
        return user.toString();
    }

    static void handleAuthAndSetUser(Request req, String jwtSecret) {
        String token = req.headers("Authorization");
        String bearer = "Bearer ";
        if (token != null && !token.isBlank() && token.startsWith(bearer)) {
            // verify token
            token = token.substring(bearer.length());
            String userId = AuthUtils.verify(token, jwtSecret);
            req.attribute("user", userId);
        }
    }

    static Driver initDriver() {
        AuthToken auth = AuthTokens.basic(System.getProperty("NEO4J_USERNAME"), System.getProperty("NEO4J_PASSWORD"));
        Driver driver = GraphDatabase.driver(System.getProperty("NEO4J_URI"), auth);
        driver.verifyConnectivity();
        return driver;
    }

    static int getServerPort() {
        return Integer.parseInt(System.getProperty("APP_PORT", "3000"));
    }

    static String getJwtSecret() {
        return System.getProperty("JWT_SECRET");
    }
}
