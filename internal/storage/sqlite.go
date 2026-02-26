package storage

import (
	"MinecraftCrawler/internal/protocol"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

func NewDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	// Restauramos PRAGMA NORMAL para seguridad de datos y añadimos todos los campos
	// Se añade UNIQUE(ip, port) para evitar duplicados
	query := `
		PRAGMA journal_mode = WAL;
		PRAGMA synchronous = NORMAL;
		CREATE TABLE IF NOT EXISTS servers (
			ip TEXT,
			port INTEGER,
			version_name TEXT,
			protocol INTEGER,
			players_online INTEGER,
			players_max INTEGER,
			whitelist BOOLEAN,
			software TEXT,
			mods TEXT,
			plugins TEXT,
			secure_chat BOOLEAN,
			timestamp DATETIME,
			UNIQUE(ip, port)
		);`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}
	return db, nil
}

// Renombramos a StartSQLiteManager para evitar colisión con buffer.go
func StartSQLiteManager(db *sql.DB, resultChan <-chan *protocol.ServerDetail, batchSize int) {
	buffer := make([]*protocol.ServerDetail, 0, batchSize)

	for res := range resultChan {
		buffer = append(buffer, res)
		if len(buffer) >= batchSize {
			if err := Flush(db, buffer); err != nil {
				log.Printf("Error flushing batch: %v", err)
			}
			buffer = buffer[:0]
		}
	}
	if len(buffer) > 0 {
		_ = Flush(db, buffer)
	}
}

func Flush(db *sql.DB, batch []*protocol.ServerDetail) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Usamos INSERT OR REPLACE para actualizar datos de servidores ya conocidos
	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO servers (
			ip, port, version_name, protocol, players_online, players_max, 
			whitelist, software, mods, plugins, secure_chat, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, s := range batch {
		modsJSON, _ := json.Marshal(s.Mods)
		pluginsJSON, _ := json.Marshal(s.Plugins)
		
		ts := s.Timestamp
		if ts.IsZero() {
			ts = time.Now()
		}

		_, err := stmt.Exec(
			s.IP, s.Port, s.VersionName, s.Protocol, s.PlayersOnline, s.PlayersMax,
			s.IsWhitelist, s.Software, string(modsJSON), string(pluginsJSON), 
			s.EnforcesSecureChat, ts,
		)
		if err != nil {
			log.Printf("Error inserting server %s: %v", s.IP, err)
			continue
		}
	}
	return tx.Commit()
}