// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package main

import (
	"bufio"
	"fmt"
	"golang.org/x/tools/cover"
	"io"
	"os"
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
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <input-profile> <output.json>\n", os.Args[0])
		os.Exit(1)
	}
	gocoverFile := os.Args[1]
	jsonFile := os.Args[2]
	profiles, err := cover.ParseProfiles(gocoverFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %q: %v", gocoverFile, err)
		os.Exit(2)
	}
	out, err := os.Create(jsonFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %q: %v", jsonFile, err)
		os.Exit(2)
	}
	if err := convertProfilesToJson(profiles, out); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
	if err := out.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing %q: %v", jsonFile, err)
		os.Exit(2)
	}
}

// lineCount is the hit count for a line, or -1 if the line has no count.
type lineCount int

const noCount lineCount = -1

func convertProfilesToJson(profiles []*cover.Profile, jsonWriter io.Writer) error {
	// The output schema is odd, in that each line is a separate attribute
	// (instead of being part of an array). This makes it hard to use Go's json
	// machinery; we just produce the output directly.
	w := bufio.NewWriter(jsonWriter)
	w.WriteString("{\n")
	w.WriteString("  \"coverage\": {")
	for fileIdx, profile := range profiles {
		var lines []lineCount
		for _, b := range profile.Blocks {
			for len(lines) <= b.EndLine {
				lines = append(lines, -1)
			}
			for i := b.StartLine; i <= b.EndLine; i++ {
				lines[i] = lineCount(b.Count)
			}
		}
		if fileIdx > 0 {
			w.WriteString(",")
		}
		w.WriteString("\n")
		fmt.Fprintf(w, "    %q: {", profile.FileName)
		first := true
		for i, count := range lines {
			if count < 0 {
				continue
			}
			if !first {
				w.WriteString(",")
			}
			first = false
			w.WriteString("\n")
			fmt.Fprintf(w, "      \"%d\": %d", i, count)
		}
		w.WriteString("\n    }")
	}
	w.WriteString("\n  }\n")
	w.WriteString("}\n")
	return w.Flush()
}
