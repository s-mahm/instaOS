package ubuntu

import (
	"fmt"
	"github.com/s-mahm/instaOS/pkg/cmd/util"
	"github.com/s-mahm/instaOS/pkg/util/flash"
	"github.com/spf13/cobra"
)

type UbuntuOptions struct {
	Flash       string
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
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete())
			util.CheckErr(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Flash, "flash", "f", "", "Required device to flash to (e.g. /dev/sda)")
	cmd.MarkFlagRequired("flash")
	cmd.Flags().StringVarP(&o.Destination, "destination", "d", "", "Optional destination directory to download iso file to")
	cmd.Flags().StringVarP(&o.Source, "source", "s", "", "Optional source iso to use")
	return cmd
}

func (o *UbuntuOptions) Complete() error {
	if err := flash.IsValidFlashDevice(o.Flash); err != nil {
		return err
	}
	_, err := flash.GetFlashDeviceInfo(o.Flash)
	if err != nil {
		return err
	}

	return nil
}

func (o *UbuntuOptions) Run(args []string) error {
	fmt.Println("")

	flash.GetFlashDevices()
	return nil
}
