package oss

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	Ali "github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOSSConfig struct {
	Endpoint  string `json:"endpoint"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Bucket    string `json:"bucket"`
	BaseURL   string `json:"baseURL"`
}

type aliOSS struct {
	client    *Ali.Client
	bucket    *Ali.Bucket
	baseURL   string
	endpoint  string
	accessKey string
	secretKey string
}

func NewAliOSS(cfg *AliOSSConfig) (OSS, error) {

	cli, err := Ali.New(cfg.Endpoint, cfg.AccessKey, cfg.SecretKey)
	if err != nil {
		return nil, err
	}

	buk, err := cli.Bucket(cfg.Bucket)

	if err != nil {
		return nil, err
	}

	return &aliOSS{client: cli, bucket: buk, baseURL: cfg.BaseURL, endpoint: cfg.Endpoint, accessKey: cfg.AccessKey, secretKey: cfg.SecretKey}, nil
}

func (S *aliOSS) Get(key string) ([]byte, error) {
	rd, err := S.bucket.GetObject(key)
	if err != nil {
		return nil, err
	}
	defer rd.Close()
	return ioutil.ReadAll(rd)
}

func (S *aliOSS) GetURL(key string) string {
	return fmt.Sprintf("%s%s", S.baseURL, key)
}

func (S *aliOSS) GetSignURL(key string, expires time.Duration) (string, error) {
	u, err := S.bucket.SignURL(key, Ali.HTTPGet, int64(expires/time.Second))
	if err != nil {
		return "", err
	}
	if S.baseURL != "" {
		i := strings.Index(u, "://")
		if i > 0 {
			b := strings.Index(u[i+3:], "/")
			if b > 0 {
				u = fmt.Sprintf("%s%s", S.baseURL, u[i+b+4:])
			}
		}
	}
	return u, nil
}

func (S *aliOSS) Put(key string, data []byte, header map[string]string) error {

	options := []Ali.Option{}

	for key, value := range header {
		if key == "Content-Type" {
			options = append(options, Ali.ContentType(value))
		} else if key == "Content-Encoding" {
			options = append(options, Ali.ContentEncoding(value))
		} else if key == "Content-Disposition" {
			options = append(options, Ali.ContentDisposition(value))
		} else {
			options = append(options, Ali.Meta(key, value))
		}
	}

	err := S.bucket.PutObject(key, bytes.NewReader(data), options...)

	if err != nil {
		return err
	}

	return nil
}

func (S *aliOSS) PutSignURL(key string, expires time.Duration, header map[string]string) (string, error) {
	u, err := S.bucket.SignURL(key, Ali.HTTPPut, int64(expires/time.Second))
	if err != nil {
		return "", err
	}
	return u, nil
}

func (S *aliOSS) PostSignURL(key string, expires time.Duration, maxSize int64, header map[string]string) (string, map[string]string, error) {

	u, err := url.Parse(S.endpoint)

	if err != nil {
		return "", nil, err
	}

	data := map[string]string{}

	data["OSSAccessKeyId"] = S.accessKey
	data["key"] = key

	policyData := map[string]interface{}{}

	policyData["expiration"] = time.Now().Add(expires).Format("2006-01-02T15:04:05Z")

	conditions := []interface{}{map[string]interface{}{"bucket": S.bucket.BucketName}, []interface{}{"eq", "$key", key}, []interface{}{"content-length-range", 0, maxSize}}

	for k, v := range header {
		conditions = append(conditions, []interface{}{"eq", fmt.Sprintf("$%s", k), v})
		data[k] = v
	}

	policyData["conditions"] = conditions

	b, _ := json.Marshal(policyData)

	policy := string(b)

	// log.Println(policy)

	v := base64.StdEncoding.EncodeToString([]byte(policy))

	data["policy"] = v

	m := hmac.New(sha1.New, []byte(S.secretKey))

	m.Write([]byte(v))

	data["Signature"] = base64.StdEncoding.EncodeToString(m.Sum(nil))

	var s string = S.baseURL

	if s == "" {
		s = fmt.Sprintf("%s://%s.%s", u.Scheme, S.bucket.BucketName, u.Host)
	}

	// log.Println(s, data)

	return s, data, nil
}

func (S *aliOSS) Del(key string) error {
	err := S.bucket.DeleteObject(key)
	if err != nil {
		return err
	}
	return nil
}

func (S *aliOSS) Has(key string) (bool, error) {
	_, err := S.bucket.GetObjectACL(key)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (S *aliOSS) Recycle() {

}
