package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const filePath = "/tmp/branches"

type Projects map[string]string

func run(args ...string) (output string, err error) {
	cmd := exec.Command(args[0], args[1:]...)
	raw, err := cmd.Output()
	if err == nil {
		output = strings.TrimSpace(string(raw))
	}
	return output, err
}

func read() (projects Projects) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return make(Projects)
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error: failed to read %s\n", filePath)
		os.Exit(1)
	}
	err = json.Unmarshal(content, &projects)
	return projects
}

func currentProject() string {
	output, err := run("git", "rev-parse", "--show-toplevel")
	if err != nil {
		fmt.Println("error: not a git repository")
		os.Exit(1)
	}
	return output
}

func currentBranch() string {
	output, err := run("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		fmt.Println("error: not a git repository")
		os.Exit(1)
	}
	return output
}

// Store current working branch.
func Store() {
	projects := read()
	project := currentProject()
	branch := currentBranch()
	projects[project] = branch

	bytes, err := json.Marshal(projects)
	err = ioutil.WriteFile(filePath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}

// Restore to last known working branch.
func Restore() {
	projects := read()
	project := currentProject()
	branch := projects[project]
	_, err := run("git", "checkout", branch)
	if err != nil {
		fmt.Println("error: invalid branch")
		os.Exit(1)
	}
}

func main() {
	if currentBranch() != "master" {
		Store()
	} else {
		Restore()
	}
}
