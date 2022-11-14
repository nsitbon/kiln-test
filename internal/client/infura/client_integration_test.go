//go:build integration
// +build integration

package infura

import (
	"context"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

var hexadecimalStringRegex = regexp.MustCompile(`^(?i)0x[\da-f]+$`)

func TestItGetGasPriceFromInfura(t *testing.T) {
	client, err := NewClient(os.Getenv("INFURA_API_KEY"), 1, &http.Client{})
	require.NoError(t, err)

	wei, err := client.GetGasPrice(context.Background())

	require.NoError(t, err)
	require.Regexp(t, hexadecimalStringRegex, string(wei))
}
