package agilemarkdown

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AddIdea(contentDirectory, ideaTitle string) (ideaPath string, ideaContent string, err error) {
	args := []string{"create-idea", "--simulate", ideaTitle}
	out, err := runAgileMarkdownCommand(contentDirectory, args)
	if len(out) > 0 {
		ideaPath = out[0]
		ideaContent = strings.Join(out[1:], "\n") + "\n"
	}
	return ideaPath, ideaContent, err
}

func AddStory(contentDirectory, projectName, storyTitle string) (storyPath string, storyContent string, err error) {
	args := []string{"create-item", "--simulate", storyTitle}
	out, err := runAgileMarkdownCommand(filepath.Join(contentDirectory, projectName), args)
	if len(out) > 0 {
		storyPath = out[0]
		storyContent = strings.Join(out[1:], "\n") + "\n"
	}
	return storyPath, storyContent, err
}

func Sync(contentDirectory string) (string, error) {
	args := []string{"sync"}
	out, err := runAgileMarkdownCommand(contentDirectory, args)
	return strings.Join(out, "\n"), err
}

func runAgileMarkdownCommand(workDir string, args []string) ([]string, error) {
	absProgramPath, _ := filepath.Abs(os.Args[0])
	agilemarkdownPath := filepath.Join(filepath.Dir(absProgramPath), "agilemarkdown")
	if _, err := os.Stat(agilemarkdownPath); err != nil {
		agilemarkdownPath = "agilemarkdown"
	}
	cmd := exec.Command(agilemarkdownPath, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	cmd.Dir = workDir
	err := cmd.Run()
	lines := make([]string, 0, 2)
	if out.Len() > 0 {
		lines = append(lines, out.String())
	}
	if stderr.Len() > 0 {
		lines = append(lines, stderr.String())
	}
	return strings.Split(strings.Join(lines, "\n"), "\n"), err
}