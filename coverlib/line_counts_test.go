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
	expect := func(lc LineCounts, m counts) {
		t.Helper()
		var expLines []int
		for i := range m {
			expLines = append(expLines, i)
		}
		sort.Ints(expLines)
		var gotLines []int
		var gotHitCounts []int

		lc.ForEach(func(lineIdx, hitCount int) {
			gotLines = append(gotLines, lineIdx)
			gotHitCounts = append(gotHitCounts, hitCount)
		})
		if !reflect.DeepEqual(expLines, gotLines) {
			t.Fatalf("invalid line set.\nexpected: %v\ngot: %v", expLines, gotLines)
		}
		for i, lineIdx := range gotLines {
			if hitCount := gotHitCounts[i]; hitCount != m[lineIdx] {
				t.Fatalf("invalid hit count for %d: expected %d, got %d", lineIdx, m[lineIdx], hitCount)
			}
		}
	}
	expect(lc, counts{})
	lc.Set(100, 10)
	expect(lc, counts{100: 10})
	lc.Set(1, 0)
	expect(lc, counts{1: 0, 100: 10})
	lc.Set(1, 2)
	expect(lc, counts{1: 2, 100: 10})
	lc.Set(1, 1)
	expect(lc, counts{1: 2, 100: 10})
	lc.Set(50, 1)
	lc.Set(500, 5)
	expect(lc, counts{1: 2, 50: 1, 100: 10, 500: 5})

	var other LineCounts
	other.CopyFrom(&lc)
	expect(other, counts{1: 2, 50: 1, 100: 10, 500: 5})

	lc.Set(1, 200)
	expect(lc, counts{1: 200, 50: 1, 100: 10, 500: 5})
	expect(other, counts{1: 2, 50: 1, 100: 10, 500: 5})

	lc.Reset()
	expect(lc, counts{})
	lc.CopyFrom(&other)
	expect(lc, counts{1: 2, 50: 1, 100: 10, 500: 5})

	other.Reset()
	other.Set(1, 0)
	other.Set(50, 50)
	other.Set(1000, 1)
	lc.MergeWith(&other)
	expect(lc, counts{1: 2, 50: 51, 100: 10, 500: 5, 1000: 1})
}
