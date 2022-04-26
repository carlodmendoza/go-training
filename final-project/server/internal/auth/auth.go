package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	gohttp "net/http"

	"github.com/carlodmendoza/go-training/final-project/server/pkg/http"
	"github.com/carlodmendoza/go-training/final-project/server/storage"
)

type userCtx string

const UserKey userCtx = "username"

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	ErrDuplicateUser = errors.New("User already exists")
	ErrEmptyFields   = errors.New("Username or password is empty")
	ErrInvalidLogin  = errors.New("Invalid username or password")
	ErrInvalidToken  = errors.New("Invalid or missing token")
)

// Signin handles a sign in request by a client.
// Upon successful sign in, a generated token is given as a cookie to client for authorizing future requests.
func Signin(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	var signinReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signinReq)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: err}
	}

	if signinReq.Username == "" || signinReq.Password == "" {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: ErrEmptyFields}
	}

	ok, err := db.AuthenticateUser(signinReq.Username, signinReq.Password)
	if !ok {
		return http.StatusError{Code: gohttp.StatusUnauthorized, Err: ErrInvalidLogin}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	token := generateSessionToken()
	err = db.CreateSession(signinReq.Username, token)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	gohttp.SetCookie(rw.Writer(), &gohttp.Cookie{
		Name:  "Token",
		Value: token,
	})
	_, _ = rw.WriteMessage("Logged in successfully!")

	return nil
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
	var signupReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: err}
	}

	if signupReq.Username == "" || signupReq.Password == "" {
		return http.StatusError{Code: gohttp.StatusBadRequest, Err: ErrEmptyFields}
	}

	exists, err := db.UserExists(signupReq.Username)
	if exists {
		return http.StatusError{Code: gohttp.StatusConflict, Err: ErrDuplicateUser}
	}
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	err = db.CreateUser(signupReq.Username, signupReq.Password)
	if err != nil {
		return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
	}

	rw.WriteHeader(gohttp.StatusCreated)
	_, _ = rw.WriteMessage("Signed up successfully!")

	return nil
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
func Authenticator(db storage.Service) func(next gohttp.Handler) gohttp.Handler {
	return func(next gohttp.Handler) gohttp.Handler {
		hfn := func(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) error {
			tokenCookie, err := r.Cookie("Token")
			if err != nil {
				return http.StatusError{Code: gohttp.StatusUnauthorized, Err: ErrInvalidToken}
			}

			session, err := db.FindSession(tokenCookie.Value)
			if session.Username == "" {
				return http.StatusError{Code: gohttp.StatusUnauthorized, Err: ErrInvalidToken}
			}
			if err != nil {
				return http.StatusError{Code: gohttp.StatusInternalServerError, Err: err}
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, session.Username)

			next.ServeHTTP(rw.Writer(), r.WithContext(ctx))
			return nil
		}
		return http.Handler{Storage: db, Function: hfn}
	}
}

// GetUser returns the username from a request context using UserKey as key.
func GetUser(r *gohttp.Request) string {
	switch v := r.Context().Value(UserKey).(type) {
	case string:
		return v
	default:
		return ""
	}
}
