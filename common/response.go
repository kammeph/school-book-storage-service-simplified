package common

import (
	"encoding/json"
	"net/http"
)

type HttpResponseModel struct {
	Error string `json:"error"`
}

func JsonResponse(w http.ResponseWriter, response any) {
	json.NewEncoder(w).Encode(response)
}

func HttpResponse(w http.ResponseWriter, response any) {
	JsonResponse(w, response)
}

func HttpResponseWithStatusCode(w http.ResponseWriter, payload any, statusCode int) {
	w.WriteHeader(statusCode)
	HttpResponse(w, payload)
}

func HttpSuccessResponse(w http.ResponseWriter) {
	response := HttpResponseModel{""}
	JsonResponse(w, response)
}

func HttpErrorResponse(w http.ResponseWriter, err string) {
	response := HttpResponseModel{err}
	JsonResponse(w, response)
}

func HttpErrorResponseWithStatusCode(w http.ResponseWriter, err string, statusCode int) {
	w.WriteHeader(statusCode)
	HttpErrorResponse(w, err)
}
