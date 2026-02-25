package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

type MasscanResult struct {
	IP    string `json:"ip"`
	Ports []struct {
		Port int `json:"port"`
	} `json:"ports"`
}

func Run(ipRange string, rate string, port int, excludeFile string, ipChan chan<- string) error {
	args := []string{
		ipRange,
		"-p", fmt.Sprintf("%d", port),
		"--rate", rate,
		"-oJ", "-", // Salida JSON por stdout
	}

	// Si se proporciona un archivo de exclusión, lo añadimos a los argumentos
	if excludeFile != "" {
		args = append(args, "--excludefile", excludeFile)
	}

	cmd := exec.Command("masscan", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			// Masscan a veces envía metadatos en el JSON, los saltamos
			if len(line) < 10 || line[0] == '[' || line[0] == ']' {
				continue
			}
			// Limpiamos la coma final si Masscan la envía
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
		_ = cmd.Wait()
		close(ipChan) // Cerramos el canal cuando masscan termina
	}()

	return nil
}