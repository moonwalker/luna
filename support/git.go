package support

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	gitRangeRx = "[0-9a-f]{7,40}...[0-9a-f]{7,40}"
	gitLastCmd = "git rev-parse --short HEAD"
	gitBaseCmd = "git rev-parse --short HEAD~2"
	gitDiffCmd = "git diff %s --name-only"
)

func GetGitDiff(compareURL string) string {
	commitRange := getCommitRange(compareURL)
	return run(fmt.Sprintf(gitDiffCmd, commitRange))
}

func getCommitRange(compareURL string) string {
	if compareURL == "" {
		baseCommit := run(gitBaseCmd)
		lastCommit := run(gitLastCmd)
		compareURL = fmt.Sprintf("%s...%s", baseCommit, lastCommit)
	}
	re := regexp.MustCompile(gitRangeRx)
	return re.FindString(compareURL)
}

func run(command string) string {
	cmd := makeCmd(command, "")
	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))

	if err != nil {
		fmt.Println(res)
		os.Exit(0)
	}

	return res
}
