package code

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.devstar.cloud/ip-systems/verification-code.git/code/legacy"
)

func VersionOf(b []byte) (Version, error) {
	if len(b) == 20 { // check if version 0
		if valid := legacy.Valid(b); !valid {
			return VERSION_ZERO, errors.New("invalid")
		}

		return VERSION_ZERO, nil
	}

	header := b[0:4]
	payload := b[4:]

	switch {
	case bytes.Compare(header, []byte("fvc1")) == 0:
		// check if version 1
		if len(payload) != 20 {
			return VERSION_ONE, errors.New("sha1s are 20 bytes long")
		}
		return VERSION_ONE, nil
	case bytes.Compare(header, []byte("fvc2")) == 0:
		// check if version 2
		if len(payload) != 32 {
			return VERSION_TWO, errors.New("sha256s are 32 bytes long")
		}
		return VERSION_TWO, nil
	default:
		return 0, errors.New(fmt.Sprintf("unpexpected header: \"%v\"", header))
	}
}

func VersionOfHex(s string) (Version, error) {
	if len(s) == 40 { // check if version 0
		if valid := legacy.ValidHex(s); !valid {
			return VERSION_ZERO, errors.New("invalid")
		}

		return VERSION_ZERO, nil
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("invalid hex: \"%s\"", s))
	}

	return VersionOf(b)
}
