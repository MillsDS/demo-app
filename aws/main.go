package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	var (
		awsRegion       = "us-east-1"
		awsBucket       = "davidrotatebucket"
		awsBucketObject = "demo.html"
	)

	awsAccessKey, ok := os.LookupEnv("MY_AWS_ACCESS_KEY_ID")
	if !ok {
		log.Fatal("missing MY_AWS_ACCESS_KEY_ID")
	}

	awsSecretAccessKey, ok := os.LookupEnv("MY_AWS_SECRET_ACCESS_KEY")
	if !ok {
		log.Fatal("missing MY_AWS_SECRET_ACCESS_KEY")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     awsAccessKey,
				SecretAccessKey: awsSecretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		// set up s3 client
		client := s3.NewFromConfig(cfg)

		// read an object from s3 bucket
		object, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
			Bucket: &awsBucket,
			Key:    &awsBucketObject,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer object.Body.Close()

		// copy the object into the response
		if _, err := io.Copy(w, object.Body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal(err)
	}
}
