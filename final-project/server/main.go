package main

import (
	"context"
	"final-project/server/data"
	"final-project/server/models"
	"final-project/server/redis"
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

	for i, trans := range data.Transactions {
		transMap := make(map[string]interface{})
		transMap["Amount"] = trans.Amount
		transMap["Date"] = trans.Date
		transMap["Notes"] = trans.Notes
		transMap["CategoryID"] = trans.CategoryID
		trKey := fmt.Sprintf("%v:%v:%v", "transactions", 1, i+1)
		redis.Client.HSet(ctx, trKey, transMap)
	}
	for i, cat := range data.Categories {
		catMap := make(map[string]interface{})
		catMap["Name"] = cat.Name
		catMap["Type"] = cat.Type
		catKey := fmt.Sprintf("%v:%v", "categories", i+1)
		redis.Client.HSet(ctx, catKey, catMap)
	}
}

// handler handles requests to the server depending on the
// request URL. It also authorizes a user to make
// requests given a request token.
func handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var transID int
		if r.URL.Path == "/signin" {
			models.Signin(w, r)
		} else if r.URL.Path == "/signup" {
			models.Signup(w, r)
		}
		// } else if r.URL.Path == "/transactions" {
		// 	uid := db.FindUidByToken(r)
		// 	if uid == -1 || uid == 0 {
		// 		w.WriteHeader(http.StatusUnauthorized)
		// 		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		// 	} else {
		// 		db.ProcessTransaction(w, r, uid)
		// 	}
		// } else if n, _ := fmt.Sscanf(r.URL.Path, "/transactions/%d", &transID); n == 1 {
		// 	uid := db.FindUidByToken(r)
		// 	if uid == -1 || uid == 0 {
		// 		w.WriteHeader(http.StatusUnauthorized)
		// 		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		// 	} else {
		// 		db.ProcessTransactionID(w, r, uid, transID)
		// 	}
		// } else if r.URL.Path == "/categories" {
		// 	uid := db.FindUidByToken(r)
		// 	if uid == -1 || uid == 0 {
		// 		w.WriteHeader(http.StatusUnauthorized)
		// 		utils.SendMessageWithBody(w, false, "Unauthorized login.")
		// 	} else {
		// 		db.ProcessCategories(w, r)
		// 	}
		// } else {
		// 	w.WriteHeader(http.StatusNotImplemented)
		// 	utils.SendMessage(w, "Invalid URL or request")
		// }
	}
}
