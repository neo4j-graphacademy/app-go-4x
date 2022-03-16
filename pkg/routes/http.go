package routes

import (
	"encoding/json"
	"net/http"
)

type withStatusCode interface {
	StatusCode() int
}

func serializeJson(writer http.ResponseWriter, result interface{}, err error) {
	if err != nil {
		serializeError(writer, err)
		return
	}
	jsonPayload, err := json.Marshal(result)
	if err != nil {
		serializeError(writer, err)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	_, _ = writer.Write(jsonPayload)
}

func serializeError(writer http.ResponseWriter, err error) {
	writer.Header().Add("Content-Type", "text/plain")
	writeStatusCode(writer, err)
	_, _ = writer.Write([]byte(err.Error()))
}

func writeStatusCode(writer http.ResponseWriter, err error) {
	if errWithCode, ok := err.(withStatusCode); ok {
		writer.WriteHeader(errWithCode.StatusCode())
	} else {
		writer.WriteHeader(500)
	}
}
