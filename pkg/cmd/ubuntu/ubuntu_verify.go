package ubuntu

import (
	"fmt"

	hkp "github.com/s-mahm/instaOS/pkg/util/openpgp-hkp"
)

const ubuntu_keyid = "0x843938DF228D22F7B3742BC0D94AA3F0EFE21092"

func VerifyISO() error {
	if err := getPublicKey(); err != nil {
		return fmt.Errorf("getting public key from server: %s", err)
	}
	return nil
}

func getPublicKey() error {
	c := hkp.Client{Host: "https://keyserver.ubuntu.com"}
	req := hkp.LookupRequest{Search: ubuntu_keyid}
	keys, err := c.Get(&req)
	if err != nil {
		return err
	}
	fmt.Println(keys)
	return nil
}
