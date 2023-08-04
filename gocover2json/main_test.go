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
	"bytes"
	"fmt"
	"golang.org/x/tools/cover"
	"os"
	"testing"

	"github.com/cockroachdb/datadriven"
)

func TestLcovToJson(t *testing.T) {
	datadriven.RunTest(t, "testdata/gocover2json", func(t *testing.T, td *datadriven.TestData) string {
		switch td.Cmd {
		case "convert":
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
			profiles, err := cover.ParseProfiles(filename)
			if err != nil {
				td.Fatalf(t, "%v", err)
			}

			out := &bytes.Buffer{}
			if err := convertProfilesToJson(profiles, out); err != nil {
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
