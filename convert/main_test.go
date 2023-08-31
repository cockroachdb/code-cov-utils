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
	"os"
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestConvert(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		dir, err := os.MkdirTemp("", "testconvert")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(dir)
		var inputFiles []string
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "input":
				var formatStr string
				td.ScanArgs(t, "fmt", &formatStr)
				filename := fmt.Sprintf("%s/%d.%s", dir, len(inputFiles)+1, formatStr)
				if err := os.WriteFile(filename, []byte(td.Input), 0666); err != nil {
					td.Fatalf(t, "%v", err)
				}
				inputFiles = append(inputFiles, filename)
				return ""

			case "convert":
				var formatStr string
				td.ScanArgs(t, "fmt", &formatStr)
				var trimPrefix string
				if td.HasArg("trim-prefix") {
					td.ScanArgs(t, "trim-prefix", &trimPrefix)
				}
				outputFile := fmt.Sprintf("%s/result.%s", dir, formatStr)
				if err := convert(inputFiles, outputFile, trimPrefix); err != nil {
					return fmt.Sprintf("Error: %s", err)
				}
				res, err := os.ReadFile(outputFile)
				if err != nil {
					td.Fatalf(t, "%v", err)
				}
				return string(res)

			default:
				td.Fatalf(t, "unknown command %s", td.Cmd)
				return ""
			}
		})
	})
}
