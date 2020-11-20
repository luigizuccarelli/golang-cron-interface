// +build fake

package handlers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"gitea-cicd.apps.aws2-dev.ocp.14west.io/cicd/trackmate-cron-interface/pkg/connectors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type Connection struct {
	Logger    *simple.Logger
	Type      string
	S3Service *FakeS3
	Force     string
}

type FakeS3 struct {
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
	if c.Force == "error" {
		return &s3.PutObjectOutput{}, errors.New("forced putobject error")
	}
	id := "testing"
	res := &s3.PutObjectOutput{VersionId: &id}
	return res, nil
}

func NewTestClient(err string, logger *simple.Logger) connectors.Client {

	conns := &Connection{S3Service: &FakeS3{}, Logger: logger, Force: err}
	return conns
}

func TestAll(t *testing.T) {

	var err error
	log := &simple.Logger{Level: "trace"}

	// need to setup files
	input, err := ioutil.ReadFile("../../tests/analytics-2020-11-11.bak")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("../../tests/analytics-2020-11-11.json", input, 0755)
	if err != nil {
		panic(err)
	}

	t.Run("ZipAndTransfer : should fail (force error)", func(t *testing.T) {
		os.Setenv("BASE_DIR", "../../tests")
		os.Setenv("TESTING", "true")
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("AWS_BUCKET", "544-trackmate-dev")
		os.Setenv("USER", "544-trackmate-dev")
		os.Setenv("AWS_ACCOUNT", "544841062556")
		conn := NewTestClient("error", log)
		res := ZipAndTransfer(conn)
		if res == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error wanted (err)", "ZipAndTransfer"))
		}
	})

	t.Run("ZipAndTransfer : should pass", func(t *testing.T) {
		os.Setenv("BASE_DIR", "../../tests")
		os.Setenv("TESTING", "true")
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("AWS_BUCKET", "544-trackmate-dev")
		os.Setenv("AWS_ACCOUNT", "544841062556")
		os.Setenv("USER", "544-trackmate-dev")
		conn := NewTestClient("normal", log)
		res := ZipAndTransfer(conn)
		if res != nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with error (%v) wanted (nil)", "ZipAndTransfer", err))
		}
	})

	t.Run("ZipAndTransfer : should fail (file not found)", func(t *testing.T) {
		os.Setenv("BASE_DIR", "../../tests")
		os.Setenv("TESTING", "true")
		os.Setenv("AWS_REGION", "us-east-1")
		os.Setenv("AWS_BUCKET", "544-trackmate-dev")
		os.Setenv("USER", "544-trackmate-dev")
		os.Setenv("AWS_ACCOUNT", "544841062556")
		conn := NewTestClient("error", log)
		res := ZipAndTransfer(conn)
		if res == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error wanted (err)", "ZipAndTransfer"))
		}
	})

}
