package scanner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os" // Importante para os.Stderr
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
		"-oJ", "-", 
	}

	if excludeFile != "" {
		args = append(args, "--excludefile", excludeFile)
	} else if ipRange == "0.0.0.0/0" {
		args = append(args, "--exclude", "255.255.255.255,127.0.0.0/8,0.0.0.0/8,224.0.0.0/4")
	}

	cmd := exec.Command("masscan", args...)
	
	// Redirigimos el stderr de masscan al stderr de nuestro programa 
	// Esto mostrará el progreso "Rate:..., 10.00% done..." en la terminal.
	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	go func() {
		defer close(ipChan)
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) < 10 || line[0] == '[' || line[0] == ']' {
				continue
			}
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
	}()

	return nil
}
