// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package fvchandler

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"gitlab.devstar.cloud/ip-systems/extract.git"
	"gitlab.devstar.cloud/ip-systems/verification-code.git/hasher"
)

// CalculateFvc is used for firstly extracting any sub packages inside the directory
// After extraction all the files are read inside the directory using 'Walk' function and added to
// the list of files after calculating their SHA256. Finally based on the flag value, the code returns
// the verification code in bytes or Hex string form
func (f *FvcUtil) CalculateFvc(h hasher.FileCollectionHasher, dirPath string) ([]byte, string) {
	f.VerbosePrint("\nOpening Directory:" + dirPath)
	//reading the entire directory
	if err := f.ReadDirectory(h, dirPath); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "error opening directory"))
		os.Exit(0)
	}
	f.VerbosePrint("Closing Directory:" + dirPath)
	//checking if the directory is empty
	if f.count < 1 {
		fmt.Fprintln(os.Stderr, "No files found in directory")
		os.Exit(0)
	}
	f.VerbosePrint("\nSorting list of files and calculating fvc of list...")
	if f.fvcBytes {
		return h.Sum(), ""
	} else {
		return []byte(nil), h.Hex()
	}
}

// VerbosePrint is a helper function used for printing messages in verbose mode
func (f *FvcUtil) VerbosePrint(str string) {
	if f.fvcVerbose {
		fmt.Fprintln(os.Stderr, str)
	}
}

// PrintResult is used to print out the final result
func (f *FvcUtil) PrintResult(hexb []byte, hexs string) {
	fmt.Println("\n---------------------------Result---------------------------")
	fmt.Println("\nFile Count =", f.count)
	if f.fvcBytes {
		fmt.Println("Verification Code(Bytes) =", hexb)
	} else {
		fmt.Println("Verification Code(String) =", hexs)
	}
}

// ReadDirectory is used to access all files inside a directory using 'Walk' function
// Checks if a file is regular and not an archive then calls ReadFile which calculates the SHA256 of
// the 'currentfile' which is open. The function returns an error on failing t read a file and nil otherwise.
func (f *FvcUtil) ReadDirectory(h hasher.FileCollectionHasher, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//check is the file is a regular file
		if info.Mode().IsRegular() == true {
			// check if the file is extractable
			if extract.IsExtractable(info.Name()) != 1.0 {
				//opening the file
				currentFile, err := os.Open(path)
				if err != nil {
					return errors.Wrapf(err, "error opening file - "+info.Name())
				}
				if err := h.ReadFile(currentFile); err != nil {
					return errors.Wrapf(err, "error reading file - "+currentFile.Name())
				}
				if f.fvcListFiles {
					fmt.Fprintln(os.Stderr, info.Name())
				}
				f.VerbosePrint("Reading and calculating SHA256 of" + currentFile.Name() + " \u2713")
				//increasing the file count and closing the current file
				f.count++
				currentFile.Close()
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
