package util

import (
	"fmt"
	"os"
	"strings"
)

const DefaultErrorExitCode = 1

var ErrExit = fmt.Errorf("exit")

func CheckErr(err error) {
	checkErr(err, fatal)
}

func fatal(msg string, code int) {
	if len(msg) > 0 {
		// add newline if needed
		if !strings.HasSuffix(msg, "\n") {
			msg += "\n"
		}
		fmt.Fprint(os.Stderr, msg)
	}
	os.Exit(code)
}

func checkErr(err error, handleErr func(string, int)) {
	if err == nil {
		return
	}
	switch {
	case err == ErrExit:
		handleErr("", DefaultErrorExitCode)
	default:
		fmt.Println(err == fmt.Errorf("exit"))
		msg := err.Error()
		if !strings.HasPrefix(msg, "error: ") {
			msg = fmt.Sprintf("error: %s", msg)
		}
		handleErr(msg, DefaultErrorExitCode)
	}
}
