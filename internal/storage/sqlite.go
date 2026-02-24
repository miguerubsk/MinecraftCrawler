package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"MinecraftCrawler/internal/protocol"
	"sync"
)

type Database struct {
	db *sql.DB
	mu sync.Mutex
}

func NewDatabase(path string) (*Database, error) {
	// Optimizaciones de PRAGMA para velocidad extrema
	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=-20000", path)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	// Crear tabla
	query := `
	CREATE TABLE IF NOT EXISTS servers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ip TEXT,
		port INTEGER,
		version_name TEXT,
		protocol_ver INTEGER,
		motd TEXT,
		players_online INTEGER,
		players_max INTEGER,
		icon BLOB,
		is_whitelist BOOLEAN,
		secure_chat BOOLEAN,
		software TEXT,
		mods_json TEXT,    -- Guardaremos los mods como JSON string
		plugins_json TEXT, -- Plugins como JSON string
		timestamp DATETIME,
		UNIQUE(ip, port)
	);`
	
	_, err = db.Exec(query)
	return &Database{db: db}, err
}

// SaveBatch inserta un grupo de servidores en una sola transacción
func (d *Database) SaveBatch(servers []*protocol.ServerDetail) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	tx, err := d.db.Begin()
	if err != nil {
		return err
	}

	stmt, _ := tx.Prepare(`
		INSERT OR REPLACE INTO servers 
		(ip, port, version_name, protocol_ver, motd, players_online, players_max, icon, is_whitelist, secure_chat, software, mods_json, plugins_json, timestamp) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	defer stmt.Close()

	for _, s := range servers {
		mods, _ := json.Marshal(s.Mods)
		plugins, _ := json.Marshal(s.Plugins)
		
		_, err := stmt.Exec(
			s.IP, s.Port, s.VersionName, s.Protocol, s.MOTD, 
			s.PlayersOnline, s.PlayersMax, s.Icon, s.IsWhitelist, 
			s.EnforcesSecureChat, s.Software, string(mods), string(plugins), s.Timestamp,
		)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}