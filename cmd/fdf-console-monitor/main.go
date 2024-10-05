package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/stefanjenkner/fdf-console-monitor/internal/fitnessmachine"
	"github.com/stefanjenkner/fdf-console-monitor/internal/serialmonitor"
)

func main() {
	name := flag.String("name", "FDF Rower", "Name of BLE device")
	port := flag.String("port", "/dev/ttyUSB0", "Serial port to use")
	flag.Parse()

	fitnessMachine := fitnessmachine.NewFitnessMachine(*name)
	serialMonitor := serialmonitor.NewSerialMonitor(*port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v\n", sig)
		serialMonitor.Stop()
	}()

	fitnessMachine.Start()
	serialMonitor.AddObserver(fitnessMachine)
	serialMonitor.Run()
	fitnessMachine.Stop()
	log.Println("Shutdown completed")
	os.Exit(0)
}
