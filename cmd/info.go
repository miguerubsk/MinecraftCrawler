package cmd

import (
	"fmt"
	"net"
	"strings"
	"time"

	"MinecraftCrawler/internal/protocol"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var InfoCmd = &cobra.Command{
	Use:   "info [objetivo]",
	Short: "Análisis profundo de un servidor Minecraft",
	Long: `Realiza un análisis exhaustivo de un servidor específico.
Resuelve registros SRV, realiza Server List Ping (SLP)
e intenta la extracción mediante protocolo UDP Query.`,
	Example: `  mccrawler info mc.hypixel.net
  mccrawler info 192.168.1.100:25565`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		PrintBanner()
		target := args[0]

		color.Cyan("[*] Analizando objetivo: %s\n", target)

		host := target
		port := 25565

		if h, p, err := net.SplitHostPort(target); err == nil {
			host = h
			if _, err := fmt.Sscanf(p, "%d", &port); err != nil {
				return fmt.Errorf("puerto inválido en el objetivo: %v", err)
			}
		} else {
			if net.ParseIP(target) != nil {
				color.Cyan("[*] IP detectada, omitiendo búsqueda SRV para %s\n", host)
			} else {
				if strings.Count(target, ":") > 1 {
					return fmt.Errorf("dirección IPv6 mal formateada: usa el formato '[addr]:port'")
				}

				color.Cyan("[*] Sin puerto especificado, buscando registro SRV para %s...\n", host)
				_, addrs, err := net.LookupSRV("minecraft", "tcp", host)
				if err == nil && len(addrs) > 0 {
					host = strings.TrimSuffix(addrs[0].Target, ".")
					port = int(addrs[0].Port)
					color.HiGreen("[+] SRV encontrado: %s:%d\n", host, port)
				} else if err != nil {
					if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
						color.HiYellow("[-] Sin registro SRV, usando puerto por defecto 25565\n")
					} else {
						return fmt.Errorf("error de resolución DNS: %v", err)
					}
				} else {
					color.HiYellow("[-] Sin registro SRV, usando puerto por defecto 25565\n")
				}
			}
		}

		if port < 1 || port > 65535 {
			return fmt.Errorf("puerto fuera de rango: %d", port)
		}

		detail, err := protocol.AnalyzeServer(host, port, 5*time.Second)
		if err != nil {
			return fmt.Errorf("análisis fallido: %v", err)
		}

		// Issue #23: UI alineada con el estilo del comando scan
		sep    := color.HiCyanString("\n" + strings.Repeat("━", 50))
		sepMid := color.HiCyanString(strings.Repeat("━", 50))
		sepEnd := color.HiCyanString(strings.Repeat("━", 50) + "\n")

		fmt.Println(sep)
		fmt.Println(color.HiWhiteString("  ANÁLISIS DE SERVIDOR"))
		fmt.Println(sepMid)

		if detail.MOTD != "" {
			fmt.Printf("  %s\n", protocol.ColorizeMOTD(detail.MOTD))
		}

		fmt.Println(sepMid)

		fmt.Printf("  %-22s %s\n", "Servidor:",  color.HiYellowString("%s:%d", host, port))
		fmt.Printf("  %-22s %s\n", "Versión:",   color.HiWhiteString(detail.VersionName))
		fmt.Printf("  %-22s %s\n", "Protocolo:", color.HiWhiteString("%d", detail.Protocol))

		playersStr := fmt.Sprintf("%d / %d", detail.PlayersOnline, detail.PlayersMax)
		if detail.PlayersOnline > 0 {
			fmt.Printf("  %-22s %s\n", "Jugadores:", color.HiGreenString(playersStr))
		} else {
			fmt.Printf("  %-22s %s\n", "Jugadores:", color.HiBlackString(playersStr))
		}

		software := detail.Software
		if software == "" {
			software = "Desconocido"
		}
		mapName := detail.MapName
		if mapName == "" {
			mapName = "N/A"
		}
		fmt.Printf("  %-22s %s\n", "Software:", color.HiWhiteString(software))
		fmt.Printf("  %-22s %s\n", "Mapa:",     color.HiWhiteString(mapName))

		fmt.Println(sepMid)

		if detail.IsWhitelist {
			fmt.Printf("  %-22s %s\n", "Lista Blanca:", color.HiRedString("ACTIVADA"))
		} else {
			fmt.Printf("  %-22s %s\n", "Lista Blanca:", color.HiBlackString("Desactivada"))
		}
		if detail.EnforcesSecureChat {
			fmt.Printf("  %-22s %s\n", "Chat Seguro:", color.HiYellowString("Obligatorio"))
		}
		if detail.RconOpen {
			fmt.Printf("  %-22s %s\n", "RCON:", color.HiRedString("ABIERTO"))
		}

		if len(detail.Mods) > 0 {
			fmt.Println(sepMid)
			fmt.Println(color.HiWhiteString("  MODS (%d)", len(detail.Mods)))
			for id, ver := range detail.Mods {
				fmt.Printf("  [+] %-30s %s\n", color.HiGreenString(id), color.HiBlackString(ver))
			}
		}

		if len(detail.Plugins) > 0 {
			fmt.Println(sepMid)
			fmt.Println(color.HiWhiteString("  PLUGINS (%d)", len(detail.Plugins)))
			for _, pl := range detail.Plugins {
				fmt.Printf("  [+] %s\n", color.HiGreenString(pl))
			}
		}

		fmt.Println(sepEnd)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
