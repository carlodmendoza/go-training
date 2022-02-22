package models

import (
	"context"
	"encoding/json"
	"final-project/server/redis"
	"final-project/server/utils"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*
	Handler program contains all fields and methods
	that make it possible to process requests from
	a client.
	Author: Carlo Mendoza
*/

// Signin handles a sign in request by a client.
// Upon successful sign in, a generated token
// is given as a cookie to client for authorizing
// future requests.
func Signin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var credentials Credentials
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if credentials.Username == "" || credentials.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		userID, ok := authenticateUser(credentials)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			utils.SendMessageWithBody(w, false, "Invalid username or password.")
			return
		}
		token := createUserSession(userID)
		http.SetCookie(w, &token)
		utils.SendMessageWithBody(w, true, "Logged in successfully!")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// Signup handles a sign up request by a client.
// It checks if an account already exists.
func Signup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var credentials Credentials
		if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if credentials.Username == "" || credentials.Password == "" {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if userExists := findCredentialsByUsername(credentials.Username); userExists {
			w.WriteHeader(http.StatusConflict)
			utils.SendMessageWithBody(w, false, "Account already exists.")
			return
		}
		createNewUser(credentials)
		w.WriteHeader(http.StatusCreated)
		utils.SendMessageWithBody(w, true, "Signed up successfully!")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessTransaction handles a transaction/ request by a client.
// The client can either get all transactions, add new transaction,
// or delete all transactions.
func ProcessTransaction(w http.ResponseWriter, r *http.Request) {
	userID, ok := authenticateToken(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		return
	}

	switch r.Method {
	case "GET":
		transactions := findTransactionsByUid(userID)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(transactions); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
		}
	case "POST":
		var transaction Transaction
		if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if transaction.Amount == 0 || transaction.Date == "" || transaction.CategoryID == 0 {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, "Missing required fields.")
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if ok := validateNewTransaction(w, r, transaction); ok {
			createNewTransaction(transaction, userID, true)
			w.WriteHeader(http.StatusCreated)
			utils.SendMessageWithBody(w, true, "Transaction added successfully!")
		}
	case "DELETE":
		if ok := deleteAllTransactionsByUid(userID); !ok {
			w.WriteHeader(http.StatusNotFound)
			utils.SendMessageWithBody(w, false, "No transactions found.")
			return
		}
		utils.SendMessageWithBody(w, true, "All transactions deleted successfully.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessTransactionID handles a transaction/id request by a client
// given a transaction ID. The client can either get, update, or
// delete a transaction.
func ProcessTransactionID(w http.ResponseWriter, r *http.Request, transID int) {
	userID, ok := authenticateToken(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		return
	}
	transaction, ok := findTransactionByTid(userID, transID)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		utils.SendMessageWithBody(w, false, "Transaction not found.")
		return
	}

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&transaction); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
		}
	case "PUT":
		newTransaction := transaction
		if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusBadRequest)
			utils.SendMessage(w, "400 Bad Request")
			return
		}
		if ok := validateNewTransaction(w, r, newTransaction); ok {
			newTransaction.TransactionID = transaction.TransactionID
			createNewTransaction(newTransaction, userID, false)
			utils.SendMessageWithBody(w, true, "Transaction updated successfully!")
		}
	case "DELETE":
		deleteTransaction(userID, transaction.TransactionID)
		utils.SendMessageWithBody(w, true, "Transaction deleted successfully.")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// ProcessCategories handles a categories/ request by a client.
// The client can get all categories.
func ProcessCategories(w http.ResponseWriter, r *http.Request) {
	if _, ok := authenticateToken(r); !ok {
		w.WriteHeader(http.StatusUnauthorized)
		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		return
	}

	switch r.Method {
	case "GET":
		categories := returnCategories()
		if err := json.NewEncoder(w).Encode(categories); err != nil {
			fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			utils.SendMessageWithBody(w, false, "500 Internal Server Error")
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		utils.SendMessage(w, "405 Method not allowed")
	}
}

// authenticateToken checks if token from a request cookie is associated
// to an existing session. If yes, it returns the corresponding
// user ID and true boolean; else, it returns 0 and false.
// If no cookie is found, it returns -1 and false.
func authenticateToken(r *http.Request) (int, bool) {
	tokenCookie, err := r.Cookie("Token")
	if err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		return -1, false
	}

	uids, _ := redis.Client.SMembers(context.Background(), "uids").Result()
	for _, uid := range uids {
		seshKey := fmt.Sprintf("%v:%v", "sessions", uid)
		sesh, _ := redis.Client.HGetAll(context.Background(), seshKey).Result()
		if tokenCookie.Value == sesh["Token"] {
			userID, _ := strconv.Atoi(uid)
			return userID, true
		}
	}
	return 0, false
}

