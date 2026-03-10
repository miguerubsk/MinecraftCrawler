package cmd_test

import (
	"MinecraftCrawler/cmd"
	"testing"

	"github.com/spf13/cobra"
)

func TestInfoCommandStructure(t *testing.T) {
	if cmd.InfoCmd.Use != "info [objetivo]" {
		t.Errorf("Command use = %s; want info [objetivo]", cmd.InfoCmd.Use)
	}
	if cmd.InfoCmd.Short == "" {
		t.Error("Command short description is empty")
	}
	if cmd.InfoCmd.RunE == nil {
		t.Error("Command RunE function is nil")
	}
	if cmd.InfoCmd.Args == nil {
		t.Error("Command Args validator is nil")
	}
	if cmd.InfoCmd.InheritedFlags().Lookup("output") != nil {
		t.Error("info command should not inherit output flag")
	}
}

func TestInfoCommandArgsValidation(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "no args", args: []string{}, wantErr: true},
		{name: "one arg", args: []string{"mc.example.net"}, wantErr: false},
		{name: "two args", args: []string{"a", "b"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.InfoCmd.Args(&cobra.Command{}, tt.args)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Args() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
