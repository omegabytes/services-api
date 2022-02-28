package handlers

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Authenticate generates a token for a known user with a valid password.
// The generated token is used in subsequent requests to versioned API endpoints.
func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	if len(name) == 0 || len(password) == 0 {
		http.Error(w, "Name and password are required", http.StatusBadRequest)
		return
	}

	// Realistically I'd have a datastore with user information
	if name == "user" && password == "pass" {
		token, err := getToken(name)
		if err != nil {
			http.Error(w, "Error generating token: "+err.Error(), http.StatusInternalServerError)
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Token: " + token))
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unable to sign in with that user or password"))
		return
	}
}

// AuthMiddleware can be used to wrap any endpoint/handler that requires authorization.
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}
		token = strings.Replace(token, "Bearer ", "", 1)
		claims, err := verifyToken(token)
		if err != nil {
			http.Error(w, "Error verifying token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		name := claims.(jwt.MapClaims)["name"].(string)
		role := claims.(jwt.MapClaims)["role"].(string)

		r.Header.Set("name", name)
		r.Header.Set("role", role)

		next.ServeHTTP(w, r)
	})
}

// getToken generates the token for a user
func getToken(name string) (string, error) {
	signingKey := []byte("signingKey")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name": name,
		"role": "member",
	})
	tokenString, err := token.SignedString(signingKey)
	return tokenString, err
}

// verifyToken validates a given token.
func verifyToken(tokenString string) (jwt.Claims, error) {
	signingKey := []byte("signingKey")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	return token.Claims, err
}
