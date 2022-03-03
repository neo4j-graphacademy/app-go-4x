package config

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io/ioutil"
)

// ReadConfig reads the application settings from config.json
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

func NewDriver(settings *Config) (neo4j.Driver, error) {
	return neo4j.NewDriver(
		settings.Uri,
		neo4j.BasicAuth(settings.Username, settings.Password, ""),
	)
}
