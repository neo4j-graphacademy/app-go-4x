package fixtures

import (
	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
	"os"
)

// ReadArray reads the content of a JSON array fixture file
// Note: error handling is a bit brutal here since fixtures will gradually be
// replaced by data coming from a Neo4j instance directly
func ReadArray(path string) (_ []map[string]interface{}, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = ioutils.DeferredClose(file, err)
	}()
	return ioutils.ReadJsonArray(file)
}

// ReadObject reads the content of a JSON object fixture file
// Note: error handling is a bit brutal here since fixtures will gradually be
// replaced by data coming from a Neo4j instance directly
func ReadObject(path string) (_ map[string]interface{}, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = ioutils.DeferredClose(file, err)
	}()
	return ioutils.ReadJson(file)
}
