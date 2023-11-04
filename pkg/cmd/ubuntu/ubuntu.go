package ubuntu

import (
	"errors"
	"fmt"
	"os"

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
	cmd.Flags().StringVarP(&o.Source, "source", "s", "", "Optional source iso to use")
	cmd.Flags().StringVarP(&o.Version, "version", "v", "22.04", "Optional Ubuntu version to install and flash\nCannot be used with source flag")
	cmd.MarkFlagsMutuallyExclusive("source", "version")
	return cmd
}

func (o *UbuntuOptions) Complete() error {
	switch o.Version {
	case "22.04", "20.04":
	default:
		return fmt.Errorf("invalid version %s", o.Version)
	}
	if err := flash.IsValidFlashDevice(o.Flash); err != nil {
		return err
	}
	_, err := flash.GetFlashDeviceInfo(o.Flash)
	if err != nil {
		return err
	}
	o.Destination = defaultDir()
	if _, err = os.ReadDir(o.Destination); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(o.Destination, os.ModePerm); err != nil {
			return fmt.Errorf("trying to create directory %s: %s", o.Destination, err)
		}
	} else {
		return err
	}

	return nil
}

func (o *UbuntuOptions) Run(args []string) error {
	iso_filename, err := DownloadUbuntuISO(o.Version, o.Destination)
	if err != nil {
		return err
	}
	if err := VerifyISO(iso_filename, o.Version, o.Destination); err != nil {
		return err
	}
	return nil
}

func defaultDir() string {
	if current_dir, err := os.Getwd(); err == nil {
		return fmt.Sprintf("%s/files", current_dir)
	}
	return ""
}
