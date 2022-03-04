package main

// tag::import[]
import (
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// end::import[]

/*
// tag::pseudo[]
driver = neo4j.newDriver(
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

// tag::createPerson[]
func helloWorld(name string) (string, error) {
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
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()
	// end::session[]

	// tag::session.writeTransaction[]
	name, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (p:Person {name: $name}) RETURN p",
			map[string]interface{}{"name": name})
		if err != nil {
			return nil, err
		}

		person := result.Record().Values[0].(neo4j.Node)

		return person.Props["name"], result.Err()
	})
	if err != nil {
		return "", err
	}
	// end::session.writeTransaction[]

	return name.(string), nil
}

// end::createPerson[]

func SessionRunExample() (string, error) {
	driver, err := neo4j.NewDriver("neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	// tag::sessionWithArgs[]
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "movies", AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	// end::sessionWithArgs[]

	// tag::session.run[]
	result, err = session.Run(
		"MATCH (p:Person {name: $name}) RETURN p",
		map[string]interface{}{"name": "Tom Hanks"})
	// end::session.run[]

	return "", nil
}

func ReadTransactionExample() (string, error) {
	driver, err := neo4j.NewDriver("neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "reviews", AccessMode: neo4j.AccessModeWrite})

	// tag::readTransaction[]
	result, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n) RETURN count(n) AS count", map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		return result, result.Err()
	})
	// end::readTransaction[]

	return "", nil
}

func ExplicitTranactionExample() (string, error) {
	driver, err := neo4j.NewDriver("neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "reviews", AccessMode: neo4j.AccessModeWrite})

	// tag::session.close[]
	defer session.Close()
	// end::session.close[]

	// tag::session.beginTransaction.Try[]
	// tag::session.beginTransaction[]
	// Begin Transaction
	tx, err := session.BeginTransaction()
	// end::session.beginTransaction[]
	if err != nil {
		return "", err
	}

	// Run a Cypher Query
	result, err = tx.Run(cypher, params)

	// If something goes wrong then rollback the transaction
	if err != nil {
		tx.Rollback()

		return "", err
	}

	// Otherwise, commit the transaction
	tx.Commit()
	// end::session.beginTransaction.Try[]

	return "", nil
}
