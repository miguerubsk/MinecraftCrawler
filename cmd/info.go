package cmd

import (
	"fmt"
	"net"
	"os"
	"github.com/spf13/cobra"
	// "github.com/fatih/color" <-- Para el Hito 03
)

var infoCmd = &cobra.Command{
	Use:   "info [target]",
	Short: "Deep analysis of a single Minecraft server",
	Long: `Performs a comprehensive analysis of a specific server.
It resolves SRV records, performs Server List Ping (SLP), 
and attempts UDP Query protocol extraction.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		
		fmt.Printf("🔍 Analyzing target: %s\n", target)

		// 1. TODO: Implement SRV Resolution (Issue #1)
		host := target
		port := 25565

		if h, p, err := net.SplitHostPort(target); err == nil {
			host = h
			fmt.Sscanf(p, "%d", &port)
		} else {
			// Intento de resolución SRV si no hay puerto explícito
			fmt.Printf("[*] No port specified, attempting SRV lookup for %s...\n", host)
			_, addrs, err := net.LookupSRV("minecraft", "tcp", host)
			if err == nil && len(addrs) > 0 {
				host = strings.TrimSuffix(addrs[0].Target, ".")
				port = int(addrs[0].Port)
				fmt.Printf("✅ SRV found: %s:%d\n", host, port)
			} else {
				fmt.Println("ℹ️ No SRV record found, using default port 25565")
			}
		}

		// 2. TODO: Trigger AnalyzeServer logic (Issue #2)
		// 3. TODO: Format rich output (Issue #3)
		
		fmt.Println("⚠️ Logic not yet implemented. Check Milestone 01.")
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}