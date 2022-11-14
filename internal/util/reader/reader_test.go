package reader

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestItSummarizeReader(t *testing.T) {
	require.Equal(t, "a more tha", Summarize(strings.NewReader("a more than 10 characters string"), 10))
}
