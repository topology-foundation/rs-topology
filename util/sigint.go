package util

import (
	"fmt"
	"os"
	"os/signal"
)

// ListenSigint catches SIGINT from terminal and sends error to error channel
func ListenSigint(ch chan error) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	<-signalChan

	ch <- fmt.Errorf("Received Ctrl+C")
}
