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
	"sort"
)

// Profiles stores LineCounts for a collection of files.
type Profiles struct {
	m map[string]*LineCounts
}

// LineCounts returns the LineCounts for the given file, adding the file to the
// collection if necessary.
func (p *Profiles) LineCounts(filename string) *LineCounts {
	lc := p.m[filename]
	if lc == nil {
		lc = &LineCounts{}
		if p.m == nil {
			p.m = make(map[string]*LineCounts)
		}
		p.m[filename] = lc
	}
	return lc
}

// Files returns all filenames in the collection (sorted).
func (p *Profiles) Files() []string {
	res := make([]string, 0, len(p.m))
	for filename := range p.m {
		res = append(res, filename)
	}
	sort.Strings(res)
	return res
}

// RenameFiles changes the names of the files in the profile.
func (p *Profiles) RenameFiles(renameFn func(filenameBefore string) string) {
	oldMap := p.m
	p.m = make(map[string]*LineCounts, len(p.m))
	for f, lc := range oldMap {
		// If a file with this name already exists, we merge the profiles. This
		// could happen if we are merging data from different profiles and some have
		// a different prefix.
		p.LineCounts(renameFn(f)).MergeWith(lc)
	}
}

func (p *Profiles) String() string {
	var buf bytes.Buffer
	for _, f := range p.Files() {
		fmt.Fprintf(&buf, "%s\n", f)
		if countsStr := p.LineCounts(f).StringWithSeparator("\n  "); countsStr != "" {
			fmt.Fprintf(&buf, "  %s\n", countsStr)
		}
	}
	return buf.String()
}

// MergeWith merges in the given profiles.
//
// For any given line, the resulting hit count is the maximum between the hit
// counts in the two profiles. This is consistent with how Bazel merges coverage
// profiles; see:
// https://github.com/bazelbuild/rules_go/blob/84d1a5964f2d92235d1677e8cb9e31eaf9b1b121/go/tools/bzltestutil/lcov.go#L104
func (p *Profiles) MergeWith(other *Profiles) {
	for _, filename := range other.Files() {
		p.LineCounts(filename).MergeWith(other.LineCounts(filename))
	}
}
