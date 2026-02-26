package protocol_test

import (
	"MinecraftCrawler/internal/protocol"
	"bytes"
	"net"
	"testing"
	"time"
)

func TestGetQueryInfo_Mock(t *testing.T) {
	// Simple mock UDP server to test the handshake and STAT request logic
	pc, err := net.ListenPacket("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not listen on udp: %v", err)
	}
	defer pc.Close()

	addr := pc.LocalAddr().(*net.UDPAddr)

	go func() {
		buf := make([]byte, 2048)
		// 1. Handshake request
		_, clientAddr, err := pc.ReadFrom(buf)
		if err != nil {
			return
		}
		
		// Respond with session ID + token string
		resp := new(bytes.Buffer)
		resp.Write([]byte{0x09})              // Type: Handshake
		resp.Write(buf[3:7])                  // Session ID
		resp.WriteString("987654321")         // Token
		resp.WriteByte(0x00)
		_, _ = pc.WriteTo(resp.Bytes(), clientAddr)

		// 2. Stat request
		_, clientAddr, err = pc.ReadFrom(buf)
		if err != nil {
			return
		}

		// Respond with Stat data
		statResp := new(bytes.Buffer)
		statResp.Write([]byte{0x00})          // Type: Stat
		statResp.Write(buf[3:7])              // Session ID
		statResp.Write([]byte("padding00"))   // 11 bytes header total including type and session id

		// KV data
		statResp.WriteString("hostname")
		statResp.WriteByte(0x00)
		statResp.WriteString("A Minecraft Server")
		statResp.WriteByte(0x00)
		statResp.WriteString("server_mod")
		statResp.WriteByte(0x00)
		statResp.WriteString("Paper (MC: 1.20.1)")
		statResp.WriteByte(0x00)
		statResp.WriteString("map")
		statResp.WriteByte(0x00)
		statResp.WriteString("world")
		statResp.WriteByte(0x00)
		statResp.Write([]byte{0x00, 0x01})
		
		// Plugins section
		statResp.Write([]byte("player_"))
		statResp.Write([]byte{0x00, 0x00})    // Termination

		_, _ = pc.WriteTo(statResp.Bytes(), clientAddr)
	}()

	res, err := protocol.GetQueryInfo("127.0.0.1", addr.Port, 1*time.Second)
	if err != nil {
		t.Errorf("GetQueryInfo() error = %v", err)
		return
	}

	if res.Software != "Paper (MC: 1.20.1)" {
		t.Errorf("expected software Paper (MC: 1.20.1), got %s", res.Software)
	}
	if res.MapName != "world" {
		t.Errorf("expected map world, got %s", res.MapName)
	}
}
