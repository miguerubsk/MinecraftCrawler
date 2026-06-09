package protocol

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"
)

func TestReadNullTerminatedString(t *testing.T) {
	data := []byte("Hello\x00World\x00")
	reader := bytes.NewReader(data)

	s1, err := readNullTerminatedString(reader)
	if err != nil || s1 != "Hello" {
		t.Errorf("Expected 'Hello', got '%s' (err: %v)", s1, err)
	}

	s2, err := readNullTerminatedString(reader)
	if err != nil || s2 != "World" {
		t.Errorf("Expected 'World', got '%s' (err: %v)", s2, err)
	}

	s3, err := readNullTerminatedString(reader)
	if err != io.EOF {
		t.Errorf("Expected EOF, got '%s' (err: %v)", s3, err)
	}
}

func TestParseBasicStatResponse(t *testing.T) {
	// Mock response for Basic Stat
	// [1] type, [4] session ID
	payload := []byte{0x00, 0x01, 0x01, 0x01, 0x01}
	
	// Data
	payload = append(payload, []byte("A Minecraft Server\x00")...)
	payload = append(payload, []byte("SMP\x00")...)
	payload = append(payload, []byte("world\x00")...)
	payload = append(payload, []byte("5\x00")...)
	payload = append(payload, []byte("20\x00")...)
	
	port := make([]byte, 2)
	binary.LittleEndian.PutUint16(port, 25565)
	payload = append(payload, port...)
	
	payload = append(payload, []byte("localhost\x00")...)

	reader := bytes.NewReader(payload[5:])
	res := &QueryResult{
		RawKV:   make(map[string]string),
		Plugins: []string{},
	}

	var err error
	res.MOTD, err = readNullTerminatedString(reader)
	if err != nil { t.Fatal(err) }
	_, _ = readNullTerminatedString(reader)
	res.MapName, _ = readNullTerminatedString(reader)
	
	numPlayers, _ := readNullTerminatedString(reader)
	maxPlayers, _ := readNullTerminatedString(reader)
	if numPlayers != "5" || maxPlayers != "20" {
		t.Errorf("Expected players 5/20, got %s/%s", numPlayers, maxPlayers)
	}

	var hostPort uint16
	_ = binary.Read(reader, binary.LittleEndian, &hostPort)
	if hostPort != 25565 {
		t.Errorf("Expected port 25565, got %d", hostPort)
	}
	
	res.HostName, _ = readNullTerminatedString(reader)
	if res.HostName != "localhost" {
		t.Errorf("Expected hostname localhost, got %s", res.HostName)
	}
}
