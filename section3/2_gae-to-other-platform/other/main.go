package main

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	subsc     = os.Getenv("PUBSUB_SUBSC")
)

func main() {
	ctx := context.Background()

	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	if err := client.Subscription(subsc).Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf("Data: %s\n", msg.Data)
		log.Printf("Attributes: %v\n", msg.Attributes)
		msg.Ack()
	}); err != nil {
		log.Printf("failed to receive message: %v\n", err)
	}
}
