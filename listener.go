package quickrelay

import (
	"crypto/tls"
	"net"
)

type ListenerOptions struct {
	// HOST:PORT of the relay server
	ServerAddr string
	// ServerInsecure is true if the relay server is insecure (non-TLS)
	ServerInsecure bool
	// ServerTLSConfig is the TLS config to connect to the relay server, overrides default TLS config
	ServerTLSConfig *tls.Config
	// ServiceName is the name of the service to register with the relay server
	ServiceName string
	// ServiceToken is the token of the service to register with the relay server
	ServiceToken string
}

type listener struct{}

// NewListener creates a new virtual listener that can be used to accept connections from relay server.
func NewListener(opts *ListenerOptions) net.Listener {
	return &listener{}
}

func (l *listener) Accept() (net.Conn, error) {
	return nil, nil
}

func (l *listener) Close() error {
	return nil
}

func (l *listener) Addr() net.Addr {
	return &net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: 0,
	}
}
