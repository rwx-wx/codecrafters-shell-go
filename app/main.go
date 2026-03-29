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

	builtin := []string{"exit", "type", "echo", "pwd", "cd"}

	for {
		fmt.Print("$ ")
		command, err:= bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Print("error", err)
		}

		cleanCommand := strings.TrimSpace(command)
		parts := parseArgs(cleanCommand) // was strings.Fields(cleanCommand)

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
		if strings.HasPrefix(command, "pwd") {
			cwd, err:= os.Getwd() 
			if err != nil {
				fmt.Print("error")
			}
			fmt.Println(cwd)
			continue
		}
		if strings.HasPrefix(command, "cd") {
			parts := strings.Fields(command)

			args := parts[1:]

			if args[0] == "~" {
				if homedir, err :=os.UserHomeDir(); err == nil {
					os.Chdir(homedir) 
					continue
				}
			}
			
			dirName, err := os.Stat(args[0])
			if err != nil {
				fmt.Print("cd: ", args[0], ": No such file or directory\n")
				continue
			}

			if dirName.IsDir() {
				os.Chdir(args[0])
			}
			continue
		}
		if parts[0] == "echo" {
			fmt.Println(strings.Join(parts[1:], " "))
			continue
		} else {
		// if strings.HasPrefix(command, "echo ") {
		// 	fmt.Print(command[5:])
		// } else {

			if len(parts) > 0 {
				cmdName := parts[0]
				
				args := parts[1:]


				if path, err:= exec.LookPath(cmdName); err == nil {
					cmd := exec.Command(cmdName, args...) // variadics very cool
					cmd.Path = path
					out, err:= cmd.CombinedOutput()
					if err != nil {
						fmt.Print(cleanCommand, ": command not found\n")
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

func parseArgs(input string) []string {
	var args []string
	var current strings.Builder
	inSingle := false
	inDouble := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		switch {
		case ch == '\'' && !inDouble:
			inSingle = !inSingle
		case ch == '"' && !inSingle:
			inDouble = !inDouble
		case ch == '\\' && inDouble && i+1 < len(input):
			// Inside double quotes, backslash only escapes specific chars
			next := input[i+1]
			if next == '"' || next == '\\' || next == '$' || next == '\n' {
				current.WriteByte(next)
				i++
			} else {
				current.WriteByte(ch)
			}
		case ch == '\\' && !inSingle && !inDouble && i+1 < len(input):
			// Outside quotes, backslash escapes the next char literally
			current.WriteByte(input[i+1])
			i++
		case (ch == ' ' || ch == '\t') && !inSingle && !inDouble:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(ch)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args
}