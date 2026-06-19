package eth

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	rc  *resty.Client
	url string
}

type rpcRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      int           `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("rpc error (code %d): %s", e.Code, e.Message)
}

type rpcResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  string    `json:"result"`
	Error   *rpcError `json:"error"`
}

func NewClient(url string) *Client {
	rc := resty.New().
		SetTimeout(10*time.Second).
		SetHeader("Content-Type", "application/json").
		SetRetryCount(3).
		SetRetryWaitTime(500 * time.Millisecond)

	return &Client{
		rc:  rc,
		url: url,
	}
}

func (c *Client) BlockNumber() (uint64, error) {
	var out rpcResponse
	resp, err := c.rc.R().
		SetBody(rpcRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "eth_blockNumber",
			Params:  []interface{}{},
		}).
		SetResult(&out).
		Post(c.url)
	if err != nil {
		return 0, fmt.Errorf("eth_blockNumber request failed: %w", err)
	}
	if resp.IsError() {
		return 0, fmt.Errorf("eth_blockNumber unexpected status: %s", resp.Status())
	}
	if out.Error != nil {
		return 0, out.Error
	}

	n, err := parseHexUint64(out.Result)
	if err != nil {
		return 0, fmt.Errorf("parse block number %q: %w", out.Result, err)
	}
	return n, nil
}

func parseHexUint64(s string) (uint64, error) {
	s = strings.TrimPrefix(s, "0x")
	if s == "" {
		return 0, fmt.Errorf("empty hex value")
	}
	n, ok := new(big.Int).SetString(s, 16)
	if !ok {
		return 0, fmt.Errorf("invalid hex value")
	}
	if !n.IsUint64() {
		return 0, fmt.Errorf("hex value overflows uint64")
	}
	return n.Uint64(), nil
}

type blockResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  *struct {
		Hash   string `json:"hash"`
		Number string `json:"number"`
	} `json:"result"`
	Error *rpcError `json:"error"`
}

func (c *Client) LatestBlockHash() (string, error) {
	var out blockResponse
	resp, err := c.rc.R().
		SetBody(rpcRequest{
			JSONRPC: "2.0",
			ID:      1,
			Method:  "eth_getBlockByNumber",
			Params:  []interface{}{"latest", false},
		}).
		SetResult(&out).
		Post(c.url)
	if err != nil {
		return "", fmt.Errorf("eth_getBlockByNumber request failed: %w", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("eth_getBlockByNumber unexpected status: %s", resp.Status())
	}
	if out.Error != nil {
		return "", out.Error
	}
	if out.Result == nil || out.Result.Hash == "" {
		return "", fmt.Errorf("eth_getBlockByNumber returned empty block")
	}
	return out.Result.Hash, nil
}

func (c *Client) RandomInRange(min, max int64) (int64, error) {
	if min > max {
		return 0, fmt.Errorf("invalid range: min %d > max %d", min, max)
	}
	hash, err := c.LatestBlockHash()
	if err != nil {
		return 0, err
	}
	return randomFromSeed(hash, min, max)
}

func (c *Client) RandomsInRange(min, max int64, count int) ([]int64, error) {
	hash, err := c.LatestBlockHash()
	if err != nil {
		return nil, err
	}
	rs := make([]int64, count)
	for i := 0; i < count; i++ {
		rs[i], err = randomFromSeed(hash, min, max)
		if err != nil {
			return nil, err
		}
	}
	return rs, nil
}

func randomFromSeed(seedHex string, min, max int64) (int64, error) {
	seedHex = strings.TrimPrefix(seedHex, "0x")
	seed, ok := new(big.Int).SetString(seedHex, 16)
	if !ok {
		return 0, fmt.Errorf("invalid seed %q", seedHex)
	}
	span := new(big.Int).Sub(big.NewInt(max), big.NewInt(min))
	span.Add(span, big.NewInt(1))

	r := new(big.Int).Mod(seed, span)
	r.Add(r, big.NewInt(min))
	return r.Int64(), nil
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func (c *Client) RandomString(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("invalid length: %d", length)
	}
	hash, err := c.LatestBlockHash()
	if err != nil {
		return "", err
	}
	return randomStringFromSeed(hash, length)
}

func randomStringFromSeed(seedHex string, length int) (string, error) {
	seedHex = strings.TrimPrefix(seedHex, "0x")
	seed, ok := new(big.Int).SetString(seedHex, 16)
	if !ok {
		return "", fmt.Errorf("invalid seed %q", seedHex)
	}

	base := big.NewInt(int64(len(charset)))
	mod := new(big.Int)
	out := make([]byte, length)
	for i := 0; i < length; i++ {
		seed.DivMod(seed, base, mod)
		out[i] = charset[mod.Int64()]
	}
	return string(out), nil
}
