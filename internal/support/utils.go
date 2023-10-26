package support

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func StringListContains(list []string, s string) bool {
	for _, e := range list {
		if e == s {
			return true
		}
	}
	return false
}

func FileExists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func IsYaml(f string) bool {
	data, err := os.ReadFile(f)
	if err != nil {
		return false
	}

	var out interface{}
	err = yaml.Unmarshal([]byte(data), &out)
	return err == nil
}

func ExpandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home := os.Getenv("HOME")
		return home + path[1:], nil
	}

	var err error
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if path == "." {
		return wd, nil
	}

	// relative path specified
	if strings.HasPrefix(path, "./") {
		return wd + path[1:], nil
	}

	// assume relative path, but not specified
	if !strings.HasPrefix(path, "/") {
		return fmt.Sprintf("%s/%s", wd, path), nil
	}

	return path, nil
}
