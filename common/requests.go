package common

import (
	"fmt"
	"net/http"
	"os"
)

var (
	corsAllowOriginProtocol = os.Getenv("CORS_ALLOW_ORIGIN_PROTOCOL")
	corsAllowOriginHost     = os.Getenv("CORS_ALLOW_ORIGIN_HOST")
	corsAllowOriginPort     = os.Getenv("CORS_ALLOW_ORIGIN_PORT")
	corsAllowOrigin         = fmt.Sprintf("%s://%s:%s", corsAllowOriginProtocol, corsAllowOriginHost, corsAllowOriginPort)
)

func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func setContentTypeJson(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
}

func Get(pattern string, handler http.HandlerFunc) {
	Request(pattern, handler, http.MethodGet)
}

func Post(pattern string, handler http.HandlerFunc) {
	Request(pattern, handler, http.MethodPost)
}

func Request(pattern string, handler http.HandlerFunc, method string) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		setContentTypeJson(&w)
		setupCORS(&w)
		if r.Method == method {
			handler(w, r)
			return
		}
		if r.Method != http.MethodOptions {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
}
