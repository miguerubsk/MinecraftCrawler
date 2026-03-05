package cmd

import (
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"os"
	"fmt"
	"strings"
)

var dbPath  string

var banner = `
   _____         _________                      .__                
  /     \   ____ \_   ___ \____________ __  _  _|  |   ___________ 
 /  \ /  \_/ ___\/    \  \/\_  __ \__  \\ \/ \/ /  | _/ __ \_  __ \
/    Y    \  \___\     \____|  | \// __ \\     /|  |_\  ___/|  | \/
\____|__  /\___  >\______  /|__|  (____  /\/\_/ |____/\___  >__|   
        \/     \/        \/            \/                 \/         
`

var rootCmd = &cobra.Command{
	Use:   "mccrawler",
	Short: color.CyanString("Un crawler de Minecraft ultra eficiente"),
	Long:  "Escanea y analiza servidores de Minecraft a gran escala usando Masscan y Go.",
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
		fmt.Println(color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "output", "o", "results.db", "Archivo SQLite de salida")
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