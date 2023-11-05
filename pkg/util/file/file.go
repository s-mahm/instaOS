package file

import (
	"os"
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
