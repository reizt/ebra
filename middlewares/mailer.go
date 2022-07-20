package middlewares

import (
	"context"
	"fmt"
	"os"

	"github.com/reizt/ebra/conf"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type Mailer struct {
}
type SendMailInput struct {
	From    string
	To      string
	Subject string
	Body    string
}

func (m *Mailer) SendMail(params *SendMailInput) (string, error) {
	conf.LoadEnv()
	ctx := context.Background()
	region := os.Getenv("AWS_REGION")
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return "", err
	}
	client := sesv2.NewFromConfig(cfg)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: &params.From,
		Destination: &types.Destination{
			ToAddresses: []string{params.To},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data: &params.Body,
					},
				},
				Subject: &types.Content{
					Data: &params.Subject,
				},
			},
		},
	}
	res, err := client.SendEmail(ctx, input)
	fmt.Println(err)
	if err != nil {
		return "", err
	}
	fmt.Println(res.MessageId)
	return *res.MessageId, nil
}
