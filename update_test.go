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

	// we first load the json payload to simulate a call to middleware
	// for now just ignore failures.
	logger.Debug(fmt.Sprintf("data %s", data))
	httpclient := NewHttpTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: code,
			// Send response to be tested

			Body: ioutil.NopCloser(bytes.NewBufferString("{\"response\":\"something\"}")),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	conns := &Connection{Http: httpclient, Logger: logger, Type: con}
	return conns
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func TestAll(t *testing.T) {

	var err error
	log := &simple.Logger{Level: "trace"}

	// create anonymous struct
	tests := []struct {
		Name     string
		Payload  string
		Handler  string
		FileName string
		Want     bool
		ErrorMsg string
	}{
		{
			"[TEST] UpdateData should pass",
			"{\"affiliate\": \"SBR-01\"},{\"affiliate\":\"Test\"}",
			"UpdateData",
			"tests/payload-example.json",
			false,
			"Handler %s returned - got (%v) wanted (%v)",
		},
		{
			"[TEST] UpdateData should fail (forced error)",
			"{\"affiliate\": \"SBR-01\"},{\"affiliate\":\"Test\"}",
			"UpdateDataError",
			"tests/payload-example.json",
			true,
			"Handler %s returned - got (%v) wanted (%v)",
		},
	}

	for _, tt := range tests {
		log.Info(fmt.Sprintf("Executing test : %s \n", tt.Name))
		switch tt.Handler {
		case "UpdateData":
			logger.Debug(fmt.Sprintf("Payload %s", tt.Payload))
			conn := NewTestClient(tt.FileName, 200, "normal", log)
			err = UpdateData(conn)
		case "UpdateDataError":
			logger.Debug(fmt.Sprintf("Payload %s", tt.Payload))
			conn := NewTestClient(tt.FileName, 500, "error", log)
			err = UpdateData(conn)
		case "UpdateDataResponseError":
			logger.Debug(fmt.Sprintf("Payload %s", tt.Payload))
			conn := NewTestClient(tt.FileName, 500, "responseError", log)
			err = UpdateData(conn)

		}

		if !tt.Want {
			if err != nil {
				t.Errorf(fmt.Sprintf(tt.ErrorMsg, tt.Handler, err, nil))
			}
		} else {
			if err == nil {
				t.Errorf(fmt.Sprintf(tt.ErrorMsg, tt.Handler, "nil", "error"))
			}
		}
		fmt.Println("")
	}
}
