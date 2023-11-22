package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func handleSession(session ssh.Session) {
	term := terminal.NewTerminal(session, "lol>>")
	io.WriteString(session, fmt.Sprintf("Successfull login, %s\n", session.User()))
	for {
		command, err := term.ReadLine()
		if err != nil {
			log.Println(err)
		}
		log.Printf("Received %s: %s", session.User(), command)
		selectCommand(session, command)

	}
}

func selectCommand(session ssh.Session, command string) {
	switch command {
	case "exit":
		return
	case "ping":
		output, err := runCommand("ping", "-c", "4", "onliner.by")
		if err != nil {
			io.WriteString(session, fmt.Sprintf("Error ping: %s\n", err))
		} else {
			io.WriteString(session, string(output))
		}
	default:
		args := strings.Fields(command)
		if len(args) > 0 {
			switch args[0] {
			case "mkdir":
				if len(args) < 2 {
					io.WriteString(session, "You must specify dir name\n")
				} else {
					err := os.Mkdir(args[1], 0755)
					if err != nil {
						io.WriteString(session, fmt.Sprintf("Failed to create dir: %s\n", err))
					} else {
						io.WriteString(session, "Dir created\n")
					}
				}
			case "rmdir":
				if len(args) < 2 {
					io.WriteString(session, "You need to specify dir to delete\n")
				} else {
					err := os.Remove(args[1])
					if err != nil {
						io.WriteString(session, fmt.Sprintf("Failed to delete: %s\n", err))
					} else {
						io.WriteString(session, "Dir successfully deleted\n")
					}
				}
			case "ls":
				files, err := os.ReadDir(".")
				if err != nil {
					io.WriteString(session, fmt.Sprintf("Error reading dir: %s\n", err))
				} else {
					var output strings.Builder
					for _, file := range files {
						output.WriteString(file.Name())
						output.WriteString("\n")
					}
					io.WriteString(session, output.String())
				}
			case "mv":
				if len(args) < 3 {
					io.WriteString(session, "You need to specify src and dst paths\n")
				} else {
					err := os.Rename(args[1], args[2])
					if err != nil {
						io.WriteString(session, fmt.Sprintf("Error while moving: %s\n", err))
					} else {
						io.WriteString(session, "File successfully moved\n")
					}
				}
			case "rm":
				if len(args) < 2 {
					io.WriteString(session, "You need to specify file to delete\n")
				} else {
					err := os.Remove(args[1])
					if err != nil {
						io.WriteString(session, fmt.Sprintf("Error while deleting: %s\n", err))
					} else {
						io.WriteString(session, "Successfully deleted\n")
					}
				}
			default:
				io.WriteString(session, "Wrong command\n")
			}
		} else {
			io.WriteString(session, "Wrong command\n")
		}
	}
}

func runCommand(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("command execution failed with: %s", err)
	}
	return output, nil
}

func main() {
	server := ssh.Server{
		Addr: "localhost:8000",
		Handler: func(s ssh.Session) {
			handleSession(s)
		},
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			log.Printf("User authenication %s", ctx.User())
			return true
		},
	}
	log.Println("SSH started at port 8000")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
