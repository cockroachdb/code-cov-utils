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
	"io"
	"strings"
)

type Format int

const (
	FormatUnset Format = iota
	FormatGoCover
	FormatLCOV
	FormatCodecovJSON
)

// FormatFromFilename determines the format from the extension of the filename.
// Supported extensions are: .gocov, .lcov, .json
func FormatFromFilename(filename string) (Format, error) {
	switch {
	case strings.HasSuffix(filename, ".gocov"):
		return FormatGoCover, nil
	case strings.HasSuffix(filename, ".lcov"):
		return FormatLCOV, nil
	case strings.HasSuffix(filename, ".json"):
		return FormatCodecovJSON, nil
	default:
		return 0, fmt.Errorf("could not determine format for filename %q; supported extensions are .out, .lcov, .json", filename)
	}
}

// Import coverage data from the given format.
func Import(format Format, reader io.Reader) (*Profiles, error) {
	switch format {
	case FormatGoCover:
		return ImportGoCover(reader)
	case FormatLCOV:
		return ImportLCOV(reader)
	case FormatCodecovJSON:
		return nil, fmt.Errorf("import from Codecov JSON not supported")
	default:
		return nil, fmt.Errorf("invalid format %d", format)
	}
}

// Export coverage data to the given format.
func Export(p *Profiles, format Format, writer io.Writer) error {
	switch format {
	case FormatGoCover:
		return fmt.Errorf("export to Go cover not supported")
	case FormatLCOV:
		return ExportLCOV(p, writer)
	case FormatCodecovJSON:
		return ExportCodecovJson(p, writer)
	default:
		return fmt.Errorf("invalid format %d", format)
	}
}
