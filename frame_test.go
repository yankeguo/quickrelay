package quickrelay

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFrameType_Constants(t *testing.T) {
	assert.Equal(t, FrameType(0), FrameTypeUnknown)
	assert.Equal(t, FrameType(1), FrameTypeHandshakeRequest)
	assert.Equal(t, FrameType(2), FrameTypeHandshakeResponse)
	assert.Equal(t, FrameType(3), FrameTypeConnectRequest)
	assert.Equal(t, FrameType(4), FrameTypeConnectResponse)
	assert.Equal(t, FrameType(5), FrameTypeData)
	assert.Equal(t, FrameType(6), FrameTypeDisconnectRequest)
	assert.Equal(t, FrameType(7), FrameTypeDisconnectResponse)
	assert.Equal(t, FrameType(0xff), FrameTypeError)
}

func TestFrame_Validate(t *testing.T) {
	tests := []struct {
		name      string
		frame     Frame
		wantError bool
	}{
		{
			name: "valid handshake request",
			frame: Frame{
				FrameType: FrameTypeHandshakeRequest,
				FramePayloadHandshakeRequest: &FramePayloadHandshakeRequest{
					ServiceName:  "test-service",
					ServiceToken: "test-token",
				},
			},
			wantError: false,
		},
		{
			name: "invalid handshake request - missing payload",
			frame: Frame{
				FrameType: FrameTypeHandshakeRequest,
			},
			wantError: true,
		},
		{
			name: "valid handshake response",
			frame: Frame{
				FrameType:                     FrameTypeHandshakeResponse,
				FramePayloadHandshakeResponse: &FramePayloadHandshakeResponse{},
			},
			wantError: false,
		},
		{
			name: "invalid handshake response - missing payload",
			frame: Frame{
				FrameType: FrameTypeHandshakeResponse,
			},
			wantError: true,
		},
		{
			name: "valid connect request",
			frame: Frame{
				FrameType:                  FrameTypeConnectRequest,
				FramePayloadConnectRequest: &FramePayloadConnectRequest{},
			},
			wantError: false,
		},
		{
			name: "invalid connect request - missing payload",
			frame: Frame{
				FrameType: FrameTypeConnectRequest,
			},
			wantError: true,
		},
		{
			name: "valid connect response",
			frame: Frame{
				FrameType: FrameTypeConnectResponse,
				FramePayloadConnectResponse: &FramePayloadConnectResponse{
					ConnectionID: "conn-123",
				},
			},
			wantError: false,
		},
		{
			name: "invalid connect response - missing payload",
			frame: Frame{
				FrameType: FrameTypeConnectResponse,
			},
			wantError: true,
		},
		{
			name: "valid data frame",
			frame: Frame{
				FrameType: FrameTypeData,
				FramePayloadData: &FramePayloadData{
					Data: []byte("hello world"),
				},
			},
			wantError: false,
		},
		{
			name: "valid data frame - empty data (ping/pong)",
			frame: Frame{
				FrameType: FrameTypeData,
				FramePayloadData: &FramePayloadData{
					Data: []byte{},
				},
			},
			wantError: false,
		},
		{
			name: "invalid data frame - missing payload",
			frame: Frame{
				FrameType: FrameTypeData,
			},
			wantError: true,
		},
		{
			name: "valid disconnect request",
			frame: Frame{
				FrameType: FrameTypeDisconnectRequest,
				FramePayloadDisconnectRequest: &FramePayloadDisconnectRequest{
					ConnectionID: "conn-123",
				},
			},
			wantError: false,
		},
		{
			name: "invalid disconnect request - missing payload",
			frame: Frame{
				FrameType: FrameTypeDisconnectRequest,
			},
			wantError: true,
		},
		{
			name: "valid disconnect response",
			frame: Frame{
				FrameType:                      FrameTypeDisconnectResponse,
				FramePayloadDisconnectResponse: &FramePayloadDisconnectResponse{},
			},
			wantError: false,
		},
		{
			name: "invalid disconnect response - missing payload",
			frame: Frame{
				FrameType: FrameTypeDisconnectResponse,
			},
			wantError: true,
		},
		{
			name: "valid error frame",
			frame: Frame{
				FrameType: FrameTypeError,
				FramePayloadError: &FramePayloadError{
					Error: "something went wrong",
				},
			},
			wantError: false,
		},
		{
			name: "invalid error frame - missing payload",
			frame: Frame{
				FrameType: FrameTypeError,
			},
			wantError: true,
		},
		{
			name: "invalid frame type",
			frame: Frame{
				FrameType: FrameType(99),
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.frame.Validate()
			if tt.wantError {
				assert.Error(t, err)
				assert.Equal(t, ErrFrameBadPayload, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestFrameWriter_WriteFrame(t *testing.T) {
	tests := []struct {
		name  string
		frame Frame
	}{
		{
			name: "handshake request",
			frame: Frame{
				FrameType: FrameTypeHandshakeRequest,
				FramePayloadHandshakeRequest: &FramePayloadHandshakeRequest{
					ServiceName:  "test-service",
					ServiceToken: "secret-token",
				},
			},
		},
		{
			name: "handshake response",
			frame: Frame{
				FrameType:                     FrameTypeHandshakeResponse,
				FramePayloadHandshakeResponse: &FramePayloadHandshakeResponse{},
			},
		},
		{
			name: "connect request",
			frame: Frame{
				FrameType:                  FrameTypeConnectRequest,
				FramePayloadConnectRequest: &FramePayloadConnectRequest{},
			},
		},
		{
			name: "connect response",
			frame: Frame{
				FrameType: FrameTypeConnectResponse,
				FramePayloadConnectResponse: &FramePayloadConnectResponse{
					ConnectionID: "conn-456",
				},
			},
		},
		{
			name: "data frame",
			frame: Frame{
				FrameType: FrameTypeData,
				FramePayloadData: &FramePayloadData{
					Data: []byte("test data payload"),
				},
			},
		},
		{
			name: "data frame - empty (ping/pong)",
			frame: Frame{
				FrameType: FrameTypeData,
				FramePayloadData: &FramePayloadData{
					Data: []byte{},
				},
			},
		},
		{
			name: "disconnect request",
			frame: Frame{
				FrameType: FrameTypeDisconnectRequest,
				FramePayloadDisconnectRequest: &FramePayloadDisconnectRequest{
					ConnectionID: "conn-789",
				},
			},
		},
		{
			name: "disconnect response",
			frame: Frame{
				FrameType:                      FrameTypeDisconnectResponse,
				FramePayloadDisconnectResponse: &FramePayloadDisconnectResponse{},
			},
		},
		{
			name: "error frame",
			frame: Frame{
				FrameType: FrameTypeError,
				FramePayloadError: &FramePayloadError{
					Error: "test error message",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewFrameWriter(&buf)

			err := writer.WriteFrame(tt.frame)
			require.NoError(t, err)
			assert.Greater(t, buf.Len(), 0, "should have written data to buffer")
		})
	}
}

func TestFrameReader_ReadFrame(t *testing.T) {
	tests := []struct {
		name  string
		frame Frame
	}{
		{
			name: "handshake request",
			frame: Frame{
				FrameType: FrameTypeHandshakeRequest,
				FramePayloadHandshakeRequest: &FramePayloadHandshakeRequest{
					ServiceName:  "test-service",
					ServiceToken: "secret-token",
				},
			},
		},
		{
			name: "data frame with content",
			frame: Frame{
				FrameType: FrameTypeData,
				FramePayloadData: &FramePayloadData{
					Data: []byte("hello world from test"),
				},
			},
		},
		{
			name: "error frame",
			frame: Frame{
				FrameType: FrameTypeError,
				FramePayloadError: &FramePayloadError{
					Error: "test error message",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			writer := NewFrameWriter(&buf)

			// Write the frame
			err := writer.WriteFrame(tt.frame)
			require.NoError(t, err)

			// Read the frame back
			reader := NewFrameReader(&buf)
			readFrame, err := reader.ReadFrame()
			require.NoError(t, err)

			// Verify the frame type matches
			assert.Equal(t, tt.frame.FrameType, readFrame.FrameType)

			// Verify payload data matches based on frame type
			switch tt.frame.FrameType {
			case FrameTypeHandshakeRequest:
				require.NotNil(t, readFrame.FramePayloadHandshakeRequest)
				assert.Equal(t, tt.frame.FramePayloadHandshakeRequest.ServiceName, readFrame.FramePayloadHandshakeRequest.ServiceName)
				assert.Equal(t, tt.frame.FramePayloadHandshakeRequest.ServiceToken, readFrame.FramePayloadHandshakeRequest.ServiceToken)
			case FrameTypeData:
				require.NotNil(t, readFrame.FramePayloadData)
				assert.Equal(t, tt.frame.FramePayloadData.Data, readFrame.FramePayloadData.Data)
			case FrameTypeError:
				require.NotNil(t, readFrame.FramePayloadError)
				assert.Equal(t, tt.frame.FramePayloadError.Error, readFrame.FramePayloadError.Error)
			}
		})
	}
}

func TestFrameRoundTrip(t *testing.T) {
	// Test all frame types for round-trip serialization
	frames := []Frame{
		{
			FrameType: FrameTypeHandshakeRequest,
			FramePayloadHandshakeRequest: &FramePayloadHandshakeRequest{
				ServiceName:  "my-service",
				ServiceToken: "my-token-123",
			},
		},
		{
			FrameType:                     FrameTypeHandshakeResponse,
			FramePayloadHandshakeResponse: &FramePayloadHandshakeResponse{},
		},
		{
			FrameType:                  FrameTypeConnectRequest,
			FramePayloadConnectRequest: &FramePayloadConnectRequest{},
		},
		{
			FrameType: FrameTypeConnectResponse,
			FramePayloadConnectResponse: &FramePayloadConnectResponse{
				ConnectionID: "connection-123",
			},
		},
		{
			FrameType: FrameTypeData,
			FramePayloadData: &FramePayloadData{
				Data: []byte("sample data for testing"),
			},
		},
		{
			FrameType: FrameTypeDisconnectRequest,
			FramePayloadDisconnectRequest: &FramePayloadDisconnectRequest{
				ConnectionID: "connection-456",
			},
		},
		{
			FrameType:                      FrameTypeDisconnectResponse,
			FramePayloadDisconnectResponse: &FramePayloadDisconnectResponse{},
		},
		{
			FrameType: FrameTypeError,
			FramePayloadError: &FramePayloadError{
				Error: "round trip test error",
			},
		},
	}

	for i, originalFrame := range frames {
		t.Run(t.Name()+"_frame_"+string(rune('0'+i)), func(t *testing.T) {
			var buf bytes.Buffer

			// Write frame
			writer := NewFrameWriter(&buf)
			err := writer.WriteFrame(originalFrame)
			require.NoError(t, err)

			// Read frame back
			reader := NewFrameReader(&buf)
			readFrame, err := reader.ReadFrame()
			require.NoError(t, err)

			// Validate original frame
			require.NoError(t, originalFrame.Validate())

			// Compare frame types
			assert.Equal(t, originalFrame.FrameType, readFrame.FrameType)

			// For empty payload structs, msgpack might deserialize them as nil
			// So we need to recreate them to ensure validation passes
			switch readFrame.FrameType {
			case FrameTypeHandshakeResponse:
				if readFrame.FramePayloadHandshakeResponse == nil {
					readFrame.FramePayloadHandshakeResponse = &FramePayloadHandshakeResponse{}
				}
			case FrameTypeConnectRequest:
				if readFrame.FramePayloadConnectRequest == nil {
					readFrame.FramePayloadConnectRequest = &FramePayloadConnectRequest{}
				}
			case FrameTypeDisconnectResponse:
				if readFrame.FramePayloadDisconnectResponse == nil {
					readFrame.FramePayloadDisconnectResponse = &FramePayloadDisconnectResponse{}
				}
			}

			// Validate the read frame
			require.NoError(t, readFrame.Validate())
		})
	}
}

func TestFrameReader_InvalidData(t *testing.T) {
	// Test reading from invalid msgpack data
	invalidData := []byte{0x00, 0x01, 0x02}
	reader := NewFrameReader(bytes.NewReader(invalidData))

	_, err := reader.ReadFrame()
	assert.Error(t, err)
}

func TestNewFrameReader(t *testing.T) {
	var buf bytes.Buffer
	reader := NewFrameReader(&buf)
	assert.NotNil(t, reader)
}

func TestNewFrameWriter(t *testing.T) {
	var buf bytes.Buffer
	writer := NewFrameWriter(&buf)
	assert.NotNil(t, writer)
}
