package cmd

import (
	"fmt"
	"strings"
	"time"

	"MinecraftCrawler/internal/protocol"
	"MinecraftCrawler/internal/resolver"

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

		resolved, err := resolver.ResolveTarget(target, nil)
		if err != nil {
			return err
		}

		if resolved.DirectIP {
			color.Cyan("[*] IP detectada, omitiendo búsqueda SRV para %s\n", resolved.Host)
		} else if resolved.SRVLookupAttempted {
			color.Cyan("[*] Sin puerto especificado, buscando registro SRV para %s...\n", target)
			if resolved.UsedSRV {
				color.HiGreen("[+] SRV encontrado: %s:%d\n", resolved.Host, resolved.Port)
			} else if resolved.SRVNotFound {
				color.HiYellow("[-] Sin registro SRV, usando puerto por defecto 25565\n")
			}
		}

		host := resolved.Host
		port := resolved.Port

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
		fmt.Fprintf(out, "  %-22s %s\n", "Versión:",   color.HiWhiteString("%s", detail.VersionName))
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
		fmt.Fprintf(out, "  %-22s %s\n", "Software:", color.HiWhiteString("%s", software))
		fmt.Fprintf(out, "  %-22s %s\n", "Mapa:",     color.HiWhiteString("%s", mapName))

		fmt.Fprintln(out, sepMid)

		if detail.IsWhitelist {
			fmt.Fprintf(out, "  %-22s %s\n", "Lista Blanca:", color.HiRedString("ACTIVADA"))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "Lista Blanca:", color.HiBlackString("Desactivada"))
		}
		if detail.EnforcesSecureChat {
			fmt.Fprintf(out, "  %-22s %s\n", "Chat Seguro:", color.HiYellowString("Obligatorio"))
		}
		if !detail.RconAttempted {
			fmt.Fprintf(out, "  %-22s %s\n", "RCON:", color.HiBlackString("No intentado"))
		} else if detail.RconOpen {
			fmt.Fprintf(out, "  %-22s %s\n", "RCON:", color.HiRedString("ABIERTO"))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "RCON:", color.HiYellowString("Cerrado"))
		}

		if !detail.QueryAttempted {
			fmt.Fprintf(out, "  %-22s %s\n", "Query UDP:", color.HiBlackString("No intentado"))
		} else if detail.QueryError != nil {
			fmt.Fprintf(out, "  %-22s %s\n", "Query UDP:", color.HiBlackString("Inactivo"))
		} else {
			fmt.Fprintf(out, "  %-22s %s\n", "Query UDP:", color.HiGreenString("CONECTADO"))
			if detail.QueryHostName != "" {
				fmt.Fprintf(out, "  %-22s %s\n", "  Host:", color.HiWhiteString("%s", detail.QueryHostName))
			}
			if detail.QueryHostPort > 0 && detail.QueryHostPort != port {
				fmt.Fprintf(out, "  %-22s %s\n", "  Port (Interno):", color.HiWhiteString("%d", detail.QueryHostPort))
			}
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
