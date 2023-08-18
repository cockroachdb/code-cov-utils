// Copyright 2023 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"fmt"
	"github.com/cockroachdb/code-cov-utils/coverlib"
	"io"
	"os"
)

// lcov2json is a program that converts from LCOV format [1] to the Codecov
// custom coverage JSON format [2].
// [1] https://ltp.sourceforge.net/coverage/lcov/geninfo.1.php
// [2] https://docs.codecov.com/docs/codecov-custom-coverage-format
//
// Example of the output format:
//
//	{
//	  "coverage": {
//	    "filename": {
//	      "1": 0,      # line 1 missed
//	      "2": 1,      # line 2 hit once
//	      "7": 5       # line 7 hit 5 times
//	    }
//	  }
//	}
func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input.dat> <output.json>\n", os.Args[0])
		os.Exit(1)
	}
	lcovFile := os.Args[1]
	jsonFile := os.Args[2]
	in, err := os.Open(lcovFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening %q: %v", lcovFile, err)
		os.Exit(2)
	}
	defer in.Close()
	out, err := os.Create(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %q: %v", jsonFile, err)
		os.Exit(2)
	}
	if err := convertLcovToJson(in, out); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	if err := out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing %q: %v", jsonFile, err)
		os.Exit(2)
	}
}

func convertLcovToJson(lcovReader io.Reader, jsonWriter io.Writer) error {
	p, err := coverlib.ImportLCOV(lcovReader)
	if err != nil {
		return err
	}
	return coverlib.ExportCodecovJson(p, jsonWriter)
}
