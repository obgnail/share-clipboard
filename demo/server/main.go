package main

import (
	clipboard "github.com/obgnail/share-clipboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("--- read config ---")
	config, err := clipboard.ReadConfig("../config.json")
	if err != nil {
		log.Errorf("read config err: %s", err)
	}
	addr := config.ServerAddr
	if len(addr) == 0 {
		log.Errorf("read config err: addr empty")
	}
	log.Infof("addr:\t %s", addr)

	log.Info("--- start serve ---")
	clipboard.ServerRun(addr)
}
