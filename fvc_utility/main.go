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
	"flag"
	"fmt"
	"os"

	"gitlab.devstar.cloud/ip-systems/verification-code.git/code"
	fvchandler "gitlab.devstar.cloud/ip-systems/verification-code.git/fvc_utility/fvc_handler"
)

var inBytes bool
var verbose bool
var example bool
var output bool
var listFile bool

func init() {
	// Flags for various modes
	flag.BoolVar(&inBytes, "b", false, "Verification Code in Bytes")
	flag.BoolVar(&verbose, "v", false, "Execution in Verbose Mode")
	flag.BoolVar(&example, "e", false, "Various Example Commands")
	flag.BoolVar(&output, "o", false, "Write output to a file")
	flag.BoolVar(&listFile, "l", false, "Lists all the files")
}

// main deals with taking in the command line arguments and
// calls various functions based on the input arguments and flags
func main() {
	// Defining the usage for -help flag
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, "Usage of go run with output on command-line:\n go run main.go [Option] <path to directory/archive>")
		fmt.Fprintln(w, "Usage of go run with output in a file:\n go run main.go -o <path to directory/archive> <name of file>")
		flag.PrintDefaults()
	}
	//parsing the flags
	flag.Parse()
	f := fvchandler.NewFvcUtil(inBytes, verbose, listFile, make([]string, 0))
	if example {
		PrintExamples()
		os.Exit(0)
	}
	path := flag.Arg(0)
	//checking the path
	if path == "" {
		fmt.Fprintln(os.Stderr, "Enter directory name as argument")
		os.Exit(0)
	}
	// prints output to file if output flag is true
	if output {
		outputDest := flag.Arg(1)
		if outputDest == "" {
			outputDest = "fvcoutput.txt"
		}
		//creating a new file at the given destinatiob
		file, err := os.Create(outputDest)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		//assigning the stdout and stderr to the new file
		os.Stdout = file
		os.Stderr = file
	}
	hasher := code.NewVersionTwo()
	f.VerbosePrint("Checking for archives...")
	path = f.CheckIfPkg(hasher, path)
	hexb, hexs := f.CalculateFvc(hasher, path)
	f.VerbosePrint("Completed sorting and calculating fvc of list  \u2713")
	if err := f.RemovePkg(); err != nil {
		fmt.Fprintln(os.Stderr, "Error removing package", err)
		os.Exit(0)
	}
	f.PrintResult(hexb, hexs)
}

// PrintExamples is used to print out various command line examples when the flag 'e' is set
func PrintExamples() {
	fmt.Fprintln(os.Stdout, "In case of go run:")
	fmt.Fprintln(os.Stdout, "$ go run main.go testdir/libarchive_.tar.gz\n$ go run main.go -v testdir/openssl-3.0.7.tar.gz\n$ go run main.go -b testdir/openssl-3.0.7.tar.gz\n$ go run main.go -o openssl-3.0.7-r2.tar.gz output.txt\n$ go run main.go -o -v openssl-3.0.7 result.txt\n$ go run main.go -e")
	fmt.Fprintln(os.Stdout, "In case of go install:")
	fmt.Fprintln(os.Stdout, "$ fvc_utility testdir/libarchive_.tar.gz\n$ fvc_utility -v testdir/openssl-3.0.7.tar.gz\n$ fvc_utility -b testdir/openssl-3.0.7.tar.gz\n$ fvc_utility -o testdir/openssl-3.0.7-r2.tar.gz output.txt\n$ fvc_utility -o -v openssl-3.0.7 result.txt\n$ fvc_utility -e")
}
