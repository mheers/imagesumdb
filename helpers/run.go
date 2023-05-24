package helpers

import (
	"bytes"
	"fmt"
	"os/exec"
)

func Run(name string, args ...string) error {
	stdout, stderr, err := RunResult(name, args...)
	if err != nil {
		fmt.Println(stdout)
		fmt.Println(stderr)
		return err
	}
	return nil
}

func RunResult(name string, args ...string) (string, string, error) {
	c := exec.Command(name, args...)

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	c.Stdout = &stdOut
	c.Stderr = &stdErr

	err := c.Run()
	if err != nil {
		return "", "", fmt.Errorf("%v: %s", err, stdErr.String())
	}

	return stdOut.String(), stdErr.String(), nil
}
