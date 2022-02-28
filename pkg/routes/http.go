package routes

import (
	"encoding/json"
	"net/http"
)

func serializeJson(writer http.ResponseWriter, result interface{}, err error) {
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	jsonPayload, err := json.Marshal(result)
	if err != nil {
		writer.WriteHeader(500)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(200)
	_, _ = writer.Write(jsonPayload)
}
