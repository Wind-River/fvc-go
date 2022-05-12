package code

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"github.com/pkg/errors"
	"gitlab.devstar.cloud/ip-systems/verification-code.git/hasher"
)

type VersionTwoHasher struct {
	fileSha256s [][]byte
	sorted      bool
}

func (h *VersionTwoHasher) ReadFile(r io.Reader) error {
	// calculate sha256 and add to list
	hasher := sha256.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return errors.Wrapf(err, "error hashing file")
	}

	h.fileSha256s = append(h.fileSha256s, hasher.Sum(nil))
	h.sorted = false
	return nil
}

func (h *VersionTwoHasher) AddSha256Hex(s string) error {
	if b, err := hex.DecodeString(s); err != nil {
		return errors.Wrapf(err, "error decoding hex string")
	} else if len(b) != 32 {
		return errors.New("sha256 hashes are 32 bytes long")
	} else {
		h.fileSha256s = append(h.fileSha256s, b)
		h.sorted = false
		return nil
	}
}

func (h *VersionTwoHasher) AddSha256(b []byte) error {
	if len(b) != 32 {
		return errors.New("sha256 hashes are 32 bytes long")
	} else {
		h.fileSha256s = append(h.fileSha256s, b)
		h.sorted = false
		return nil
	}
}

func (h VersionTwoHasher) Sum() []byte {
	if !h.sorted {
		sort.Slice(h.fileSha256s, func(i int, j int) bool {
			return bytes.Compare(h.fileSha256s[i], h.fileSha256s[j]) < 0
		})
		h.sorted = true
	}

	hasher := sha256.New()
	for _, v := range h.fileSha256s {
		hasher.Write(v)
	}
	return append([]byte("FVC2\000"), hasher.Sum(nil)...)
}

func (h VersionTwoHasher) Hex() string {
	return hex.EncodeToString(h.Sum())
}

func NewVersionTwo() hasher.FileCollectionHasher {
	return &VersionTwoHasher{
		fileSha256s: make([][]byte, 0),
	}
}
