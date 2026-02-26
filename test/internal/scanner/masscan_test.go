package scanner_test

import (
	"MinecraftCrawler/internal/scanner"
	"reflect"
	"testing"
)

func TestBuildArguments(t *testing.T) {
	tests := []struct {
		name        string
		ipRange     string
		rate        string
		port        int
		excludeFile string
		want        []string
	}{
		{
			name:        "Basic args",
			ipRange:     "192.168.1.1/24",
			rate:        "1000",
			port:        25565,
			excludeFile: "",
			want: []string{
				"192.168.1.1/24",
				"-p", "25565",
				"--rate", "1000",
				"-oJ", "-",
			},
		},
		{
			name:        "With Exclude File",
			ipRange:     "10.0.0.0/8",
			rate:        "500",
			port:        80,
			excludeFile: "exclude.txt",
			want: []string{
				"10.0.0.0/8",
				"-p", "80",
				"--rate", "500",
				"-oJ", "-",
				"--excludefile", "exclude.txt",
			},
		},
		{
			name:        "Global Scan Excludes",
			ipRange:     "0.0.0.0/0",
			rate:        "10000",
			port:        25565,
			excludeFile: "",
			want: []string{
				"0.0.0.0/0",
				"-p", "25565",
				"--rate", "10000",
				"-oJ", "-",
				"--exclude", "255.255.255.255,127.0.0.0/8,0.0.0.0/8,224.0.0.0/4",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scanner.BuildArguments(tt.ipRange, tt.rate, tt.port, tt.excludeFile)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildArguments() = %v, want %v", got, tt.want)
			}
		})
	}
}

