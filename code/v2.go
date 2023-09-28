// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package code

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"sort"

	"github.com/Wind-River/fvc-go.git/hasher"
	"github.com/pkg/errors"
)

type VersionTwoHasher struct {
	fileSha256s [][]byte
	sorted      bool
}

// ReadFile is used to take in a parameter of type io.Reader as input, calculate the SHA256
// of the file and adding it to the list, returns an error
// in case of failure
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

// AddSha256Hex is used to add a file's SHA256 from the hex string when the data used
// is the metadata from the database instead of the directory
// of file, returns an error in case of failure.
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

// AddSha256 is used to add a file's SHA256 from the slice of bytes when the data used
// is the metadata from the database instead of the directory
// of file, returns an error in case of failure.
func (h *VersionTwoHasher) AddSha256(b []byte) error {
	if len(b) != 32 {
		return errors.New("sha256 hashes are 32 bytes long")
	} else {
		h.fileSha256s = append(h.fileSha256s, b)
		h.sorted = false
		return nil
	}
}

// Sum is used to sort the list of files and calculating the SHA1
// of the list and returns the fvc as a slice of bytes
func (h VersionTwoHasher) Sum() []byte {
	//sorts the list of files if not already sorted
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

// Hex converts the fvc from slice of bytes to hex string
func (h VersionTwoHasher) Hex() string {
	return hex.EncodeToString(h.Sum())
}

func NewVersionTwo() hasher.FileCollectionHasher {
	return &VersionTwoHasher{
		fileSha256s: make([][]byte, 0),
	}
}
