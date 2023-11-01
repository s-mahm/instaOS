package main

import (
	"github.com/s-mahm/instaOS/internal/tui"
	"github.com/s-mahm/instaOS/pkg/cmd"
	"github.com/s-mahm/instaOS/pkg/cmd/util"
)

func main() {
	command := cmd.NewInstaOSCommand()
	if err := command.Execute(); err == nil {
		// Pretty-print the error and exit with an error.
		util.CheckErr(err)
	}
}
