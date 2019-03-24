package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	credpb "google.golang.org/genproto/googleapis/iam/credentials/v1"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	"cloud.google.com/go/storage"
)

var (
	projectID      = os.Getenv("GOOGLE_CLOUD_PROJECT")
	serviceAccount = fmt.Sprintf("%s@appspot.gserviceaccount.com", projectID)
	bucket         = os.Getenv("GCS_BUCKET")
)

func main() {
	http.HandleFunc("/uploadurl", func(w http.ResponseWriter, r *http.Request) {
		fileName := r.FormValue("name")
		if fileName == "" {
			http.Error(w, "name parameter is required", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		client, err := credentials.NewIamCredentialsClient(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Close()

		// 署名付きURLの生成
		uploadURL, err := storage.SignedURL(bucket, fileName, &storage.SignedURLOptions{
			GoogleAccessID: serviceAccount,
			Method:         http.MethodPut,
			ContentType:    "image/png",
			SignBytes: func(b []byte) ([]byte, error) {
				// Cloud IAM APIでサービスアカウントを利用して署名を生成する
				resp, err := client.SignBlob(ctx, &credpb.SignBlobRequest{
					Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount),
					Payload: b,
				})
				if err != nil {
					log.Printf("Failed to signe blob: %v", err)
					return nil, err
				}
				return resp.GetSignedBlob(), nil
			},
			Expires: time.Now().Add(15 * time.Minute),
		})
		if err != nil {
			log.Printf("Failed to signe url: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		body, _ := json.Marshal(map[string]string{
			"url": uploadURL,
		})
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
