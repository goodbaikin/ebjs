package encode

import (
	"bufio"
	"fmt"
	"os/exec"
)

func execute(name string, args []string, logger func(string)) error {
	cmd := exec.Command(name, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		logger(fmt.Sprintf("%s %s", name, scanner.Text()))
	}

	return cmd.Wait()
}
