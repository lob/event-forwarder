package sqs

import (
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/pkg/errors"
)

// The Options struct allows you to provide options to modify the behavior of a
// new SQS struct returned by the New() function.
type Options struct {
	QueueURL string
}

// SQS encapsulates the client used to send events to a SQS queue.
type SQS struct {
	queueURL string
	client   sqsiface.SQSAPI
}

// New returns an instance of SQS with a valid AWS SQS client.
func New(opts *Options) (*SQS, error) {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	sess, err := session.NewSession(aws.NewConfig().WithHTTPClient(c))
	if err != nil {
		return &SQS{}, errors.Wrap(err, "sqs")
	}

	client := sqs.New(sess)

	sqs := &SQS{
		client: client,
	}

	if opts != nil {
		sqs.queueURL = opts.QueueURL
	}

	return sqs, nil
}

// SendMessage takes in a msg and a map of tags and sends the message to the
// configured SQS queue.
func (s *SQS) SendMessage(msg string, tags map[string]string) error {
	attributes := make(map[string]*sqs.MessageAttributeValue)

	for k, v := range tags {
		attributes[k] = &sqs.MessageAttributeValue{
			StringValue: aws.String(v),
			DataType:    aws.String("String"),
		}
	}

	_, err := s.client.SendMessage(&sqs.SendMessageInput{
		MessageBody:       &msg,
		MessageAttributes: attributes,
		QueueUrl:          &s.queueURL,
	})

	return err
}
