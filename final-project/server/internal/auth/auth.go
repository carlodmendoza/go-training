package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/storage"
)

type userCtx string

const UserKey userCtx = "username"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrDuplicateUser = errors.New("User already exists")
	ErrInvalidLogin  = errors.New("Invalid username or password")
	ErrInvalidToken  = errors.New("Invalid or missing token")
)

// Signin handles a sign in request by a client.
// Upon successful sign in, a generated token is given as a cookie to client for authorizing future requests.
func Signin(db storage.Service, w http.ResponseWriter, r *http.Request) {
	var signinReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signinReq)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if signinReq.Username == "" || signinReq.Password == "" {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
		http.Error(w, "Missing required fields.", http.StatusBadRequest)
		return
	}

	ok, err := db.AuthenticateUser(signinReq.Username, signinReq.Password)
	if !ok {
		http.Error(w, ErrInvalidLogin.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	token := generateSessionToken()
	err = db.CreateSession(signinReq.Username, token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Token",
		Value: token,
	})
	_, _ = w.Write([]byte("Logged in successfully!"))
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(db storage.Service, w http.ResponseWriter, r *http.Request) {
	var signupReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if signupReq.Username == "" || signupReq.Password == "" {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
		http.Error(w, "Missing required fields.", http.StatusBadRequest)
		return
	}

	exists, err := db.UserExists(signupReq.Username)
	if exists {
		http.Error(w, ErrDuplicateUser.Error(), http.StatusConflict)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = db.CreateUser(signupReq.Username, signupReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Signed up successfully!"))
}

// GenerateSessionToken returns a token for authorizing client requests.
func generateSessionToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// Authenticator is a middleware that checks if token from a request cookie is associated to an existing session.
// If yes, it passes the username to the request context then proceeds with the request.
// If not, or no cookie is found, an error response is sent.
func Authenticator(db storage.Service) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			tokenCookie, err := r.Cookie("Token")
			if err != nil {
				fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}

			session, err := db.FindSession(tokenCookie.Value)
			if session.Username == "" {
				http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
				return
			}
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, session.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

// GetUser returns the username from a request context using UserKey as key.
func GetUser(r *http.Request) string {
	switch v := r.Context().Value(UserKey).(type) {
	case string:
		return v
	default:
		return ""
	}
}
