package ethclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

const (
	DefaultApiVersion = "2.0"
)

const (
	getBlockNumberMethod       = "eth_blockNumber"
	getBlockByNumberMethod     = "eth_getBlockByNumber"
	getTransactionByHashMethod = "eth_getTransactionByHash"
)

type Client interface {
	// Info returns the client info.
	Info() (*Info, error)
	// BlockNumber returns the current block number.
	// It calls the eth_blockNumber RPC method.
	BlockNumber() (int, error)
	// BlockByNumber returns the block with the given block number.
	// It calls the eth_getBlockByNumber RPC method.
	BlockByNumber(blockNumber int) (*Block, error)
}

type client struct {
	endpoint   string
	apiVersion string
}

type Info struct {
	ApiVersion string
	Endpoint   string
}

type Block struct {
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Hash        string `json:"hash"`
	BlockNumber string `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
}

type requestBody struct {
	Jsonrpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	ID      int         `json:"id"`
	Params  interface{} `json:"params"`
}

type baseResponseBody struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
}

type responseBodyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type blockResult struct {
	Hash         string        `json:"hash"`
	Number       string        `json:"number"`
	Transactions []Transaction `json:"transactions"`
}

type blockNumberResponseBody struct {
	baseResponseBody
	Result string            `json:"result"`
	Error  responseBodyError `json:"error"`
}

type blockResponseBody struct {
	baseResponseBody
	Result blockResult       `json:"result"`
	Error  responseBodyError `json:"error"`
}

func New(endpoint string) Client {
	return &client{
		endpoint:   endpoint,
		apiVersion: DefaultApiVersion,
	}
}

// Info returns the client info.
func (c *client) Info() (*Info, error) {
	return &Info{
		ApiVersion: c.apiVersion,
		Endpoint:   c.endpoint,
	}, nil
}

// BlockNumber returns the current block number.
// It calls the eth_blockNumber RPC method.
func (c *client) BlockNumber() (int, error) {
	// create body request
	reqBody := makeReqBody(c.apiVersion, getBlockNumberMethod, []interface{}{})
	// do request
	res, err := c.doReq(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to do request: %v", err)
	}

	// decode response
	var resBody blockNumberResponseBody
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return 0, fmt.Errorf("failed to decode response body: %v", err)
	}

	// validate the result is a hexadecimal number
	if !strings.HasPrefix(resBody.Result, "0x") {
		return 0, fmt.Errorf("result is not a valid quantitity")
	}
	// extract block number from result
	// the first 2 characters are "0x" so we ignore them
	blockNumber, err := strconv.ParseInt(resBody.Result[2:], 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse block number: %v", err)
	}
	// return result
	return int(blockNumber), nil
}

// BlockByNumber returns the block with the given block number.
// It calls the eth_getBlockByNumber RPC method.
func (c *client) BlockByNumber(blockNumber int) (*Block, error) {
	// create body request
	reqBody := makeReqBody(
		c.apiVersion,
		getBlockByNumberMethod,
		[]interface{}{"0x" + strconv.FormatInt(int64(blockNumber), 16), true},
	)
	// do post request
	res, err := c.doReq(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %v", err)
	}
	// decode response
	var resBody blockResponseBody
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %v", err)
	}

	return &Block{
		Number:       resBody.Result.Number,
		Transactions: resBody.Result.Transactions,
	}, nil
}

func (c *client) doReq(reqBody requestBody) (*http.Response, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %v", err)
	}
	// do post request
	res, err := http.Post(c.endpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to do post request: %v", err)
	}
	// check status code
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response status: %v", res.Status)
	}

	return res, nil
}

func makeReqBody(apiVersion, method string, params interface{}) requestBody {
	return requestBody{
		Jsonrpc: apiVersion,
		Method:  method,
		ID:      rand.Int(),
		Params:  params,
	}
}
