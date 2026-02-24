package storage

import (
	"log"
	"MinecraftCrawler/internal/protocol"
	"time"
)

// StartManager inicia el proceso de guardado en base de datos
// Retorna un canal 'done' que se cerrará cuando termine de procesar
func StartManager(db *Database, resultChan <-chan *protocol.ServerDetail, batchSize int) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		var batch []*protocol.ServerDetail
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case res, ok := <-resultChan:
				if !ok {
					// Canal cerrado, guardar lo que quede y salir
					if len(batch) > 0 {
						if err := db.SaveBatch(batch); err != nil {
							log.Printf("Error guardando batch final: %v", err)
						}
					}
					return
				}
				batch = append(batch, res)
				if len(batch) >= batchSize {
					if err := db.SaveBatch(batch); err != nil {
						log.Printf("Error guardando batch: %v", err)
					}
					batch = nil // Reiniciar batch después de guardar
				}
			case <-ticker.C:
				if len(batch) > 0 {
					if err := db.SaveBatch(batch); err != nil {
						log.Printf("Error guardando batch por tiempo: %v", err)
					}
					batch = nil // Reiniciar batch después de guardar
				}
			}
		}
	}()
	return done
}