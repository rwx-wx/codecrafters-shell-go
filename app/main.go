package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)


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
			} else if path, err := exec.LookPath(cleanCommand[5:]); err == nil {
				fmt.Print(cleanCommand[5:], " is ", path + "\n")
			} else {
				fmt.Print(cleanCommand[5:], ": not found\n")
			}
			continue
		}
		if strings.HasPrefix(command, "echo ") {
			fmt.Print(command[5:])
		} else {

			parts := strings.Fields(cleanCommand)

			if len(parts) > 0 {
				cmdName := parts[0]
				args := parts[1:]

				if path, err:= exec.LookPath(cmdName); err == nil {
					cmd := exec.Command(cmdName, args...) // variadics very cool
					cmd.Path = path
					out, err:= cmd.CombinedOutput()
					if err != nil {
						fmt.Println(cleanCommand, ": command not found\n")
					} else {
						fmt.Print(string(out))
					}
				} else {
					fmt.Print(cleanCommand, ": command not found\n")
				}
			}
		}		
	}
}
