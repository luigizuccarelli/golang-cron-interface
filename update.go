package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func UpdateData(c Client) error {

	req, _ := http.NewRequest("POST", os.Getenv("URL"), nil)
	req.Header.Set("X-Api-Key", os.Getenv("APIKEY"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.Do(req)
	if err != nil {
		c.Error(fmt.Sprintf("Function UpdateData http request %v", err))
		return err
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		c.Error(fmt.Sprintf("Function UpdateData %v", e))
		return err
	}
	c.Debug(fmt.Sprintf("Function UpdateData response from server %s", string(body)))
	return nil
}
