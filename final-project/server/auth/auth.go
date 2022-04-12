package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"server/storage"
)

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
// Upon successful sign in, a generated token
// is given as a cookie to client for authorizing
// future requests.
func Signin(db storage.Service, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
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
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(db storage.Service, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
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

		exists, err := db.FindUser(signupReq.Username)
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
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// GenerateSessionToken returns a token for authorizing
// client requests.
func generateSessionToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// AuthenticateToken checks if token from a request cookie is associated
// to an existing session. If yes, it returns the corresponding
// username. If not, or no cookie is found, it returns an empty string.
func AuthenticateToken(db storage.Service, w http.ResponseWriter, r *http.Request) string {
	tokenCookie, err := r.Cookie("Token")
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err)
		http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
		return ""
	}

	username, err := db.FindSession(tokenCookie.Value)
	if username == "" {
		http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
		return ""
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return ""
	}
	return username
}
