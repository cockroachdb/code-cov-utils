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
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestLcovToJson(t *testing.T) {
	datadriven.RunTest(t, "testdata/gocover2json", func(t *testing.T, td *datadriven.TestData) string {
		switch td.Cmd {
		case "convert":
			var trimPrefix string
			for _, arg := range td.CmdArgs {
				switch arg.Key {
				case "trim-prefix":
					trimPrefix = arg.Vals[0]
				default:
					td.Fatalf(t, "unknown argument %q", arg.Key)
				}
			}

			f, err := os.CreateTemp("", "profile")
			if err != nil {
				td.Fatalf(t, "%v", err)
			}
			filename := f.Name()
			if _, err := f.WriteString(td.Input); err != nil {
				td.Fatalf(t, "%v", err)
			}
			if err := f.Close(); err != nil {
				td.Fatalf(t, "%v", err)
			}

			out := &bytes.Buffer{}
			if err := convertGocoverToJson(filename, out, trimPrefix); err != nil {
				return fmt.Sprintf("error: %v", err)
			}
			result := out.String()
			return result

		default:
			td.Fatalf(t, "unknown command %s", td.Cmd)
			return ""
		}
	})
}
