package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
)

type Redirect struct {
	fd   int    // 1 = stdout, 2 = stderr
	path string
	append bool
}

func findRedirect(parts []string) (cmdParts []string, redirect *Redirect) {
	for i, p := range parts {
		var fd int
		var appendMode bool

		switch p {
		case ">", "1>":
			fd = 1
		case ">>", "1>>":
			fd = 1
			appendMode = true
		case "2>":
			fd = 2
		case "2>>":
			fd = 2
			appendMode = true
		default:
			continue
		}

		if i+1 >= len(parts) {
			fmt.Println("syntax error: missing file after redirect")
			return parts, nil
		}
		return parts[:i], &Redirect{fd: fd, path: parts[i+1], append: appendMode}
	}
	return parts, nil
}

func openRedirectFile(r *Redirect) (*os.File, error) {
	flags := os.O_WRONLY | os.O_CREATE
	if r.append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}
	return os.OpenFile(r.path, flags, 0644)
}

func runCommand(parts []string, redirect *Redirect) {
	if len(parts) == 0 {
		return
	}

	cmdName := parts[0]
	args := parts[1:]

	path, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Println(cmdName + ": command not found")
		return
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Path = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if redirect != nil {
		f, err := openRedirectFile(redirect)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()

		switch redirect.fd {
		case 1:
			cmd.Stdout = f
		case 2:
			cmd.Stderr = f
		}
	}
	cmd.Run()

	// if err := cmd.Run(); err != nil {
	// 	fmt.Fprintf(os.Stderr, "%v\n", err)
	// }
}

func main() {
	builtin := []string{"exit", "type", "echo", "pwd", "cd"}

	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Print("error", err)
		}

		cleanCommand := strings.TrimSpace(command)
		parts := parseArgs(cleanCommand)

		if len(parts) == 0 {
			continue
		}

		if cleanCommand == "exit" {
			return
		}

		if parts[0] == "type" {
			if len(parts) < 2 {
				fmt.Println("type: missing argument")
				continue
			}
			arg := parts[1]
			if slices.Contains(builtin, arg) {
				fmt.Println(arg + " is a shell builtin")
			} else if path, err := exec.LookPath(arg); err == nil {
				fmt.Println(arg + " is " + path)
			} else {
				fmt.Println(arg + ": not found")
			}
			continue
		}

		if parts[0] == "pwd" {
			cwd, err := os.Getwd()
			if err != nil {
				fmt.Println("error:", err)
			} else {
				fmt.Println(cwd)
			}
			continue
		}

		if parts[0] == "cd" {
			if len(parts) < 2 || parts[1] == "~" {
				if homedir, err := os.UserHomeDir(); err == nil {
					os.Chdir(homedir)
				}
				continue
			}
			info, err := os.Stat(parts[1])
			if err != nil || !info.IsDir() {
				fmt.Printf("cd: %s: No such file or directory\n", parts[1])
				continue
			}
			os.Chdir(parts[1])
			continue
		}

		// Parse redirect before running anything
		cmdParts, redirect := findRedirect(parts)

		if len(cmdParts) > 0 && cmdParts[0] == "echo" {
			output := strings.Join(cmdParts[1:], " ") + "\n"
			if redirect != nil {
				f, err := openRedirectFile(redirect)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				if redirect.fd == 1 {
					fmt.Fprint(f, output)
				} else {
					fmt.Print(output) 
				}
			} else {
				fmt.Print(output)
			}
			continue
		}

		runCommand(cmdParts, redirect)
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