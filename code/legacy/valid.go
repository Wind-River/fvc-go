package legacy

import "encoding/hex"

func Valid(b []byte) bool {
	if len(b) != 20 { // length of sha1
		return false
	}

	return true
}

func ValidHex(s string) bool {
	if len(s) != 40 { // length of hex sha1
		return false
	}

	if _, err := hex.DecodeString(s); err != nil { // not a hax-encoded string
		return false
	}

	return true
}
