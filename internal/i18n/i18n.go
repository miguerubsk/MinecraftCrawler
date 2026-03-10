package i18n

import (
	"fmt"
	"os"
	"strings"
)

var currentLanguage = "es"

var dictionary = map[string]map[string]string{
	"es": {
		"error.prefix":                          "Error: %v",
		"scan.start":                            "[*] Iniciando escaneo en %s (Puerto: %d, Workers: %d, Rate: %s)\n",
		"scan.masscan_error":                    "Error ejecutando Masscan: %v\n",
		"scan.logfile_error":                    "Error al crear archivo de log: %v",
		"scan.summary.title":                    "  RESUMEN DEL ESCANEO",
		"scan.summary.total":                    "Total Encontrados:",
		"scan.summary.duration":                 "Tiempo Total:",
		"scan.summary.speed":                    "Velocidad Media:",
		"scan.summary.database":                 "Base de Datos:",
		"scan.verbose.limit":                    "[*] Límite de %d logs alcanzado. Continuando escaneo silencioso en base de datos...\n",
		"info.analyzing_target":                 "[*] Analizando objetivo: %s\n",
		"info.srv.lookup":                       "[*] Sin puerto especificado, buscando registro SRV para %s...\n",
		"info.srv.found":                        "[+] SRV encontrado: %s:%d\n",
		"info.srv.not_found":                    "[-] Sin registro SRV, usando puerto por defecto 25565\n",
		"info.ip.detected":                      "[*] IP detectada, omitiendo búsqueda SRV para %s\n",
		"info.analysis_failed":                  "análisis fallido: %v",
		"info.title":                            "  ANÁLISIS DE SERVIDOR",
		"info.label.server":                     "Servidor:",
		"info.label.version":                    "Versión:",
		"info.label.protocol":                   "Protocolo:",
		"info.label.players":                    "Jugadores:",
		"info.label.software":                   "Software:",
		"info.label.map":                        "Mapa:",
		"info.label.whitelist":                  "Lista Blanca:",
		"info.label.secure_chat":                "Chat Seguro:",
		"info.label.rcon":                       "RCON:",
		"info.label.query":                      "Query UDP:",
		"info.value.unknown":                    "Desconocido",
		"info.value.na":                         "N/A",
		"info.value.whitelist_enabled":          "ACTIVADA",
		"info.value.whitelist_disabled":         "Desactivada",
		"info.value.secure_chat_required":       "Obligatorio",
		"info.value.rcon_not_attempted":         "No intentado",
		"info.value.rcon_open":                  "ABIERTO",
		"info.value.rcon_closed":                "Cerrado",
		"info.value.query_not_attempted":        "No intentado",
		"info.value.query_inactive":             "Inactivo",
		"info.value.query_connected":            "CONECTADO",
		"info.mods.title":                       "  MODS (%d)",
		"info.plugins.title":                    "  PLUGINS (%d)",
	},
	"en": {
		"error.prefix":                          "Error: %v",
		"scan.start":                            "[*] Starting scan on %s (Port: %d, Workers: %d, Rate: %s)\n",
		"scan.masscan_error":                    "Error running Masscan: %v\n",
		"scan.logfile_error":                    "Error creating log file: %v",
		"scan.summary.title":                    "  SCAN SUMMARY",
		"scan.summary.total":                    "Total Found:",
		"scan.summary.duration":                 "Total Time:",
		"scan.summary.speed":                    "Average Speed:",
		"scan.summary.database":                 "Database:",
		"scan.verbose.limit":                    "[*] Reached %d log lines. Continuing scan silently to database...\n",
		"info.analyzing_target":                 "[*] Analyzing target: %s\n",
		"info.srv.lookup":                       "[*] No port specified, looking up SRV record for %s...\n",
		"info.srv.found":                        "[+] SRV found: %s:%d\n",
		"info.srv.not_found":                    "[-] No SRV record found, using default port 25565\n",
		"info.ip.detected":                      "[*] IP detected, skipping SRV lookup for %s\n",
		"info.analysis_failed":                  "analysis failed: %v",
		"info.title":                            "  SERVER ANALYSIS",
		"info.label.server":                     "Server:",
		"info.label.version":                    "Version:",
		"info.label.protocol":                   "Protocol:",
		"info.label.players":                    "Players:",
		"info.label.software":                   "Software:",
		"info.label.map":                        "Map:",
		"info.label.whitelist":                  "Whitelist:",
		"info.label.secure_chat":                "Secure Chat:",
		"info.label.rcon":                       "RCON:",
		"info.label.query":                      "UDP Query:",
		"info.value.unknown":                    "Unknown",
		"info.value.na":                         "N/A",
		"info.value.whitelist_enabled":          "ENABLED",
		"info.value.whitelist_disabled":         "Disabled",
		"info.value.secure_chat_required":       "Required",
		"info.value.rcon_not_attempted":         "Not attempted",
		"info.value.rcon_open":                  "OPEN",
		"info.value.rcon_closed":                "Closed",
		"info.value.query_not_attempted":        "Not attempted",
		"info.value.query_inactive":             "Inactive",
		"info.value.query_connected":            "CONNECTED",
		"info.mods.title":                       "  MODS (%d)",
		"info.plugins.title":                    "  PLUGINS (%d)",
	},
}

func Init(preferred string) {
	lang := strings.ToLower(strings.TrimSpace(preferred))
	if lang == "es" || lang == "en" {
		currentLanguage = lang
		return
	}
	currentLanguage = detectFromEnvironment()
}

func Current() string {
	return currentLanguage
}

func T(key string, args ...any) string {
	table, ok := dictionary[currentLanguage]
	if !ok {
		table = dictionary["es"]
	}
	msg, ok := table[key]
	if !ok {
		msg = key
	}
	if len(args) == 0 {
		return msg
	}
	return fmt.Sprintf(msg, args...)
}

func detectFromEnvironment() string {
	candidates := []string{"LC_ALL", "LC_MESSAGES", "LANGUAGE", "LANG"}
	for _, envName := range candidates {
		if val := os.Getenv(envName); val != "" {
			lang := normalizeLocale(val)
			if lang == "es" || lang == "en" {
				return lang
			}
		}
	}
	return "es"
}

func normalizeLocale(raw string) string {
	v := strings.ToLower(strings.TrimSpace(raw))
	if v == "" {
		return ""
	}
	if idx := strings.Index(v, "."); idx > 0 {
		v = v[:idx]
	}
	if idx := strings.Index(v, ":"); idx > 0 {
		v = v[:idx]
	}
	v = strings.ReplaceAll(v, "-", "_")
	parts := strings.Split(v, "_")
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}
