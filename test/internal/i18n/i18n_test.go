package i18n_test

import (
	"MinecraftCrawler/internal/i18n"
	"testing"
)

func TestInitPreferredLanguage(t *testing.T) {
	i18n.Init("en")
	if i18n.Current() != "en" {
		t.Fatalf("Current() = %q, want %q", i18n.Current(), "en")
	}

	i18n.Init("es")
	if i18n.Current() != "es" {
		t.Fatalf("Current() = %q, want %q", i18n.Current(), "es")
	}
}

func TestTranslateUsesCurrentLanguage(t *testing.T) {
	i18n.Init("en")
	if got := i18n.T("info.label.server"); got != "Server:" {
		t.Fatalf("T(info.label.server) = %q, want %q", got, "Server:")
	}

	i18n.Init("es")
	if got := i18n.T("info.label.server"); got != "Servidor:" {
		t.Fatalf("T(info.label.server) = %q, want %q", got, "Servidor:")
	}
}
