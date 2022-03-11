package challenges_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/neo4j-graphacademy/neoflix/pkg/config"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestErrors(outer *testing.T) {
	// Load Settings
	settings, err := config.ReadConfig("../../config.json")
	assertNilError(outer, err)

	// Init Driver
	driver, err := config.NewDriver(settings)
	assertNilError(outer, err)

	defer func() {
		assertNilError(outer, driver.Close())
	}()

	session := driver.NewSession(neo4j.SessionConfig{})

	// tag::handle[]
	result, err := session.Run("MTCH (n) RETURN x(n)", map[string]interface{}{})
	// end::handle[]

	assertNil(outer, result)
	assertNotNil(outer, err)

	// tag::handle[]

	/**
		Neo4jError: Neo.ClientError.Statement.SyntaxError (Invalid input 'T': expected 'a/A' or 'e/E' (line 1, column 2 (offset: 1))
		"MTCH (n) RETURN x(n)"
	  	 ^)]
	*/

	parts := strings.SplitN(err.Error(), " ", 3)

	fmt.Println(parts[0]) // <1> Neo4jError:
	fmt.Println(parts[1]) // <2> Neo.ClientError.Statement.SyntaxError
	fmt.Println(parts[2]) // <3> (Invalid input 'T':...
	// end::handle[]

}
