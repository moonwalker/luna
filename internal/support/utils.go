package support

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

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