// authenticateUser returns the user ID and true if given
// Credentials is correct. If not, it returns 0 and false.
func authenticateUser(credentials Credentials) (int, bool) {
	uids, _ := redis.Client.SMembers(context.Background(), "uids").Result()
	for _, uid := range uids {
		credKey := fmt.Sprintf("%v:%v", "credentials", uid)
		creds, _ := redis.Client.HGetAll(context.Background(), credKey).Result()
		if credentials.Username == creds["Username"] && credentials.Password == creds["Password"] {
			userID, _ := strconv.Atoi(uid)
			return userID, true
		}
	}
	return 0, false
}

// createUserSession saves or updates session data in Redis and
// returns an http.Cookie given an ID of authenticated user.
func createUserSession(uid int) http.Cookie {
	seshKey := fmt.Sprintf("%v:%v", "sessions", uid)

	seshMap := make(map[string]interface{})
	token := utils.GenerateSessionToken()
	seshMap["Token"] = token
	seshMap["Timestamp"] = time.Now().Unix()
	redis.Client.HSet(context.Background(), seshKey, seshMap)

	cookie := http.Cookie{Name: "Token", Value: token}
	return cookie
}

// findCredentialsByUsername returns true if a given username
// already has an existing account. Otherwise, it returns false.
func findCredentialsByUsername(username string) bool {
	uids, _ := redis.Client.SMembers(context.Background(), "uids").Result()
	for _, uid := range uids {
		credKey := fmt.Sprintf("%v:%v", "credentials", uid)
		creds, _ := redis.Client.HGetAll(context.Background(), credKey).Result()
		if username == creds["Username"] {
			return true
		}
	}
	return false
}

// createNewUser creates a new user with user ID
// and Credentials, and stores it in Redis.
func createNewUser(credentials Credentials) {
	nextUid, _ := redis.Client.Incr(context.Background(), "nextUserID").Result()
	redis.Client.SAdd(context.Background(), "uids", nextUid)
	credsKey := fmt.Sprintf("%v:%v", "credentials", nextUid)
	redis.Client.HSet(context.Background(), credsKey, map[string]interface{}{"Username": credentials.Username, "Password": credentials.Password})
}

// findTransactionsByUid returns a list of Transaction given
// a user ID. If there are no existing user transactions,
// it returns an empty list.
func findTransactionsByUid(uid int) []Transaction {
	transactions := []Transaction{}
	trUsrKey := fmt.Sprintf("%v:%v", "transactions", uid)
	trids, _ := redis.Client.SMembers(context.Background(), trUsrKey).Result()
	for _, trid := range trids {
		trKey := fmt.Sprintf("%v:%v", trUsrKey, trid)
		trans, _ := redis.Client.HGetAll(context.Background(), trKey).Result()
		transID, _ := strconv.Atoi(trans["TransactionID"])
		transAmount, _ := strconv.ParseFloat(trans["Amount"], 64)
		transCategoryID, _ := strconv.Atoi(trans["CategoryID"])
		transaction := Transaction{TransactionID: transID, Amount: transAmount, Date: trans["Date"], Notes: trans["Notes"], CategoryID: transCategoryID}
		transactions = append(transactions, transaction)
	}
	return transactions
}

