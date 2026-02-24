package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var cfgFile string
var dbPath  string

var rootCmd = &cobra.Command{
	Use:   "mccrawler",
	Short: "Un crawler de Minecraft ultra eficiente",
	Long:  `Escanea y analiza servidores de Minecraft a gran escala usando Masscan y Go.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "output", "o", "results.db", "Archivo SQLite de salida")
}