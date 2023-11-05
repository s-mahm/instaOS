package ubuntu

type UserData struct {
	Autoinstall struct {
		Version          int         `yaml:"version"`
		RefreshInstaller interface{} `yaml:"refresh-installer"`
		Update           string      `yaml:"update"`
		Storage          struct {
			Grub struct {
				ReorderUefi bool `yaml:"reorder_uefi"`
			} `yaml:"grub"`
			Layout struct {
				Name  string `yaml:"name"`
				Match struct {
					Ssd  string `yaml:"ssd"`
					Size string `yaml:"size"`
				} `yaml:"match"`
			} `yaml:"layout"`
			Config []struct {
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
			} `yaml:"config"`
		} `yaml:"storage"`
		Keyboard struct {
			Layout string `yaml:"layout"`
		} `yaml:"keyboard"`
		Identity struct {
			Hostname string `yaml:"hostname"`
			Password string `yaml:"password"`
			Username string `yaml:"username"`
		} `yaml:"identity"`
		SSH struct {
			AllowPw       bool `yaml:"allow-pw"`
			InstallServer bool `yaml:"install-server"`
		} `yaml:"ssh"`
		LateCommands []string `yaml:"late-commands"`
	} `yaml:"autoinstall"`
}
