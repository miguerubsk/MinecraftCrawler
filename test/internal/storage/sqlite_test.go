package storage_test

import (
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/storage"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func TestNewDatabase_InMemory(t *testing.T) {
	db, err := storage.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	defer db.Close()

	// Verify table exists
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='servers'").Scan(&name)
	if err != nil {
		t.Errorf("failed to find servers table: %v", err)
	}
	if name != "servers" {
		t.Errorf("expected table name 'servers', got '%s'", name)
	}
}

func TestStartSQLiteManager_InMemory(t *testing.T) {
	// 1. Setup in-memory DB
	db, err := storage.NewDatabase(":memory:")
	if err != nil {
		t.Fatalf("failed to create in-memory database: %v", err)
	}
	defer db.Close()

	// 2. Setup channel and Start Manager
	resultChan := make(chan *protocol.ServerDetail, 5)
	done := make(chan bool)

	// Run in background
	go func() {
		storage.StartSQLiteManager(db, resultChan, 2)
		done <- true
	}()

	// 3. Send test data
	testServer := &protocol.ServerDetail{
		IP:            "127.0.0.1",
		Port:          25565,
		VersionName:   "1.20.4",
		Protocol:      765,
		PlayersOnline: 10,
		PlayersMax:    100,
		Software:      "Paper",
		Timestamp:     time.Now(),
		Mods:          make(map[string]string),
	}
	resultChan <- testServer

	// Close channel to flush buffer
	close(resultChan)

	// Wait for manager to finish
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("StartSQLiteManager timed out")
	}

	// 4. Verify data in DB
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count rows: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 row, got %d", count)
	}

	var storedIP string
	var storedPort int
	err = db.QueryRow("SELECT ip, port FROM servers LIMIT 1").Scan(&storedIP, &storedPort)
	if err != nil {
		t.Fatalf("Failed to query row: %v", err)
	}

	if storedIP != testServer.IP {
		t.Errorf("Expected IP %s, got %s", testServer.IP, storedIP)
	}
	if storedPort != testServer.Port {
		t.Errorf("Expected Port %d, got %d", testServer.Port, storedPort)
	}
}

