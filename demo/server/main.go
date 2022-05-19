package main

import (
	"github.com/juju/errors"
	"github.com/obgnail/share-clipboard"
	log "github.com/sirupsen/logrus"
)

func main() {
	addr := "127.0.0.1:8899"
	if err := share_clipboard.NewServer(addr).Serve(); err != nil {
		log.Error(errors.ErrorStack(err))
	}
}
