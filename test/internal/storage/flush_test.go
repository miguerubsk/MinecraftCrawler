package storage_test

import (
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/storage"
	"testing"
	"time"
)

const (
	storageTestVersionLatest = "1.20.4"
	storageTestIPUpsert      = "127.0.0.20"
)

func TestFlushInsertsAndAssignsTimestampForZeroValue(t *testing.T) {
	db := newSingleConnTestDB(t)

	batch := []*protocol.ServerDetail{
		{
			IP:            "127.0.0.10",
			Port:          25565,
			VersionName:   storageTestVersionLatest,
			Protocol:      765,
			PlayersOnline: 4,
			PlayersMax:    20,
			Software:      "Paper",
			Mods:          map[string]string{"fabric-api": "0.100.0"},
			Plugins:       []string{"EssentialsX", "LuckPerms"},
			// Timestamp intentionally zero to validate fallback behavior.
		},
	}

	if err := storage.Flush(db, batch); err != nil {
		t.Fatalf("Flush() error = %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&count); err != nil {
		t.Fatalf("failed to count rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 row, got %d", count)
	}

	var ts string
	var software string
	if err := db.QueryRow("SELECT timestamp, software FROM servers WHERE ip = ? AND port = ?", "127.0.0.10", 25565).Scan(&ts, &software); err != nil {
		t.Fatalf("failed to query inserted row: %v", err)
	}
	if ts == "" {
		t.Fatalf("timestamp should not be empty for zero-value input timestamp")
	}
	if software != "Paper" {
		t.Fatalf("software = %q, want %q", software, "Paper")
	}
}

func TestFlushUpsertReplacesExistingServerByUniqueKey(t *testing.T) {
	db := newSingleConnTestDB(t)

	first := []*protocol.ServerDetail{
		{
			IP:            storageTestIPUpsert,
			Port:          25565,
			VersionName:   "1.20.1",
			Protocol:      763,
			PlayersOnline: 1,
			PlayersMax:    10,
			Software:      "Vanilla",
			Mods:          map[string]string{},
			Plugins:       []string{},
			Timestamp:     time.Now().Add(-time.Hour),
		},
	}

	second := []*protocol.ServerDetail{
		{
			IP:            storageTestIPUpsert,
			Port:          25565,
			VersionName:   storageTestVersionLatest,
			Protocol:      765,
			PlayersOnline: 7,
			PlayersMax:    50,
			Software:      "Paper",
			Mods:          map[string]string{"fabric-api": "0.100.1"},
			Plugins:       []string{"LuckPerms"},
			Timestamp:     time.Now(),
		},
	}

	if err := storage.Flush(db, first); err != nil {
		t.Fatalf("Flush(first) error = %v", err)
	}
	if err := storage.Flush(db, second); err != nil {
		t.Fatalf("Flush(second) error = %v", err)
	}

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM servers WHERE ip = ? AND port = ?", storageTestIPUpsert, 25565).Scan(&count); err != nil {
		t.Fatalf("failed to count upsert rows: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 upserted row, got %d", count)
	}

	var version string
	var protocolNum int
	var players int
	var software string
	if err := db.QueryRow(
		"SELECT version_name, protocol, players_online, software FROM servers WHERE ip = ? AND port = ?",
		storageTestIPUpsert, 25565,
	).Scan(&version, &protocolNum, &players, &software); err != nil {
		t.Fatalf("failed to query upserted row: %v", err)
	}

	if version != storageTestVersionLatest {
		t.Fatalf("version_name = %q, want %q", version, storageTestVersionLatest)
	}
	if protocolNum != 765 {
		t.Fatalf("protocol = %d, want %d", protocolNum, 765)
	}
	if players != 7 {
		t.Fatalf("players_online = %d, want %d", players, 7)
	}
	if software != "Paper" {
		t.Fatalf("software = %q, want %q", software, "Paper")
	}
}
