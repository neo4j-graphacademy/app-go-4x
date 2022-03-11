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
	// Defer the closing of the Driver
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

	// tag::session.readTransaction[]
	result, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (n) RETURN count(n) AS count", map[string]interface{}{})
		if err != nil {
			return nil, err
		}

		return result, result.Err()
	})
	// end::session.readTransaction[]

	return "", nil
}

func ExplicitTranactionExample() (string, error) {
	driver, err := neo4j.NewDriver("neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	defer driver.Close()

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

// tag::getActors[]
func GetActors() (nil, error) {
	// <1> Initiate Driver
	driver, err := neo4j.NewDriver("neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", "letmein", ""))

	// <2> Check for driver instantiation error
	if err != nil {
		return "", err
	}

	// <3> Defer closing of the driver
	defer driver.Close()

	// <4> Create a new Session
	session := driver.NewSession(neo4j.SessionConfig{DatabaseName: "movies", AccessMode: neo4j.AccessModeWrite})

	// <5> Defer closing the session
	defer session.Close()

	// <6> Execute Cypher and get Result
	// tag::run[]
	result, queryErr := session.Run(
		"MATCH (p:Person)-[r:ACTED_IN]->(m:Movie {title: $title}) RETURN p, r, m",
		map[string]interface{}{"title": "Arthur"})
	// end::run[]

	// <7> Handle Query Errors
	if queryErr != nil {
		// Problem executing the query, maybe a syntax error?
		return "", queryErr
	}

	// <8> For each Record in the Result
	for result.Next() {
		// <9> Get the next record
		record := result.Record()

		// <10> Interact with the record object
		// tag::keys[]
		fmf.Println(record.Keys) // ['p', 'r', 'm']
		// end::keys[]
		// tag::index[]
		// Access a value by its index
		fmf.Println(record.Values[0].(neo4j.Node)) // The Person node
		// end::index[]
		// tag::alias[]
		// Access a value by its alias
		fmf.Println(record.Values["movie"].(neo4j.Node)) // The Movie node
		// end::alias[]
	}

	return nil, nil
}

// end::getActors[]

/*
Shortform examples

// tag::Single[]
// Get the first and only result from the stream.
first, err := record.Single()
// end::Single[]

// tag::Next[]

// .Next() returns false upon error
for result.Next() {
    record := result.Record()
    handleRecord(record)
}
// Err returns the error that caused Next to return false
if err = result.Err(); err != nil {
    handleError(err)
}

// end::Next[]


// tag::NextRecord[]
for result.NextRecord(&record) {
    fmf.Println(record.Keys)
}
// end::NextRecord[]

// tag::Consume[]
summary := result.Consume()

// Time in milliseconds before receiving the first result
fmt.Println(summary.ResultAvailableAfter())

// Time in milliseconds once the final result was consumed
fmt.Println(summary.ResultConsumedAfter())
// end::Consume[]


// tag::Collect[]
remaining, remainingErr := result.Collect()
// end::Collect[]



*/
