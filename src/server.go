package src

import (
	"bytes"
	"fmt"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	laddr string
	body  []byte
}

func NewServer(laddr string) *Server {
	return &Server{laddr: laddr}
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.laddr)
	log.Infof("Serve: %s", lis.Addr().String())
	if err != nil {
		return errors.Trace(err)
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Error(errors.Trace(err))
			continue
		}
		log.Debugf("accept: %s -> %s", conn.LocalAddr().String(), conn.RemoteAddr().String())
		go func() {
			if err := s.handleConn(conn); err != nil {
				log.Error(errors.ErrorStack(err))
			}
		}()
	}
}

func (s *Server) handleConn(conn net.Conn) error {
	buf := GetBuffer()
	defer PutBuffer(buf)

	// [ActionType]Body
	header := make([]byte, 1)
	n, err := conn.Read(header)
	if err != nil {
		return errors.Trace(err)
	}
	if n != 1 {
		return fmt.Errorf("read type error: %s", string(header))
	}

	resp := bytes.NewBufferString(StatusOK)
	switch action := header[0]; action {
	case ActionTypeCopy:
		n, err := conn.Read(buf)
		if err != nil {
			return errors.Trace(err)
		}
		if n >= defaultBufSize {
			return fmt.Errorf("too large body")
		}
		s.body = buf[:n]
	case ActionTypePaste:
		resp.Write(s.body)
	default:
		return fmt.Errorf("error type: %v", action)
	}
	if _, err := conn.Write(resp.Bytes()); err != nil {
		return errors.Trace(err)
	}
	return nil
}
