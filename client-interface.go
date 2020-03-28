package main

import (
	"net/http"
)

type Client interface {
	Error(string, ...interface{})
	Info(string, ...interface{})
	Debug(string, ...interface{})
	Trace(string, ...interface{})
	Do(req *http.Request) (*http.Response, error)
}
