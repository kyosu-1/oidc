package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

const (
	healthEndpoint       = "/health"
	discoveryEndpoint    = "/.well-known/openid-configuration"
	authorizationEndpoin = "/oauth2/authorize"
	tokenEndpoint        = "/oauth2/token"
	userinfoEndpoint     = "/userinfo"
	jwksEndpoint         = "/jwks.json"
)

func main() {
	r := chi.NewRouter()
	r.Get(healthEndpoint, health)
	r.Get(discoveryEndpoint, discovery)
	r.Get(authorizationEndpoin, authorize)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("Server started")

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
