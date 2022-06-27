package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateMediaType(t *testing.T) {
	cases := []struct {
		path string
		mt    string
	}{
		{ "testdata/resource.txt", "text/plain; charset=utf-8" },
		{ "testdata/resource.csv", "text/csv" },
		{ "testdata/resource.dat", "application/octet-stream" },
		{ "testdata/resource.json", "application/json" },
		{ "testdata/resource.pdf", "application/pdf" },
	}

	for _, tc := range cases {
		t.Run(tc.mt, func(t *testing.T) {
			data, err := os.ReadFile(tc.path)
			require.NoError(t, err)

			detected := DetectMediaType(data)

			require.Equal(t, tc.mt, detected)
		})
	}
}
