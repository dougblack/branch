package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const filePath = "/tmp/working"

type Projects map[string]string

func (p Projects) branch(project string) string {
	return p[project]
}

func On(branch string) {
	projects := read()
	projects[currentProject()] = branch
	dl, err := json.Marshal(projects)
	err = ioutil.WriteFile(filePath, dl, 0644)
	if err != nil {
		panic(err)
	}
}

func read() (projects Projects) {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return make(Projects)
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, &projects)
	return projects
}

func resolve(project string) (branch string) {
	return read().branch(project)
}

func currentProject() (project string) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	raw, err := cmd.Output()
	if err != nil {
		fmt.Println("error: not a git repository")
		os.Exit(1)
	}
	output := strings.TrimSpace(string(raw))
	return output
}

func Switch() {
	branch := resolve(currentProject())
	cmd := exec.Command("git", "checkout", branch)
	err := cmd.Run()
	if err != nil {
		fmt.Println("error: invalid branch")
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) > 1 {
		branch := os.Args[1]
		On(branch)
	} else {
		Switch()
	}
}
