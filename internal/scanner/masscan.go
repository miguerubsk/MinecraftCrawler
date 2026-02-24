package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

// Estructura para parsear el JSON de Masscan línea a línea
type MasscanResult struct {
	IP    string `json:"ip"`
	Ports []struct {
		Port int `json:"port"`
	} `json:"ports"`
}

// Run lanza masscan y envía las IPs encontradas al canal de análisis
func Run(ipRange string, rate string, port int, exclusions string, ipChan chan<- string) error {
	// Argumentos: -oJ - (salida JSON por stdout), --rate, etc.
	args := []string{
		ipRange,
		"-p", fmt.Sprintf("%d", port),
		"--rate", rate,
		"--exclude 255.255.255.255", // Excluir broadcast por defecto
		"-oJ", "-", // Salida JSON a stdout
	}

	if exclusions != "" {
		args = append(args, "--excludefile", exclusions)
	}

	cmd := exec.Command("masscan", args...)
	cmd.Stderr = os.Stderr // Redirigir stderr para ver logs de masscan
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Ejecutar procesamiento en el hilo actual (bloqueante)
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Bytes()
		// Masscan a veces envía comas o corchetes al inicio/final del JSON
		if len(line) < 10 || line[0] == '[' || line[0] == ']' {
			continue
		}
		// Eliminar la coma final si existe
		if line[len(line)-1] == ',' {
			line = line[:len(line)-1]
		}

		var res MasscanResult
		if err := json.Unmarshal(line, &res); err == nil {
			if len(res.Ports) > 0 {
				ipChan <- res.IP
			}
		}
	}
	return cmd.Wait()
}

// ReadFromFile lee un archivo JSON generado previamente por masscan
func ReadFromFile(filePath string, ipChan chan<- string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Masscan genera un array JSON [{}, {}]. Usamos un decoder para no cargar todo en RAM
	decoder := json.NewDecoder(file)
	
	// Leer el corchete inicial
	if _, err := decoder.Token(); err != nil { return err }

	for decoder.More() {
		var res MasscanResult
		if err := decoder.Decode(&res); err == nil {
			ipChan <- res.IP
		}
	}
	return nil
}