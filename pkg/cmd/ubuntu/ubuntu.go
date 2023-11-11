package ubuntu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/s-mahm/instaOS/pkg/cmd/distro"
	"github.com/s-mahm/instaOS/pkg/cmd/util"
	"github.com/s-mahm/instaOS/pkg/flash"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type UbuntuDistro struct {
	distro.Distro
	Flash      string
	ExtraFiles string
	UserData   UserData
	NoVerify   bool
}

func NewUbuntuDistro() *UbuntuDistro {
	return &UbuntuDistro{}
}

func NewCmdUbuntu() *cobra.Command {
	o := NewUbuntuDistro()
	cmd := &cobra.Command{
		Use:   "ubuntu",
		Short: "Create and flash Ubuntu unattended installation ISO",
		Run: func(cmd *cobra.Command, args []string) {
			util.CheckErr(o.Complete(cmd))
			util.CheckErr(o.Run(args))
		},
	}
	cmd.Flags().StringVarP(&o.Flash, "flash", "f", "", "Required device to flash to (e.g. /dev/sda)")
	cmd.MarkFlagRequired("flash")
	cmd.Flags().StringVarP(&o.ISOSource, "source", "s", "", "Source iso to use. By default the latest daily ISO for Ubuntu will be downloaded and saved in a \"files\" directory (created if not exists)")
	cmd.Flags().StringVarP(&o.ExtraFiles, "extra-files", "x", "", "Additional file or files to include in the image.")
	cmd.Flags().StringVarP(&o.Version, "version", "v", "22.04", "Ubuntu version to install and flash. Cannot be used with source flag")
	cmd.MarkFlagsMutuallyExclusive("source", "version")
	cmd.Flags().BoolVarP(&o.NoVerify, "no-verify", "k", false,
		"Disable GPG verification of the source ISO file. By default SHA256SUMS and SHA256SUMS.gpg from releases will be used to verify the authenticity and integrity of the source ISO file")
	cmd.Flags().StringP("userdata", "u", "", "user-data file to add. More info at https://ubuntu.com/server/docs/install/autoinstall-reference")
	return cmd
}

func (o *UbuntuDistro) Complete(cmd *cobra.Command) error {
	o.DistroType = distro.Ubuntu
	_, err := exec.LookPath("xorriso")
	if err != nil {
		return fmt.Errorf("package xorriso not found")
	}
	userdata_path, err := cmd.Flags().GetString("userdata")
	if err != nil {
		return err
	}
	if len(userdata_path) == 0 {
		o.UserData = DefaultUserData()
	} else {
		yaml_file, err := os.ReadFile(userdata_path)
		if err != nil {
			return fmt.Errorf("opening yaml file: %s", err)
		}
		err = yaml.Unmarshal(yaml_file, &o.UserData)
		if err != nil {
			return fmt.Errorf("un-marshaling user-data: %s", err)
		}
	}
	switch o.Version {
	case "22.04", "20.04":
	default:
		return fmt.Errorf("invalid version %s", o.Version)
	}
	if err := flash.IsValidFlashDevice(o.Flash); err != nil {
		return err
	}
	_, err = flash.GetFlashDeviceInfo(o.Flash)
	if err != nil {
		return err
	}
	o.ISODestination = distro.DefaultDir()
	if _, err = os.ReadDir(o.ISODestination); errors.Is(err, os.ErrNotExist) {
		if err := os.Mkdir(o.ISODestination, os.ModePerm); err != nil {
			return fmt.Errorf("trying to create directory %s: %s", o.ISODestination, err)
		}
	} else {
		return err
	}

	return nil
}

func (o *UbuntuDistro) Run(args []string) error {
	iso_filename, err := DownloadUbuntuISO(o.Version, o.ISODestination)
	if err != nil {
		return err
	}
	if !o.NoVerify {
		if err := VerifyISO(iso_filename, o.Version, o.ISODestination); err != nil {
			return err
		}
	}
	err = CreateInstaISO(iso_filename, o.Version, o.ISODestination, o.UserData)
	if err != nil {
		return err
	}
	return nil
}
