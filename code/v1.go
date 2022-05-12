package code

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"sort"

	"github.com/pkg/errors"
	"gitlab.devstar.cloud/ip-systems/verification-code.git/hasher"
)

type VersionOneHasher struct {
	fileSha1Hexs []string
	sorted       bool
}

func (h *VersionOneHasher) ReadFile(r io.Reader) error {
	// calculate sha1 and add to list
	hasher := sha1.New()
	if _, err := io.Copy(hasher, r); err != nil {
		return errors.Wrapf(err, "error hashing file")
	}

	h.fileSha1Hexs = append(h.fileSha1Hexs, hex.EncodeToString(hasher.Sum(nil)))
	h.sorted = false
	return nil
}

func (h *VersionOneHasher) AddSha1Hex(s string) error {
	if b, err := hex.DecodeString(s); err != nil {
		return errors.Wrapf(err, "error decoding hex string")
	} else if len(b) != 20 {
		return errors.New("sha1 hashes are 20 bytes long")
	}

	h.fileSha1Hexs = append(h.fileSha1Hexs, s)
	return nil
}

func (h VersionOneHasher) Sum() []byte {
	if !h.sorted {
		sort.Sort(sort.StringSlice(h.fileSha1Hexs))
		h.sorted = true
	}

	hasher := sha1.New()
	for _, v := range h.fileSha1Hexs {
		hasher.Write([]byte(v))
	}
	return append([]byte("FVC1\000"), hasher.Sum(nil)...)
}

func (h VersionOneHasher) Hex() string {
	return hex.EncodeToString(h.Sum())
}

func NewVersionOne() hasher.FileCollectionHasher {
	return &VersionOneHasher{
		fileSha1Hexs: make([]string, 0),
	}
}
