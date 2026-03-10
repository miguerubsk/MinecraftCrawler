package protocol

import (
	"strings"

	"github.com/fatih/color"
)

// MinecraftColorMap mapea los códigos § a secuencias ANSI básicas
var MinecraftColorMap = map[rune]string{
	'0': "\033[30m", // Negro
	'1': "\033[34m", // Azul Oscuro
	'2': "\033[32m", // Verde Oscuro
	'3': "\033[36m", // Cian Oscuro
	'4': "\033[31m", // Rojo Oscuro
	'5': "\033[35m", // Púrpura
	'6': "\033[33m", // Dorado (Amarillo)
	'7': "\033[37m", // Gris
	'8': "\033[90m", // Gris Oscuro
	'9': "\033[94m", // Azul
	'a': "\033[92m", // Verde
	'b': "\033[96m", // Cian
	'c': "\033[91m", // Rojo
	'd': "\033[95m", // Rosa
	'e': "\033[93m", // Amarillo
	'f': "\033[97m", // Blanco
	'l': "\033[1m",  // Negrita
	'm': "\033[9m",  // Tachado
	'n': "\033[4m",  // Subrayado
	'o': "\033[3m",  // Itálica
	'r': "\033[0m",  // Reset
}

// ColorizeMOTD convierte códigos § de Minecraft a ANSI para terminal
func ColorizeMOTD(motd string) string {
	var result strings.Builder
	runes := []rune(motd)
	hasColor := false
	
	for i := 0; i < len(runes); i++ {
		if runes[i] == '§' && i+1 < len(runes) {
			code := runes[i+1]
			if ansi, ok := MinecraftColorMap[code]; ok {
				if color.NoColor {
					// Solo eliminamos el código si es un código válido y estamos en modo sin color
					i++ 
					continue
				}
				result.WriteString(ansi)
				hasColor = true
				i++ // Saltamos el código de color
				continue
			}
		}
		result.WriteRune(runes[i])
	}
	
	if hasColor && !color.NoColor {
		// Aseguramos que el color se resetee al final
		result.WriteString("\033[0m")
	}
	return result.String()
}
