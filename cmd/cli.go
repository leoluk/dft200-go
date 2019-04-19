package main

import (
	"github.com/muka/go-bluetooth/api"
	log "github.com/sirupsen/logrus"
	"os"
)

const tagAddress = "C1:55:0A:91:1B:9E"

func main() {
	manager, err := api.NewManager()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	err = manager.RefreshState()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	dev, err := api.GetDeviceByAddress(tagAddress)
	if err != nil {
		log.Fatal(err)
	}

	log.Infof("device (dev): %v", dev)


}
