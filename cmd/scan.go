package cmd

import (
	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/scanner"
	"MinecraftCrawler/internal/storage"
	"io"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
)

var (
	ipRange     string
	rate        string
	port        int
	workers     int
	verbose     int
	excludeFile string
)

var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Inicia el escaneo y análisis",
	Example: `  mccrawler scan --range 1.1.1.0/24 --rate 1000
  mccrawler scan -r 1.2.3.4/32 -w 500 -v
  mccrawler scan --range 192.168.1.0/24 --exclude samples/exclude.txt`,
	Run: func(cmd *cobra.Command, args []string) {


		// 1. Configurar Logger dual (Archivo + Consola)
		logFile, err := os.OpenFile("crawler.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			color.Red("Error al crear archivo de log: %v", err)
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
		// El manager de almacenamiento siempre usará un batch de 500 para eficiencia
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

						if verbose > 0 && int(count) <= verbose {
							color.Green("[+] %-15s | %-15s | P: %d/%d | WL: %t",
								detail.IP, detail.VersionName, detail.PlayersOnline, detail.PlayersMax, detail.IsWhitelist)
						} else if verbose > 0 && int(count) == verbose + 1 {
							color.Yellow("[*] Límite de %d logs alcanzado. Continuando escaneo silencioso en base de datos...", verbose)
						}
						
						resultChan <- detail
					}
				}
			}()
		}

		// 5. Ejecutar Masscan
		color.Cyan("[*] Iniciando escaneo en %s (Puerto: %d, Workers: %d, Rate: %s)", ipRange, port, workers, rate)
		
		err = scanner.Run(ipRange, rate, port, excludeFile, ipChan)
		if err != nil {
			color.Red("Error ejecutando Masscan: %v", err)
			os.Exit(1)
		}

		// Esperar a que los workers terminen
		wg.Wait()
		close(resultChan)
		
		// Pequeña pausa para asegurar que el storage manager termine de escribir el último batch
		time.Sleep(1 * time.Second)
		color.Green("\n[*] Escaneo finalizado. Total encontrados: %d. Datos en: %s", atomic.LoadInt32(&foundCount), dbPath)
	},
}

func init() {
	ScanCmd.Flags().StringVarP(&ipRange, "range", "r", "", "Rango CIDR (ej: 1.1.1.0/24)")
	ScanCmd.Flags().StringVarP(&rate, "rate", "p", "1000", "PPS de Masscan")
	ScanCmd.Flags().IntVar(&port, "port", 25565, "Puerto objetivo")
	ScanCmd.Flags().IntVarP(&workers, "workers", "w", 1000, "Goroutines concurrentes")
	ScanCmd.Flags().IntVarP(&verbose, "verbose", "v", 0, "Muestra detalles de cada servidor encontrado (opcional: límite de líneas, default 500)")
	ScanCmd.Flags().Lookup("verbose").NoOptDefVal = "500"
	ScanCmd.Flags().StringVar(&excludeFile, "exclude", "", "Archivo de exclusiones (rangos de IP a evitar)")
	
	ScanCmd.MarkFlagRequired("range")
	
	rootCmd.AddCommand(ScanCmd)
}


