package sns

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	pkgConfig "github.com/juanmalvarez3/twit/pkg/config"
)

func Provide(ctx context.Context) (*sns.Client, error) {
	cfg, err := pkgConfig.New()
	credProvider := credentials.NewStaticCredentialsProvider(
		cfg.AWS.AccessKey,
		cfg.AWS.SecretKey,
		"",
	)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cfg.AWS.Endpoint,
			SigningRegion: cfg.AWS.Region,
		}, nil
	})

	awsCfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credProvider),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		return nil, err
	}
	return sns.NewFromConfig(awsCfg), nil
}
