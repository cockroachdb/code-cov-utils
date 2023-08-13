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
	"encoding/json"
	"io"
)

// ExportCodecovJson exports profile data to the Codecov custom coverage JSON
// format (https://docs.codecov.com/docs/codecov-custom-coverage-format).
//
// Sample output (with added comments):
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
func ExportCodecovJson(p *Profiles, writer io.Writer) error {
	type fileCoverage map[int]int
	o := struct {
		Coverage map[string]fileCoverage `json:"coverage"`
	}{
		Coverage: make(map[string]fileCoverage),
	}
	for _, filename := range p.Files() {
		counts := make(fileCoverage)
		p.LineCounts(filename).ForEach(func(lineIdx, hitCount int) {
			counts[lineIdx] = hitCount
		})
		o.Coverage[filename] = counts
	}
	marshalled, err := json.MarshalIndent(&o, "", "  ")
	if err != nil {
		return err
	}
	_, err = writer.Write(marshalled)
	return err
}
