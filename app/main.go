package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {

	builtin := []string{"exit", "type", "echo"}

	for {
		fmt.Print("$ ")
		command, err:= bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Print("error", err)
		}

		cleanCommand := strings.TrimSpace(command)

		if cleanCommand == "exit" {
			return
		}
		if strings.HasPrefix(command, "type ") {
			if slices.Contains(builtin, cleanCommand[5:]) {
				fmt.Print(cleanCommand[5:], " is a shell builtin\n")
			} 
			if path, err := exec.LookPath(cleanCommand[5:]); err!= nil {
				fmt.Print(cleanCommand[5:], "is", path + "\n")
			} else {
				fmt.Print(cleanCommand[5:], ": not found\n")
			}
			continue
		}
		if strings.HasPrefix(command, "echo ") {
			fmt.Print(command[5:])
		} else {
			fmt.Print(cleanCommand, ": command not found\n")
		}		
	}
}
