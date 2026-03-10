package cmd

import (
	"fmt"
	"strings"
	"time"

	"MinecraftCrawler/internal/i18n"
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
		const (
			infoLineFmt   = "  %-22s %s\n"
			rconLabelKey  = "info.label.rcon"
			queryLabelKey = "info.label.query"
		)

		PrintBanner()
		target := args[0]

		color.Cyan(i18n.T("info.analyzing_target"), target)

		resolved, err := resolver.ResolveTarget(target, nil)
		if err != nil {
			return err
		}

		if resolved.DirectIP {
			color.Cyan(i18n.T("info.ip.detected"), resolved.Host)
		} else if resolved.SRVLookupAttempted {
			color.Cyan(i18n.T("info.srv.lookup"), target)
			if resolved.UsedSRV {
				color.HiGreen(i18n.T("info.srv.found"), resolved.Host, resolved.Port)
			} else if resolved.SRVNotFound {
				color.HiYellow(i18n.T("info.srv.not_found"))
			}
		}

		host := resolved.Host
		port := resolved.Port

		detail, err := protocol.AnalyzeServer(host, port, 5*time.Second)
		if err != nil {
			return fmt.Errorf(i18n.T("info.analysis_failed"), err)
		}

		// Issue #23: UI alineada con el estilo del comando scan
		const separatorWidth = 50
		out := color.Output
		sep    := color.HiCyanString("\n" + strings.Repeat("━", separatorWidth))
		sepMid := color.HiCyanString(strings.Repeat("━", separatorWidth))
		sepEnd := color.HiCyanString(strings.Repeat("━", separatorWidth) + "\n")

		fmt.Fprintln(out, sep)
		fmt.Fprintln(out, color.HiWhiteString(i18n.T("info.title")))
		fmt.Fprintln(out, sepMid)

		if detail.MOTD != "" {
			fmt.Fprintf(out, "  %s\n", protocol.ColorizeMOTD(detail.MOTD))
			fmt.Fprintln(out, sepMid)
		}

		fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.server"), color.HiYellowString("%s:%d", host, port))
		fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.version"), color.HiWhiteString("%s", detail.VersionName))
		fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.protocol"), color.HiWhiteString("%d", detail.Protocol))

		playersStr := fmt.Sprintf("%d / %d", detail.PlayersOnline, detail.PlayersMax)
		if detail.PlayersOnline > 0 {
			fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.players"), color.HiGreenString("%s", playersStr))
		} else {
			fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.players"), color.HiBlackString("%s", playersStr))
		}

		software := detail.Software
		if software == "" {
			software = i18n.T("info.value.unknown")
		}
		mapName := detail.MapName
		if mapName == "" {
			mapName = i18n.T("info.value.na")
		}
		fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.software"), color.HiWhiteString("%s", software))
		fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.map"), color.HiWhiteString("%s", mapName))

		fmt.Fprintln(out, sepMid)

		if detail.IsWhitelist {
			fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.whitelist"), color.HiRedString(i18n.T("info.value.whitelist_enabled")))
		} else {
			fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.whitelist"), color.HiBlackString(i18n.T("info.value.whitelist_disabled")))
		}
		if detail.EnforcesSecureChat {
			fmt.Fprintf(out, infoLineFmt, i18n.T("info.label.secure_chat"), color.HiYellowString(i18n.T("info.value.secure_chat_required")))
		}
		if !detail.RconAttempted {
			fmt.Fprintf(out, infoLineFmt, i18n.T(rconLabelKey), color.HiBlackString(i18n.T("info.value.rcon_not_attempted")))
		} else if detail.RconOpen {
			fmt.Fprintf(out, infoLineFmt, i18n.T(rconLabelKey), color.HiRedString(i18n.T("info.value.rcon_open")))
		} else {
			fmt.Fprintf(out, infoLineFmt, i18n.T(rconLabelKey), color.HiYellowString(i18n.T("info.value.rcon_closed")))
		}

		if !detail.QueryAttempted {
			fmt.Fprintf(out, infoLineFmt, i18n.T(queryLabelKey), color.HiBlackString(i18n.T("info.value.query_not_attempted")))
		} else if detail.QueryError != nil {
			fmt.Fprintf(out, infoLineFmt, i18n.T(queryLabelKey), color.HiBlackString(i18n.T("info.value.query_inactive")))
		} else {
			fmt.Fprintf(out, infoLineFmt, i18n.T(queryLabelKey), color.HiGreenString(i18n.T("info.value.query_connected")))
		}

		if len(detail.Mods) > 0 {
			fmt.Fprintln(out, sepMid)
			fmt.Fprintln(out, color.HiWhiteString(i18n.T("info.mods.title"), len(detail.Mods)))
			for id, ver := range detail.Mods {
				fmt.Fprintf(out, "  [+] %-30s %s\n", color.HiGreenString("%s", id), color.HiBlackString("%s", ver))
			}
		}

		if len(detail.Plugins) > 0 {
			fmt.Fprintln(out, sepMid)
			fmt.Fprintln(out, color.HiWhiteString(i18n.T("info.plugins.title"), len(detail.Plugins)))
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
