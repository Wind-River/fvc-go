// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//       http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.

package extract

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"

	"github.com/pkg/errors"
)

/*
#cgo pkg-config: libarchive
#cgo pkg-config: libcrypto
#include "lib.h"
#include "stdlib.h"
#include "stdio.h"
*/
import "C"

var Extensions []string
var MabyeExtensions []string

const (
	archiveEOF    = 1   // Found end of archive.
	archiveOK     = 0   // Operation was successful.
	archiveRetry  = -10 // Retry might succeed.
	archiveWarn   = -20 // Partial success.
	archiveFailed = -25 // Current operation cannot complete.
	archiveFatal  = -30 // No more operations are possible.

	archiveOpenError  = -100 // Error on opening file
	archiveNextError  = -200 // Error on reading next header
	archiveCopyError  = -300 // Error from copy_data
	archiveWriteError = -400 // Error finishing archive write
)

var Logger zerolog.Logger

func init() {
	Extensions = []string{".ar", ".arj", ".cpio", ".dump", ".jar", ".7z", ".zip", ".pack", ".pack2000", ".tar", ".bz2", ".gz", ".lzma", ".snz", ".xz", ".z", ".tgz", ".rpm", ".gem", ".deb", ".whl", ".apk"}
	Logger = zerolog.New(os.Stderr)
}

// IsExtractable looks at the file extension, and possibly the context of files around it, to guess whether that file is an extractable file
func IsExtractable(file string) float64 {
	fileName, _, ext := SplitExt(file)

	if ext == ".pack" {
		var hasIdx bool
		var inObjectsDir bool

		if _, err := os.Stat(fileName + ".idx"); err == nil {
			hasIdx = true
		}
		parentPath, _ := filepath.Split(file)
		parent := path.Base(parentPath)

		inObjectsDir = parent == "objects"

		ret := 1.0
		if hasIdx {
			ret -= 0.5
		}
		if inObjectsDir {
			ret -= 0.5
		}
		return ret
	}

	for _, v := range Extensions {
		if v == ext {
			return 1.0
		}
	}

	return 0.0
}

// RecognizeExtension checks a file's extension against known archive extensions
func RecognizeExtension(file string) bool {
	_, _, ext := SplitExt(file)

	for _, v := range Extensions {
		if v == ext {
			return true
		}
	}

	return false
}

// SplitExt returns the base file name, full extension (.tar + compression), and extension
// If extension is not a compressed tar, full extension and extension will be identical
func SplitExt(s string) (string, string, string) {
	ext := filepath.Ext(s)
	fileName := strings.TrimSuffix(s, ext)
	var fullExt string

	if filepath.Ext(fileName) == ".tar" {
		fullExt = ".tar" + ext
		fileName = strings.TrimSuffix(fileName, ".tar")
	} else {
		fullExt = ext
	}

	return fileName, fullExt, ext
}

type Extract struct {
	source     string // filepath to archive
	filename   string // filename that is different from filepath; extension is used to aid extraction
	target     string // filepath to extract archive to
	isEnclosed bool   // if true, extracts to a sha256-named directory at target
}

// Getters
// //
func (e Extract) Source() string {
	return e.source
}
func (e Extract) Target() string {
	return e.target
}
func (e Extract) IsEnclosed() bool {
	return e.isEnclosed
}

////

// New creates a new extractor set-up to extract source to current directory
func New(source, filename string) (*Extract, error) {
	if filename == "" {
		filename = source
	}
	if _, err := os.Stat(source); os.IsNotExist(err) {
		return nil, errors.Wrapf(err, "extract.New(%s, %s).Stat(%s)", source, filename, source)
	}

	source, err := filepath.Abs(source)
	if err != nil {
		return nil, errors.Wrapf(err, "extract.New(%s, %s).Abs(%s)", source, filename, source)
	}
	e := Extract{source, filename, "", false}

	return &e, nil
}

// NewAt creates a new extractor set-up to extract source to target directory
func NewAt(source, filename, target string) (*Extract, error) {
	Logger.Debug().Str(zerolog.CallerFieldName, "extract.NewAt()").Str("source", source).Str("filename", filename).Str("target", target).Msg("Making Target Extractor")
	if filename == "" {
		filename = source
	}
	e, err := New(source, filename)
	if err != nil {
		return e, err
	}

	target, err = filepath.Abs(target)
	if err != nil {
		return e, errors.Wrapf(err, "extract.NewAt(%s, %s, %s).Abs(%s)", source, filename, target, target)
	}
	e.target = target

	Logger.Debug().Str(zerolog.CallerFieldName, "extract.NewAt()").Str("source", source).Str("filename", filename).Str("target", target).Str("e.target", e.target).Msg("Making Target Directory")
	if err = os.MkdirAll(e.target, 0755); err != nil && !os.IsExist(err) {
		return e, errors.Wrapf(err, "extract.NewAt(%s, %s, %s).Mkdir(%s)", source, filename, target, target)
	}

	return e, nil
}

// Create a directory at target where the name is the hex-encoded sha256 of the source
func (e *Extract) Enclose() error {
	if e.isEnclosed {
		return nil
	}
	if e.source == "" {
		return errors.New("No source file provided")
	}

	f, err := os.Open(e.source)
	if err != nil {
		return errors.Wrapf(err, "Extract.Enclose().Open(%s)", e.source)
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return errors.Wrapf(err, "EXtract.Enclose().Copy(sha256.New(), %s)", e.source)
	}
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	targetDir := filepath.Join(e.target, hash)
	if err = os.MkdirAll(targetDir, 0755); err != nil && !os.IsExist(err) {
		return errors.Wrapf(err, "Extract.Enclose().MkdirAll(%s)", targetDir)
	}
	e.target = targetDir
	e.isEnclosed = true

	return nil
}

// Extract extracts the source to target
func (e Extract) Extract() (string, error) {
	Logger.Debug().Str(zerolog.CallerFieldName, "extract.Extract{}.Extract()").
		Str("source", e.source).
		Str("filename", e.filename).
		Str("target", e.target).
		Bool("isEnclosed", e.isEnclosed).
		Msg("extracting")

	var cTarget *C.char
	if e.target != "" {
		cTarget = C.CString(e.target)
	}

	cs := C.CString(e.source)
	n, f, ex := SplitExt(e.filename)
	var exit C.status
	if f != ex || (len(ex) > 0 && ex[1] == 'c') { // f contains tar
		exit = C.extractOne(cs, nil, cTarget)
	} else {
		cn := C.CString(n)
		exit = C.extractOne(cs, cn, cTarget)
	}
	defer C.status_free(exit)
	if exit == nil {
		Logger.Debug().Str(zerolog.CallerFieldName, "extract.Extract{}.Extract()").Msg("exit is null")
	}
	Logger.Debug().Str(zerolog.CallerFieldName, "extract.Extract{}.Extract()").Str("code", fmt.Sprintf("%d", exit.code)).Msg("Exited")

	if exit.code < 0 {
		return e.target, errors.New(fmt.Sprintf("extract returned with status: %d\n%s\n%s\n", exit.code, C.GoString(exit.message), C.GoString(exit.tag)))
	}

	// ret := C.GoString(exit.tag)
	ret := e.Target()
	Logger.Debug().Str(zerolog.CallerFieldName, "extract.Extract{}.Extract()").Interface("extract", e).Str("target", ret).Msg("extracted")
	return ret, nil
}
