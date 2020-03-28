// +build !test

package main

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/microlib/simple"
)

type Connection struct {
	Http   *http.Client
	Logger *simple.Logger
}

func NewClientConnection(logger *simple.Logger) Client {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}
	return &Connection{Http: httpClient, Logger: logger}
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
	return c.Http.Do(req)
}
