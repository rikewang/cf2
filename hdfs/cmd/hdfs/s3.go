package main

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/colinmarc/hdfs/v2"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type req struct {
	key   string
	path  string
	isDir bool
}

type resp struct {
	key string
	err error
}

func uploadFile(sess *session.Session, bucket string, q chan req, p chan resp, client *hdfs.Client, wg *sync.WaitGroup) {
	uploader := s3manager.NewUploader(sess)

	for r := range q {
		if r.key == "" {
			continue
		}
		var rst resp
		var err error
		rst.key = r.key
		if !r.isDir {
			remote, err := client.Open(r.path)
			if err != nil {
				rst.err = err
			}
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(r.key),
				Body:   remote,
			})
			remote.Close()
		} else {
			buffer := &aws.WriteAtBuffer{}
			remote := bytes.NewReader(buffer.Bytes())
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(r.key),
				Body:   remote,
			})
		}

		rst.err = err
		p <- rst
	}
	wg.Done()
}

func startClients(workerNums uint, q chan req, p chan resp, client *hdfs.Client, ak, sk, bucket, endpoint string) {
	cfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(ak, sk, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String("default"),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
		//LogLevel:         aws.LogLevel(aws.LogDebug),
	}
	sess := session.New(cfg)
	wg := &sync.WaitGroup{}
	var i uint
	for i = 0; i < workerNums; i++ {
		wg.Add(1)
		go uploadFile(sess, bucket, q, p, client, wg)
	}
	wg.Wait()
	close(p)
}

func s3put(paths []string, ak, sk, bucket, s3url string, workerNums uint, exitQuick bool) {
	sources, nn, err := normalizePaths(paths[0:1])
	if err != nil {
		fatal(err)
	}

	source := sources[0]
	u, err := url.Parse(bucket)
	if err != nil || u.Scheme != "s3" {
		fatal("bucket url error")
	}

	s3Bucket := u.Host
	s3Path := strings.TrimPrefix(u.Path, "/")

	client, err := getClient(nn)
	if err != nil {
		fatal(err)
	}

	reqQueue := make(chan req, 1000)
	resqQueue := make(chan resp, 1000)

	go func() {
		err = client.Walk(source, func(p string, fi os.FileInfo, err error) error {
			if err != nil {
				fatal(err)
			}
			var key string
			if p == source && !fi.IsDir() {
				key = filepath.Join(s3Path, filepath.Base(source))
			} else {
				key = filepath.Join(s3Path, strings.TrimPrefix(p, source))
			}

			if key == "" {
				return nil
			}
			if key[0] == '/' {
				key = key[1:]
			}
			r := req{
				key:  key,
				path: p,
			}

			if fi.IsDir() {
				r.isDir = true
			}
			reqQueue <- r
			return nil

		})
		if err != nil {
			fatal(err)
		}
		close(reqQueue)
	}()

	go startClients(workerNums, reqQueue, resqQueue, client, *s3ak, *s3sk, s3Bucket, *s3endpoint)

	for rst := range resqQueue {
		if rst.err != nil {
			fmt.Fprintf(os.Stderr, "key: %s, error: %s\n", rst.key, rst.err)
			if exitQuick {
				os.Exit(1)
			}
		}
	}
}
