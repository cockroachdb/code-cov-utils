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

package coverlib

import (
	"bytes"
	"fmt"
	"github.com/cockroachdb/datadriven"
	"strings"
	"testing"
)

func TestCoverlib(t *testing.T) {
	datadriven.Walk(t, "testdata", func(t *testing.T, path string) {
		var p Profiles
		datadriven.RunTest(t, path, func(t *testing.T, td *datadriven.TestData) string {
			switch td.Cmd {
			case "set":
				if len(td.CmdArgs) != 1 {
					td.Fatalf(t, "usage: set <filename>")
				}
				file := td.CmdArgs[0].String()
				lc := p.LineCounts(file)
				if td.Input != "" {
					for _, l := range strings.Split(td.Input, "\n") {
						var lineIdx, hitCount int
						if _, err := fmt.Sscanf(l, "%d %d", &lineIdx, &hitCount); err != nil {
							td.Fatalf(t, "%v", err)
						}
						lc.Set(lineIdx, hitCount)
					}
				}
				return p.String()

			case "export":
				var formatStr string
				td.ScanArgs(t, "fmt", &formatStr)
				format, err := FormatFromFilename("some-name." + formatStr)
				if err != nil {
					td.Fatalf(t, "%v", err)
				}
				var buf bytes.Buffer
				if err := Export(&p, format, &buf); err != nil {
					td.Fatalf(t, "%v", err)
				}
				return buf.String()

			case "import":
				var formatStr string
				td.ScanArgs(t, "fmt", &formatStr)
				format, err := FormatFromFilename("some-name." + formatStr)
				if err != nil {
					td.Fatalf(t, "%v", err)
				}
				res, err := Import(format, strings.NewReader(td.Input))
				if err != nil {
					td.Fatalf(t, "%v", err)
				}
				if td.HasArg("merge") {
					p.MergeWith(res)
				} else {
					p = *res
				}
				return p.String()

			default:
				td.Fatalf(t, "unknown command %s", td.Cmd)
				return ""
			}
		})
	})
}
