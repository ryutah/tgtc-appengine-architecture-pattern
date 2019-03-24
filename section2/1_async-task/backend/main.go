package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		// `X-AppEngine-TaskName`はCloud Tasks以外では設定不可な
		// リクエストヘッダのため送信元の検証として利用することができる。
		taskName := r.Header.Get("X-AppEngine-TaskName")
		if taskName == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		payload := make(map[string]interface{})
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			return
		}
		for k, v := range payload {
			log.Printf("%v: %v\n", k, v)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
