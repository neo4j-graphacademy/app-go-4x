package fixtures

import (
	"os"
	"path/filepath"

	"github.com/neo4j-graphacademy/neoflix/pkg/ioutils"
)

type FixtureLoader struct {
	Prefix string
}

// ReadArray reads the content of a JSON array fixture file
// Note: error handling is a bit brutal here since fixtures will gradually be
// replaced by data coming from a Neo4j instance directly
func (fl *FixtureLoader) ReadArray(fixture string) (_ []map[string]interface{}, err error) {
	newPath := filepath.Join(fl.Prefix, fixture)

	file, err := os.Open(newPath)
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
func (fl *FixtureLoader) ReadObject(fixture string) (_ map[string]interface{}, err error) {
	newPath := filepath.Join(fl.Prefix, fixture)

	file, err := os.Open(newPath)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = ioutils.DeferredClose(file, err)
	}()
	return ioutils.ReadJson(file)
}
