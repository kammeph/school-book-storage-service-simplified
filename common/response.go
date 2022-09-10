package common

import (
	"encoding/json"
	"net/http"
)

type HttpResponseModel struct {
	Error string `json:"error"`
}

type HttpGenericResponseModel struct {
	Payload any `json:"payload"`
}

func JsonResponse(w http.ResponseWriter, payload any) {
	json.NewEncoder(w).Encode(payload)
}

func HttpResponse(w http.ResponseWriter, payload any) {
	response := HttpGenericResponseModel{Payload: payload}
	JsonResponse(w, response)
}

func HttpResponseWithStatusCode(w http.ResponseWriter, payload any, statusCode int) {
	w.WriteHeader(statusCode)
	HttpResponse(w, payload)
}

func HttpErrorResponse(w http.ResponseWriter, err string) {
	response := HttpResponseModel{err}
	JsonResponse(w, response)
}

func HttpErrorResponseWithStatusCode(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	HttpErrorResponse(w, err)
}
