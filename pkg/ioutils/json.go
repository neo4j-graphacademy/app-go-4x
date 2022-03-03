package ioutils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ReadJson(r io.Reader) (map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func ReadJsonArray(r io.Reader) ([]map[string]interface{}, error) {
	body, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var results []map[string]interface{}
	if err = json.Unmarshal(body, &results); err != nil {
		return nil, err
	}
	return results, nil
}
