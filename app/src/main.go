package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"io"
	"net/http"
	"os"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

var (
	ctx context.Context
	rdb *redis.Client
)

func init() {
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
		TLSConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		},
	})
}

func main() {
	http.HandleFunc("/set_key", SetKeyHandler)
	http.HandleFunc("/get_key", GetKeyHandler)
	http.HandleFunc("/del_key", DelKeyHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Forbidden", http.StatusForbidden)
	})

	http.ListenAndServe(":8080", nil)
}

func SetKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method must be POST", http.StatusMethodNotAllowed)
		return
	}

	// read json data
	var keyValue KeyValue
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&keyValue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	key := keyValue.Key
	value := keyValue.Value

	if key == "" || value == "" {
		http.Error(w, "Both key and value must be provided", http.StatusBadRequest)
		return
	}

	// set data to db
	err := rdb.Set(ctx, key, value, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Operation SET is successful\n")
}

func GetKeyHandler(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}(r.Body)

	key := r.FormValue("key")

	if key == "" {
		http.Error(w, "Key must be provided\n", http.StatusBadRequest)
		return
	}

	value, err := rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		http.Error(w, "Cannot find key: "+key, http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate response
	response := map[string]string{
		"value": value,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DelKeyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method must be DELETE", http.StatusMethodNotAllowed)
		return
	}

	// read json data
	var keyValue KeyValue
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&keyValue); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if keyValue.Key == "" {
		http.Error(w, "Key must be provided", http.StatusBadRequest)
		return
	}

	// delete by key
	value, err := rdb.GetDel(ctx, keyValue.Key).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	fmt.Fprint(w, "DELETE "+value+"\n")
}
