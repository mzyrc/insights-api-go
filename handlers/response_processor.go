package handlers

import (
	"encoding/json"
	"net/http"
)

type httpResponsePayload struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

func respondWithError(writer http.ResponseWriter, statusCode int, errorMessage string) {
	errorResponse := httpResponsePayload{Error: errorMessage}
	processHTTPResponse(writer, statusCode, errorResponse)
}

func respondWithSuccess(writer http.ResponseWriter, statusCode int, payload interface{}) {
	successResponse := httpResponsePayload{Data: payload}
	processHTTPResponse(writer, statusCode, successResponse)
}

func processHTTPResponse(writer http.ResponseWriter, statusCode int, responsePayload httpResponsePayload) {
	response, _ := json.Marshal(responsePayload)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(response)
}
