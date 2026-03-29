package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
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
		if strings.HasPrefix(command, "echo ") {
			fmt.Print(command[5:])
		} else {
			fmt.Print(cleanCommand, ": command not found\n")
		}		
	}
}