// createNewTransaction creates a new Transaction, associates it to
// a user ID, and stores it in Redis. If isNew is true, a new transaction
// is created; else, it is only updated.
func createNewTransaction(transaction Transaction, uid int, isNew bool) {
	var trid int64
	trUsrKey := fmt.Sprintf("%v:%v", "transactions", uid)
	transMap := make(map[string]interface{})
	if isNew {
		trid, _ = redis.Client.Incr(context.Background(), "nextTransactionID").Result()
		redis.Client.SAdd(context.Background(), trUsrKey, trid)
	} else {
		trid = int64(transaction.TransactionID)
	}
	transMap["TransactionID"] = trid
	transMap["Amount"] = transaction.Amount
	transMap["Date"] = transaction.Date
	transMap["Notes"] = transaction.Notes
	transMap["CategoryID"] = transaction.CategoryID
	trKey := fmt.Sprintf("%v:%v", trUsrKey, trid)
	redis.Client.HSet(context.Background(), trKey, transMap)
}

// findTransactionByTid returns a Transaction and true if given
// transaction ID exists in given user ID. Otherwise, it returns
// an empty Transaction and false.
func findTransactionByTid(uid, tid int) (Transaction, bool) {
	trUsrKey := fmt.Sprintf("%v:%v", "transactions", uid)
	if isMember, _ := redis.Client.SIsMember(context.Background(), trUsrKey, tid).Result(); isMember {
		trKey := fmt.Sprintf("%v:%v", trUsrKey, tid)
		trans, _ := redis.Client.HGetAll(context.Background(), trKey).Result()
		transID, _ := strconv.Atoi(trans["TransactionID"])
		transAmount, _ := strconv.ParseFloat(trans["Amount"], 64)
		transCategoryID, _ := strconv.Atoi(trans["CategoryID"])
		transaction := Transaction{TransactionID: transID, Amount: transAmount, Date: trans["Date"], Notes: trans["Notes"], CategoryID: transCategoryID}
		return transaction, true
	}
	return Transaction{}, false
}

// deleteAllTransactionsByUid deletes all transactions
// of a user from Redis given a user ID. If successful,
// it returns true; else, it returns false.
func deleteAllTransactionsByUid(uid int) bool {
	trUsrKey := fmt.Sprintf("%v:%v", "transactions", uid)
	trids, _ := redis.Client.SMembers(context.Background(), trUsrKey).Result()
	if len(trids) == 0 {
		return false
	}
	redis.Client.Del(context.Background(), trUsrKey)
	for _, trid := range trids {
		trKey := fmt.Sprintf("%v:%v", trUsrKey, trid)
		redis.Client.Del(context.Background(), trKey)
	}
	return true
}

// deleteTransaction deletes a transaction of a user
// from Redis given a user ID and transaction ID.
func deleteTransaction(uid, tid int) {
	trUsrKey := fmt.Sprintf("%v:%v", "transactions", uid)
	redis.Client.SRem(context.Background(), trUsrKey, tid)
	trKey := fmt.Sprintf("%v:%v", trUsrKey, tid)
	redis.Client.Del(context.Background(), trKey)
}

// returnCategories gets all category data from Redis
// and returns an array of Category.
func returnCategories() []Category {
	categories := []Category{}
	catids, _ := redis.Client.SMembers(context.Background(), "catids").Result()
	for _, catid := range catids {
		catKey := fmt.Sprintf("%v:%v", "categories", catid)
		cat, _ := redis.Client.HGetAll(context.Background(), catKey).Result()
		catID, _ := strconv.Atoi(cat["CategoryID"])
		category := Category{CategoryID: catID, Name: cat["Name"], Type: cat["Type"]}
		categories = append(categories, category)
	}
	return categories
}

// findCategoryByCid returns true if a given category ID
// exists. Otherwise, it returns false.
func findCategoryByCid(cid int) bool {
	isMember, _ := redis.Client.SIsMember(context.Background(), "catids", cid).Result()
	return isMember
}

// validateNewTransaction validates a POST or PUT transaction request.
// It sends a message to the client if it is a bad request.
func validateNewTransaction(w http.ResponseWriter, r *http.Request, trans Transaction) bool {
	if categoryExists := findCategoryByCid(trans.CategoryID); !categoryExists {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, "Category doesn't exist.")
		w.WriteHeader(http.StatusNotFound)
		utils.SendMessageWithBody(w, false, "Category doesn't exist.")
		return false
	}
	if _, err := time.Parse("01-02-2006", trans.Date); err != nil {
		fmt.Printf("Error in %s: %s\n", r.URL.Path, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		utils.SendMessageWithBody(w, false, "Invalid date format.")
		return false
	}
	return true
}
