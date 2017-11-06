package archive

import (
	"testing"
)

func TestPackagePathname(t *testing.T) {
	cases := []struct {
		in, out string
	}{
		{"https://nccwsc.usgs.gov/sites/default/files/files/ACCCNRS_Report_2015.pdf", "/nccwsc.usgs.gov/sites/default/files/files/ACCCNRS_Report_2015.pdf"},
	}

	for i, c := range cases {
		got := PackagePathName(c.in)
		if got != c.out {
			t.Errorf("case %d result mismatch. expected: '%s', got: '%s'", i, c.out, got)
			continue
		}
	}
}
