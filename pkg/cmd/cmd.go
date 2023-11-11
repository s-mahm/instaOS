package cmd

import (
	"github.com/s-mahm/instaOS/pkg/cmd/ubuntu"
	"github.com/s-mahm/instaOS/pkg/cmd/util/templates"
	"github.com/spf13/cobra"
)

type OS int64

const (
	Linux OS = iota
	Mac
	Windows
)

func (os OS) String() string {
	switch os {
	case Linux:
		return "linux"
	case Mac:
		return "mac"
	case Windows:
		return "windows"
	}
	return "unknown"
}

type Distro struct {
	OS             OS
	Name           string
	Version        string
	ISOSource      string
	ISODestination string
}

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
