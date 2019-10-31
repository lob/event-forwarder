package sqs

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	sqsiface.SQSAPI

	messageBody       string
	messageAttributes map[string]string
}

func (m mockClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	if *input.MessageBody != m.messageBody {
		return nil, fmt.Errorf("unexpected message body found; found %q; expected %q", *input.MessageBody, m.messageBody)
	}

	for k, v := range m.messageAttributes {
		if input.MessageAttributes[k] != nil {
			if *input.MessageAttributes[k].StringValue != v {
				return nil, fmt.Errorf("unexpected message attribute value found for key %q. found %q; expected %q", k, *input.MessageAttributes[k].StringValue, v)
			}

			if *input.MessageAttributes[k].DataType != "String" {
				return nil, fmt.Errorf("unexpected message attribute data type found for key %q. found %q; expected %q", k, *input.MessageAttributes[k].DataType, "String")
			}
		}
	}

	return nil, nil
}

func TestNew(t *testing.T) {
	opts := &Options{
		QueueURL: "test",
	}
	cfg, err := New(opts)
	assert.NoError(t, err)

	assert.NotNil(t, cfg, "returned assets shouldn't be nil")
	assert.Equal(t, "test", cfg.queueURL)
}

func TestSendMessage(t *testing.T) {
	sqs := &SQS{
		client: &mockClient{
			messageBody: `{"message": "hello world"}`,
			messageAttributes: map[string]string{
				"type":        "cloudwatch",
				"environment": "test",
			},
		},
	}

	err := sqs.SendMessage(
		`{"message": "hello world"}`,
		map[string]string{
			"type":        "cloudwatch",
			"environment": "test",
		},
	)
	assert.NoError(t, err)
}
