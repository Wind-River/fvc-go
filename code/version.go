package code

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
	"gitlab.devstar.cloud/ip-systems/verification-code.git/code/legacy"
)

func VersionOf(b []byte) (*Version, error) {
	if len(b) == 20 { // check if version 0
		if valid := legacy.Valid(b); !valid {
			return nil, errors.New("invalid")
		}

		v := VERSION_ZERO
		return &v, nil
	}

	header := b[0:5]
	payload := b[5:]

	switch {
	case bytes.Compare(header, []byte("FVC1\000")) == 0:
		v := VERSION_ONE
		// verify version 1
		if len(payload) != 20 {
			return &v, errors.New("sha1s are 20 bytes long")
		}
		return &v, nil
	case bytes.Compare(header, []byte("FVC2\000")) == 0:
		v := VERSION_TWO
		// verify version 2
		if len(payload) != 32 {
			return &v, errors.New("sha256s are 32 bytes long")
		}
		return &v, nil
	default:
		return nil, errors.New(fmt.Sprintf("unpexpected header: \"%v\"", header))
	}
}

func VersionOfHex(s string) (*Version, error) {
	if len(s) == 40 { // check if version 0
		if valid := legacy.ValidHex(s); !valid {
			return nil, errors.New("invalid")
		}

		v := VERSION_ZERO
		return &v, nil
	}

	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid hex: \"%s\"", s))
	}

	return VersionOf(b)
}
