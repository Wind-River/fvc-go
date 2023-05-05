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

// Struct for holding some flags, file count and list of extracted packages
type FvcUtil struct {
	fvcBytes         bool
	fvcVerbose       bool
	fvcListFiles     bool
	count            int
	listExtractedPkg []string
}

// CheckIfPkg checks if a given path is a directory or an archive package
// In case it is an archive package, the archive is extracted and its name is added
// to the list of extracted packages. Returns the updated path incase an archive has beeen extracted.
func (f *FvcUtil) CheckIfPkg(h hasher.FileCollectionHasher, path string) string {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	//check is the file is a regular file
	if info.Mode().IsRegular() {
		//check is the file is extractable
		if extract.IsExtractable(info.Name()) == 1.0 {
			destPath := filepath.Dir(path)
			// removing the extensions from the file name
			fileName, _, _ := extract.SplitExt(info.Name())
			if err1 := f.ExtractPkg(path, info.Name(), destPath+"/"+fileName+"-archive"); err1 != nil {
				fmt.Fprintln(os.Stderr, errors.Wrapf(err1, "error extracting archive"))
				os.Exit(0)
			}
			// setting the path to newly extracted package
			path = destPath + "/" + fileName + "-archive"
			f.listExtractedPkg = append(f.listExtractedPkg, path)
		}
	}
	//extracting sub packages
	if err := f.ExtractSubPkg(h, path); err != nil {
		fmt.Println(errors.Wrapf(err, "error opening directory"))
		os.Exit(0)
	}
	return path
}

// ExtractSubPkg walks through a directory looking for any achives
// and extracts all such archives that it encounters if they are extractable.
// During exctraction if it fails then we treat them as normal files and call ReadFile
// Returns an error if extraction fails and nil otherwise.
func (f *FvcUtil) ExtractSubPkg(h hasher.FileCollectionHasher, dirPath string) error {
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//check is the file is a regular file
		if info.Mode().IsRegular() {
			//check is the file is extractable
			if extract.IsExtractable(info.Name()) == 1.0 {
				destPath := filepath.Dir(path)
				fileName, _, _ := extract.SplitExt(info.Name())
				err := f.ExtractPkg(path, info.Name(), destPath+"/"+fileName+"-archive")
				// if extraction fails then we treat it as a file and call readfile on it
				if err != nil {
					file, err := os.Open(path)
					if err != nil {
						return err
					}
					if err := h.ReadFile(file); err != nil {
						return errors.Wrapf(err, "error reading file - "+file.Name())
					}
					if f.fvcListFiles {
						fmt.Fprintln(os.Stderr, info.Name())
					}
					f.count++
					file.Close()
				}
				// appending the name of the extracted package to the list for removing at the end
				f.listExtractedPkg = append(f.listExtractedPkg, destPath+"/"+fileName+"-archive")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// ExtractPkg makes use of the extract library to extract the package at the
// given 'path' having the given 'name' onto the given 'dest' address.
// Returns an error on failure and nil otherwise.
func (f *FvcUtil) ExtractPkg(path string, name string, dest string) error {
	f.VerbosePrint("\nExtracting the archive:" + name + "...")
	// using the NewAt function from the extract library for creating the new directory for extraction
	extractor, err := extract.NewAt(path, name, dest)
	if err != nil {
		return err
	}
	// using the Extract function from the extract library for extracting the archive
	if extracted, err := extractor.Extract(); err != nil {
		return err
	} else {
		f.VerbosePrint("Completed Extracting the archive at:" + extracted + " \u2713")
	}
	return nil
}

// RemovePkg uses the list of extracted packages to remove each package
// in order of them being extracted. The function returns an error when failure to
// removing a file is encountered and nil otherwise
func (f *FvcUtil) RemovePkg() error {
	// check if there are any packages that have been extracted
	if len(f.listExtractedPkg) != 0 {
		for i := len(f.listExtractedPkg) - 1; i >= 0; i-- {
			f.VerbosePrint("\nRemoving the extracted archive:" + f.listExtractedPkg[i] + "...")
			// removing the extracted package
			if err := os.RemoveAll(f.listExtractedPkg[i]); err != nil {
				return err
			}
			f.VerbosePrint("Completed removing archive:" + f.listExtractedPkg[i] + " \u2713")
		}
	}
	return nil
}

func NewFvcUtil(inBytes bool, verbose bool, listFiles bool, listExtracted []string) FvcUtil {
	return FvcUtil{fvcBytes: inBytes, fvcVerbose: verbose, fvcListFiles: listFiles, count: 0, listExtractedPkg: listExtracted}
}
