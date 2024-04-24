package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// Code here
	cf := "config.yaml"

	// Load the configuration
	// Resolve the configuration
	p := NewProvider(cf)

	col, err := NewCollector(context.Background(), p)
	if err != nil {
		fmt.Printf("Error creating collector: %v\n", err)
		return
	}

	// context from signal
	ctx := context.Background()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		col.Shutdown()
	}()

	col.Run(ctx)
}
