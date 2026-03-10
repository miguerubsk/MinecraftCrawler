package storage

import (
	"MinecraftCrawler/internal/protocol"
	"database/sql"
	"encoding/json"
	"fmt"
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
			map_name TEXT,
			mods TEXT,
			plugins TEXT,
			secure_chat BOOLEAN,
			timestamp DATETIME,
			UNIQUE(ip, port)
		);`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	if err := ensureServersColumn(db, "map_name", "TEXT"); err != nil {
		return nil, err
	}
	return db, nil
}

func ensureServersColumn(db *sql.DB, columnName string, columnType string) error {
	rows, err := db.Query("PRAGMA table_info(servers)")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name string
		var ctype string
		var notNull int
		var defaultValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notNull, &defaultValue, &pk); err != nil {
			return err
		}
		if name == columnName {
			return nil
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	query := fmt.Sprintf("ALTER TABLE servers ADD COLUMN %s %s", columnName, columnType)
	_, err = db.Exec(query)
	return err
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
			whitelist, software, map_name, mods, plugins, secure_chat, timestamp
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
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
			s.IsWhitelist, s.Software, s.MapName, string(modsJSON), string(pluginsJSON), 
			s.EnforcesSecureChat, ts,
		)
		if err != nil {
			log.Printf("Error inserting server %s: %v", s.IP, err)
			continue
		}
	}
	return tx.Commit()
}