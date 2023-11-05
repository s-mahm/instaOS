package ubuntu

func DefaultUserData() UserData {
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
	default_network := Network{
		Version:    2,
		Renderers:  "network-manager",
		Activators: "network-manager",
		Ethernets: map[string]struct {
			DHCP4 bool `yaml:"dhcp4"`
		}{
			"eth0": {
				DHCP4: true,
			},
		},
	}
	default_identity := Identity{
		Hostname: "ubuntu",
		Username: "user",
		Password: "\"$1$xcqce3Uk$Xlpomuf6zmspPe2d4piIi1\"",
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
		}{{
			Arches: []string{"default"},
			Uri:    "http://us.archive.ubuntu.com/ubuntu/",
		},
		},
	}
	default_packages := []string{
		"build-essential",
		"network-manager",
		"dkms",
		"vim",
		// "ubuntu-desktop-minimal",
	}
	default_autoinstall := Autoinstall{
		Version:          2,
		RefreshInstaller: nil,
		Update:           true,
		Storage:          default_storage,
		Network:          default_network,
		Keyboard:         default_keyboad,
		Identity:         default_identity,
		SSH:              default_ssh,
		Apt:              default_apt,
		Packages:         default_packages,
		PackageUpdate:    true,
		PackageUpgrade:   true,
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
	Network          Network     `yaml:"network"`
	Storage          Storage     `yaml:"storage"`
	Keyboard         Keyboard    `yaml:"keyboard"`
	Identity         Identity    `yaml:"identity"`
	SSH              SSH         `yaml:"ssh"`
	Apt              Apt         `yaml:"apt"`
	Packages         []string    `yaml:"packages"`
	PackageUpdate    bool        `ymal:"package_update"`
	PackageUpgrade   bool        `ymal:"package_upgrade"`
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
	Version    int    `yaml:"version"`
	Renderers  string `yaml:"renderers"`
	Activators string `yaml:"activators"`
	Ethernets  map[string]struct {
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

// default_config_layout := []StorageConfig{
// 		{
// 			GrubDevice: true,
// 			ID:         "disk-sda",
// 			Path:       "/dev/sda",
// 			Ptable:     "gpt",
// 			Type:       "string",
// 			Wipe:       "superblock-recursive",
// 		},
// 		{
// 			Device: "disk-sda",
// 			Flag:   "bios_grub",
// 			ID:     "partiton-0",
// 			Number: 1,
// 			Size:   1048576,
// 			Type:   "partiton",
// 		},
// 		{
// 			Device: "disk-sda",
// 			ID:     "partiton-1",
// 			Number: 2,
// 			Size:   -1,
// 			Type:   "partition",
// 			Wipe:   "superblock",
// 		},
// 		{
// 			Fstype: "ext4",
// 			ID:     "format-0",
// 			Type:   "format",
// 			Volume: "partition-1",
// 		},
// 		{
// 			Device: "format-0",
// 			ID:     "mount-0",
// 			Path:   "/",
// 			Type:   "mount",
// 		},
