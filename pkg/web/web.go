package web

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func HttpClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: timeout * time.Second,
	}
	return client
}

func GetRequest(client *http.Client, endpoint string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if response.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("404 not found: %s", endpoint)
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s error code: %s", response.Status, endpoint)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DownloadFile(client *http.Client, endpoint string, dirpath string) error {
	var (
		fileName string
		fullPath string
	)

	if _, err := os.ReadDir(dirpath); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	} else if response.StatusCode == http.StatusNotFound {
		return fmt.Errorf("404 not found: %s", endpoint)
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s error code: %s", response.Status, endpoint)
	}

	defer response.Body.Close()

	fileURL, err := url.Parse(endpoint)
	path := fileURL.Path

	segments := strings.Split(path, "/")
	fileName = segments[len(segments)-1]
	fullPath = fmt.Sprintf("%s/%s", dirpath, fileName)
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		if os.IsTimeout(err) {
			return fmt.Errorf("timeout exceeded")
		}
		return err
	}
	return nil
}
