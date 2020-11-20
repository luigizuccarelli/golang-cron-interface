// +build real

package connectors

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microlib/simple"
)

type Connection struct {
	Logger *simple.Logger
}

func NewClientConnection(logger *simple.Logger) Client {
	return &Connection{Logger: logger}
}

func (c *Connection) Error(msg string, val ...interface{}) {
	c.Logger.Error(fmt.Sprintf(msg, val...))
}

func (c *Connection) Info(msg string, val ...interface{}) {
	c.Logger.Info(fmt.Sprintf(msg, val...))
}

func (c *Connection) Debug(msg string, val ...interface{}) {
	c.Logger.Debug(fmt.Sprintf(msg, val...))
}

func (c *Connection) Trace(msg string, val ...interface{}) {
	c.Logger.Trace(fmt.Sprintf(msg, val...))
}

func (c *Connection) PutObject(input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	arn := "arn:aws:iam::" + os.Getenv("AWS_ACCOUNT") + ":role/" + os.Getenv("AWS_USER") + "-sts"
	sess := session.Must(session.NewSession())
	c.Trace("Function ZipAndTransfer session %v ", sess)
	creds := stscreds.NewCredentials(sess, arn)
	c.Trace("Function ZipAndTransfer creds %v ", creds)
	region := os.Getenv("AWS_REGION")
	config := aws.NewConfig().WithCredentials(creds).WithRegion(region).WithMaxRetries(10)
	svc := s3.New(sess, config)
	return svc.PutObject(input)
}
