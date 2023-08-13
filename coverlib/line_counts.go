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

const noCount int = -1

// LineCounts stores the hit counts for a file.
type LineCounts struct {
	hitCounts []int
}

// Set the hit count for a line. If the line already has a hit count, the larger
// value is used.
func (lc *LineCounts) Set(lineIdx, hitCount int) {
	for lineIdx >= len(lc.hitCounts) {
		lc.hitCounts = append(lc.hitCounts, noCount)
	}
	if lc.hitCounts[lineIdx] < hitCount {
		lc.hitCounts[lineIdx] = hitCount
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
