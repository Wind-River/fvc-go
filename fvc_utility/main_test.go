// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.
package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"gitlab.devstar.cloud/ip-systems/verification-code.git/code"
	fvchandler "gitlab.devstar.cloud/ip-systems/verification-code.git/fvc_utility/fvc_handler"
)

// TestCheckIfPkg tests the CheckIfPkg function with the input
// of an archive and compares the result with the
// given string and throws an error in case of failure
func TestCheckIfPkg(t *testing.T) {
	Bytes := false
	Verbose := false
	hasher := code.NewVersionTwo()
	listFiles := false
	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, make([]string, 0))
	result := f.CheckIfPkg(hasher, "testdir/libarchive_.tar.gz")
	if result != "testdir/libarchive_-archive" {
		t.Error("Expected testdir/libarchive_-archive but got", result)
	}
}

// TestCalculateFvc tests the CalculateFvc function with the input
// of a directory and compares the result with the
// given string and throws an error in case of failure
func TestCalculateFvc(t *testing.T) {
	Bytes := false
	Verbose := false
	hasher := code.NewVersionTwo()
	listFiles := false
	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, make([]string, 0))
	resBytes, resString := f.CalculateFvc(hasher, "testdir/openmp_5.0.0")
	if resBytes != nil || resString != "4656433200de6ba2e40e22cac22fe6b608e6ec26241268069f8314bb29b49e09c6c9ee9240" {
		t.Error("Expected nil,4656433200de6ba2e40e22cac22fe6b608e6ec26241268069f8314bb29b49e09c6c9ee9240 but got: ", resBytes, resString)
	}
}

// TestExtractSubPkg tests the ExtractSubPkg function with the input
// of an archive and check if the result is an error
//  or nil and throws an error in case of error
func TestExtractSubPkg(t *testing.T) {
	Bytes := false
	Verbose := false
	hasher := code.NewVersionTwo()
	listFiles := false
	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, make([]string, 0))
	result := f.ExtractSubPkg(hasher, "testdir/openssl-3.0.7-r2.tar.gz")
	if result != nil {
		t.Error("Expected nil but got", result)
	}
}

// TestReadDirectory tests the ReadDirectory function with the input
// of a directory and check if the result is an error
//  or nil and throws an error in case of error
func TestReadDirectory(t *testing.T) {
	Bytes := false
	Verbose := false
	hasher := code.NewVersionTwo()
	listFiles := false
	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, make([]string, 0))
	result := f.ReadDirectory(hasher, "testdir/openmp_5.0.0")
	if result != nil {
		t.Error("Expected nil but got", result)
	}
}

// TestExtractPkg tests the ExtractPkg function with the input
// of an archive, archive name and destination directory and check if the result
// is an error or nil and throws an error in case of error
func TestExtractPkg(t *testing.T) {
	Bytes := false
	Verbose := false
	listFiles := false
	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, make([]string, 0))
	result := f.ExtractPkg("testdir/openssl-3.0.7-r2.tar.gz", "openssl-3.0.7-r2.tar.gz", "testdir/openssl")
	if result != nil {
		t.Error("Expected nil but got", result)
	}
}

// TestRemovePkg tests the RemovePkg function with the input
// of a list of directories and check if the result is an error
//  or nil and throws an error in case of error
func TestRemovePkg(t *testing.T) {
	Bytes := false
	Verbose := false
	listFiles := false
	var testList []string
	testList = append(testList, "testdir/libarchive_-archive")
	testList = append(testList, "testdir/openssl-3.0.7-r2-archive")
	testList = append(testList, "testdir/openssl")

	f := fvchandler.NewFvcUtil(Bytes, Verbose, listFiles, testList)
	result := f.RemovePkg()
	if result != nil {
		t.Error("Expected nil but got", result)
	}
}

// TestMain checks the main function by passing in command-line arguments
// and calculate the file verification code which is compared to the given
// string and returns an error in case of failure
func TestMain(t *testing.T) {
	os.Args = []string{"fvc_utility", "testdir/libarchive_.tar.gz"}
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()
	w.Close()
	os.Stdout = old
	out := <-outC
	expected := "\n---------------------------Result---------------------------\n\nFile Count = 627\nVerification Code(String) = 46564332007d96db439742c53b1094915b6cf65901a408ff13a42ae77bee3cb2851480fafc\n"
	if out != expected {
		t.Errorf("Expected %s but got %s", expected, out)
	}

}
