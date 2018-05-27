package agilemarkdown

import (
	"bytes"
	"os/exec"
	"strings"
)

func AddIdea(contentDirectory, ideaTitle string) (ideaPath string, ideaContent string, err error) {
	args := []string{"create-idea", "--simulate", ideaTitle}
	out, err := runAgileMarkdownCommand(contentDirectory, args)
	if len(out) > 0 {
		ideaPath = out[0]
		ideaContent = strings.Join(out[1:], "\n")
	}
	return ideaPath, ideaContent, err
}

func runAgileMarkdownCommand(workDir string, args []string) ([]string, error) {
	cmd := exec.Command("./agilemarkdown", args...)
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
	return strings.Split(strings.TrimSpace(strings.Join(lines, "\n")), "\n"), err
}
