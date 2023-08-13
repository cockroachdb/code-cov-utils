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
	"reflect"
	"sort"
	"testing"
)

func TestLineCounts(t *testing.T) {
	var lc LineCounts

	type counts map[int]int
	expect := func(m counts) {
		t.Helper()
		var expLines []int
		for i := range m {
			expLines = append(expLines, i)
		}
		sort.Ints(expLines)
		var gotLines []int
		lc.ForEach(func(lineIdx, hitCount int) {
			gotLines = append(gotLines, lineIdx)
			if hitCount != m[lineIdx] {
				t.Fatalf("invalid hit count for %d: expected %d, got %d", lineIdx, m[lineIdx], hitCount)
			}
		})
		if !reflect.DeepEqual(expLines, gotLines) {
			t.Fatalf("invalid line set.\nexpected: %v\ngot: %v", expLines, gotLines)
		}
	}
	expect(counts{})
	lc.Set(100, 10)
	expect(counts{100: 10})
	lc.Set(1, 0)
	expect(counts{1: 0, 100: 10})
	lc.Set(1, 2)
	expect(counts{1: 2, 100: 10})
	lc.Set(1, 1)
	expect(counts{1: 2, 100: 10})
	lc.Set(50, 1)
	lc.Set(500, 5)
	expect(counts{1: 2, 50: 1, 100: 10, 500: 5})
}
