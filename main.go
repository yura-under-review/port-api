package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/yura-under-review/port-api/server"
)

const (
	DEBUG = "DEBUG"
	INFO  = "INFO"
	ERROR = "ERROR"
	WARN  = "WARN"
)

func main() {

	// TODO: get config from envs

	initLogger("DEBUG")

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	setupGracefulShutdown(cancel)

	srv := server.New(":8080", "./static-html/root.html")

	if err := srv.Run(ctx, &wg); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	wg.Wait()
}

func initLogger(lvl string) {

	switch strings.ToUpper(lvl) {
	case INFO:
		log.SetLevel(log.InfoLevel)

	case WARN:
		log.SetLevel(log.WarnLevel)

	case ERROR:
		log.SetLevel(log.ErrorLevel)

	default:
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.JSONFormatter{PrettyPrint: false})
	log.SetOutput(os.Stderr)
}

func setupGracefulShutdown(cancel context.CancelFunc) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Info("interrupt signal")
		cancel()
	}()
}
