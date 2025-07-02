package quickrelay

import (
	"errors"
	"io"

	"github.com/vmihailenco/msgpack/v5"
)

var (
	ErrFrameBadPayload = errors.New("frame payload does not match frame type")
)

type FrameType uint8

const (
	FrameTypeUnknown FrameType = iota

	FrameTypeHandshakeRequest  // client -> server, contains the service name and service token
	FrameTypeHandshakeResponse // server -> client

	FrameTypeConnectRequest  // server -> client, a new connection is about to be established
	FrameTypeConnectResponse // client -> server, contains the connection ID

	FrameTypeData // client <-> server, contains the data, empty data means ping/pong

	FrameTypeDisconnectRequest  // client <-> server, contains the connection ID
	FrameTypeDisconnectResponse // server <-> client

	FrameTypeError FrameType = 0xff
)

type FramePayloadHandshakeRequest struct {
	ServiceName  string `msgpack:"sn"`
	ServiceToken string `msgpack:"st"`
}

type FramePayloadHandshakeResponse struct {
}

type FramePayloadConnectRequest struct{}

type FramePayloadConnectResponse struct {
	ConnectionID string `msgpack:"cid"`
}

type FramePayloadData struct {
	Data []byte `msgpack:"d"`
}

type FramePayloadDisconnectRequest struct {
	ConnectionID string `msgpack:"cid"`
}

type FramePayloadDisconnectResponse struct{}

type FramePayloadError struct {
	Error string `msgpack:"e"`
}

type Frame struct {
	FrameType                       FrameType `msgpack:"ft"`
	*FramePayloadHandshakeRequest   `msgpack:"hr,omitempty"`
	*FramePayloadHandshakeResponse  `msgpack:"hs,omitempty"`
	*FramePayloadConnectRequest     `msgpack:"cr,omitempty"`
	*FramePayloadConnectResponse    `msgpack:"cs,omitempty"`
	*FramePayloadData               `msgpack:"data,omitempty"`
	*FramePayloadDisconnectRequest  `msgpack:"dr,omitempty"`
	*FramePayloadDisconnectResponse `msgpack:"ds,omitempty"`
	*FramePayloadError              `msgpack:"err,omitempty"`
}

func (f *Frame) Validate() error {
	switch f.FrameType {
	case FrameTypeHandshakeRequest:
		if f.FramePayloadHandshakeRequest == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeHandshakeResponse:
		if f.FramePayloadHandshakeResponse == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeConnectRequest:
		if f.FramePayloadConnectRequest == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeConnectResponse:
		if f.FramePayloadConnectResponse == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeData:
		if f.FramePayloadData == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeDisconnectRequest:
		if f.FramePayloadDisconnectRequest == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeDisconnectResponse:
		if f.FramePayloadDisconnectResponse == nil {
			return ErrFrameBadPayload
		}
	case FrameTypeError:
		if f.FramePayloadError == nil {
			return ErrFrameBadPayload
		}
	default:
		return ErrFrameBadPayload
	}
	return nil
}

type FrameReader interface {
	ReadFrame() (f Frame, err error)
}

type frameReader struct {
	r   io.Reader
	dec *msgpack.Decoder
}

func NewFrameReader(r io.Reader) FrameReader {
	return &frameReader{r: r, dec: msgpack.NewDecoder(r)}
}

func (r *frameReader) ReadFrame() (f Frame, err error) {
	err = r.dec.Decode(&f)
	return
}

type FrameWriter interface {
	WriteFrame(f Frame) error
}

type frameWriter struct {
	w   io.Writer
	enc *msgpack.Encoder
}

func NewFrameWriter(w io.Writer) FrameWriter {
	return &frameWriter{w: w, enc: msgpack.NewEncoder(w)}
}

func (w *frameWriter) WriteFrame(f Frame) error {
	return w.enc.Encode(&f)
}
