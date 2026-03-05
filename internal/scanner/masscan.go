package scanner

import (
	"bufio"
	"fmt"
	"os" // Importante para os.Stderr
	"os/exec"
	"strings"
)

type MasscanResult struct {
	IP    string `json:"ip"`
	Ports []struct {
		Port int `json:"port"`
	} `json:"ports"`
}

// BuildArguments constructs the arguments for masscan
func BuildArguments(ipRange string, rate string, port int, excludeFile string) []string {
	args := []string{
		ipRange,
		"-p", fmt.Sprintf("%d", port),
		"--rate", rate,
		"-oL", "-",
	}

	if excludeFile != "" {
		args = append(args, "--excludefile", excludeFile)
	} else if ipRange == "0.0.0.0/0" {
		args = append(args, "--exclude", "255.255.255.255,127.0.0.0/8,0.0.0.0/8,224.0.0.0/4")
	}
	return args
}

func Run(ipRange string, rate string, port int, excludeFile string, ipChan chan<- string) error {
	args := BuildArguments(ipRange, rate, port, excludeFile)

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
			line := scanner.Text()
			// Formato: "open tcp 25565 1.2.3.4 123456789"
			parts := strings.Fields(line)
			if len(parts) >= 4 && parts[0] == "open" {
				ipChan <- parts[3]
			}
		}
		_ = cmd.Wait()
	}()

	return nil
}
