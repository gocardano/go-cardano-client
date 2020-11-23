package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gocardano/go-cardano-client/shelley"
	log "github.com/sirupsen/logrus"
)

const (
	readTimeoutMs  = 5000
	writeTimeoutMs = 5000
)

var version string = "-"

func main() {

	socketFilename := flag.String("socket", "", "Socket filename")
	showVersion := flag.Bool("version", false, "Display version information")
	debug := flag.Bool("debug", false, "Enable debug level logging")
	trace := flag.Bool("trace", false, "Enable trace level logging")

	flag.Parse()

	if *showVersion {
		fmt.Printf("go-cardano-client version: %s\n", version)
		os.Exit(0)
	}

	if *socketFilename == "" {
		flag.Usage()
		os.Exit(1)
	}

	info, err := os.Stat(*socketFilename)
	if err != nil && os.IsNotExist(err) {
		log.Errorf("File [%s] does not exists", *socketFilename)
		os.Exit(1)
	} else if err != nil {
		log.WithError(err).Errorf("Unknown error with file [%s]", *socketFilename)
		os.Exit(1)
	} else if info.IsDir() {
		log.Errorf("[%s] is a directory, expecting a unix file socket to cardano-node", *socketFilename)
		os.Exit(1)
	}

	log.SetLevel(log.ErrorLevel)
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	if *trace {
		log.SetLevel(log.TraceLevel)
	}

	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(false)

	log.Infof("Starting application version: %s", version)

	client, err := shelley.NewClient(*socketFilename)
	if err != nil {
		log.WithError(err).Error("Error creating shelley client")
		os.Exit(1)
	}

	err = client.Handshake()
	if err != nil {
		log.WithError(err).Error("Error negotiating handshake protocol")
		os.Exit(1)
	}

	slotNumber, hash, blockNumber, err := client.QueryTip()
	if err != nil {
		log.WithError(err).Error("Error querying tip block header hash")
	}

	fmt.Println("SlotNumber  : ", slotNumber)
	fmt.Println("Hash        : ", fmt.Sprintf("%x", hash))
	fmt.Println("BlockNumber : ", blockNumber)

	// Disconnect
	err = client.Disconnect()
}
