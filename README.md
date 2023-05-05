# File Verification Code

File Verification Code is a tool to find whether two given packages/archives are equivalent, i.e, the files inside the two should be exactly same. It traverses the entire package/archive and calculates the SHA256 of each file. These are added to a list which is sorted and finally the SHA256 of this list is calculated. Any package/archive containing exactly the same files will have the same verification code no matter the structure or name of the package/archive. File verification code is useful as an identifier or unique id for a given package/archive. It can also be used in cases where packages/archives must be compared or handled for duplicacy.  


A Golang Library and Utility for calculating File Verification Codes can be found in this repository. Read More about them:

- [Library](code/README.md)
- [Utility](fvc_utility/README.md)

# Getting started

## What's contained in this project

- code - contains the source code for Fvc library
- fvc_utility - contains the source code for Fvc utility
- hasher - Contains the interface used for Fvc

## Dependencies

File Verification code uses the extract library for extracting archives and sub-archives. The extract library makes use of the libarchive package which makes it dependent on CGO. The extract library is provided as a vendor in the vendor directory.



