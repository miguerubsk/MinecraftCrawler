package protocol_test

import (
	"MinecraftCrawler/internal/protocol"
	"bytes"
	"encoding/json"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestReadVarIntSafe(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    int
		wantErr bool
	}{
		{"SingleByte", []byte{0x05}, 5, false},
		{"TwoBytes", []byte{0xac, 0x02}, 300, false},
		{"Overflow", []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0, true},
		{"Incomplete", []byte{0x80}, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bytes.NewReader(tt.input)
			got, err := protocol.ReadVarIntSafe(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadVarIntSafe() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ReadVarIntSafe() = %v, want %v", got, tt.want)
			}
		})

	}
}

func mockMCServer(t *testing.T, response protocol.StatusResponse) (string, func()) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Read Handshake packet length
		_, err = protocol.ReadVarIntSafe(conn)
		if err != nil { return }
		
		// Read Handshake Packet ID (should be 0x00)
		_, _ = protocol.ReadVarIntSafe(conn)
		
		// Skip verify content for test simplicity, assuming client correct
		// Just consume rest of handshake + request packet
		// In reality we should parse it properly, but skipping bytes works for mocking
		buffer := make([]byte, 1024)
		_ = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		_, _ = conn.Read(buffer) // Consume status request 0x00

		// Prepare response
		jsonBytes, _ := json.Marshal(response)
		
		var payload bytes.Buffer
		_ = protocol.WriteVarInt(&payload, 0x00) // Packet ID: JSON Response
		_ = protocol.WriteVarInt(&payload, len(jsonBytes)) // String length
		payload.Write(jsonBytes)

		var packet bytes.Buffer
		_ = protocol.WriteVarInt(&packet, payload.Len())
		packet.Write(payload.Bytes())

		_, _ = conn.Write(packet.Bytes())
	}()

	return l.Addr().String(), func() { l.Close() }
}

func TestGetServerStatus(t *testing.T) {
	expected := protocol.StatusResponse{}
	expected.Version.Name = "1.20.4"
	expected.Version.Protocol = 765
	expected.Players.Max = 100
	expected.Players.Online = 5
	expected.Description = "A Minecraft Server"

	addr, cleanup := mockMCServer(t, expected)
	defer cleanup()

	host, portStr, _ := net.SplitHostPort(addr)
	port, _ := strconv.Atoi(portStr)

	status, err := protocol.GetServerStatus(host, port, 2*time.Second)
	if err != nil {
		t.Fatalf("GetServerStatus failed: %v", err)
	}

	if status.Version.Name != expected.Version.Name {
		t.Errorf("Version Name = %s, want %s", status.Version.Name, expected.Version.Name)
	}
	if status.Players.Online != expected.Players.Online {
		t.Errorf("Players Online = %d, want %d", status.Players.Online, expected.Players.Online)
	}
}

