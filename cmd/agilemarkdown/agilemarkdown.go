package agilemarkdown

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var syncMutex sync.Mutex

func AddIdea(contentDirectory, ideaTitle string, jwtToken string) (ideaPath string, ideaContent string, err error) {
	user, _ := getUserInfoFromJwtToken(contentDirectory, jwtToken)

	args := []string{"create-idea", "--simulate"}
	if user != "" {
		args = append(args, "--user", user)
	}
	args = append(args, ideaTitle)
	out, err := runAgileMarkdownCommand(contentDirectory, args)
	if len(out) > 0 {
		ideaPath = out[0]
		ideaContent = strings.Join(out[1:], "\n") + "\n"
	}
	return ideaPath, ideaContent, err
}

func AddStory(contentDirectory, projectName, storyTitle string, jwtToken string) (storyPath string, storyContent string, err error) {
	user, _ := getUserInfoFromJwtToken(contentDirectory, jwtToken)

	args := []string{"create-item", "--simulate"}
	if user != "" {
		args = append(args, "--user", user)
	}
	args = append(args, storyTitle)
	out, err := runAgileMarkdownCommand(filepath.Join(contentDirectory, projectName), args)
	if len(out) > 0 {
		storyPath = out[0]
		storyContent = strings.Join(out[1:], "\n") + "\n"
	}
	return storyPath, storyContent, err
}

func Sync(contentDirectory, jwtToken string) (string, error) {
	syncMutex.Lock()
	defer syncMutex.Unlock()

	userName, userEmail := getUserInfoFromJwtToken(contentDirectory, jwtToken)

	args := []string{"sync"}
	if userName != "" {
		userEmail = getUserEmail(userName, userEmail)
		args = append(args, "--author", fmt.Sprintf("%s <%s>", userName, userEmail))

	}
	out, err := runAgileMarkdownCommand(contentDirectory, args)
	return strings.Join(out, "\n"), err
}

func CreateUserIfNotExist(workDir, jwtToken string) {
	getUserInfoFromJwtToken(workDir, jwtToken)
}

func createUserIfNotExist(workDir, name, email string) {
	args := []string{"create-user", "--name", name, "--email", email}
	runAgileMarkdownCommand(workDir, args)
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
	output := strings.TrimSpace(strings.Join(lines, "\n"))
	return strings.Split(output, "\n"), createVerboseError(err, output)
}

func getUserInfoFromJwtToken(workDir, jwtToken string) (name, email string) {
	token, _ := jwt.Parse(jwtToken, nil)
	if token == nil {
		return "", ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok {
		return "", ""
	} else {

		var name, email string
		if rawName, ok := claims["name"]; ok {
			name = rawName.(string)
		}
		if rawEmail, ok := claims["email"]; ok {
			email = rawEmail.(string)
		}

		if name == "" {
			if rawSub, ok := claims["sub"]; ok {
				name = rawSub.(string)
			}
		}

		if email == "" && name != "" {
			email = getUserEmail(name, email)
		}

		if name != "" {
			createUserIfNotExist(workDir, name, email)
		}

		return name, email
	}
}

func getUserEmail(name, email string) string {
	if email != "" {
		return email
	}
	return fmt.Sprintf("%s@unknown.org", strings.Replace(name, " ", ".", -1))
}

func createVerboseError(err error, out string) error {
	if err == nil {
		return nil
	}

	if out == "" {
		return err
	}

	return fmt.Errorf("%s\n%s", err.Error(), out)
}
