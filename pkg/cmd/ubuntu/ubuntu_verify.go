package ubuntu

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	hkp "github.com/s-mahm/instaOS/pkg/util/openpgp-hkp"
	"github.com/s-mahm/instaOS/pkg/web"
)

const ubuntu_keyid = "0x843938DF228D22F7B3742BC0D94AA3F0EFE21092"

func VerifyISO(filename string, version string, destination string) error {
	checksums, err := verifySignature(version, destination)
	if err != nil {
		return fmt.Errorf("verifying signature: %s", err)
	}
	fmt.Println("GOOD Signature")
	checksum_pattern, err := regexp.Compile(fmt.Sprintf(`(?m)(.*) \*%s`, filename))
	matches := checksum_pattern.FindStringSubmatch(checksums)
	if err != nil || len(matches) == 0 {
		return fmt.Errorf("cannot find iso in SHA256SUMS")
	}
	valid_checksum := matches[1]
	if err = CompareChecksums(filename, valid_checksum, destination); err != nil {
		return fmt.Errorf("comparing checksums: %s", err)
	}
	fmt.Println("Checksums matched")
	return nil

}

func verifySignature(version string, destination string) (string, error) {
	client := web.HttpClient(10)
	keyring, err := getKeyRing()
	if err != nil {
		return "", fmt.Errorf("getting key ring: %s", err)
	}
	signature, err := web.GetRequest(client, fmt.Sprintf("%s/%s/SHA256SUMS.gpg", iso_base_url, version))
	if err != nil {
		return "", fmt.Errorf("downloading SHA256SUMS.gpg: %s", err)
	}
	verification_target, err := web.GetRequest(client, fmt.Sprintf("%s/%s/SHA256SUMS", iso_base_url, version))
	if err != nil {
		return "", fmt.Errorf("downloading SHA256SUMS: %s", err)
	}
	entity, err := openpgp.CheckArmoredDetachedSignature(keyring, bytes.NewReader(verification_target), bytes.NewReader(signature), &packet.Config{})
	if err != nil {
		return "", fmt.Errorf("checking detached signature: %s", err)
	}
	for _, identity := range entity.Identities {
		if identity.Name != "Ubuntu CD Image Automatic Signing Key (2012) <cdimage@ubuntu.com>" {
			return "", fmt.Errorf("BAD signature from: %s", identity.Name)
		}
	}
	return string(verification_target), nil
}

func getKeyRing() (openpgp.EntityList, error) {
	var key_buffer bytes.Buffer
	c := hkp.Client{Host: "https://keyserver.ubuntu.com"}
	req := hkp.LookupRequest{Search: ubuntu_keyid}
	keys, err := c.Get(&req)
	if err != nil {
		return nil, err
	}
	err = keys[0].Serialize(&key_buffer)
	if err != nil {
		return nil, err
	}
	key_bytes := key_buffer.Bytes()
	keyring, err := openpgp.ReadKeyRing(bytes.NewReader(key_bytes))

	return keyring, nil
}

func CompareChecksums(filename string, valid_checksum string, destintation string) error {
	iso_file, err := os.Open(fmt.Sprintf("%s/%s", destintation, filename))
	if err != nil {
		return fmt.Errorf("opening iso file: %s", err)
	}
	defer iso_file.Close()
	iso_hash := sha256.New()
	if _, err := io.Copy(iso_hash, iso_file); err != nil {
		return fmt.Errorf("generating checksum for iso: %s", err)
	}
	iso_checksum := fmt.Sprintf("%x", iso_hash.Sum(nil))
	if iso_checksum != valid_checksum {
		return fmt.Errorf("checksums do not match")
	}
	return nil
}
