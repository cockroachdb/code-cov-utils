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
	"fmt"
	"strings"
)

const noCount int = -1

// LineCounts stores the hit counts for a file.
type LineCounts struct {
	hitCounts []int
}

// Set the hit count for a line. If the line already has a hit count, the larger
// value is used.
func (lc *LineCounts) Set(lineIdx, hitCount int) {
	lc.ensureSize(lineIdx + 1)
	if lc.hitCounts[lineIdx] < hitCount {
		lc.hitCounts[lineIdx] = hitCount
	}
}

func (lc *LineCounts) ensureSize(n int) {
	// Safety guard in case of corrupt data.
	if n > 10000000 {
		panic(fmt.Sprintf("desired size too large: %d", n))
	}
	for len(lc.hitCounts) < n {
		lc.hitCounts = append(lc.hitCounts, noCount)
	}
}

// ForEach runs the given function for each line that has a hit count (in
// increasing lineIdx order).
func (lc *LineCounts) ForEach(fn func(lineIdx, hitCount int)) {
	for i, c := range lc.hitCounts {
		if c != noCount {
			fn(i, c)
		}
	}
}

// Reset deletes all counts.
func (lc *LineCounts) Reset() {
	lc.hitCounts = lc.hitCounts[:0]
}

// CopyFrom copies the given counts.
func (lc *LineCounts) CopyFrom(other *LineCounts) {
	lc.hitCounts = append(lc.hitCounts[:0], other.hitCounts...)
}

// MergeWith merges in the given line counts.
//
// For any given line, the resulting hit count is the sum between the two hit
// counts.
func (lc *LineCounts) MergeWith(other *LineCounts) {
	lc.ensureSize(len(other.hitCounts))
	other.ForEach(func(lineIdx, hitCount int) {
		if lc.hitCounts[lineIdx] == noCount {
			lc.hitCounts[lineIdx] = 0
		}
		lc.hitCounts[lineIdx] += hitCount
	})
}

func (lc *LineCounts) String() string {
	return lc.StringWithSeparator(", ")
}

// StringWithSeparator generates a string representation showing groups of lines
// with the same hit count, separated by the given separator.
func (lc *LineCounts) StringWithSeparator(sep string) string {
	// We will RLE compress the counts and emit each "block" as a string.
	var lastStart, lastEnd, lastCount int
	var blocks []string
	maybeEmit := func() {
		if lastStart == 0 {
			return
		}
		var str string
		if lastStart == lastEnd {
			str = fmt.Sprintf("%d:%d", lastStart, lastCount)
		} else {
			str = fmt.Sprintf("%d-%d:%d", lastStart, lastEnd, lastCount)
		}
		blocks = append(blocks, str)
	}
	lc.ForEach(func(line, count int) {
		if lastStart != 0 && lastEnd == line-1 && lastCount == count {
			lastEnd = line
			return
		}
		maybeEmit()
		lastStart = line
		lastEnd = line
		lastCount = count
	})
	maybeEmit()
	return strings.Join(blocks, sep)
}
