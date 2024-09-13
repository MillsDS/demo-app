package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/storage"

	"google.golang.org/api/option"

	"golang.org/x/oauth2/google"
)

func main() {
	ctx := context.Background()

	// fetch credentials from environment & decode them
	credentialsJSON, ok := os.LookupEnv("MY_GCP_CREDENTIALS")
	if !ok {
		log.Fatal("missing MY_GCP_CREDENTIALS")
	}

	credentials, err := google.CredentialsFromJSON(ctx, []byte(credentialsJSON), secretmanager.DefaultAuthScopes()...)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// set up a GCS (Google Cloud Storage) client
		client, err := storage.NewClient(ctx, option.WithCredentials(credentials))
		if err != nil {
			log.Fatal(err)
		}

		// read from my-cool-test-bucket & copy it to response
		bucket := client.Bucket("my-cool-test-bucket")
		object := bucket.Object(objectLabel())

		b, err := object.NewReader(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer b.Close()

		if _, err := io.Copy(w, b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal(err)
	}
}

func objectLabel() string {
	if inKubernetes() {
		return "data-k8s"
	}
	return "data"
}

func inKubernetes() bool {
	_, ok := os.LookupEnv("KUBERNETES_SERVICE_HOST")
	return ok
}
