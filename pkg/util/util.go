package util

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

func ReplaceTextInFiles(paths []string, search_text string, replace_text string) error {
	for _, path := range paths {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		newContents := strings.Replace(string(read), search_text, replace_text, -1)
		err = os.WriteFile(path, []byte(newContents), 0)
		if err != nil {
			return err
		}
	}
	return nil
}

func ReplaceLineInFiles(paths []string, search_pattern string, replacement_line string) error {
	pattern, err := regexp.Compile(fmt.Sprintf(`%s`, search_pattern))
	if err != nil {
		return err
	}
	for _, path := range paths {
		read, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(read), "\n")
		for i, line := range lines {
			if matched := pattern.MatchString(line); matched == true {
				lines[i] = replacement_line
			}
		}
		output := strings.Join(lines, "\n")
		err = os.WriteFile(path, []byte(output), 0)
		if err != nil {
			return err
		}
	}
	return nil
}
