package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	credentials "cloud.google.com/go/iam/credentials/apiv1"
	credpb "google.golang.org/genproto/googleapis/iam/credentials/v1"
)

const (
	issure   = "sample_frontend_service"
	audience = "sample_micro_service"
)

var (
	projectID      = os.Getenv("GOOGLE_CLOUD_PROJECT")
	serviceAccount = fmt.Sprintf("%s@appspot.gserviceaccount.com", projectID)
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := newClient().Get(fmt.Sprintf("%s://%s/microservice/foo", scheme(), r.Host))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Response from microservice is: %s", body)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func scheme() string {
	isDev := os.Getenv("DEV_SERVER")
	if isDev != "" {
		return "http"
	}
	return "https"
}

func newClient() *http.Client {
	return &http.Client{
		Transport: newTransport(),
	}
}

type transport struct {
	rt http.RoundTripper
}

func newTransport() *transport {
	return &transport{
		rt: http.DefaultTransport,
	}
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	token, err := t.createAuthToken(req.Context())
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return t.rt.RoundTrip(req)
}

func (*transport) createAuthToken(ctx context.Context) (string, error) {
	client, err := credentials.NewIamCredentialsClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	payload, _ := json.Marshal(jwt.StandardClaims{
		Issuer:    issure,
		Audience:  audience,
		ExpiresAt: time.Now().Add(10 * time.Second).Unix(),
		IssuedAt:  time.Now().Unix(),
	})
	resp, err := client.SignJwt(ctx, &credpb.SignJwtRequest{
		Name:    fmt.Sprintf("projects/-/serviceAccounts/%s", serviceAccount),
		Payload: string(payload),
	})
	if err != nil {
		return "", err
	}
	return resp.GetSignedJwt(), nil
}
