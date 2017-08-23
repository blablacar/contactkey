package utils

import (
	"io"
	"os/exec"
)

func StreamCombinedOutput(c *exec.Cmd) (io.ReadCloser, error) {
	reader, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}
	c.Stderr = c.Stdout

	c.Start()
	return reader, nil
}
