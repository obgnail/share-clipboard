package main

import (
	"github.com/juju/errors"
	"github.com/obgnail/share-clipboard"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	addr := "127.0.0.1:8899"
	p, _ := share_clipboard.NewPeer(addr)

	go func() {
		for {
			p.Send()
			p.Load(func(data []byte) { log.Println("return::", string(data)) })
			time.Sleep(time.Second * 3)
		}
	}()

	if err := p.Run(); err != nil {
		err = errors.Trace(err)
		log.Error(errors.ErrorStack(err))
	}
}
