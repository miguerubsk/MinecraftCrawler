package protocol

import (
	"net"
	"testing"
	"time"
)

const (
	loopbackHost      = "127.0.0.1"
	loopbackDynamic   = "127.0.0.1:0"
	canonicalRCONAddr = "127.0.0.1:25575"
)

func startMockRCONServer(t *testing.T) (host string, port int, cleanup func()) {
	t.Helper()

	l, err := net.Listen("tcp", loopbackDynamic)
	if err != nil {
		t.Fatalf("failed to start mock RCON server: %v", err)
	}

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

	addr := l.Addr().(*net.TCPAddr)
	return loopbackHost, addr.Port, func() { _ = l.Close() }
}

func TestProbeRconSuccess(t *testing.T) {
	host, port, cleanup := startMockRCONServer(t)
	defer cleanup()

	if err := probeRcon(host, port, time.Second); err != nil {
		t.Fatalf("probeRcon() error = %v, want nil", err)
	}
}

func TestProbeRconFailure(t *testing.T) {
	l, err := net.Listen("tcp", loopbackDynamic)
	if err != nil {
		t.Fatalf("failed to reserve random port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	err = probeRcon(loopbackHost, port, 200*time.Millisecond)
	if err == nil {
		t.Fatalf("probeRcon() error = nil, want non-nil")
	}
}

func TestAnalyzeRconSuccess(t *testing.T) {
	host, port, cleanup := startMockRCONServer(t)
	defer cleanup()

	detail := &ServerDetail{IP: host, Port: port}
	got, err := analyzeRcon(detail, time.Second)
	if err != nil {
		t.Fatalf("analyzeRcon() error = %v", err)
	}
	if got == nil {
		t.Fatal("analyzeRcon() returned nil detail")
	}
	if !got.RconAttempted {
		t.Fatalf("RconAttempted = false, want true")
	}
	if !got.RconOpen {
		t.Fatalf("RconOpen = false, want true")
	}
	if got.Software != "RCON Service" {
		t.Fatalf("Software = %q, want %q", got.Software, "RCON Service")
	}
}

func TestAnalyzeRconFailureMarksAttempt(t *testing.T) {
	l, err := net.Listen("tcp", loopbackDynamic)
	if err != nil {
		t.Fatalf("failed to reserve random port: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	_ = l.Close()

	detail := &ServerDetail{IP: loopbackHost, Port: port}
	got, err := analyzeRcon(detail, 200*time.Millisecond)
	if err == nil {
		t.Fatalf("analyzeRcon() error = nil, want non-nil")
	}
	if got != nil {
		t.Fatalf("analyzeRcon() detail = %#v, want nil on error", got)
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