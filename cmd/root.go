package cmd

import (
	"MinecraftCrawler/internal/i18n"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"os"
	"fmt"
	"strings"
)

var dbPath  string
var language string

var banner = `
___  ____                            __ _     _____                    _           
|  \/  (_)                          / _| |   /  __ \                  | |          
| .  . |_ _ __   ___  ___ _ __ __ _| |_| |_  | /  \/_ __ __ ___      _| | ___ _ __ 
| |\/| | | '_ \ / _ \/ __| '__/ _` + "`" + ` |  _| __| | |   | '__/ _` + "`" + ` \ \ /\ / / |/ _ \ '__|
| |  | | | | | |  __/ (__| | | (_| | | | |_  | \__/\ | | (_| |\ V  V /| |  __/ |   
\_|  |_/_|_| |_|\___|\___|_|  \__,_|_|  \__|  \____/_|  \__,_| \_/\_/ |_|\___|_|   
`

var rootCmd = &cobra.Command{
	Use:   "mccrawler",
	Short: color.CyanString("Un crawler de Minecraft ultra eficiente"),
	Long:  "Escanea y analiza servidores de Minecraft a gran escala usando Masscan y Go.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		i18n.Init(language)
	},
}

func PrintBanner() {
	fmt.Fprintln(color.Output, color.GreenString(banner))
}

func Execute() {
	cobra.AddTemplateFunc("colorYellow", func(s string) string { return color.YellowString(s) })
	cobra.AddTemplateFunc("colorGreen", func(s string) string { return color.GreenString(s) })
	cobra.AddTemplateFunc("colorCyan", func(s string) string { return color.CyanString(s) })
	cobra.AddTemplateFunc("pad", func(s string, padding int) string {
		template := fmt.Sprintf("%%-%ds", padding)
		return fmt.Sprintf(template, s)
	})
	cobra.AddTemplateFunc("trimTrailingWhitespaces", func(s string) string {
		return strings.TrimRight(s, " \t\n\r")
	})
	cobra.AddTemplateFunc("banner", func() string {
		return color.GreenString(banner)
	})
	
	rootCmd.SetHelpTemplate(rootHelpTemplate)
	
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString(i18n.T("error.prefix"), err))
		os.Exit(1)
	}
}

func init() {
	i18n.Init("auto")
	rootCmd.PersistentFlags().StringVar(&language, "lang", "auto", "Language for CLI output (auto|es|en)")
}

var rootHelpTemplate = `{{banner}}
{{if .Long}}{{.Long}}{{else}}{{.Short}}{{end}}

{{colorYellow "USO:"}}
  {{colorCyan .UseLine}}

{{if .HasAvailableSubCommands}}{{colorYellow "COMANDOS DISPONIBLES:"}}
{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}  {{colorGreen (pad .Name .NamePadding)}} {{.Short}}{{end}}
{{end}}{{end}}
{{if .HasAvailableLocalFlags}}{{colorYellow "PARÁMETROS:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasAvailableInheritedFlags}}{{colorYellow "PARÁMETROS GLOBALES:"}}
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}

{{if .HasExample}}{{colorYellow "EJEMPLOS:"}}
{{.Example}}{{end}}
Usa "{{colorCyan "mccrawler [comando] --help"}}" para más información sobre un comando.
`