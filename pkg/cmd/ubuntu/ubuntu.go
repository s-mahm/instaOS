package ubuntu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/s-mahm/instaOS/pkg/cmd/util"
	"github.com/s-mahm/instaOS/pkg/util/flash"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type UbuntuOptions struct {
	Flash        string
	ExtraFiles   string
	UserDataPath string
	UserData     UserData
	Version      string
	Destination  string
	Source       string
	NoVerify     bool
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
	cmd.Flags().StringVarP(&o.Source, "source", "s", "", "Source iso to use. By default the latest daily ISO for Ubuntu will be downloaded and saved in a \"files\" directory (created if not exists)")
	cmd.Flags().StringVarP(&o.Version, "version", "v", "22.04", "Ubuntu version to install and flash. Cannot be used with source flag")
	cmd.MarkFlagsMutuallyExclusive("source", "version")
	cmd.Flags().BoolVarP(&o.NoVerify, "no-verify", "k", false,
		"Disable GPG verification of the source ISO file. By default SHA256SUMS and SHA256SUMS.gpg from releases will be used to verify the authenticity and integrity of the source ISO file")
	cmd.Flags().StringVarP(&o.UserDataPath, "userdata", "u", "", "user-data file to add. More info at https://ubuntu.com/server/docs/install/autoinstall-reference")
	return cmd
}

func (o *UbuntuOptions) Complete() error {
	_, err := exec.LookPath("xorriso")
	if err != nil {
		return fmt.Errorf("package xorriso not found")
	}
	if len(o.UserDataPath) == 0 {
		o.UserData = DefaultUserData()
	} else {
		yaml_file, err := os.ReadFile(o.UserDataPath)
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
	if !o.NoVerify {
		if err := VerifyISO(iso_filename, o.Version, o.Destination); err != nil {
			return err
		}
	}
	err = CreateInstaISO(iso_filename, o.Version, o.Destination, o.UserData)
	if err != nil {
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
