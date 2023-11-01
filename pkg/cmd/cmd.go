package cmd

import (
	"github.com/s-mahm/instaOS/pkg/cmd/ubuntu"
	"github.com/s-mahm/instaOS/pkg/cmd/util/templates"
	"github.com/spf13/cobra"
)

func NewInstaOSCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "instaOS",
		Short:   "instaOS creates a bootable flash drive for unattended operating system installation",
	}

	cmd.AddCommand(ubuntu.NewCmdUbuntu())

	templates.GenerateTemplates(cmd)

	return cmd
}
