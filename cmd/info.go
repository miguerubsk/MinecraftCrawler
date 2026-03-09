package cmd

import (
	"fmt"
	"net"
	"strings"
	"time"
	"MinecraftCrawler/internal/protocol"

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

		fmt.Printf("Analyzing target: %s\n", target)

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
		detail, err := protocol.AnalyzeServer(host, port, 5*time.Second)
		if err != nil {
			return fmt.Errorf("analysis failed: %v", err)
		}

		// 3. TODO: Format rich output (Issue #23)
		fmt.Printf("\n--- Technical Data ---\n")
		fmt.Printf("IP: %s | Port: %d\n", detail.IP, detail.Port)
		fmt.Printf("MOTD: %s\n", detail.MOTD)
		fmt.Printf("Version: %s (Protocol: %d)\n", detail.VersionName, detail.Protocol)
		fmt.Printf("Players: %d/%d\n", detail.PlayersOnline, detail.PlayersMax)
		fmt.Printf("Software: %s | Map: %s\n", detail.Software, detail.MapName)
		
		if detail.IsWhitelist {
			fmt.Println("Whitelist: Enabled")
		}
		if detail.EnforcesSecureChat {
			fmt.Println("Secure Chat: Enforced")
		}
		if detail.RconOpen {
			fmt.Println("RCON: Open")
		}

		if len(detail.Mods) > 0 {
			fmt.Printf("Mods (%d): ", len(detail.Mods))
			mods := make([]string, 0, len(detail.Mods))
			for id, ver := range detail.Mods {
				mods = append(mods, fmt.Sprintf("%s (%s)", id, ver))
			}
			fmt.Println(strings.Join(mods, ", "))
		}

		if len(detail.Plugins) > 0 {
			fmt.Printf("Plugins (%d): %s\n", len(detail.Plugins), strings.Join(detail.Plugins, ", "))
		}

		fmt.Println("\nDeep analysis complete. Visual formatting pending (Issue #23).")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
