package config

// tag::import[]
import (
	"encoding/json"
	"io/ioutil"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// end::import[]

/**
 * ReadConfig reads the application settings from config.json
 */
// tag::readConfig[]
func ReadConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	if err = json.Unmarshal(file, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// end::readConfig[]

type Config struct {
	Uri      string `json:"NEO4J_URI"`
	Username string `json:"NEO4J_USERNAME"`
	Password string `json:"NEO4J_PASSWORD"`

	Port       int    `json:"APP_PORT"`
	JwtSecret  string `json:"JWT_SECRET"`
	SaltRounds int    `json:"SALT_ROUNDS"`
}

/**
 * Initiate the Neo4j Driver
 *
 * @param {Config} config   Config struct loaded from config.json
 * @returns {neo4j.Driver}	A new Driver instance
 */
// tag::initDriver[]
func NewDriver(settings *Config) (neo4j.Driver, error) {
	// Create new Driver instance
	driver, err := neo4j.NewDriver(
		settings.Uri,
		neo4j.BasicAuth(settings.Username, settings.Password, ""),
	)

	// Check error in driver instantiation
	if err != nil {
		return nil, err
	}

	// Verify Connectivity
	err = driver.VerifyConnectivity()

	// If connectivity fails, handle the error
	if err != nil {
		return nil, err
	}

	return driver, nil
}

// end::initDriver[]
