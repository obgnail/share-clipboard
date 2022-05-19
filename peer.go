package share_clipboard

import (
	"bytes"
	"fmt"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"net"
)

type Peer struct {
	raddr  *net.TCPAddr
	events chan *event
}

func NewPeer(raddr string) (*Peer, error) {
	rAddr, err := net.ResolveTCPAddr("tcp4", raddr)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Peer{raddr: rAddr, events: make(chan *event, 1)}, nil
}

func (p *Peer) Send() {
	p.events <- defaultCopyEvent
}

func (p *Peer) Load(callback func(data []byte)) {
	p.events <- pasteEvent(callback)
}

func (p *Peer) DefaultLoad() {
	p.events <- defaultPasteEvent
}

func (p *Peer) Run() error {
	buf := GetBuffer()
	defer PutBuffer(buf)

	limitChan := make(chan struct{}, 1)
	for event := range p.events {
		limitChan <- struct{}{}
		conn, err := net.DialTCP("tcp", nil, p.raddr)
		if err != nil {
			return errors.Trace(err)
		}
		log.Debugf("Dail: %s -> %s", conn.LocalAddr().String(), conn.RemoteAddr().String())
		msg, err := event.build()
		if err != nil {
			return errors.Trace(err)
		}
		if _, err := conn.Write(msg); err != nil {
			return errors.Trace(err)
		}
		n, err := conn.Read(buf)
		if err != nil {
			return errors.Trace(err)
		}
		Len := len(StatusOK)
		if n < Len || string(buf[:Len]) != StatusOK {
			return fmt.Errorf("status error")
		}
		if event.success != nil {
			if err := event.success(buf[Len:n]); err != nil {
				return errors.Trace(err)
			}
		}
		if event.extra != nil {
			event.extra(buf[Len:n])
		}
		<-limitChan
	}
	return nil
}

var (
	defaultCopyEvent  = copyEvent()
	defaultPasteEvent = pasteEvent(nil)
)

type event struct {
	typ     int
	build   func() ([]byte, error)
	success func(data []byte) error
	extra   func(data []byte)
}

func copyEvent() *event {
	return &event{
		typ: ActionTypeCopy,
		build: func() ([]byte, error) {
			data, err := GetBytesFromClip()
			if err != nil {
				return nil, errors.Trace(err)
			}
			w := bytes.Buffer{}
			w.WriteByte(ActionTypeCopy)
			w.Write(data)
			return w.Bytes(), nil
		},
		success: nil,
		extra:   nil,
	}
}

func pasteEvent(callback func(data []byte)) *event {
	return &event{
		typ: ActionTypePaste,
		build: func() ([]byte, error) {
			w := bytes.Buffer{}
			w.WriteByte(ActionTypePaste)
			return w.Bytes(), nil
		},
		success: func(b []byte) error {
			if err := SetClipBytes(b); err != nil {
				return errors.Trace(err)
			}
			return nil
		},
		extra: callback,
	}
}
