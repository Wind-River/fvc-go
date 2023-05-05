# Verification Code - Library

A Golang Library for calculating File Verification Codes.

- Version 2 - Version 2 is a SHA256 with the file signature 'FVC2' + Null. So a hex-encoded Version 2 file verification code using the empty SHA256 would be:
```
4656433200e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
```

# Getting started

## What's contained in this project

- legacy - contains the following:
    - legacy.go - src code for version 0
    - upgrade.go - upgrading version 0 to version 1
    - valid.go - src code for validating
- enum.go - custom type for versions
- v1.go - src code for version 1
- v2.go - src code for version 2
- version.go - checking version of Hex string or bytes

## Algorithm

- Collect a list of every regular file in your collection.
    - Ignore: directories, Named Pipes, Sockets, Block Device Files, Character Device Files, or Symbolic Links.
- Calculate the sha256 of every file.
- Sort, ascending, the list of sha256s.
- Calculate the sha256 of the list of sha256s.
- Append 'FVC2' + Null to the front of this sha256 and return.


