package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"server/storage"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Signin handles a sign in request by a client.
// Upon successful sign in, a generated token
// is given as a cookie to client for authorizing
// future requests.
func Signin(db storage.StorageService, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var signinReq AuthRequest

		err := json.NewDecoder(r.Body).Decode(&signinReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if signinReq.Username == "" || signinReq.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			http.Error(w, "Missing required fields.", http.StatusBadRequest)
			return
		}

		uid, ok := db.AuthenticateUser(signinReq.Username, signinReq.Password)
		if !ok {
			http.Error(w, "Invalid username or password.", http.StatusUnauthorized)
			return
		}

		session := db.CreateSession(uid)
		http.SetCookie(w, &http.Cookie{
			Name:  "Token",
			Value: session.Token,
		})
		_, _ = w.Write([]byte("Logged in successfully!"))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(db storage.StorageService, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var signupReq AuthRequest

		err := json.NewDecoder(r.Body).Decode(&signupReq)
		if err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if signupReq.Username == "" || signupReq.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			http.Error(w, "Missing required fields.", http.StatusBadRequest)
			return
		}

		if db.FindUser(signupReq.Username) {
			http.Error(w, "Account already exists.", http.StatusConflict)
			return
		}

		db.CreateUser(signupReq.Username, signupReq.Password)
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("Signed up successfully!"))
	default:
		http.Error(w, "", http.StatusMethodNotAllowed)
	}
}

// GenerateSessionToken returns a token for authorizing
// client requests.
func GenerateSessionToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

// AuthenticateToken checks if token from a request cookie is associated
// to an existing session. If yes, it returns the corresponding
// user ID and true boolean; else, it returns 0 and false.
// If no cookie is found, it returns -1 and false.
func AuthenticateToken(db storage.StorageService, r *http.Request) (int, bool) {
	tokenCookie, err := r.Cookie("Token")
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		return -1, false
	}

	uid := db.FindSession(tokenCookie.Value)
	if uid == 0 {
		return uid, false
	}
	return uid, true
}
