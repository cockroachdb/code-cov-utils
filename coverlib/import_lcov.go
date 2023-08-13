package coverlib

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// ImportLCOV imports profile data from LCOV format (see
// https://ltp.sourceforge.net/coverage/lcov/geninfo.1.php).
func ImportLCOV(reader io.Reader) (*Profiles, error) {
	p := &Profiles{}
	lcov := bufio.NewScanner(reader)
	var currentCounts *LineCounts
	for lcov.Scan() {
		l := lcov.Text()
		if l == "end_of_record" {
			if currentCounts == nil {
				return nil, errors.New("end_of_record with no file path")
			}
			currentCounts = nil
			continue
		}
		idx := strings.Index(l, ":")
		if idx == -1 {
			// Don't know how to parse this line; skip.
			fmt.Fprintf(os.Stderr, "Warning: cannot parse %q\n", l)
			continue
		}
		key, val := l[:idx], l[idx+1:]
		switch key {
		case "SF":
			currentCounts = p.LineCounts(val)

		case "DA":
			var line, count int
			_, err := fmt.Sscanf(val, "%d,%d", &line, &count)
			if err != nil {
				return nil, fmt.Errorf("error parsing DA line: %v", err)
			}
			currentCounts.Set(line, count)
		}
	}
	if currentCounts != nil {
		return nil, errors.New("unfinished record")
	}
	if err := lcov.Err(); err != nil {
		return nil, err
	}
	return p, nil
}
