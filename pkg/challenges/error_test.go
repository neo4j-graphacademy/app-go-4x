package challenges_test

import (
	"fmt"
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
	result, err := session.Run(
		"MTCH (n) RETURN x(n)", nil)
	// end::handle[]

	assertNil(outer, result)
	assertNotNil(outer, err)

	// tag::handle[]

	/**
		Neo4jError: Neo.ClientError.Statement.SyntaxError (Invalid input 'T': expected 'a/A' or 'e/E' (line 1, column 2 (offset: 1))
		"MTCH (n) RETURN x(n)"
	  	 ^)]
	*/

	// Cast error as a Neo4j Error
	neo4jError := err.(*neo4j.Neo4jError)

	fmt.Println(neo4jError.Code) // <1> Neo.ClientError.Statement.SyntaxError
	fmt.Println(neo4jError.Msg)  // <2> (Invalid input 'T':...

	// The error code can be further broken down into the following parts:
	fmt.Println(neo4jError.Classification()) // ClientError
	fmt.Println(neo4jError.Category())       // Statement
	fmt.Println(neo4jError.Title())          // SyntaxError
	// end::handle[]
}

func TestConstraintErrors(outer *testing.T) {
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

	// Create a new constraint
	result, err := session.Run(
		"CREATE CONSTRAINT IF NOT EXISTS ON (t:Test) ASSERT t.value IS UNIQUE",
		map[string]interface{}{})

	assertNotNil(outer, result)
	assertNilError(outer, err)

	// Fail constraint validation
	result, resultErr := session.Run(
		"UNWIND range(0, 2) AS row CREATE (t:Test {value: 'notunique'})",
		map[string]interface{}{})

	assertNil(outer, result)

	// tag::constrainterror[]
	// Cast as a Neo4jError
	neo4jError, ok := resultErr.(*neo4j.Neo4jError)

	// Check the Error Code Title
	if ok && neo4jError.Title() == "ConstraintValidationFailed" {
		// Handle the error here
	}
	// end::constrainterror[]

	assertEquals(outer, neo4jError.Title(), "ConstraintValidationFailed")

}
