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
func Signin(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) (int, error) {
	var signinReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signinReq)
	if err != nil {
		return gohttp.StatusBadRequest, err
	}

	if signinReq.Username == "" || signinReq.Password == "" {
		return gohttp.StatusBadRequest, ErrEmptyFields
	}

	ok, err := db.AuthenticateUser(signinReq.Username, signinReq.Password)
	if !ok {
		return gohttp.StatusUnauthorized, ErrInvalidLogin
	}
	if err != nil {
		return gohttp.StatusInternalServerError, err
	}

	token := generateSessionToken()
	err = db.CreateSession(signinReq.Username, token)
	if err != nil {
		return gohttp.StatusInternalServerError, err
	}

	gohttp.SetCookie(rw.Writer(), &gohttp.Cookie{
		Name:  "Token",
		Value: token,
	})
	_, _ = rw.WriteMessage("Logged in successfully!")

	return gohttp.StatusOK, nil
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) (int, error) {
	var signupReq AuthRequest

	err := json.NewDecoder(r.Body).Decode(&signupReq)
	if err != nil {
		return gohttp.StatusBadRequest, err
	}

	if signupReq.Username == "" || signupReq.Password == "" {
		return gohttp.StatusBadRequest, ErrEmptyFields
	}

	exists, err := db.UserExists(signupReq.Username)
	if exists {
		return gohttp.StatusConflict, ErrDuplicateUser
	}
	if err != nil {
		return gohttp.StatusInternalServerError, err
	}

	err = db.CreateUser(signupReq.Username, signupReq.Password)
	if err != nil {
		return gohttp.StatusInternalServerError, err
	}

	_, _ = rw.WriteMessage("Signed up successfully!")

	return gohttp.StatusCreated, nil
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
		hfn := func(db storage.Service, rw *http.ResponseWriter, r *gohttp.Request) (int, error) {
			tokenCookie, err := r.Cookie("Token")
			if err != nil {
				return gohttp.StatusUnauthorized, ErrInvalidToken
			}

			session, err := db.FindSession(tokenCookie.Value)
			if session.Username == "" {
				return gohttp.StatusUnauthorized, ErrInvalidToken
			}
			if err != nil {
				return gohttp.StatusInternalServerError, err
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, session.Username)

			next.ServeHTTP(rw.Writer(), r.WithContext(ctx))
			return gohttp.StatusContinue, nil
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
