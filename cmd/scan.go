package cmd

import (
	"github.com/spf13/cobra"
	"MinecraftCrawler/internal/scanner"
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/storage"
	"log"
	"sync"
	"time"
)

var (
	ipRange    string
	rate       string
	port       int
	exclusions string
	workers    int
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Inicia un escaneo completo (Masscan + Analysis)",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Inicializar DB
		db, _ := storage.NewDatabase(dbPath)
		
		// 2. Canales
		ipChan := make(chan string, 10000)
		resultChan := make(chan *protocol.ServerDetail, 1000)

		// 3. Iniciar Storage y Workers
		storageDone := storage.StartManager(db, resultChan, 500)
		
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ip := range ipChan {
					// Si el puerto es 25575 (RCON), solo hacemos estadística básica
					if port == 25575 {
						// Lógica simplificada para RCON aquí
						continue
					}
					detail, err := protocol.AnalyzeServer(ip, port, 4*time.Second)
					if err == nil {
						log.Printf("[+] Found: %s:%d (%s) - Players: %d/%d - Whitelist: %t", detail.IP, detail.Port, detail.VersionName, detail.PlayersOnline, detail.PlayersMax, detail.IsWhitelist)
						resultChan <- detail
					}
				}
			}()
		}

		// 4. Ejecutar Masscan
		log.Printf("[*] Starting Masscan on %s with rate %s...", ipRange, rate)
		// Ahora bloquea hasta que termina y cerramos el canal
		if err := scanner.Run(ipRange, rate, port, exclusions, ipChan); err != nil {
			// En caso de error crítico en masscan, cerramos para no bloquear
			close(ipChan)
			// Podríamos loguear el error fatal aquí
		} else {
			close(ipChan)
		}

		// 5. Esperar finalización ordenada
		wg.Wait()
		close(resultChan)
		<-storageDone
	},
}

func init() {
	scanCmd.Flags().StringVarP(&ipRange, "range", "r", "", "Rango de IPs (CIDR)")
	scanCmd.Flags().StringVarP(&rate, "rate", "p", "1000", "Ratio de Masscan (pps)")
	scanCmd.Flags().IntVar(&port, "port", 25565, "Puerto a escanear (25565 o 25575)")
	scanCmd.Flags().StringVar(&exclusions, "exclude", "", "Archivo de exclusiones")
	scanCmd.Flags().IntVarP(&workers, "workers", "w", 1000, "Número de goroutines concurrentes")
	rootCmd.AddCommand(scanCmd)
}