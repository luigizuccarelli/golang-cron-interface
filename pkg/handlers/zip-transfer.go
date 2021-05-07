package handlers

import (
	"compress/gzip"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"lmzsoftware.com/lzuccarelli/golang-cron-interface/pkg/connectors"
)

func ZipAndTransfer(c connectors.Client) error {

	year, month, day := time.Now().Date()
	m := int(month)
	prevDay := day - 1

	if os.Getenv("TESTING") != "" && os.Getenv("TESTING") == "true" {
		year = 2020
		m = 11
		prevDay = 11
	}

	d, err := ioutil.ReadFile(os.Getenv("BASE_DIR") + "/analytics-" + strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay) + ".json")
	c.Debug("Function ZipAndTransfer key %s ", "analytics-"+strconv.Itoa(year)+"-"+strconv.Itoa(m)+"-"+strconv.Itoa(prevDay))
	if err != nil {
		c.Error("Function ZipAndTransfer %v", err)
		return err
	}

	f, _ := os.Create(os.Getenv("BASE_DIR") + "/analytics-" + strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay) + ".gz")
	w := gzip.NewWriter(f)
	defer w.Close()
	_, err = w.Write(d)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Body:   aws.ReadSeekCloser(strings.NewReader(os.Getenv("BASE_DIR") + "/analytics-" + strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay) + ".gz")),
		Bucket: aws.String(os.Getenv("AWS_BUCKET")),
		Key:    aws.String(strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay)),
	}

	_, e := c.PutObject(input)
	if e != nil {
		c.Error("Function ZipAndTransfer %v", e)
		return e
	}

	e = os.Remove(os.Getenv("BASE_DIR") + "/analytics-" + strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay) + ".json")
	if e != nil {
		c.Error("Function ZipAndTransfer %v", e)
		return e
	}

	e = os.Remove(os.Getenv("BASE_DIR") + "/analytics-" + strconv.Itoa(year) + "-" + strconv.Itoa(m) + "-" + strconv.Itoa(prevDay) + ".gz")
	if e != nil {
		c.Error("Function ZipAndTransfer %v", e)
		return e
	}

	return nil
}
