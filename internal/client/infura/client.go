package infura

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"interview-test/internal/util/reader"
)

const RPCVersion2 = "2.0"

type Wei = string

type Client interface {
	GetGasPrice(context.Context) (Wei, error)
}

type client struct {
	httpClient                *http.Client
	getGasPriceRequest        *http.Request
	getGasPriceRequestPayload string
}

func NewClient(apiKey string, projectID uint32, httpClient *http.Client) (Client, error) {
	if getGasPriceRequest, payload, err := prepareGetGasPriceRequest(apiKey, projectID); err != nil {
		return nil, err
	} else {
		return &client{
			httpClient:                httpClient,
			getGasPriceRequest:        getGasPriceRequest,
			getGasPriceRequestPayload: payload}, nil
	}
}

func prepareGetGasPriceRequest(apiKey string, projectID uint32) (*http.Request, string, error) {
	payload, err := NewJSONRPCPayloadBuilder(RPCVersion2, projectID).WithMethod("eth_gasPrice").Build()

	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create get gas price request payload")
	}

	rawPayload, err := json.Marshal(payload)

	if err != nil {
		return nil, "", errors.Wrap(err, "failed to marshal get gas price request payload")
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("https://mainnet.infura.io/v3/%s", apiKey), nil)

	if err != nil {
		return nil, "", errors.Wrap(err, "failed to create get gas price request")
	}

	req.Header["Content-Type"] = []string{"application/json"}
	req.ContentLength = int64(len(rawPayload))
	return req, string(rawPayload), nil
}

type JSONRPCPayload struct {
	RPCVersion string `json:"jsonrpc"`
	Method     string `json:"method"`
	Params     []any  `json:"params"`
	Id         uint32 `json:"id"`
}

type JSONRPCResult[T any] struct {
	Result T `json:"result"`
}

type JSONRPCPayloadBuilder struct {
	rpcVersion string
	method     string
	id         uint32
}

func NewJSONRPCPayloadBuilder(rpcVersion string, id uint32) *JSONRPCPayloadBuilder {
	return &JSONRPCPayloadBuilder{rpcVersion: rpcVersion, id: id}
}

func (b *JSONRPCPayloadBuilder) WithMethod(method string) *JSONRPCPayloadBuilder {
	b.method = method
	return b
}

func (b JSONRPCPayloadBuilder) Build() (*JSONRPCPayload, error) {
	if len(b.method) == 0 {
		return nil, errors.New("failed to build JSONRPCPayload: no method provided")
	}

	return &JSONRPCPayload{
		RPCVersion: b.rpcVersion,
		Method:     b.method,
		Params:     []any{},
		Id:         b.id,
	}, nil
}

func (c client) GetGasPrice(ctx context.Context) (Wei, error) {
	if resp, err := c.httpClient.Do(c.getGasPriceRequestWithContext(ctx)); err != nil {
		return "", errors.Wrap(err, "failed to send get gas price request")
	} else {
		defer reader.DrainAndClose(resp.Body)
		return extractGasPriceFromResponse(resp)
	}
}

func (c client) getGasPriceRequestWithContext(ctx context.Context) *http.Request {
	r := c.getGasPriceRequest.WithContext(ctx)
	r.Body = io.NopCloser(strings.NewReader(c.getGasPriceRequestPayload))
	return r
}

func extractGasPriceFromResponse(resp *http.Response) (Wei, error) {
	if resp.StatusCode >= 300 {
		return "", errors.Errorf("get gas price request failed with status [%s] and body [%s]", resp.Status, reader.Summarize(resp.Body, 1000))
	}

	var result JSONRPCResult[Wei]
	err := json.NewDecoder(resp.Body).Decode(&result)
	return result.Result, errors.Wrap(err, "failed to decode result")
}
