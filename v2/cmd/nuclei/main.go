package main

import (
	"github.com/SecuriWiser/nuclei/v2/internal/firebase"
	"github.com/SecuriWiser/nuclei/v2/internal/mongo"
	"github.com/SecuriWiser/nuclei/v2/internal/scanner"
	"github.com/projectdiscovery/gologger"
	"os"
	"os/signal"
)

func waitForShutdown(c chan os.Signal) {
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			gologger.Info().Msgf("CTRL+C pressed: Exiting\n")
			mongo.Disconnect()
			firebase.Disconnect()
			os.Exit(1)
		}
	}()
}

func main() {
	mongo.Connect()
	firebase.Connect()

	c := make(chan os.Signal, 1)
	defer close(c)
	waitForShutdown(c)

	scanner.WaitForScanning()
}
