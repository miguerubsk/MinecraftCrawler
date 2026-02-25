package cmd

import (
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/scanner"
	"MinecraftCrawler/internal/storage"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"sync"
	"time"
)

var (
	ipRange    string
	rate       string
	port       int
	workers    int
	verbose    bool
	excludeFile string
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Inicia el escaneo y análisis",
	Run: func(cmd *cobra.Command, args []string) {
		// 1. Inicializar DB con buffer optimizado
		db, _ := storage.NewDatabase(dbPath)
		
		ipChan := make(chan string, 10000)
		resultChan := make(chan *protocol.ServerDetail, 1000)

		// 2. Storage Manager (Escritura en disco optimizada)
		// El batchSize de 500 reduce drásticamente los IOPS en disco
		go storage.StartSQLiteManager(db, resultChan, 500)

		// 3. Worker Pool de Análisis
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ip := range ipChan {
					detail, err := protocol.AnalyzeServer(ip, port, 4*time.Second)
					if err == nil {
						if verbose {
							// Formato de log mejorado y legible
							log.Printf("[+] %-15s | %-10s | P: %d/%d | WL: %t", 
								detail.IP, detail.VersionName, detail.PlayersOnline, detail.PlayersMax, detail.IsWhitelist)
						}
						resultChan <- detail
					}
				}
			}()
		}

		// 4. Ejecutar Masscan
		fmt.Printf("[*] Iniciando escaneo en %s (Puerto: %d, Workers: %d)\n", ipRange, port, workers)
		err := scanner.Run(ipRange, rate, port, excludeFile, ipChan)
		if err != nil {
			log.Fatal(err)
		}

		wg.Wait()
		close(resultChan)
		fmt.Println("\n[*] Escaneo finalizado. Los datos se han guardado en", dbPath)
	},
}

func init() {
	scanCmd.Flags().StringVarP(&ipRange, "range", "r", "", "Rango CIDR (ej: 1.1.1.0/24)")
	scanCmd.Flags().StringVarP(&rate, "rate", "p", "1000", "PPS de Masscan")
	scanCmd.Flags().IntVar(&port, "port", 25565, "Puerto objetivo")
	scanCmd.Flags().IntVarP(&workers, "workers", "w", 1000, "Goroutines concurrentes")
	scanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Muestra detalles de cada servidor encontrado")
	scanCmd.Flags().StringVar(&excludeFile, "exclude", "", "Archivo de exclusiones (rangos de IP a evitar)")
	rootCmd.AddCommand(scanCmd)
}