package main

import (
	"log"
	"time"

	"colormind/cmd"

	"github.com/getsentry/sentry-go"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://be656b4e9f3645cd8b4f5f55787ef37e@o1387515.ingest.sentry.io/6708903",
		Debug:            true,
		Release:          "colomind-cli@0.0.0",
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(5 * time.Second)

	cmd.Execute()
}
