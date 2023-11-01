package ubuntu

import (
	"fmt"
	"github.com/s-mahm/instaOS/pkg/cmd/util"
	"github.com/spf13/cobra"
)

type UbuntuOptions struct {
	Device      string
	Destination string
	ExtraFile   string
	MetaData    string
	NoMD5       bool
	UserData    string
	Verbose     bool
	Version     string
	Source      string
}

func NewUbuntuOptions() *UbuntuOptions {
	return &UbuntuOptions{}
}

func NewCmdUbuntu() *cobra.Command {
	o := NewUbuntuOptions()
	cmd := &cobra.Command{
		Use:   "ubuntu",
		Short: "Create and flash Ubuntu unattended installation ISO",
		Args: func(cmd *cobra.Command, args []string) error {
			if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
				util.CheckErr(fmt.Errorf("Run %s", err))
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete())
			util.CheckErr(o.Run(args))
		},
	}
	return cmd
}

func (o *UbuntuOptions) Complete() error {
	return nil
}

func (o *UbuntuOptions) Run(args []string) error {
	return nil
}
