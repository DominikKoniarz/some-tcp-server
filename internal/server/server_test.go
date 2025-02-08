package server

import (
	"testing"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  []byte
		expected ParsedRequest
	}{
		{
			name:     "valid request",
			request:  []byte("0001Idb1\x00user1"),
			expected: ParsedRequest{ProtocolVersion: "0001", MessageType: "I", Data: "db1\x00user1"},
		},
		{
			name:     "invalid protocol version",
			request:  []byte("0002Idb1\x00user1"),
			expected: ParsedRequest{},
		},
		{
			name:     "invalid data format",
			request:  []byte("0001I"),
			expected: ParsedRequest{},
		},
	}

	srv := NewServer()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			parsedRequest, err := srv.ParseRequest(tt.request)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if parsedRequest.ProtocolVersion != tt.expected.ProtocolVersion {
				t.Errorf("Expected protocol version %v, got %v", tt.expected.ProtocolVersion, parsedRequest.ProtocolVersion)
			}

			if parsedRequest.MessageType != tt.expected.MessageType {
				t.Errorf("Expected message type %v, got %v", tt.expected.MessageType, parsedRequest.MessageType)
			}

			if parsedRequest.Data != tt.expected.Data {
				t.Errorf("Expected data %v, got %v", tt.expected.Data, parsedRequest.Data)
			}
		})
	}
}
