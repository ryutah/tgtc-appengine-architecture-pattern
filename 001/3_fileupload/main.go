package main

import (
	"encoding/json"
	"fmt"
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

		uploadURL, err := storage.SignedURL(bucket, fileName, &storage.SignedURLOptions{
			GoogleAccessID: serviceAccount,
			Method:         http.MethodPut,
			SignBytes: func(b []byte) ([]byte, error) {
				resp, err := client.SignBlob(ctx, &credpb.SignBlobRequest{
					Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount),
					Payload: b,
				})
				if err != nil {
					return nil, err
				}
				return resp.GetSignedBlob(), nil
			},
			Expires: time.Now().Add(15 * time.Minute),
		})
		if err != nil {
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
