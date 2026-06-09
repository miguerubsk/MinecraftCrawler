package protocol_test

import (
	"MinecraftCrawler/internal/protocol"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	queryTestLoopbackAddr   = "127.0.0.1:0"
	queryTestLoopbackHost   = "127.0.0.1"
	errListenUDPTemplate    = "could not listen on udp: %v"
)

func TestGetQueryInfoShortHandshakeResponse(t *testing.T) {
	pc, err := net.ListenPacket("udp", queryTestLoopbackAddr)
	if err != nil {
		t.Fatalf(errListenUDPTemplate, err)
	}
	defer pc.Close()

	addr := pc.LocalAddr().(*net.UDPAddr)

	go func() {
		buf := make([]byte, 2048)
		_, clientAddr, err := pc.ReadFrom(buf)
		if err != nil {
			return
		}
		// < 5 bytes forces "no query response"
		_, _ = pc.WriteTo([]byte{0x09, 0x00, 0x00, 0x00}, clientAddr)
	}()

	_, err = protocol.GetQueryInfo(queryTestLoopbackHost, addr.Port, time.Second)
	if err == nil || !strings.Contains(err.Error(), "no query response") {
		t.Fatalf("GetQueryInfo() error = %v, want containing %q", err, "no query response")
	}
}

func TestGetQueryInfoInvalidTokenResponse(t *testing.T) {
	pc, err := net.ListenPacket("udp", queryTestLoopbackAddr)
	if err != nil {
		t.Fatalf(errListenUDPTemplate, err)
	}
	defer pc.Close()

	addr := pc.LocalAddr().(*net.UDPAddr)

	go func() {
		buf := make([]byte, 2048)
		n, clientAddr, err := pc.ReadFrom(buf)
		if err != nil || n < 7 {
			return
		}

		resp := make([]byte, 0, 16)
		resp = append(resp, 0x09)
		resp = append(resp, buf[3:7]...)
		resp = append(resp, []byte("abc\x00")...)
		_, _ = pc.WriteTo(resp, clientAddr)
	}()

	_, err = protocol.GetQueryInfo(queryTestLoopbackHost, addr.Port, time.Second)
	if err == nil {
		t.Fatalf("GetQueryInfo() error = nil, want non-nil for invalid token")
	}
}

func TestGetQueryInfoShortStatResponse(t *testing.T) {
	pc, err := net.ListenPacket("udp", queryTestLoopbackAddr)
	if err != nil {
		t.Fatalf(errListenUDPTemplate, err)
	}
	defer pc.Close()

	addr := pc.LocalAddr().(*net.UDPAddr)

	go func() {
		buf := make([]byte, 2048)
		
		// 1. First attempt (Full Query)
		// Handshake request
		n, clientAddr, err := pc.ReadFrom(buf)
		if err != nil || n < 7 { return }

		handshakeResp := make([]byte, 0, 32)
		handshakeResp = append(handshakeResp, 0x09)
		handshakeResp = append(handshakeResp, buf[3:7]...)
		handshakeResp = append(handshakeResp, []byte(strconv.Itoa(12345))...)
		handshakeResp = append(handshakeResp, 0x00)
		_, _ = pc.WriteTo(handshakeResp, clientAddr)

		// Stat request (will fail)
		_, clientAddr, err = pc.ReadFrom(buf)
		if err != nil { return }
		_, _ = pc.WriteTo([]byte{0x00, 0x01, 0x02, 0x03, 0x04}, clientAddr)

		// 2. Second attempt (Basic Query fallback)
		// Handshake request
		_, clientAddr, err = pc.ReadFrom(buf)
		if err != nil { return }
		_, _ = pc.WriteTo(handshakeResp, clientAddr)

		// Stat request (will fail)
		_, clientAddr, err = pc.ReadFrom(buf)
		if err != nil { return }
		_, _ = pc.WriteTo([]byte{0x00, 0x01, 0x02, 0x03, 0x04}, clientAddr)
	}()

	_, err = protocol.GetQueryInfo(queryTestLoopbackHost, addr.Port, time.Second)
	if err == nil || !strings.Contains(err.Error(), "stat request failed") {
		t.Fatalf("GetQueryInfo() error = %v, want containing %q", err, "stat request failed")
	}
}
