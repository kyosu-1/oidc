package main

import (
	"encoding/json"
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

// DiscoveryResponse is the response for the discovery endpoint
// ref: https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type discoveryResponse struct {
	Issuer                            string   `json:"issuer"`
	AuthEndpoint                      string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserInfoEndpoint                  string   `json:"userinfo_endpoint"`
	JWKSURI                           string   `json:"jwks_uri"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	GrantTypesSupported               []string `json:"grant_types_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
}

// supported scopes
var (
	SupportedResponseTypes            = []string{"code"}
	SupportedGrantTypes               = []string{"authorization_code"}
	SupportedScopes                   = []string{"openid"}
	SupportedTokenEndpointAuthMethods = []string{"client_secret_post", "client_secret_basic"}
	SupportedSubjectTypes             = []string{"public"}
	SupportedIDTokenSigningAlgs       = []string{"RS256"}
	SupportedClaims                   = []string{"aud", "exp", "iat", "iss", "sub"}
)

func newDicoveryResponse(host string) *discoveryResponse {
	return &discoveryResponse{
		Issuer:                            formalURL(host, ""),
		AuthEndpoint:                      formalURL(host, authorizationEndpoin),
		TokenEndpoint:                     formalURL(host, tokenEndpoint),
		UserInfoEndpoint:                  formalURL(host, userinfoEndpoint),
		JWKSURI:                           formalURL(host, jwksEndpoint),
		ScopesSupported:                   SupportedScopes,
		ResponseTypesSupported:            SupportedResponseTypes,
		GrantTypesSupported:               SupportedGrantTypes,
		SubjectTypesSupported:             SupportedSubjectTypes,
		IDTokenSigningAlgValuesSupported:  SupportedIDTokenSigningAlgs,
		ClaimsSupported:                   SupportedClaims,
		TokenEndpointAuthMethodsSupported: SupportedTokenEndpointAuthMethods,
	}
}

// https formal url
func formalURL(host, path string) string {
	return fmt.Sprintf("https://%s%s", host, path)
}

// discovery endpoint
func discovery(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, newDicoveryResponse(r.Host))
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
