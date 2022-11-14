package infura

import (
	"context"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var hexadecimalStringRegex = regexp.MustCompile(`^(?i)0x[\da-f]+$`)

func TestItGetGasPriceFromInfura(t *testing.T) {
	client, err := NewClient(os.Getenv("INFURA_API_KEY"), 1, &http.Client{})
	require.NoError(t, err)

	wei, err := client.GetGasPrice(context.Background())

	require.NoError(t, err)
	require.Regexp(t, hexadecimalStringRegex, wei)
}

type MockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m MockRoundTripper) RoundTrip(_ *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestItGetGasPriceFromInfuraMocked(t *testing.T) {
	client, err := NewClient(
		os.Getenv("INFURA_API_KEY"),
		1,
		&http.Client{
			Transport: MockRoundTripper{
				Response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"result": "0xaf234"}`))}}})
	require.NoError(t, err)

	wei, err := client.GetGasPrice(context.Background())

	require.NoError(t, err)
	require.Equal(t, "0xaf234", wei)
}
