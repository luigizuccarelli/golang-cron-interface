package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/microlib/simple"
	"github.com/robfig/cron"
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

	conn := NewClientConnection(logger)

	cr := cron.New()
	cr.AddFunc(os.Getenv("CRON"),
		func() {
			UpdateData(conn)
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

func cleanup(c *cron.Cron) {
	logger.Warn("Cleanup resources")
	logger.Info("Terminating")
	c.Stop()
}
