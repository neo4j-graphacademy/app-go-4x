package main

// tag::import[]
import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// end::import[]

/*
// tag::pseudo[]
driver = GraphDatabase.driver(
  connectionString, // <1>
  auth=(username, password), // <2>
  **configuration // <3>
)
// end::pseudo[]

// tag::connection[]
  address of server
          ↓
neo4j://localhost:7687
  ↑                ↑
scheme        port number
// end::connection[]
*/

func basicAuth() (neo4j.AuthToken, error) {
	username := "neo4j"
	password := "letmein"

	auth :=
		// tag::auth[]
		neo4j.BasicAuth(username, password, "")
	// end::auth[]

	return auth, nil
}

func helloWorld() (string, error) {
	// tag::driver[]
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}
	// end::driver[]

	// tag::close[]
	defer driver.Close()
	// end::close[]

	// tag::verifyConnectivity[]
	err = driver.VerifyConnectivity()
	if err != nil {
		return "", err
	}
	// end::verifyConnectivity[]

	// tag::session[]
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	// end::session[]

	// tag::writeTransaction[]
	greeting, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (a:Greeting) SET a.message = $message RETURN a.message + ', from node ' + id(a)",
			map[string]interface{}{"message": "hello, world"})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return "", err
	}
	// end::writeTransaction[]

	return greeting.(string), nil
}
