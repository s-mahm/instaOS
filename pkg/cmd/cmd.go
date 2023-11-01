package cmd

import (
	"github.com/spf13/cobra"
)

func NewInstaOSCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmd := &cobra.Command{
		Version: "v0.0.1",
		Use:     "instaOS",
		Short:   "instaOS creates a bootable flash drive for unattended operating system installation",
		Long: "instaOS creates a bootable flash drive for unattended operating system installation\n\n" +
			"Find more information at: " +
			"https://github.com/s-mahm/instaOS",
	}
	cmd.SilenceUsage = true
	return cmd
}
