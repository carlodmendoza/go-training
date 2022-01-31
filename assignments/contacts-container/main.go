package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	if err := http.ListenAndServe(":8080", db.handler()); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

// handler handles requests to the server
// depending on the request URL
func (db *Database) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		if r.URL.Path == "/contacts" {
			db.process(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/contacts/%d", &id); n == 1 {
			db.processID(id, w, r)
		} else {
			fmt.Fprintln(w, "Invalid URL or request")
		}
	}
}

// process handles the request to "/contacts"
func (db *Database) process(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(db.records); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	case "POST":
		var contact Contact
		if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
			http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
			return
		}
		if con, ok := db.isDuplicateContact(contact); ok {
			w.WriteHeader(http.StatusConflict)
			w.Header().Set("Content-Type", "application/json")
			respBody := formatResponseBody(false, "Contact already exists.", *con)
			fmt.Fprintln(w, respBody)
			return
		}
		db.mu.Lock()
		db.nextID++
		contact.ID = db.nextID
		db.records = append(db.records, contact)
		db.mu.Unlock()

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		respBody := formatResponseBody(true, "Contact added successfully!", contact)
		fmt.Fprintln(w, respBody)
	case "PUT":
		http.Error(w, "Error: method not allowed", http.StatusMethodNotAllowed)
	case "DELETE":
		http.Error(w, "Error: method not allowed", http.StatusMethodNotAllowed)
	}
}

// processID handles the request to "/contacts/id"
func (db *Database) processID(id int, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if con, _, ok := db.retrieveContactById(id); ok {
			if err := json.NewEncoder(w).Encode(&con); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Error: record not found", http.StatusNotFound)
		}
	case "POST":
		http.Error(w, "Error: method not allowed", http.StatusMethodNotAllowed)
	case "PUT":
		if con, index, ok := db.retrieveContactById(id); ok {
			var contact Contact
			contact.ID = con.ID
			if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
				http.Error(w, "Error: "+err.Error(), http.StatusBadRequest)
				return
			}
			db.mu.Lock()
			db.records[index] = contact
			db.mu.Unlock()

			w.Header().Set("Content-Type", "application/json")
			respBody := formatResponseBody(true, "Contact updated successfully!", contact)
			fmt.Fprintln(w, respBody)
		} else {
			http.Error(w, "Error: record not found", http.StatusNotFound)
		}
	case "DELETE":
		if _, _, ok := db.retrieveContactById(id); ok {
			db.deleteContactById(id)
			w.Header().Set("Content-Type", "application/json")
			respBody := fmt.Sprintf(
				"{\"success\": %t,\"message\": \"%s\"}",
				true, "Contact deleted successfully.")
			fmt.Fprintln(w, respBody)
		} else {
			http.Error(w, "Error: record not found", http.StatusNotFound)
		}
	}
}

// deleteContactById deletes a record from the
// database given the record ID
func (db *Database) deleteContactById(id int) {
	db.mu.Lock()
	for i, rec := range db.records {
		if rec.ID == id {
			db.records = append(db.records[:i], db.records[i+1:]...)
		}
	}
	db.mu.Unlock()
}

// retrieveContactById retrieves a record from the
// database given the record ID. It returns a pointer
// to the found record, its index, and a boolean (true
// if record was found; otherwise, false)
func (db *Database) retrieveContactById(id int) (*Contact, int, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()
	for i, rec := range db.records {
		if rec.ID == id {
			return &rec, i, true
		}
	}
	return nil, -1, false
}

// isDuplicateContact checks for duplicate records given
// a Contact. It returns the pointer to the found record
// and a boolean (true if record already exists; otherwise, false)
func (db *Database) isDuplicateContact(c Contact) (*Contact, bool) {
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

// formatResponseBody returns a formatted string given
// a boolean, message, and Contact converted to JSON format
func formatResponseBody(success bool, message string, data Contact) string {
	body, _ := json.Marshal(data)
	return fmt.Sprintf(
		"{\"success\": %t,"+
			"\"message\": \"%s\","+
			"\"data\": %s\n}",
		success, message, string(body))
}
