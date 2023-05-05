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
	"encoding/hex"
	"fmt"

	"github.com/pkg/errors"
)

// VersionOf checks the slice of bytes to figure out the
// version of fvc used and returns the version
// and error
func VersionOf(b []byte) (*Version, error) {
	header := b[0:5]
	payload := b[5:]

	switch {
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

// VersionOfHex checks the hex string to figure out the
// version of fvc used and returns the version
// and error
func VersionOfHex(s string) (*Version, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("invalid hex: \"%s\"", s))
	}

	return VersionOf(b)
}
