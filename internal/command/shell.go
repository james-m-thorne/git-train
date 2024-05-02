package command

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os/exec"
)

var PrintBold = color.New(color.Bold).PrintlnFunc()

// Exec is a simple wrapper for exec.Command("sh", "-c", ...).
func Exec(command string) error {
	var stdout, stderr bytes.Buffer
	fmt.Println()
	PrintBold(command)
	cmd := exec.Command("sh", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err, stderr.String())
	}

	fmt.Print(stdout.String())
	return nil
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
