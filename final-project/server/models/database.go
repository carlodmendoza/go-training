package models

import (
	"encoding/json"
	"final-project/server/utils"
	"fmt"
	"net/http"
	"sync"
)

type Database struct {
	Users       []User              `json:"users"`
	NextUserID  int                 `json:"nextUserID"`
	Categories  map[string][]string `json:"categories"`
	Mu          sync.Mutex
	CurrentUser User
}

func (db *Database) Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if creds.Username == "" || creds.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if user := db.authenticateUser(creds); user != nil {
				utils.SendMessageWithBody(w, true, "Logged in successfully!")
				db.CurrentUser = *user
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				utils.SendMessageWithBody(w, false, "Invalid username or password.")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if user.Name == "" || user.Credentials.Username == "" || user.Credentials.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
		} else {
			if tempUser := db.findUser(user.Credentials); tempUser != nil {
				w.WriteHeader(http.StatusConflict)
				utils.SendMessageWithBody(w, false, "Account already exists.")
			} else {
				db.Mu.Lock()
				db.NextUserID++
				user.UserID = db.NextUserID
				db.Users = append(db.Users, user)
				db.Mu.Unlock()

				w.WriteHeader(http.StatusCreated)
				utils.SendMessageWithBody(w, true, "Signed up successfully!")
			}
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

func (db *Database) authenticateUser(creds Credentials) *User {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, user := range db.Users {
		if creds.Username == user.Credentials.Username && creds.Password == user.Credentials.Password {
			return &user
		}
	}
	return nil
}

func (db *Database) findUser(creds Credentials) *User {
	db.Mu.Lock()
	defer db.Mu.Unlock()
	for _, user := range db.Users {
		if creds.Username == user.Credentials.Username {
			return &user
		}
	}
	return nil
}
