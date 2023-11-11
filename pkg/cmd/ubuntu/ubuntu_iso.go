package ubuntu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/s-mahm/instaOS/pkg/util/file"
	"github.com/s-mahm/instaOS/pkg/web"
)

const iso_base_url = "http://releases.ubuntu.com"
const iso_fallback_url = "http://cdimage.ubuntu.com/releases"

func DownloadUbuntuISO(version string, destination string) (string, error) {
	iso_url := iso_base_url
	client := web.HttpClient(300)
	_, err := web.GetRequest(client, fmt.Sprintf("%s", iso_url))
	if err != nil {
		iso_url = iso_fallback_url
	}
	iso_pattern, err := regexp.Compile(fmt.Sprintf(`(?m)ubuntu-%s.*server-amd64.iso`, version))
	if err != nil {
		return "", fmt.Errorf("downloading iso: %s", err)
	}
	releases_out, err := web.GetRequest(client, fmt.Sprintf("%s/%s", iso_url, version))
	if err != nil {
		return "", fmt.Errorf("finding ubuntu %s iso download url: %s", version, err)
	}
	matches := iso_pattern.FindStringSubmatch(string(releases_out))
	if len(matches) == 0 {
		return "", fmt.Errorf("cannot find ubuntu %s iso download url", version)
	}
	iso_name := matches[0]
	download_url := fmt.Sprintf("%s/%s/%s", iso_url, version, iso_name)
	if isoAlreadyExists(destination, iso_name) {
		fmt.Printf("File %s already downloded, using existing file\n", iso_name)
	} else {
		err = web.DownloadFile(client, download_url, destination, true)
		if err != nil {
			return "", fmt.Errorf("downloading iso from url: %s", err)
		}
	}
	return iso_name, nil
}

func isoAlreadyExists(destination string, iso_name string) bool {
	_, err := os.Stat(fmt.Sprintf("%s/%s", destination, iso_name))
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func CreateInstaISO(filename string, version string, destination string, userdata UserData) error {
	tempdir, err := os.MkdirTemp("files", fmt.Sprintf("ubuntu-%s-", version))
	if err != nil {
		return fmt.Errorf("creating temp dir: %s", err)
	}
	defer os.RemoveAll(tempdir)
	if err = extractISOToDirectory(filename, destination, tempdir); err != nil {
		return err
	}
	files_to_edit := func(prefix string, files []string) []string {
		for i, filename := range files {
			files[i] = fmt.Sprintf("%s/%s", prefix, filename)
		}
		return files
	}
	if err = file.ReplaceTextInFiles(files_to_edit(tempdir, []string{"boot/grub/grub.cfg", "boot/grub/loopback.cfg"}), "---", `  autoinstall    quiet    ds=nocloud\;s=/cdrom/nocloud/  ---`); err != nil {
		return fmt.Errorf("error editing iso files: %s", err)
	}
	if err = file.ReplaceTextInFiles(files_to_edit(tempdir, []string{"boot/grub/grub.cfg"}), "timeout=30", "timeout=1"); err != nil {
		return fmt.Errorf("error editing iso files: %s", err)
	}
	if err = os.Mkdir(fmt.Sprintf("%s/nocloud", tempdir), os.ModePerm); err != nil {
		return fmt.Errorf("creating nocloud dircetory")
	}
	if err = AddUserData(userdata, fmt.Sprintf("%s/nocloud", tempdir)); err != nil {
		return err
	}
	if _, err = os.Create(fmt.Sprintf("%s/nocloud/meta-data", tempdir)); err != nil {
		return err
	}
	if err = createISOFromDirectory(fmt.Sprintf("%s/%s", destination, filename), tempdir, fmt.Sprintf("%s/ubuntu-%s-autoinstall.iso", destination, version)); err != nil {
		return err
	}
	return nil
}

func extractISOToDirectory(filename string, source string, destination string) error {
	xorriso_args := fmt.Sprintf("-osirrox on -indev %s/%s -extract / %s", source, filename, destination)
	cmd := exec.Command("xorriso", strings.Split(xorriso_args, " ")...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("xorriso error: %s", err)
	}
	err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, os.FileMode(0755))
	})
	if err != nil {
		return fmt.Errorf("setting directory permissions: %s", err)
	}
	return nil
}

func createISOFromDirectory(reference string, source string, isoname string) error {
	xorriso_create_args := "-as mkisofs -r -V \"ubuntu-autoinstall\""
	xorriso_ref_args := fmt.Sprintf("-indev %s -report_el_torito as_mkisofs", reference)
	cmd := exec.Command("xorriso", strings.Split(xorriso_ref_args, " ")...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xorriso error: %s", err)
	}
	flag_pattern, err := regexp.Compile(`(?m)^((?:-|--).*)$`)
	matches := flag_pattern.FindAllStringSubmatch(string(out), -1)
	for _, match := range matches {
		if !strings.HasPrefix(match[0], "-V") {
			xorriso_create_args = xorriso_create_args + " " + match[0]
		}
	}
	xorriso_create_args = xorriso_create_args + " " + fmt.Sprintf("-o %s %s", isoname, source)
	xorriso_create_args = strings.ReplaceAll(xorriso_create_args, "'", "")
	cmd = exec.Command("xorriso", strings.Split(xorriso_create_args, " ")...)
	out, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("xorriso error: %s", err)
	}

	return nil
}

func AddUserData(userdata UserData, destination string) error {
	userdata_yaml, err := yaml.Marshal(&userdata)
	if err != nil {
		return fmt.Errorf("marshaling user-data: %s", err)
	}
	final_data := []byte("#cloud-config\n" + string(userdata_yaml))
	err = os.WriteFile(fmt.Sprintf("%s/user-data", destination), final_data, 0755)
	if err != nil {
		return fmt.Errorf("writing to user-data: %s", err)
	}
	return nil
}
