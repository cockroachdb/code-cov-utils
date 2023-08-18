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
	"io"

	"golang.org/x/tools/cover"
)

// ImportGoCover imports go cover profile data.
func ImportGoCover(reader io.Reader) (*Profiles, error) {
	profiles, err := cover.ParseProfilesFromReader(reader)
	if err != nil {
		return nil, err
	}
	p := &Profiles{}
	for _, profile := range profiles {
		lineCounts := p.LineCounts(profile.FileName)
		for _, b := range profile.Blocks {
			for i := b.StartLine; i <= b.EndLine; i++ {
				lineCounts.Set(i, b.Count)
			}
		}
	}
	return p, nil
}
