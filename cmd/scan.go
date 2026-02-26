package cmd

import (
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/scanner"
	"MinecraftCrawler/internal/storage"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var (
	ipRange     string
	rate        string
	port        int
	workers     int
	verbose     bool
	excludeFile string
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Inicia el escaneo y análisis",
	Run: func(cmd *cobra.Command, args []string) {


		// 1. Configurar Logger dual (Archivo + Consola)
		logFile, err := os.OpenFile("crawler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Error al crear archivo de log: %v\n", err)
			return
		}
		defer logFile.Close()

		// MultiWriter envía los logs a ambos destinos
		multiWriter := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(multiWriter)

		// 2. Inicializar DB
		db, err := storage.NewDatabase(dbPath)
		if err != nil {
			log.Fatalf("Error al abrir la base de datos: %v", err)
		}

		ipChan := make(chan string, 10000)
		resultChan := make(chan *protocol.ServerDetail, 1000)
		
		// Contador para limitar la salida a 500 servidores
		var foundCount int32

		// 3. Storage Manager (Escritura en disco optimizada)
		go storage.StartSQLiteManager(db, resultChan, 500)

		// 4. Worker Pool de Análisis
		var wg sync.WaitGroup
		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for ip := range ipChan {
					detail, err := protocol.AnalyzeServer(ip, port, 4*time.Second)
					if err == nil {
						// Incrementamos el contador de forma segura entre hilos
						count := atomic.AddInt32(&foundCount, 1)

						if verbose && count <= 500 {
							log.Printf("[+] %-15s | %-15s | P: %d/%d | WL: %t",
								detail.IP, detail.VersionName, detail.PlayersOnline, detail.PlayersMax, detail.IsWhitelist)
						} else if count == 501 {
							log.Println("[*] Límite de 500 logs alcanzado. Continuando escaneo silencioso en base de datos...")
						}
						
						resultChan <- detail
					}
				}
			}()
		}

		// 5. Ejecutar Masscan
		log.Printf("[*] Iniciando escaneo en %s (Puerto: %d, Workers: %d, Rate: %s)\n", ipRange, port, workers, rate)
		
		err = scanner.Run(ipRange, rate, port, excludeFile, ipChan)
		if err != nil {
			log.Fatalf("Error ejecutando Masscan: %v", err)
		}

		// Esperar a que los workers terminen
		wg.Wait()
		close(resultChan)
		
		// Pequeña pausa para asegurar que el storage manager termine de escribir el último batch
		time.Sleep(1 * time.Second)
		log.Printf("\n[*] Escaneo finalizado. Total encontrados: %d. Datos en: %s\n", atomic.LoadInt32(&foundCount), dbPath)
	},
}

func init() {
	ScanCmd.Flags().StringVarP(&ipRange, "range", "r", "", "Rango CIDR (ej: 1.1.1.0/24)")
	ScanCmd.Flags().StringVarP(&rate, "rate", "p", "1000", "PPS de Masscan")
	ScanCmd.Flags().IntVar(&port, "port", 25565, "Puerto objetivo")
	ScanCmd.Flags().IntVarP(&workers, "workers", "w", 1000, "Goroutines concurrentes")
	ScanCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Muestra detalles de cada servidor encontrado")
	ScanCmd.Flags().StringVar(&excludeFile, "exclude", "", "Archivo de exclusiones (rangos de IP a evitar)")
	rootCmd.AddCommand(ScanCmd)
}


