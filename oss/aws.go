package oss

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AwsOSSConfig struct {
	Region    string `json:"region"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	BaseURL   string `json:"baseURL"`
}

type awsOSS struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	baseURL       string
	bucket        string
}

func Unwrap(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			return err
		}
		err = e
	}
	return err
}

func NewAwsOSS(cfg *AwsOSSConfig) (OSS, error) {

	config, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(cfg.Region),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     cfg.AccessKey,
				SecretAccessKey: cfg.SecretKey,
			},
		}),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(config)

	return &awsOSS{client: client, presignClient: s3.NewPresignClient(client), bucket: cfg.Bucket, baseURL: cfg.BaseURL}, nil

}

func (S *awsOSS) Get(key string) ([]byte, error) {
	ctx := context.Background()
	rs, err := S.client.GetObject(ctx, &s3.GetObjectInput{Bucket: &S.bucket, Key: &key})
	if err != nil {
		_, ok := Unwrap(err).(*types.NoSuchKey)
		if ok {
			return nil, ErrNoSuchKey
		}
		return nil, err
	}
	defer rs.Body.Close()
	return ioutil.ReadAll(rs.Body)
}

func (S *awsOSS) GetURL(key string) string {
	return fmt.Sprintf("%s%s", S.baseURL, key)
}

func (S *awsOSS) GetSignURL(key string, expires time.Duration) (string, error) {
	ctx := context.Background()
	rs, err := S.presignClient.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: &S.bucket, Key: &key}, s3.WithPresignExpires(expires))
	if err != nil {
		_, ok := Unwrap(err).(*types.NoSuchKey)
		if ok {
			return "", ErrNoSuchKey
		}
		return "", err
	}
	return rs.URL, nil
}

func (S *awsOSS) Put(key string, data []byte, header map[string]string) error {

	ctx := context.Background()

	input := &s3.PutObjectInput{Bucket: &S.bucket, Key: &key, Body: bytes.NewReader(data)}

	for key, value := range header {
		if key == "Content-Type" {
			{
				s := value
				input.ContentType = &s
			}
		} else if key == "Content-Encoding" {
			{
				s := value
				input.ContentEncoding = &s
			}
		} else if key == "Content-Disposition" {
			{
				s := value
				input.ContentDisposition = &s
			}
		}
	}

	_, err := S.client.PutObject(ctx, input)

	if err != nil {
		return err
	}

	return nil
}

func (S *awsOSS) PutSignURL(key string, expires time.Duration, header map[string]string) (string, error) {

	ctx := context.Background()

	input := &s3.PutObjectInput{Bucket: &S.bucket, Key: &key}

	for key, value := range header {
		if key == "Content-Type" {
			{
				s := value
				input.ContentType = &s
			}
		} else if key == "Content-Encoding" {
			{
				s := value
				input.ContentEncoding = &s
			}
		} else if key == "Content-Disposition" {
			{
				s := value
				input.ContentDisposition = &s
			}
		}
	}

	rs, err := S.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(expires))
	if err != nil {
		return "", err
	}
	return rs.URL, nil
}

func (S *awsOSS) PostSignURL(key string, expires time.Duration, maxSize int64, header map[string]string) (string, map[string]string, error) {
	return "", nil, io.ErrUnexpectedEOF
}

func (S *awsOSS) Del(key string) error {
	ctx := context.Background()
	_, err := S.client.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &S.bucket, Key: &key})
	return err
}

func (S *awsOSS) Has(key string) (bool, error) {
	ctx := context.Background()
	_, err := S.client.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &S.bucket, Key: &key})
	if err != nil {
		_, ok := Unwrap(err).(*types.NoSuchKey)
		if ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (S *awsOSS) Recycle() {

}
