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
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cockroachdb/code-cov-utils/coverlib"
	"golang.org/x/tools/cover"
)

// gocover2json is a program that converts from go cover profile format to the
// Codecov JSON format
// (https://docs.codecov.com/docs/codecov-custom-coverage-format).
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
	var trimPrefix string
	flag.StringVar(&trimPrefix, "trim-prefix", "", "trim prefix from filenames")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Converts a go cover profile to a codecov json file\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] <input-profile> <output.json>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	gocoverFile := flag.Arg(0)
	jsonFile := flag.Arg(1)
	out, err := os.Create(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %q: %v\n", jsonFile, err)
		os.Exit(2)
	}
	if err := convertGocoverToJson(gocoverFile, out, trimPrefix); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting %q: %v\n", gocoverFile, err)
		os.Exit(2)
	}
	if err := out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing %q: %v\n", jsonFile, err)
		os.Exit(2)
	}
}

func convertGocoverToJson(gocoverFile string, jsonWriter io.Writer, trimPrefix string) error {
	profiles, err := cover.ParseProfiles(gocoverFile)
	if err != nil {
		return err
	}
	// The output schema is odd, in that each line is a separate attribute
	// (instead of being part of an array). This makes it hard to use Go's json
	// machinery; we just produce the output directly.
	w := bufio.NewWriter(jsonWriter)
	w.WriteString("{\n")
	w.WriteString("  \"coverage\": {")
	for fileIdx, profile := range profiles {
		var lineCounts coverlib.LineCounts
		for _, b := range profile.Blocks {
			for i := b.StartLine; i <= b.EndLine; i++ {
				lineCounts.Set(i, b.Count)
			}
		}
		if fileIdx > 0 {
			w.WriteString(",")
		}
		w.WriteString("\n")
		fileName := strings.TrimPrefix(profile.FileName, trimPrefix)
		fmt.Fprintf(w, "    %q: {", fileName)
		first := true
		lineCounts.ForEach(func(lineIdx, hitCount int) {
			if !first {
				w.WriteString(",")
			}
			first = false
			w.WriteString("\n")
			fmt.Fprintf(w, "      \"%d\": %d", lineIdx, hitCount)
		})
		w.WriteString("\n    }")
	}
	w.WriteString("\n  }\n")
	w.WriteString("}\n")
	return w.Flush()
}
