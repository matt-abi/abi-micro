package aws

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ability-sh/abi-lib/json"
	"github.com/ability-sh/abi-micro/micro"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfigdata"
	"gopkg.in/yaml.v2"
)

func getConfigObject(content []byte, contentType string) (interface{}, error) {

	var config interface{} = nil

	if strings.Contains(contentType, "json") {
		err := json.Unmarshal(content, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	if strings.Contains(contentType, "yaml") {
		err := yaml.Unmarshal(content, &config)
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return nil, fmt.Errorf("unknown content type %s", contentType)
}

func SetAppConfig(c context.Context, p micro.Payload) (interface{}, error) {

	sess, err := session.NewSession(aws.NewConfig().
		WithRegion(os.Getenv("AWS_DEFAULT_REGION")).
		WithCredentials(credentials.NewEnvCredentials()))

	if err != nil {
		return nil, err
	}

	// 创建一个新的AppConfigData客户端
	svc := appconfigdata.New(sess)

	svc.GetLatestConfiguration(&appconfigdata.GetLatestConfigurationInput{})

	e_app := os.Getenv("AWS_APPCONFIG_APP")
	e_env := os.Getenv("AWS_APPCONFIG_ENV")
	e_config := os.Getenv("AWS_APPCONFIG_CONFIG")

	input, err := svc.StartConfigurationSession(&appconfigdata.StartConfigurationSessionInput{
		ApplicationIdentifier:          &e_app,
		EnvironmentIdentifier:          &e_env,
		ConfigurationProfileIdentifier: &e_config,
	})

	if err != nil {
		return nil, err
	}

	rs, err := svc.GetLatestConfiguration(&appconfigdata.GetLatestConfigurationInput{ConfigurationToken: input.InitialConfigurationToken})

	if err != nil {
		return nil, err
	}

	v, err := getConfigObject(rs.Configuration, *rs.ContentType)

	if err != nil {
		return nil, err
	}

	err = p.SetConfig(v)

	if err != nil {
		return nil, err
	}

	configurationToken := rs.NextPollConfigurationToken
	nextPollIntervalInSeconds := *rs.NextPollIntervalInSeconds

	go func() {

		for {
			select {
			case <-c.Done():
				return
			default:

				time.Sleep(time.Second * time.Duration(nextPollIntervalInSeconds))

				rs, err := svc.GetLatestConfiguration(&appconfigdata.GetLatestConfigurationInput{ConfigurationToken: configurationToken})

				if err != nil {
					log.Println("svc.GetLatestConfiguration", err)
					continue
				}

				configurationToken = rs.NextPollConfigurationToken
				nextPollIntervalInSeconds = *rs.NextPollIntervalInSeconds

				if len(rs.Configuration) == 0 {
					continue
				}

				v, err := getConfigObject(rs.Configuration, *rs.ContentType)

				if err != nil {
					log.Println("getConfigObject", err, string(rs.Configuration))
					continue
				}

				err = p.SetConfig(v)

				if err != nil {
					log.Println("p.SetConfig", err)
					continue
				}
			}
		}

	}()

	return v, nil
}
