package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"flag"
)

// Check if a string is empty
func isEmptyString(query string) bool {
	return query == ""
}

// Write a string line by line
func write(data string) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		fmt.Println(line)
	}
}

// Run a shell command and return the output
func runCommandOld(command string) (string, error) {
        parts := strings.Fields(command)
        cmd := exec.Command(parts[0], parts[1:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Run a shell command and return the output
func runCommand(gitprog string,gitcmd []string) (string, error) {
        cmd := exec.Command(gitprog, gitcmd[0:]...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// Get the current branch of the Git repository
func getCurrentBranch() (string, error) {
	output, err := runCommandOld("git symbolic-ref --short HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v", err)
	}
	return strings.TrimSpace(output), nil
}

// Get the top-level directory of the Git repository
func getTopLevelDir() (string, error) {
	output, err := runCommandOld("git rev-parse --show-toplevel")
	if err != nil {
		return "", fmt.Errorf("failed to get top-level directory: %v", err)
	}
	return strings.TrimSpace(output), nil
}

// Check if a Git command is allowed based on the current branch
func isCommandAllowed(gitBranch string) bool {
	allowedBranches := []string{"trunk", "main", "master"}
	for _, branch := range allowedBranches {
		if gitBranch == branch {
			return false
		}
	}
	return true
}

// Perform a Git fetch and pull
func update() error {
	_, err := runCommandOld("git fetch --all -p -t")
	if err != nil {
		return fmt.Errorf("failed to fetch: %v", err)
	}
	_, err = runCommandOld("git pull")
	if err != nil {
		return fmt.Errorf("failed to pull: %v", err)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: go run main.go <command> [args]")
	}

	// Make sure Git is available
	_, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Git is not installed")
	}

	// cmd := os.Args[1]

	// Set up the global variables
	gitDir, err := getTopLevelDir()
	if err != nil {
		log.Fatalf("failed to get top-level directory: %v", err)
	}
	gitBranch, err := getCurrentBranch()
	if err != nil {
		log.Fatalf("failed to get current branch: %v", err)
	}

	cmd := []string{os.Args[1]}
	// Process the command
	switch os.Args[1] {
	case "check":
		fmt.Printf("%s in %s\n", gitBranch, gitDir)
		cmd = []string{"rev-parse","HEAD"}
	case "checkout", "co":
		cmd = []string{"checkout"}
	case "update", "u":
		err := update()
		if err != nil {
			log.Fatalf("failed to update: %v", err)
		}
		os.Exit(0)
	case "log", "l":
		cmd = []string{"log", "--oneline", "--graph"}
	case "add", "a":
		if !isCommandAllowed(gitBranch) {
			log.Fatalf("command not allowed in %s", gitBranch)
		}
		cmd = []string{"add"}
	case "commit", "c":
		if !isCommandAllowed(gitBranch) {
			log.Fatalf("command not allowed in %s", gitBranch)
		}
		cmd = []string{"commit"}
	case "push", "p":
		if !isCommandAllowed(gitBranch) {
			log.Fatalf("command not allowed in %s", gitBranch)
		}
		cmd = []string{"push"}
	case "originpush", "op", "og":
		if !isCommandAllowed(gitBranch) {
			log.Fatalf("command not allowed in %s", gitBranch)
		}
		cmd = []string{"push", "-u", "origin", "$git_branch"}
	case "current_hash", "hash":
		cmd = []string{"rev-parse", "HEAD"}
	case "grep", "gg":
		cmd = []string{"grep","-n"}
	case "clone":
		if isEmptyString(gitDir) && isEmptyString(gitBranch) {
			cmd = []string{os.Args[1]}
		} else {
			log.Fatal("No!")
		}
	default:
		cmd = []string{os.Args[1]}
	}

	// Run the Git command
	flag.Parse()
	output, err := runCommand("git",append(cmd,flag.Args()[1:]...))
	if err != nil {
		fmt.Fprintln(os.Stderr,output)
		log.Fatalf("failed to uun command: %v", err)
	}

	// Write the output
	write(output)

	os.Exit(0)
}
