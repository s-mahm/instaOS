package ubuntu

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

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
		fmt.Println("Downloading ", download_url)
		err = web.DownloadFile(client, download_url, destination)
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

func CreateInstaISO(filename string, version string, destination string) error {
	temp, err := os.MkdirTemp("files", fmt.Sprintf("ubuntu-%s-", version))
	if err != nil {
		return fmt.Errorf("creating temp dir: %s", err)
	}
	// defer os.RemoveAll(temp)
	extractISOToDirectory(filename, destination, temp)
	files_to_edit := func(prefix string, files []string) []string {
		for i, filename := range files {
			files[i] = fmt.Sprintf("%s/%s", prefix, filename)
		}
		return files
	}
	if err = file.ReplaceTextInFiles(files_to_edit(temp, []string{"boot/grub/grub.cfg", "boot/grub/loopback.cfg"}), "---", `autoinstall ds=nocloud\;s=/cdrom/server/  ---`); err != nil {
		return fmt.Errorf("error editing iso files: %s", err)
	}
	return nil
}

func extractISOToDirectory(filename string, destination string, dirpath string) error {
	xorriso_args := fmt.Sprintf("-osirrox on -indev %s/%s -extract / %s", destination, filename, dirpath)
	cmd := exec.Command("xorriso", strings.Split(xorriso_args, " ")...)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("xorriso error: %s", err)
	}
	err = filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
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

func getUserData() {
	return
}
