package protocol

import (
	"net"
	"testing"
	"time"
)

const (
	loopbackHost      = "127.0.0.1"
	canonicalRCONAddr = "127.0.0.1:25575"
)

func TestAnalyzeServerTimeout(t *testing.T) {
	// Use an unroutable IP for timeout test
	_, err := AnalyzeServer("192.0.2.1", 25565, 100*time.Millisecond)
	if err == nil {
		t.Fatal("AnalyzeServer() error = nil, want timeout error")
	}
}

func TestAnalyzeServerRconFail(t *testing.T) {
	// TCP port closed should still attempt Query if it's the default port
	detail, err := AnalyzeServer(loopbackHost, 25565, 1*time.Second)
	if err != nil {
		t.Fatalf("AnalyzeServer() error = %v", err)
	}
	if !detail.RconAttempted {
		t.Fatalf("RconAttempted = false, want true")
	}
	if detail.RconOpen {
		t.Fatalf("RconOpen = true, want false")
	}
}

func TestAnalyzeServerRconPortPath(t *testing.T) {
	// Try binding to canonical RCON port; skip test if unavailable.
	l, err := net.Listen("tcp", canonicalRCONAddr)
	if err != nil {
		t.Skip("port 25575 unavailable on this environment; skipping dedicated path test")
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		_ = conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		buf := make([]byte, 64)
		_, _ = conn.Read(buf)
	}()

	detail, err := AnalyzeServer(loopbackHost, 25575, time.Second)
	if err != nil {
		t.Fatalf("AnalyzeServer() error = %v", err)
	}
	if !detail.RconAttempted || !detail.RconOpen {
		t.Fatalf("RCON flags not set as expected: attempted=%v open=%v", detail.RconAttempted, detail.RconOpen)
	}
	if detail.QueryAttempted {
		t.Fatalf("QueryAttempted = true, want false for dedicated RCON path")
	}
	if detail.Port != 25575 {
		t.Fatalf("Port = %d, want 25575", detail.Port)
	}
}
