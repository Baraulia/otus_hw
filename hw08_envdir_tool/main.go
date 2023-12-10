package main

import (
	"fmt"
	"os"
)

func main() {
	path := os.Args[1]
	command := os.Args[2:]
	mapEnvs, err := ReadDir(path)
	if err != nil {
		fmt.Printf("error while reading dir: %s", err)
		return
	}

	code := RunCmd(command, mapEnvs)

	if code != 0 {
		fmt.Printf("command execution error. code:%d", code)
	}
}
