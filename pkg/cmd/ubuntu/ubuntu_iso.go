package ubuntu

import (
	"fmt"
	"regexp"

	"github.com/s-mahm/instaOS/pkg/web"
)

const iso_base_url = "http://releases.ubuntu.com"
const iso_fallback_url = "http://cdimage.ubuntu.com/releases"

func DownloadUbuntuISO(version string, destination string) error {
	client := web.HttpClient(300)
	iso_pattern, err := regexp.Compile(fmt.Sprintf(`(?m)ubuntu-%s.*server-amd64.iso`, version))
	if err != nil {
		return fmt.Errorf("downloading iso: %s", err)
	}
	releases_out, err := web.GetRequest(client, fmt.Sprintf("%s/%s", iso_base_url, version))
	if err != nil {
		return fmt.Errorf("finding ubuntu %s iso download url - %s", version, err)
	}
	matches := iso_pattern.FindStringSubmatch(string(releases_out))
	if len(matches) == 0 {
		return fmt.Errorf("cannot find ubuntu %s iso download url", version)
	}
	download_url := fmt.Sprintf("%s/%s/%s", iso_base_url, version, matches[0])
	fmt.Println("Downloading ", download_url)
	err = web.DownloadFile(client, download_url, destination)
	if err != nil {
		return fmt.Errorf("downloading iso from url - %s", err)
	}
	return nil
}
