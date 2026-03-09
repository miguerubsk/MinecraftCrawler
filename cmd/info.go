package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"
	// TODO: Integrate colored console output (milestone 03) using github.com/fatih/color
)

var InfoCmd = &cobra.Command{
	Use:   "info [target]",
	Short: "Deep analysis of a single Minecraft server",
	Long: `Performs a comprehensive analysis of a specific server.
It resolves SRV records, performs Server List Ping (SLP), 
and attempts UDP Query protocol extraction.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]

		fmt.Printf("🔍 Analyzing target: %s\n", target)

		// 1. TODO: Implement SRV Resolution (Issue #21)
		host := target
		port := 25565

		if h, p, err := net.SplitHostPort(target); err == nil {
			host = h
			if _, err := fmt.Sscanf(p, "%d", &port); err != nil {
				return fmt.Errorf("invalid port in target: %v", err)
			}
		} else {
			// Comprobamos si es una IP literal (v4 o v6) antes de cualquier otra lógica
			if net.ParseIP(target) != nil {
				fmt.Printf("[*] IP address detected, skipping SRV lookup for %s\n", host)
			} else {
				// Si no es un IP válida y tiene varios ":", es un IPv6 mal formateado (le faltan los [])
				if strings.Count(target, ":") > 1 {
					return fmt.Errorf("malformed IPv6 address: use '[addr]:port' format for IPv6 with literal ports")
				}

				// Intento de resolución SRV si no hay puerto explícito ni es IP
				fmt.Printf("[*] No port specified, attempting SRV lookup for %s...\n", host)
				_, addrs, err := net.LookupSRV("minecraft", "tcp", host)
				if err == nil && len(addrs) > 0 {
					host = strings.TrimSuffix(addrs[0].Target, ".")
					port = int(addrs[0].Port)
					fmt.Printf("SRV found: %s:%d\n", host, port)
				} else if err != nil {
					// Diferenciamos errores de DNS
					if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
						fmt.Println("No SRV record found, using default port 25565")
					} else {
						return fmt.Errorf("DNS resolution error: %v", err)
					}
				} else {
					fmt.Println("No SRV record found, using default port 25565")
				}
			}
		}

		if port < 1 || port > 65535 {
			return fmt.Errorf("port out of range: %d", port)
		}

		// 2. TODO: Trigger AnalyzeServer logic (Issue #22)
		// 3. TODO: Format rich output (Issue #23)

		fmt.Printf("[*] Target resolved: %s:%d\n", host, port)
		fmt.Println("Logic for deep analysis (Issue #22) and rich UI (Issue #23) pending.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
