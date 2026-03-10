package protocol_test

import (
	"MinecraftCrawler/internal/protocol"
	"strings"
	"testing"

	"github.com/fatih/color"
)

func withNoColor(t *testing.T, noColor bool) {
	t.Helper()
	original := color.NoColor
	color.NoColor = noColor
	t.Cleanup(func() {
		color.NoColor = original
	})
}

func TestColorizeMOTDWithColorCodesAndEdgeCases(t *testing.T) {
	withNoColor(t, false)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "no formatting codes",
			in:   "Servidor vanilla",
			want: "Servidor vanilla",
		},
		{
			name: "lowercase color code",
			in:   "§aHola",
			want: "\033[92mHola\033[0m",
		},
		{
			name: "uppercase color code",
			in:   "§AHola",
			want: "\033[92mHola\033[0m",
		},
		{
			name: "uppercase style code",
			in:   "§LNegrita",
			want: "\033[1mNegrita\033[0m",
		},
		{
			name: "obfuscated code is consumed as no-op",
			in:   "Cuidado §k106",
			want: "Cuidado 106",
		},
		{
			name: "unknown code is preserved",
			in:   "§xTexto",
			want: "§xTexto",
		},
		{
			name: "dangling section sign is preserved",
			in:   "Texto§",
			want: "Texto§",
		},
		{
			name: "mixed color and explicit reset",
			in:   "§aVerde§r normal",
			want: "\033[92mVerde\033[0m normal\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := protocol.ColorizeMOTD(tt.in)
			if got != tt.want {
				t.Fatalf("ColorizeMOTD() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestColorizeMOTDNoColorStripsValidCodes(t *testing.T) {
	withNoColor(t, true)

	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "lowercase color code stripped",
			in:   "§aHola",
			want: "Hola",
		},
		{
			name: "uppercase color code stripped",
			in:   "§BHola",
			want: "Hola",
		},
		{
			name: "uppercase style code stripped",
			in:   "§ONice",
			want: "Nice",
		},
		{
			name: "obfuscated code stripped",
			in:   "abc§kdef",
			want: "abcdef",
		},
		{
			name: "mixed valid codes stripped from text",
			in:   "§aHola §LMundo§r!",
			want: "Hola Mundo!",
		},
		{
			name: "unknown code remains untouched",
			in:   "§xHola",
			want: "§xHola",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := protocol.ColorizeMOTD(tt.in)
			if got != tt.want {
				t.Fatalf("ColorizeMOTD() = %q, want %q", got, tt.want)
			}

			if strings.Contains(got, "\033[") {
				t.Fatalf("ColorizeMOTD() emitted ANSI escape in NoColor mode: %q", got)
			}
		})
	}
}
