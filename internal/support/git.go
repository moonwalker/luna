package support

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	gitRangeRx = "[0-9a-f]{7,40}...[0-9a-f]{7,40}"
	gitPrevCmd = "git rev-parse HEAD~1"
	gitCurrCmd = "git rev-parse HEAD"
	gitDiffCmd = "git diff --name-only %s"
)

func getGitDiff(compareRange string) []string {
	commitRange := getCommitRange(compareRange)
	diff := run(fmt.Sprintf(gitDiffCmd, commitRange))
	return strings.Split(diff, "\n")
}

func getCommitRange(compareRange string) string {
	if compareRange == "" {
		prevCommit := run(gitPrevCmd)
		currCommit := run(gitCurrCmd)
		compareRange = fmt.Sprintf("%s...%s", prevCommit, currCommit)
	}
	fmt.Println(compareRange)
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
