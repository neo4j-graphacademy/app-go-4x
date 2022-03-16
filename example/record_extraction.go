package main

import (
	"fmt"
	"reflect"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// tag::single[]
func recordExtractSingleExample(queryResult neo4j.Result) (string, error) {
	singleRecord, err := queryResult.Single()
	if err != nil {
		// oh no! 0 or 2+ results
		return "", err
	}
	value, found := singleRecord.Get("some-column")
	if !found {
		// probable typo: some-column is not in the returned row
		return "", fmt.Errorf("some-column not found")
	}
	// let's say we get a string
	result, ok := value.(string)
	if !ok {
		// oh no! it's not a string
		return "", err
	}
	return result, nil
}

// end::single[]

// tag::collect[]
func recordExtractCollectExample(queryResult neo4j.Result) ([]bool, error) {
	// buffers everything in memory
	records, err := queryResult.Collect()
	if err != nil {
		// oh no! sth went wrong when fetching one of the results
		return nil, err
	}
	// this time, all the results are boolean
	// we know the size in advance, so let's allocate everything now!
	results := make([]bool, len(records))
	for i, record := range records {
		value, found := record.Get("some-boolean-column")
		if !found {
			// probable typo: some-boolean-column is not in the returned row
			return nil, fmt.Errorf("some-boolean-column not found in record number %d", i)
		}
		result, ok := value.(bool)
		if !ok {
			// oh no! it's not a bool
			return nil, fmt.Errorf("expected boolean, got: %v", reflect.TypeOf(result))
		}
		results[i] = result
	}
	return results, nil
}

// end::collect[]

// tag::next[]
func recordExtractNextRecordExample(queryResult neo4j.Result) ([]neo4j.Duration, error) {
	// this time, we do not know the size in advance, we'll allocate and grow the slice as we go
	var results []neo4j.Duration
	// alternatively to loop below:
	//var record *neo4j.Record
	//for queryResult.NextRecord(&record) {
	//	// ...
	//}
	var i int
	for queryResult.Next() {
		i++
		// get the current record
		record := queryResult.Record()
		value, found := record.Get("some-duration-column")
		if !found {
			// probable typo: some-duration-column is not in the returned row
			return nil, fmt.Errorf("some-duration-column not found in record number %d", i)
		}
		result, ok := value.(neo4j.Duration)
		if !ok {
			// oh no! it's not a neo4j.Duration
			return nil, fmt.Errorf("expected neo4j.Duration, got: %v", reflect.TypeOf(result))
		}
		results[i] = result
	}
	return results, nil
}

// end::next[]
