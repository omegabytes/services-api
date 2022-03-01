package handlers

import (
	"encoding/json"
	"fmt"
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
		resBadRequest(w, "Name and password are required")
		return
	}

	// Realistically I'd have a datastore with user information
	if name == "user" && password == "pass" {
		token, err := getToken(name)
		if err != nil {
			resInternalError(w, "Error generating token: "+err.Error())
		} else {
			w.Header().Set("Authorization", "Bearer "+token)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("{Token:%s}", token)))
		}
	} else {
		resUnauthorized(w, "Unable to sign in with that user or password")
		return
	}
}

// AuthMiddleware can be used to wrap any endpoint/handler that requires authorization.
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if len(token) == 0 {
			resUnauthorized(w, "Missing authorization header")
			return
		}
		token = strings.Replace(token, "Bearer ", "", 1)
		claims, err := verifyToken(token)
		if err != nil {
			resUnauthorized(w, "Error verifying token: "+err.Error())
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

func resUnauthorized(w http.ResponseWriter, message string) {
	response := map[string]interface{}{
		"status":  http.StatusUnauthorized,
		"message": message,
	}
	w.WriteHeader(http.StatusUnauthorized)
	resp, _ := json.Marshal(response)
	w.Write(resp)
}
