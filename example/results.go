package main

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"reflect"
)

func GetActorsForMovie(movie string) ([]string, error) {
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return nil, err
	}
	// end::driver[]

	// tag::close[]
	// Defer the closing of the Driver
	defer driver.Close()
	// end::close[]

	// tag::verifyConnectivity[]
	err = driver.VerifyConnectivity()
	if err != nil {
		return nil, err
	}
	// end::verifyConnectivity[]

	// tag::get_actors[]
	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	names, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
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

		var record *neo4j.Record
		var actorNames []string

		// tag::keys[]
		fmt.Println(result.Keys()) // ['p', 'm', 'a', 'path']
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
			movie, _ := record.Get("m")
			movieByAlias := movie.(neo4j.Node) // m, a :Movie node as type neo4j.Node()
			// end::alias[]
			fmt.Println(movieByAlias) // ['p', 'm', 'a', 'path']

			// Add the Person name property to the array
			actorNames[i] = personByIndex.Props["name"].(string)

			// tag::cast[]
			// Get a Relationship and check assertion
			actedInRelationship, _ := record.Get("r")
			actedIn, ok := actedInRelationship.(neo4j.Relationship)

			if !ok {
				// Value is not a relationship
				return nil, fmt.Errorf("expected a neo4j.Relationship, got: %v", reflect.TypeOf(result))
			}
			// end::cast[]

			// tag::node[]
			personNode, _ := record.Get("p")
			person := personNode.(neo4j.Node)

			fmt.Println(person.Id)     // <1>
			fmt.Println(person.Labels) // <2>
			fmt.Println(person.Props)  // <3>
			// end::node[]

			// tag::relationship[]
			actedInRelationship, _ = record.Get("p")
			actedIn = actedInRelationship.(neo4j.Relationship)

			fmt.Println(actedIn.Id)      // <1>
			fmt.Println(actedIn.Type)    // <2>
			fmt.Println(actedIn.Props)   // <3>
			fmt.Println(actedIn.StartId) // <4>
			fmt.Println(actedIn.EndId)   // <5>
			// end::relationship[]

			// tag::path[]
			returnedPath, _ := record.Get("path")
			path := returnedPath.(neo4j.Path)

			nodes := path.Nodes
			relationships := path.Relationships

			for node := range nodes {
				fmt.Println(node) // neo4j.Node
			}
			for relationship := range relationships {
				fmt.Println(relationship) // neo4j.Relationship
			}
			// end::path[]

			// Add the Person name property to the array
			actorNames[i] = personByIndex.Props["name"].(string)
			// tag::nextrecord[]
		}
		// end::nextrecord[]

		return actorNames, nil

	})
	if err != nil {
		return nil, err
	}
	// end::get_actors[]
	return names.([]string), nil

}

func TimeExamples(result neo4j.Result) {
	record := result.Record()

	// tag::time[]
	timeProperty, _ := record.Get("time")
	time := timeProperty.(neo4j.Time).Time()

	fmt.Println(time.Year())  // 2022
	fmt.Println(time.Month()) // January
	fmt.Println(time.Day())   // 4

	// For Time, DateTime,
	fmt.Println(time.Day()) // 4
	// end::time[]

}

func DurationExample() (neo4j.Duration, error) {
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))

	if err != nil {
		return neo4j.Duration{}, err
	}

	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{})

	defer session.Close()

	result, err := session.Run("RETURN duration('P1Y2M3DT12H34M56S.9876') AS duration", map[string]interface{}{})
	if err != nil {
		return neo4j.Duration{}, err
	}
	record := result.Record()

	// tag::duration[]
	// duration('P1Y2M3DT12H34M56S')
	// 1 year, 2 months, 3 days; 12 hours, 34 minutes, 56 seconds
	durationProperty, _ := record.Get("duration")
	duration := durationProperty.(neo4j.Duration)

	fmt.Println(duration.Months)  // 14 (1 year, 2 months = 14 months)
	fmt.Println(duration.Days)    // 3
	fmt.Println(duration.Seconds) // 45296
	fmt.Println(duration.Nanos)   // 987600000
	// end::duration[]

	return duration, nil
}

func PointExample() error {
	driver, err := neo4j.NewDriver("neo4j+s://dbhash.databases.neo4j.io",
		neo4j.BasicAuth("neo4j", "letmein", ""))
	if err != nil {
		return err
	}
	defer driver.Close()

	session := driver.NewSession(neo4j.SessionConfig{})
	defer session.Close()

	result, err := session.Run(`RETURN
		point({longitude:20, latitude:10}) AS wgs842D,
		point({longitude:20, latitude:10, height:30}) AS wgs843D,
		point({x:20, y:10}) AS cartesian2D,
		point({x:20, y:10, z:30}) AS cartesian3D`, map[string]interface{}{})
	if err != nil {
		return err
	}

	record := result.Record()

	// tag::point2d[]
	wgs842DResult, _ := record.Get("wgs842D")
	wgs842D := wgs842DResult.(neo4j.Point2D)
	// {SpatialRefId: xxxx, X: 10, y: 20}

	cartesian2DResult, _ := record.Get("cartesian2D")
	cartesian2D := cartesian2DResult.(neo4j.Point2D)
	// {SpatialRefId: xxxx, X: 10, y: 20}
	// end::point2d[]
	fmt.Printf("%v\n", wgs842D)
	fmt.Printf("%v\n", cartesian2D)

	// tag::point3d[]
	wgs843DResult, _ := record.Get("wgs843D")
	wgs843D := wgs843DResult.(neo4j.Point3D)
	// {SpatialRefId: xxxx, X: 10, y: 20, z: 30}
	cartesian3DResult, _ := record.Get("cartesian3D")
	cartesian3D := cartesian3DResult.(neo4j.Point3D)
	// {SpatialRefId: xxxx, X: 10, y: 20, z: 30}
	// tag::point3d[]
	fmt.Printf("%v\n", wgs843D)
	fmt.Printf("%v\n", cartesian3D)
	return nil
}
