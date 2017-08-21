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

func getGitDiff(compareRange string) []string {
	commitRange := getCommitRange(compareRange)
	diff := run(fmt.Sprintf(gitDiffCmd, commitRange))
	return strings.Split(diff, "\n")
}

func getCommitRange(compareRange string) string {
	if compareRange == "" {
		baseCommit := run(gitBaseCmd)
		lastCommit := run(gitLastCmd)
		compareRange = fmt.Sprintf("%s...%s", baseCommit, lastCommit)
	}
	re := regexp.MustCompile(gitRangeRx)
	return re.FindString(compareRange)
}

func run(command string) string {
	cmd := MakeCmd(command, "")
	out, err := cmd.CombinedOutput()
	res := strings.TrimSpace(string(out))

	if err != nil {
		fmt.Println(res)
		os.Exit(0)
	}

	return res
}
