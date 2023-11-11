package ubuntu

import (
	"golang.org/x/crypto/bcrypt"
)

func DefaultUserData() UserData {
	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	default_storage_layout := StorageLayout{
		Name: "lvm",
		Match: struct {
			Ssd  bool   `yaml:"ssd,omitempty"`
			Size string `yaml:"size"`
		}{
			Size: "largest",
		},
	}
	default_storage := Storage{
		Grub: struct {
			ReorderUefi bool `yaml:"reorder_uefi"`
		}{
			ReorderUefi: false,
		},
		Layout: default_storage_layout,
	}
	default_identity := Identity{
		Hostname: "instaos",
		Username: "admin",
		Password: string(password),
	}
	default_keyboad := Keyboard{
		Layout: "us",
	}
	default_ssh := SSH{
		AllowPw:       true,
		InstallServer: true,
	}
	default_apt := Apt{
		Primary: []struct {
			Arches []string `ymal:"arches"`
			Uri    string   `yaml:"uri"`
		}{
			{
				Arches: []string{"default"},
				Uri:    "http://repo.internal/",
			},
			{
				Arches: []string{"i386", "amd64"},
				Uri:    "http://archive.ubuntu.com/ubuntu",
			},
		},
	}

	default_packages := []string{
		"build-essential",
		"network-manager",
		"dkms",
		"vim",
		"ubuntu-desktop-minimal",
	}
	default_latecommands := []string{"shutdown -h now"}
	default_autoinstall := Autoinstall{
		Version:          1,
		RefreshInstaller: nil,
		Update:           true,
		Storage:          default_storage,
		Locale:           "en_US",
		Keyboard:         default_keyboad,
		Identity:         default_identity,
		SSH:              default_ssh,
		Apt:              default_apt,
		Packages:         default_packages,
		LateCommands:     default_latecommands,
	}
	return UserData{
		Autoinstall: default_autoinstall,
	}
}

type UserData struct {
	Autoinstall Autoinstall `yaml:"autoinstall"`
}

type Autoinstall struct {
	Version          int         `yaml:"version"`
	RefreshInstaller interface{} `yaml:"refresh-installer"`
	Update           bool        `yaml:"update"`
	EarlyCommands    []string    `yaml:"early-commands,omitempty"`
	Network          Network     `yaml:"network,omitempty"`
	Storage          Storage     `yaml:"storage,omitempty"`
	Locale           string      `yaml:"locale,omitempty"`
	Keyboard         Keyboard    `yaml:"keyboard,omitempty"`
	Identity         Identity    `yaml:"identity,omitempty"`
	SSH              SSH         `yaml:"ssh,omitempty"`
	Apt              Apt         `yaml:"apt,omitempty"`
	Packages         []string    `yaml:"packages,omitempty"`
	LateCommands     []string    `yaml:"late-commands,omitempty"`
}

type StorageLayout struct {
	Name     string `yaml:"name"`
	Password string `yaml:"password,omitempty"`
	Match    struct {
		Ssd  bool   `yaml:"ssd,omitempty"`
		Size string `yaml:"size"`
	} `yaml:"match"`
}

type StorageConfig struct {
	GrubDevice bool   `yaml:"grub_device,omitempty"`
	ID         string `yaml:"id"`
	Path       string `yaml:"path,omitempty"`
	Ptable     string `yaml:"ptable,omitempty"`
	Type       string `yaml:"type"`
	Wipe       string `yaml:"wipe,omitempty"`
	Device     string `yaml:"device,omitempty"`
	Flag       string `yaml:"flag,omitempty"`
	Number     int    `yaml:"number,omitempty"`
	Size       int    `yaml:"size,omitempty"`
	Fstype     string `yaml:"fstype,omitempty"`
	Volume     string `yaml:"volume,omitempty"`
}

type Storage struct {
	Grub struct {
		ReorderUefi bool `yaml:"reorder_uefi"`
	} `yaml:"grub"`
	Layout StorageLayout   `yaml:"layout"`
	Config []StorageConfig `yaml:"config,omitempty"`
}

type Network struct {
	Version   int `yaml:"version"`
	Ethernets map[string]struct {
		DHCP4 bool `yaml:"dhcp4"`
	} `yaml:"ethernets"`
}

type Keyboard struct {
	Layout string `yaml:"layout"`
}

type Identity struct {
	Hostname string `yaml:"hostname"`
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}

type SSH struct {
	AllowPw       bool `yaml:"allow-pw"`
	InstallServer bool `yaml:"install-server"`
}

type Apt struct {
	Primary []struct {
		Arches []string `ymal:"arches"`
		Uri    string   `yaml:"uri"`
	} `yaml:"primary"`
}
