package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppendString(t *testing.T) {
	list := []string{"foo", "bar"}

	appended := AppendUniqueString(list, "baz")
	require.Equal(t, appended[2], "baz")
	appended = AppendUniqueString(list, "baz")
	require.Equal(t, len(appended), 3)
}
