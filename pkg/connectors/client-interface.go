package connectors

import (
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	PutObject(*s3.PutObjectInput) (*s3.PutObjectOutput, error)
}
