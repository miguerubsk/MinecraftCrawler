package storage

import (
	"MinecraftCrawler/internal/protocol"
	"database/sql"
	"log"
	"time"
	_ "modernc.org/sqlite"
)

func NewDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	_, _ = db.Exec(`
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = OFF;
		PRAGMA temp_store = MEMORY;
		CREATE TABLE IF NOT EXISTS servers (
			ip TEXT,
			port INTEGER,
			version_name TEXT,
			protocol INTEGER,
			players_online INTEGER,
			players_max INTEGER,
			whitelist BOOLEAN,
			timestamp DATETIME
		);
	`)
	return db, nil
}

func StartManager(db *sql.DB, resultChan <-chan *protocol.ServerDetail, batchSize int) {
	buffer := make([]*protocol.ServerDetail, 0, batchSize)

	for res := range resultChan {
		buffer = append(buffer, res)
		if len(buffer) >= batchSize {
			flush(db, buffer)
			buffer = buffer[:0]
		}
	}
	if len(buffer) > 0 {
		flush(db, buffer)
	}
}

func flush(db *sql.DB, batch []*protocol.ServerDetail) {
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error Tx: %v", err)
		return
	}

	stmt, _ := tx.Prepare(`INSERT INTO servers (ip, port, version_name, protocol, players_online, players_max, whitelist, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	defer stmt.Close()

	for _, s := range batch {
		// Aseguramos que el timestamp sea actual si viene vacío
		ts := s.Timestamp
		if ts.IsZero() {
			ts = time.Now()
		}
		_, _ = stmt.Exec(s.IP, s.Port, s.VersionName, s.Protocol, s.PlayersOnline, s.PlayersMax, s.IsWhitelist, ts)
	}
	_ = tx.Commit()
}