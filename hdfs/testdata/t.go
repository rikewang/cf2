package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
)

func main() {
	ak, sk := "TJ1EAD2CNV9ZSQURG29P", "Qs6NSUWSf5gRYZpjB573qk2e1Jphc3KSbaNnJbmR"
	endpoint := "http://9.37.2.110"
	bucket := "xixi"
	path := os.Args[1]
	cfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("default"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		LogLevel:         aws.LogLevel(aws.LogDebug),
	}

	sess := session.New(cfg)
	uploader := s3manager.NewUploader(sess)

	remote, err := os.Open(path)
	if err != nil {
		log.Fatal("error open file")
	}
	defer remote.Close()

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
		Body:   remote,
	})

	fmt.Println("success")
}
