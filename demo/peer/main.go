package main

import (
	"github.com/juju/errors"
	clipboard "github.com/obgnail/share-clipboard"
	log "github.com/sirupsen/logrus"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) }

func fn() {
	log.Info("--- read config ---")
	config, err := clipboard.ReadConfig("../config.json")
	if err != nil {
		log.Errorf("read config err: %s", err)
	}
	addr := config.ServerAddr
	sendHK := config.SendClipboardHotKey
	loadHK := config.LoadClipboardHotKey
	if len(addr) == 0 || len(sendHK) == 0 || len(loadHK) == 0 {
		log.Errorf("read config err: sendHK/loadHK empty")
	}

	log.Infof("addr:\t %s", addr)
	log.Infof("send:\t %s", sendHK)
	log.Infof("load:\t %s", loadHK)

	log.Info("--- start ---")
	err = clipboard.ListenHotKey(sendHK, func() {
		clipboard.PeerSendClipboard(addr)
	})
	if err != nil {
		log.Error(errors.ErrorStack(err))
	}
	err = clipboard.ListenHotKey(loadHK, func() {
		clipboard.PeerLoadClipboard(addr)
	})
	if err != nil {
		log.Error(errors.ErrorStack(err))
	}

	forever := make(chan struct{}, 1)
	<-forever
}
