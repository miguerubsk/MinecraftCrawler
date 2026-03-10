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
		const separatorWidth = 50
		out := color.Output
		sep    := color.HiCyanString("\n" + strings.Repeat("━", separatorWidth))
		sepMid := color.HiCyanString(strings.Repeat("━", separatorWidth))
		sepEnd := color.HiCyanString(strings.Repeat("━", separatorWidth) + "\n")

		fmt.Fprintln(out, sep)
		fmt.Fprintln(out, color.HiWhiteString("  ANÁLISIS DE SERVIDOR"))
		fmt.Fprintln(out, sepMid)

		if detail.MOTD != "" {
			fmt.Fprintf(out, "  %s\n", protocol.ColorizeMOTD(detail.MOTD))
			fmt.Fprintln(out, sepMid)
		}

		fmt.Fprintf(out, "  %-22s %s\n", "Servidor:",  color.HiYellowString("%s:%d", host, port))
		fmt.Fprintf(out, "  %-22s %s\n", "Versión:",   color.HiWhiteString(detail.VersionName))
		fmt.Fprintf(out, "  %-22s %s\n", "Protocolo:", color.HiWhiteString("%d", detail.Protocol))

		playersStr := fmt.Sprintf("%d / %d", detail.PlayersOnline, detail.PlayersMax)
		if detail.PlayersOnline > 0 {
			fmt.Fprintf(out, "  %-22s %s\n", "Jugadores:", color.HiGreenString("%s", playersStr))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "Jugadores:", color.HiBlackString("%s", playersStr))
		}

		software := detail.Software
		if software == "" {
			software = "Desconocido"
		}
		mapName := detail.MapName
		if mapName == "" {
			mapName = "N/A"
		}
		fmt.Fprintf(out, "  %-22s %s\n", "Software:", color.HiWhiteString(software))
		fmt.Fprintf(out, "  %-22s %s\n", "Mapa:",     color.HiWhiteString(mapName))

		fmt.Fprintln(out, sepMid)

		if detail.IsWhitelist {
			fmt.Fprintf(out, "  %-22s %s\n", "Lista Blanca:", color.HiRedString("ACTIVADA"))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "Lista Blanca:", color.HiBlackString("Desactivada"))
		}
		if detail.EnforcesSecureChat {
			fmt.Fprintf(out, "  %-22s %s\n", "Chat Seguro:", color.HiYellowString("Obligatorio"))
		}
		if detail.RconOpen {
			fmt.Fprintf(out, "  %-22s %s\n", "RCON:", color.HiRedString("ABIERTO"))
		}

		if detail.QueryError != nil {
			fmt.Fprintf(out, "  %-22s %s\n", "Query UDP:", color.HiBlackString("Inactivo"))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "Query UDP:", color.HiGreenString("CONECTADO"))
		}

		if len(detail.Mods) > 0 {
			fmt.Fprintln(out, sepMid)
			fmt.Fprintln(out, color.HiWhiteString("  MODS (%d)", len(detail.Mods)))
			for id, ver := range detail.Mods {
				fmt.Fprintf(out, "  [+] %-30s %s\n", color.HiGreenString("%s", id), color.HiBlackString("%s", ver))
			}
		}

		if len(detail.Plugins) > 0 {
			fmt.Fprintln(out, sepMid)
			fmt.Fprintln(out, color.HiWhiteString("  PLUGINS (%d)", len(detail.Plugins)))
			for _, pl := range detail.Plugins {
				fmt.Fprintf(out, "  [+] %s\n", color.HiGreenString("%s", pl))
			}
		}

		fmt.Fprintln(out, sepEnd)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(InfoCmd)
}
