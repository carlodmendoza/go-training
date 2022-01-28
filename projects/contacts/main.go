package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Contact struct {
	ID       int
	First    string
	Last     string
	Company  string
	Address  string
	Country  string
	Position string
}

type Database struct {
	nextID  int
	mu      sync.Mutex
	records []Contact
}

func main() {
	fmt.Println("Server running in port 8080")
	db := &Database{records: []Contact{}}
	http.ListenAndServe("localhost:8080", db.handler())
}

func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.URL.Path == "/contacts" {
			db.process(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/contacts/%d", &id); n == 1 {
			// db.processID(id, w, r)
		} else {
			fmt.Fprintln(w, "Invalid URL or request")
		}
	}
}

func (db *Database) process(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var contact Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		if con, ok := db.isDuplicateContact(&contact); ok {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			respBody := formatResponseBody(false, "Contact already exists.", con)
			fmt.Fprintf(w, respBody)
			return
		}
		db.mu.Lock()
		db.nextID++
		contact.ID = db.nextID
		db.records = append(db.records, contact)
		db.mu.Unlock()
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		respBody := formatResponseBody(true, "Contact added successfully!", &contact)
		fmt.Fprintf(w, respBody)
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db.records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "PUT":
		http.Error(w, "Error: method not allowed", 405)
	case "DELETE":
		http.Error(w, "Error: method not allowed", 405)
	}
}

func (db *Database) isDuplicateContact(c *Contact) (*Contact, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for _, rec := range db.records {
		if rec.First == c.First &&
			rec.Last == c.Last &&
			rec.Company == c.Company &&
			rec.Address == c.Address &&
			rec.Country == c.Country &&
			rec.Position == c.Position {
			return &rec, true
		}
	}
	return nil, false
}

func formatResponseBody(success bool, message string, data *Contact) string {
	body, _ := json.Marshal(data)
	return fmt.Sprintf(
		"{\"success\": %t,"+
			"\"message\": \"%s\","+
			" \"data\": %s\n}",
		success, message, string(body))
}
