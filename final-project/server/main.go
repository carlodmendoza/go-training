package main

import (
	"context"
	"final-project/server/data"
	"final-project/server/models"
	"final-project/server/redis"
	"final-project/server/utils"
	"fmt"
	"log"
	"net/http"
)

/*
	Main program for running the server, handling requests,
	and initializing data in Redis.
	Author: Carlo Mendoza
*/

func main() {
	fmt.Println("Server running on port 8080")
	initRedis()
	if err := http.ListenAndServe("localhost:8080", handler()); err != nil {
		log.Fatalf("Error ListenAndServe(): %s", err.Error())
	}
}

// initRedis stores the initial data in Redis.
func initRedis() {
	ctx := context.Background()
	redis.Client.FlushDB(ctx)

	redis.Client.SAdd(ctx, "uids", 1)
	redis.Client.Set(ctx, "nextUserID", 1, 0)
	redis.Client.Set(ctx, "nextTransactionID", 10, 0)

	credsKey := fmt.Sprintf("%v:%v", "credentials", 1)
	redis.Client.HSet(ctx, credsKey, map[string]interface{}{"Username": "cmendoza", "Password": "123"})

	trUsrKey := fmt.Sprintf("%v:%v", "transactions", 1)
	for _, trans := range data.Transactions {
		redis.Client.SAdd(ctx, trUsrKey, trans.TransactionID)
		transMap := make(map[string]interface{})
		transMap["TransactionID"] = trans.TransactionID
		transMap["Amount"] = trans.Amount
		transMap["Date"] = trans.Date
		transMap["Notes"] = trans.Notes
		transMap["CategoryID"] = trans.CategoryID
		trKey := fmt.Sprintf("%v:%v", trUsrKey, trans.TransactionID)
		redis.Client.HSet(ctx, trKey, transMap)
	}
	for _, cat := range data.Categories {
		redis.Client.SAdd(ctx, "catids", cat.CategoryID)
		catMap := make(map[string]interface{})
		catMap["CategoryID"] = cat.CategoryID
		catMap["Name"] = cat.Name
		catMap["Type"] = cat.Type
		catKey := fmt.Sprintf("%v:%v", "categories", cat.CategoryID)
		redis.Client.HSet(ctx, catKey, catMap)
	}
}

// handler handles requests to the server depending on the
// request URL. It also authorizes a user to make
// requests given a request token.
func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var transID int
		if r.URL.Path == "/signin" {
			models.Signin(w, r)
		} else if r.URL.Path == "/signup" {
			models.Signup(w, r)
		} else if r.URL.Path == "/categories" {
			models.ProcessCategories(w, r)
		} else if r.URL.Path == "/transactions" {
			models.ProcessTransaction(w, r)
		} else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
			models.ProcessTransactionID(w, r, transID)
		} else {
			w.WriteHeader(http.StatusNotImplemented)
			utils.SendMessage(w, "Invalid URL or request")
		}
	}
}
