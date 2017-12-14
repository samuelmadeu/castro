package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	// third party packages
	"github.com/auth0-community/go-auth0"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gopkg.in/square/go-jose.v2/jwt"
)

type Response struct {
	Message string `json:"message"`
}

var AUTH0_API_AUDIENCE = []string{"https://instamatches-backend.herokuapp.com"}

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr := ":" + os.Getenv("PORT")

	r := mux.NewRouter()

	// This route is always accessible
	r.Handle("/api/public", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Message: "Hello from a public endpoint! You don't need to be authenticated to see this.",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))

	// This route is only accessible if the user has a valid access_token with the read:messages scope
	// We are wrapping the checkJwt middleware around the handler function which will check for a
	// valid token and scope.
	r.Handle("/api/private", checkJwt(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{
			Message: "Hello from a private endpoint! You need to be authenticated and have a scope of read:messages to see this.",
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	})))

	http.HandleFunc("/", hello)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World from Go")
}

func checkJwt(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Validate the access_token
		if err != nil {
			// Handle invalid token case
		} else {
			// Ensure the token has the correct scope
			result := checkScope(r, validator, token)
			if result == true {
				// If the token is valid and we have the right scope, we'll pass through the middleware
				h.ServeHTTP(w, r)
			} else {
				response := Response{
					Message: "You do not have the read:messages scope.",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
			}
		}
	})
}

func checkScope(r *http.Request, validator *auth0.JWTValidator, token *jwt.JSONWebToken) bool {
	claims := map[string]interface{}{}
	err := validator.Claims(r, token, &claims)

	if err != nil {
		fmt.Println(err)
		return false
	}

	if strings.Contains(claims["scope"].(string), "read:messages") {
		return true
	} else {
		return false
	}
}
