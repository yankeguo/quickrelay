package quickrelay

import (
	"errors"
	"fmt"
	"net"
)

var (
	ErrServiceUnknown = errors.New("service unknown")
)

type ServiceOptions struct {
	ServiceName  string
	ServiceToken string
	ServicePort  int
}

type ServerOptions struct {
	GetServiceOptions func(name string) (ServiceOptions, error)
}

type Server interface {
	Run(l net.Listener) error
}

type server struct {
	opts ServerOptions
}

func NewServer(opts ServerOptions) Server {
	return &server{opts: opts}
}

func (s *server) serve(conn net.Conn) {
	defer conn.Close()

	var f Frame
	var err error

	fr := NewFrameReader(conn)
	fw := NewFrameWriter(conn)

	if f, err = fr.ReadFrame(); err != nil {
		fw.WriteFrame(Frame{
			FrameType: FrameTypeError,
			FramePayloadError: &FramePayloadError{
				Error: fmt.Sprintf("invalid frame: %v", err),
			},
		})
		return
	}

	if f.FrameType != FrameTypeHandshakeRequest {
		fw.WriteFrame(Frame{
			FrameType: FrameTypeError,
			FramePayloadError: &FramePayloadError{
				Error: "invalid handshake request",
			},
		})
		return
	}

	var so ServiceOptions

	if so, err = s.opts.GetServiceOptions(f.ServiceName); err != nil {
		fw.WriteFrame(Frame{
			FrameType: FrameTypeError,
			FramePayloadError: &FramePayloadError{
				Error: fmt.Sprintf("failed to get service options: %v", err),
			},
		})
		return
	}

	if f.ServiceToken != so.ServiceToken {
		fw.WriteFrame(Frame{
			FrameType: FrameTypeError,
			FramePayloadError: &FramePayloadError{
				Error: "invalid service token",
			},
		})
		return
	}

}

func (s *server) Run(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.serve(conn)
	}
}
