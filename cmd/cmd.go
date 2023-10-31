package cmd

import (
	"github.com/spf13/cobra"
)

func NewInstaOSCommand() *cobra.Command {
	// Parent command to which all subcommands are added.
	cmd := &cobra.Command{
		Use:   "instaos",
		Short: "instaos creates a bootable flash drive for unattended operating system installation",
		Long: "instaos creates a bootable flash drive for unattended operating system installation\n\n" +
			"Find more information at: " +
			"https://github.com/s-mahm/instaOS",
	}
	return cmd
}
