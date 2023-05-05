# Verification Code - Utility

A Golang Utility for calculating File Verification Codes which makes use of the 'FVC' and the 'extract' libraries. 
It processes command line arguments from the cli. The extract library is used to extract the input archive or any archives present in the input directory/archive. The FVC library is used to read each file inside the given package and to calculate the SHA256 of each file. Next, it makes use of the FVC library to sort the list of files and to calculate the verification code in Hex string or byte form. Finally, it removes any packages that were extracted during execution and prints the file count and verification code of the given package.


# Getting started

## What's contained in this project

- main.go - contains the source code for the FVC Utility
    - main() - handles the cli and calls all functions
    - PrintExamples() - prints a standard set of examples
    - CheckIfPkg() - checks if initial path is an archive and extracts it
    - CalculateFvc() - handles the process of extraction, reading files and calculating FVC
    - ExtractSubPkg() - checks for sub packages/archives and extracts it
    - ReadDirectory() - read all the files inside the directory
    - ExtractPkg() - extracts the given archive at the given destination
    - RemovePkg() - removes all the packages extracted during execution 
- main_test.go - contains the tests for everything in main.go
- testdir - contains various packages for testing

## Build

```shell
go build -mod=vendor
```

## Run

```shell
$ go run main.go [OPTION] <path>
```
or
```shell
$ fvc_utility [OPTION] <path>
```
### For getting FVC in bytes
```shell
$ go run main.go -b <path>
```
### For executing in verbose mode
```shell
$ go run main.go -v <path>
```
### For viewing example commands
```shell
$ go run main.go -e
```
Note: 'path' can be of a directory or an archive. In case of bad archives, the archives are treated as files instead.

### For Windows
Use
```shell
GOOS="windows" GOARCH="amd64" CC=x86_64-w64-mingw32-gcc
```
before any go command such as:
```shell
GOOS="windows" GOARCH="amd64" CC=x86_64-w64-mingw32-gcc go build main.go
```

## Test

### For testing
```shell
$ go test
```
### For testing in verbose mode
```shell
$ go test -v
```
## Output

### File Count

The total number of files present in the given archive/directory.

### Verification Code

The Verification code for the given archive/directory in bytes or Hex string depending upon choosen flags.