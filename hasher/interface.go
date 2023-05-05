// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package hasher

import (
	"io"
)

//Interface for the fvc library
type FileCollectionHasher interface {
	// ReadFile is used to take in a parameter of type io.Reader as input, calculate the SHA256
	// of the file and adding it to the list, returns an error
	// in case of failure
	ReadFile(io.Reader) error
	// Sum is used to sort the list of files and calculating the SHA1
	// of the list and returns the fvc as a slice of bytes
	Sum() []byte
	// Hex converts the fvc from slice of bytes to hex string
	Hex() string
}
