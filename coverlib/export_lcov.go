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
	"bufio"
	"fmt"
	"io"
)

// ExportLCOV exports profile data to the LCOV format (see
// https://ltp.sourceforge.net/coverage/lcov/geninfo.1.php).
func ExportLCOV(p *Profiles, writer io.Writer) error {
	w := bufio.NewWriter(writer)

	var err error
	emit := func(s string) {
		if err == nil {
			_, err = w.WriteString(s)
		}
	}
	// Note: this code is similar to Bazel's LCOV converter:
	// https://github.com/bazelbuild/rules_go/blob/84d1a5964f2d92235d1677e8cb9e31eaf9b1b121/go/tools/bzltestutil/lcov.go#L117
	for _, filename := range p.Files() {
		lineCounts := p.LineCounts(filename)
		emit(fmt.Sprintf("SF:%s\n", filename))

		numLines := 0
		numCovered := 0
		lineCounts.ForEach(func(lineIdx, hitCount int) {
			emit(fmt.Sprintf("DA:%d,%d\n", lineIdx, hitCount))
			numLines++
			if hitCount > 0 {
				numCovered++
			}
		})
		emit(fmt.Sprintf("LH:%d\nLF:%d\nend_of_record\n", numCovered, numLines))
	}
	if err == nil {
		err = w.Flush()
	}
	return err
}

type writer struct {
	w   *bufio.Writer
	err error
}

func newWriter(w io.Writer) *writer {
	return &writer{
		w: bufio.NewWriter(w),
	}
}

func (w *writer) Emit(s string) {
	if w.err != nil {
		return
	}
	if _, err := w.w.WriteString(s); err != nil {
		w.err = err
	}
}

func (w *writer) Finish() error {
	if w.err != nil {
		return w.err
	}
	return w.w.Flush()
}
