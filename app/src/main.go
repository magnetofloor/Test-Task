package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"net/http"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})

	router := mux.NewRouter()
	router.HandleFunc("/set_key/{key}/{value}", SetKeyHandler).Methods("POST")
	router.HandleFunc("/get_key/{key}", GetKeyHandler).Methods("GET")
	router.HandleFunc("/del_key/{key}", DelKeyHandler).Methods("DELETE")

	http.Handle("/", router)
	http.ListenAndServe(":8080", nil)
}

func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK Set"+key+" "+value)
}

func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		value = "Key not found"
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, value)
}

func DelKeyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK del")
}
