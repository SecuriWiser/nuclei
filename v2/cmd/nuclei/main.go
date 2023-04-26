package main

import (
	"github.com/apex/log"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v2/internal/firebase"
	"github.com/projectdiscovery/nuclei/v2/internal/scanner"
	"os"
	"os/signal"
)

func waitForShutdown(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			firebase.Disconnect()
			os.Exit(1)
		}
	}()
}

func main() {
	firebase.Connect()

	c := make(chan os.Signal, 1)
	defer close(c)
	waitForShutdown(c)

	url := os.Getenv("URL")
	riskID := os.Getenv("RISK_ID")
	if url != "" && riskID != "" {
		scanner.StartScanning(url, riskID)
	} else {
		log.Info("Url or riskID is empty")
	}
}
