package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

func UpdateData(c Client) error {

	list := strings.Split(os.Getenv("PAYLOAD"), ",")

	for x, _ := range list {
		// send the payload
		c.Info(fmt.Sprintf("Retrieving data for  %s", list[x]))
		req, _ := http.NewRequest("POST", os.Getenv("URL"), bytes.NewBuffer([]byte(list[x])))
		req.Header.Set("X-Api-Key", os.Getenv("APIKEY"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := c.Do(req)
		if err != nil {
			c.Error(fmt.Sprintf("Http request %v", err))
			continue
		}

		defer resp.Body.Close()
		body, e := ioutil.ReadAll(resp.Body)
		if e != nil {
			c.Error(fmt.Sprintf("Cron service updatePrices %v", e))
			continue
		}
		c.Debug(fmt.Sprintf("Response from server %s", string(body)))
		time.Sleep(2 * time.Second)
	}

	return nil
}
