package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/microlib/simple"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

type Connection struct {
	Http   *http.Client
	Logger *simple.Logger
	Type   string
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

func (c *Connection) Do(req *http.Request) (*http.Response, error) {
	if c.Type == "error" {
		return nil, errors.New("Forced error")
	}
	return c.Http.Do(req)
}

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewHttpTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func NewTestClient(data string, code int, con string, logger *simple.Logger) Client {

	// we first load the json payload to simulate a call
	// for now just ignore failures.
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString(data)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &Connection{Http: httpclient, Logger: logger, Type: con}
	return conns
}

func TestAll(t *testing.T) {

	var err error
	log := &simple.Logger{Level: "trace"}

	t.Run("UpdateData : should pass", func(t *testing.T) {
		conn := NewTestClient("{\"name\":\"test\"}", 200, "normal", log)
		res := UpdateData(conn)
		if res != nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with error (%v) wanted (nil)", "UpdateData", err))
		}
	})

	t.Run("UpdateData : should fail (force error)", func(t *testing.T) {
		conn := NewTestClient("test", 500, "error", log)
		res := UpdateData(conn)
		if res == nil {
			t.Errorf(fmt.Sprintf("Handler %s returned with no error wanted (err)", "UpdateData"))
		}
	})

}
