func GetActorsForMovie(movie string): ([]string, error) {
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
		// tag::run[]
		result, err := transaction.Run(
			"MATCH path = (p:Person)-[r:ACTED_IN]->(m:Movie {title: $title}) RETURN p, r, m, path",
			map[string]interface{}{"title": "Arthur"})

		if err != nil {
			return nil, err
		}
		// end::run[]

		// Prepare an array to return
		i := 0
		actorNames = make([]string)

		// tag::keys[]
		fmf.PrintLn(record.Keys) // ['p', 'm', 'a', 'path']
		// end::keys[]

		// tag::nextrecord[]
		for result.NextRecord(&record) {
		// end::nextrecord[]
			i++

			// tag::index[]
			// Get the first value, in this case `p`
			personByIndex := record.Values[0].(neo4j.Node) // p, a :Person node as type neo4j.Node()
			// end::index[]

			// tag::alias[]
			// Get the value of `m`
			movieByAlias := record.Values["m"].(neo4j.Node) // m, a :Movie node as type neo4j.Node()
			// end::alias[]

			// Add the Person name property to the array
			actorNames[i] = personByIndex.Props["name"]


			// tag::cast[]
			// Get a Relationship and check assertion
			actedIn, ok := record.Values["r"].(neo4j.Relationship)

			if !ok {
				// Value is not a relationship
				return nil, fmt.Errorf("expected a neo4j.Relationship, got: %v", reflect.TypeOf(result))
			}
			// end::cast[]


			// tag::node[]
			person := record.Values["p"]

			fmf.PrintLn(person.Id) // <1>
			fmf.PrintLn(person.Labels) // <2>
			fmf.PrintLn(person.Props) // <3>
			// end::node[]

			// tag::relationship[]
			actedIn := record.Values["p"]

			fmf.PrintLn(actedIn.Id) // <1>
			fmf.PrintLn(actedIn.Type) // <2>
			fmf.PrintLn(actedIn.Props) // <3>
			fmf.PrintLn(actedIn.StartId) // <4>
			fmf.PrintLn(actedIn.EndId) // <5>
			// end::relationship[]

			// tag::path[]
			path := record.Values["path"]

			nodes := path.Nodes.([]neo4j.Node)
			relationships := path.Relationships.([]neo4j.Relationship)

			for segment := range path {
				fmf.PrintLn(segment) // neo4j.Relationship
			}
			// end::path[]

			// Add the Person name property to the array
			actorNames[i] = personByIndex.Props["name"]



		// tag::nextrecord[]
		}
		// end::nextrecord[]

		return actorNames, nil

	})
	if err != nil {
		return [], err
	}

	// end::get_actors[]
}

func TimeExamples(result neo4j.Result) {
	record := result.Record()

	// tag::time[]
	time := record.Values["time"]

	fmf.PrintLn(time.Year())  // 2022
	fmf.PrintLn(time.Month()) // January
	fmf.PrintLn(time.Day())   // 4

	// For Time, DateTime,
	fmf.PrintLn(time.Day())   // 4
	// end::time[]


}

func DurationExample() (neo4j.Duration, error) {
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	defer driver.Close()

	session := driver.Session()

	defer session.Close()

	result := session.Run("RETURN duration('P1Y2M3DT12H34M56S.9876') AS duration")
	record := result.Record()

	// tag::duration[]
	// duration('P1Y2M3DT12H34M56S')
	// 1 year, 2 months, 3 days; 12 hours, 34 minutes, 56 seconds
	duration := record.Values["duration"].(neo4j.Duration)

	fmf.PrintLn(duration.Months)  // 14 (1 year, 2 months = 14 months)
	fmf.PrintLn(duration.Days)	  // 3
	fmf.PrintLn(duration.Seconds) // 45296
	fmf.PrintLn(duration.Nanos)	  // 987600000
	// end::duration[]

	return duration, nil
}

func PointExample() {
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return "", err
	}

	defer driver.Close()

	session := driver.Session()

	defer session.Close()

	result := session.Run("RETURN
		point({longitude:20, latitude:10}) AS wgs842D,
		point({longitude:20, latitude:10, height:30}) AS wgs843D,
		point({x:20, y:10}) AS cartesian2D,
		point({x:20, y:10, z:30}) AS cartesian3D")
	record := result.Record()

	// tag::point2d[]
	wgs842D := result.Values["wgs842D"]
	// {SpatialRefId: xxxx, X: 10, y: 20}

	cartesian2D := result.Values["cartesian2D"]
	// {SpatialRefId: xxxx, X: 10, y: 20}
	// end::point2d[]

	// tag::point3d[]
	wgs843D := result.Values["wgs843D"]
	// {SpatialRefId: xxxx, X: 10, y: 20, z: 30}
	cartesian3D := result.Values["cartesian3D"]
	// {SpatialRefId: xxxx, X: 10, y: 20, z: 30}
	// tag::point3d[]


}