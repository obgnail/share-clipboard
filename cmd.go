package share_clipboard

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func PeerSendClipboard(addr string) {
	p, _ := NewPeer(addr)
	p.SendOnce()
	if err := p.Run(); err != nil {
		err = errors.Trace(err)
		log.Error(errors.ErrorStack(err))
	}
}

func PeerDefaultLoadClipboard(addr string) {
	p, _ := NewPeer(addr)
	p.LoadOnce(nil)
	if err := p.Run(); err != nil {
		err = errors.Trace(err)
		log.Error(errors.ErrorStack(err))
	}
}

func PeerLoadClipboard(addr string) {
	p, _ := NewPeer(addr)
	p.LoadOnce(func(data []byte) { log.Println("return::", string(data)) })
	if err := p.Run(); err != nil {
		err = errors.Trace(err)
		log.Error(errors.ErrorStack(err))
	}
}

func SyncClipboard(addr string) {
	p, _ := NewPeer(addr)
	p.SendLoop()
	p.LoadOnce(nil)
	if err := p.Run(); err != nil {
		err = errors.Trace(err)
		log.Error(errors.ErrorStack(err))
	}
}

func ServerRun(addr string) {
	if err := NewServer(addr).Serve(); err != nil {
		log.Error(errors.ErrorStack(err))
	}
}
