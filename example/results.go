func GetActorsForMovie(movie string) {
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

	// tag::get_actors[]
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

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

	// end::get_actors[]
}