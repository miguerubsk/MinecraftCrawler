package cmd_test

import (
	"MinecraftCrawler/cmd"
	"testing"
)

func TestScanFlags(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		expected string
	}{
		{"Rate", "rate", "1000"},
		{"Port", "port", "25565"},
		{"Workers", "workers", "1000"},
		{"Verbose", "verbose", "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := cmd.ScanCmd.Flags().Lookup(tt.flag)
			if flag == nil {
				t.Errorf("Flag %s not found", tt.flag)
				return
			}
			if flag.DefValue != tt.expected {
				t.Errorf("Flag %s default value = %s; want %s", tt.flag, flag.DefValue, tt.expected)
			}
		})
	}
}

func TestScanCommandStructure(t *testing.T) {
	if cmd.ScanCmd.Use != "scan" {
		t.Errorf("Command use = %s; want scan", cmd.ScanCmd.Use)
	}
	if cmd.ScanCmd.Short == "" {
		t.Error("Command short description is empty")
	}
	if cmd.ScanCmd.Run == nil {
		t.Error("Command Run function is nil")
	}
}

