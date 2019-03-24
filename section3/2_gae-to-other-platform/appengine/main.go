package main

import (
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	topic     = os.Getenv("PUBSUB_TOPIC")
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := pubsub.NewClient(ctx, projectID)
		defer client.Close()

		id, err := client.Topic(topic).Publish(ctx, &pubsub.Message{
			Data: []byte("Message"),
			Attributes: map[string]string{
				"foo": "bar",
			},
		}).Get(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Publish messages: %v\n", id)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
