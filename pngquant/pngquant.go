package pngquant

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

func Compress(finput io.Reader, out io.Writer, speed string) error {
	if err := speedCheck(speed); err != nil {
		return err
	}
	cmd := exec.Command("pngquant", "-", "--speed", speed)
	cmd.Stdin = finput
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func speedCheck(speed string) error {
	speedInt, err := strconv.Atoi(speed)
	if err != nil {
		return err
	}

	if speedInt > 10 {
		return fmt.Errorf("speed cannot exceed value of 10")
	}

	return nil
}
