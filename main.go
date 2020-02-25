package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/microlib/simple"
	"gopkg.in/robfig/cron.v2"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	logger *simple.Logger
)

func main() {

	logger = &simple.Logger{Level: os.Getenv("LOG_LEVEL")}

	err := ValidateEnvars(logger)
	if err != nil {
		os.Exit(1)
	}

	cr := cron.New()
	cr.AddFunc(os.Getenv("CRON"),
		func() {
			updatePrices(logger)
		})
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	go func() {
		<-c
		cleanup(cr)
		os.Exit(1)
	}()

	cr.Start()

	for {
		logger.Debug(fmt.Sprintf("NOP sleeping for %s seconds\n", os.Getenv("SLEEP")))
		s, _ := strconv.Atoi(os.Getenv("SLEEP"))
		time.Sleep(time.Duration(s) * time.Second)
	}
}

func updatePrices(logger *simple.Logger) error {
	// set up http object
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	// send the payload from the scanner routine
	req, _ := http.NewRequest("POST", os.Getenv("URL"), bytes.NewBuffer([]byte(os.Getenv("PAYLOAD"))))
	req.Header.Set("X-Api-Key", os.Getenv("APIKEY"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Http request %v", err))
		return err
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		logger.Error(fmt.Sprintf("Cron service updatePrices %v", e))
		return e
	}
	logger.Debug(fmt.Sprintf("Response from server %s", string(body)))

	return nil
}

func cleanup(c *cron.Cron) {
	logger.Warn("Cleanup resources")
	logger.Info("Terminating")
	c.Stop()
}
