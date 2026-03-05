package cmd

import (
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"os"
	"fmt"
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
	Long:  color.GreenString(banner) + "\nEscanea y analiza servidores de Minecraft a gran escala usando Masscan y Go.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(color.RedString("Error: %v", err))
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&dbPath, "output", "o", "results.db", "Archivo SQLite de salida")
}