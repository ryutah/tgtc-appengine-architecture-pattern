package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	adminpb "google.golang.org/genproto/googleapis/iam/admin/v1"

	jwt "github.com/dgrijalva/jwt-go"

	admin "cloud.google.com/go/iam/admin/apiv1"
)

var keys = map[string]interface{}{}

var (
	projectID      = os.Getenv("GOOGLE_CLOUD_PROJECT")
	serviceAccount = fmt.Sprintf("%s@appspot.gserviceaccount.com", projectID)
)

func main() {
	http.HandleFunc("/microservice/foo", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		admClient, err := admin.NewIamClient(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer admClient.Close()

		if err := authorizeRequest(ctx, admClient, r); err != nil {
			http.Error(w, "forbidden request", http.StatusForbidden)
			return
		}

		fmt.Fprintln(w, "success to call api!!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

func authorizeRequest(ctx context.Context, client *admin.IamClient, r *http.Request) error {
	// Authorization Bearerトークンの抽出
	authParts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
	if len(authParts) != 2 || strings.ToLower(strings.TrimSpace(authParts[0])) != "bearer" {
		return errors.New("invalid authorization header")
	}
	payload := authParts[1]

	// JWTの検証(github.com/dgrijalva/jwt-goを利用) JWTのパースをしつつ、トークンの有効期限や署名の検証を自動で行ってくれる
	_, err := jwt.Parse(payload, func(token *jwt.Token) (interface{}, error) {
		kid := token.Header["kid"].(string)
		if key, ok := keys[kid]; ok {
			return key, nil
		}

		// Cloud IAM APIを利用し公開鍵を取得する
		resp, err := client.GetServiceAccountKey(ctx, &adminpb.GetServiceAccountKeyRequest{
			Name:          fmt.Sprintf("projects/%s/serviceAccounts/%s/keys/%s", projectID, serviceAccount, kid),
			PublicKeyType: adminpb.ServiceAccountPublicKeyType_TYPE_RAW_PUBLIC_KEY,
		})
		if err != nil {
			return nil, err
		}
		block, _ := pem.Decode(resp.GetPublicKeyData())
		if block == nil {
			return nil, errors.New("failed to parse pem block")
		}
		key, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		keys[kid] = key
		return key, nil
	})
	return err
}
