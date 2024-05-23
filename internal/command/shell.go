package command

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"os/exec"
	"strings"
)

var PrintBold = color.New(color.Bold).PrintlnFunc()

func PrintFatalError(format string, a ...any) {
	color.Red(format, a...)
	os.Exit(1)
}

// Exec is a simple wrapper for exec.Command("sh", "-c", ...).
func Exec(command string, exitOnError bool, dryRun bool) {
	fmt.Println()
	PrintBold(command)
	if dryRun {
		return
	}

	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		if exitOnError {
			PrintFatalError("%s: %s", command, err)
		} else {
			color.Red("%s: %s", command, err)
		}
	}
}

func GetOutput(command string) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%s: %s", err, stderr.String())
	}

	output := stdout.String()
	output = output[:len(output)-1] // trim the new line char
	return output, nil
}

func GetOutputFatal(shell string) string {
	result, err := GetOutput(shell)
	if err != nil {
		PrintFatalError("%s: %s", shell, err)
	}
	return result
}

func YesNoInput(question string) (bool, error) {
	fmt.Printf("\n%s (y/n)\n", question)

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y", nil
}
