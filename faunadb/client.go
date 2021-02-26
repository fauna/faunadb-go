package faunadb

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

const (
	apiVersion        = "3"
	defaultEndpoint   = "https://db.fauna.com"
	requestTimeout    = 60 * time.Second
	headerTxnTime     = "X-Txn-Time"
	headerLastSeenTxn = "X-Last-Seen-Txn"
	headerFaunaDriver = "go"
)

var resource = ObjKey("resource")

// ClientConfig is the base type for the configuration parameters of a FaunaClient.
type ClientConfig func(*FaunaClient)

// Endpoint configures the FaunaDB URL for a FaunaClient.
func Endpoint(url string) ClientConfig { return func(cli *FaunaClient) { cli.endpoint = url } }

// HTTP allows the user to override the http.Client used by a FaunaClient.
func HTTP(http *http.Client) ClientConfig { return func(cli *FaunaClient) { cli.http = http } }

/*
EnableTxnTimePassthrough configures the FaunaClient to keep track of the last seen transaction time.
The last seen transaction time is used to avoid reading stale data from outdated replicas when
reading and writing from different nodes at the same time.

Deprecated: This function is deprecated since this feature is enabled by default.
*/
func EnableTxnTimePassthrough() ClientConfig {
	return func(cli *FaunaClient) { cli.isTxnTimeEnabled = true }
}

// QueryTimeoutMS sets the server timeout for ALL queries issued with this client.
// This is not the http request timeout.
func QueryTimeoutMS(millis uint64) ClientConfig {
	return func(cli *FaunaClient) { cli.queryTimeoutMs = millis }
}

/*
DisableTxnTimePassthrough configures the FaunaClient to not keep track of the last seen transaction time.
The last seen transaction time is used to avoid reading stale data from outdated replicas when
reading and writing from different nodes at the same time.

Disabling this option might lead to data inconsistencies and is not recommended. If don't know what you're
doing leave this alone. Use at your own risk.
*/
func DisableTxnTimePassthrough() ClientConfig {
	return func(cli *FaunaClient) { cli.isTxnTimeEnabled = false }
}

//QueryConfig is the base type for query specific configuration parameters.
type QueryConfig func(*faunaRequest)

// TimeoutMS sets the server timeout for a specific query.
// This is not the http request timeout.
func TimeoutMS(millis uint64) QueryConfig {
	return func(req *faunaRequest) {
		req.headers["X-Query-Timeout"] = strconv.FormatUint(millis, 10)
	}
}

type faunaRequest struct {
	headers map[string]string
}

// ObserverCallback is the callback type for requests.
type ObserverCallback func(*QueryResult)

// Observer configures a callback function called for every query executed.
func Observer(observer ObserverCallback) ClientConfig {
	return func(cli *FaunaClient) { cli.observer = observer }
}

/*
FaunaClient provides methods for performing queries on a FaunaDB cluster.

This structure should be reused as much as possible. Avoid copying this structure.
If you need to create a client with a different secret, use the NewSessionClient method.
*/
type FaunaClient struct {
	basicAuth        string
	endpoint         string
	http             *http.Client
	isTxnTimeEnabled bool
	lastTxnTime      int64
	queryTimeoutMs   uint64
	observer         ObserverCallback
	headers          map[string]string
}

// QueryResult is a structure containing the result context for a given FaunaDB query.
type QueryResult struct {
	Client     *FaunaClient
	Query      Expr
	Result     Value
	StatusCode int
	Headers    map[string][]string
	StartTime  time.Time
	EndTime    time.Time
}

/*
NewFaunaClient creates a new FaunaClient structure. Possible configuration options:
	Endpoint: sets a specific FaunaDB url. Default: https://db.fauna.com
	HTTP: sets a specific http.Client. Default: a new net.Client with 60 seconds timeout.
*/
func NewFaunaClient(secret string, configs ...ClientConfig) *FaunaClient {
	client := &FaunaClient{basicAuth: basicAuth(secret), isTxnTimeEnabled: true}

	for _, config := range configs {
		config(client)
	}

	if client.endpoint == "" {
		client.endpoint = defaultEndpoint
	}

	if client.http == nil {
		dial := func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		}
		if len(client.endpoint) > 5 && client.endpoint[:5] == "https" {
			dial = nil
		}
		transport := &http2.Transport{
			DialTLS:   dial,
			AllowHTTP: true,
		}

		client.http = &http.Client{
			Transport: transport,
			Timeout: requestTimeout,
		}
	}

	if client.observer == nil {
		client.observer = func(queryResult *QueryResult) {}
	}

	client.headers = map[string]string{
		"Content-Type":          "application/json; charset=utf-8",
		"X-FaunaDB-API-Version": apiVersion,
		"X-Fauna-Driver":        headerFaunaDriver,
	}

	if client.queryTimeoutMs > 0 {
		client.headers["X-Query-Timeout"] = strconv.FormatUint(client.queryTimeoutMs, 10)
	} else {
		client.queryTimeoutMs = uint64(requestTimeout.Milliseconds())
		client.headers["X-Query-Timeout"] = strconv.FormatUint(uint64(requestTimeout.Milliseconds()), 10)
	}

	return client
}

// QueryResult run and return the cost headers associated with this query.
func (client *FaunaClient) QueryResult(expr Expr) (value Value, headers map[string][]string, err error) {
	value, err = client.NewWithObserver(func(queryResult *QueryResult) {
		headers = queryResult.Headers
	}).Query(expr)

	return
}

// BatchQueryResult run and return the cost headers associated with this query.
func (client *FaunaClient) BatchQueryResult(expr []Expr) (value []Value, headers map[string][]string, err error) {
	value, err = client.NewWithObserver(func(queryResult *QueryResult) {
		headers = queryResult.Headers
	}).BatchQuery(expr)

	return
}

// Query is the primary method used to send a query language expression to FaunaDB.
func (client *FaunaClient) Query(expr Expr, configs ...QueryConfig) (value Value, err error) {
	startTime := time.Now()
	response, err := client.performRequest(expr, configs)

	if response != nil {
		defer func() {
			_, _ = io.Copy(ioutil.Discard, response.Body) // Discard remaining bytes so the connection can be reused
			_ = response.Body.Close()
		}()
	}

	if err == nil {
		if err = checkForResponseErrors(response); err == nil {
			value, err = client.parseResponse(response, expr, startTime)
		}
	}

	return
}

// BatchQuery will sends multiple simultaneous queries to FaunaDB. values are returned in the same order
// as the queries.
func (client *FaunaClient) BatchQuery(exprs []Expr) (values []Value, err error) {
	arr := make(unescapedArr, len(exprs))

	for i, expr := range exprs {
		arr[i] = expr
	}

	var res Value

	if res, err = client.Query(arr); err == nil {
		err = res.Get(&values)
	}

	return
}

// NewSessionClient creates a new child FaunaClient with a new secret. The returned client reuses its parent's internal http resources.
func (client *FaunaClient) NewSessionClient(secret string) *FaunaClient {
	return client.newClient(basicAuth(secret), client.observer)
}

// NewWithObserver creates a new FaunaClient with a specific observer callback. The returned client reuses its parent's internal http resources.
func (client *FaunaClient) NewWithObserver(observer ObserverCallback) *FaunaClient {
	return client.newClient(client.basicAuth, observer)
}

// GetLastTxnTime gets the freshest timestamp reported to this client.
func (client *FaunaClient) GetLastTxnTime() int64 {
	if client.isTxnTimeEnabled {
		return client.lastTxnTime
	}
	return 0
}

// SyncLastTxnTime syncs the freshest timestamp seen by this client.
// This has no effect if more stale than the currently stored timestamp.
// WARNING: This should be used only when coordinating timestamps across
//          multiple clients. Moving the timestamp arbitrarily forward into
//          the future will cause transactions to stall.
func (client *FaunaClient) SyncLastTxnTime(newTxnTime int64) {
	if client.isTxnTimeEnabled {
		for {
			oldTxnTime := atomic.LoadInt64(&client.lastTxnTime)
			if oldTxnTime >= newTxnTime ||
				atomic.CompareAndSwapInt64(&client.lastTxnTime, oldTxnTime, newTxnTime) {
				break
			}
		}
	}
}

func (client *FaunaClient) newClient(basicAuth string, observer ObserverCallback) *FaunaClient {
	return &FaunaClient{
		basicAuth:        basicAuth,
		endpoint:         client.endpoint,
		headers:          client.headers,
		http:             client.http,
		isTxnTimeEnabled: client.isTxnTimeEnabled,
		queryTimeoutMs:   client.queryTimeoutMs,
		lastTxnTime:      client.lastTxnTime,
		observer:         observer,
	}
}

func (client *FaunaClient) performRequest(expr Expr, configs []QueryConfig) (response *http.Response, err error) {
	var request *http.Request
	if request, err = client.prepareRequest(expr, configs); err == nil {
		response, err = client.http.Do(request)
	}

	return
}

func (client *FaunaClient) prepareRequest(expr Expr, configs []QueryConfig) (request *http.Request, err error) {
	var body []byte

	if body, err = json.Marshal(expr); err == nil {
		if request, err = http.NewRequest("POST", client.endpoint, bytes.NewReader(body)); err == nil {
			request.Header.Add("Authorization", client.basicAuth)
			for k, v := range client.headers {
				request.Header.Add(k, v)
			}

			if len(configs) > 0 {
				req := &faunaRequest{
					headers: map[string]string{},
				}
				for _, config := range configs {
					config(req)
				}
				for k, v := range req.headers {
					request.Header.Add(k, v)
				}
			}

			client.addLastTxnTimeHeader(request)
		}
	}

	return
}

func (client *FaunaClient) parseResponse(response *http.Response, expr Expr, startTime time.Time) (value Value, err error) {
	var parsedResponse Value

	if err = client.storeLastTxnTime(response.Header); err == nil {
		if parsedResponse, err = parseJSON(response.Body); err == nil {
			value, err = parsedResponse.At(resource).GetValue()
			client.callObserver(response, expr, value, startTime)
		}
	}

	return
}

func (client *FaunaClient) callObserver(response *http.Response, expr Expr, value Value, startTime time.Time) {
	queryResult := &QueryResult{
		client,
		expr,
		value,
		response.StatusCode,
		response.Header,
		startTime,
		time.Now(),
	}

	client.observer(queryResult)
}

func (client *FaunaClient) addLastTxnTimeHeader(request *http.Request) {
	if client.isTxnTimeEnabled {
		if lastSeen := atomic.LoadInt64(&client.lastTxnTime); lastSeen != 0 {
			request.Header.Add(headerLastSeenTxn, strconv.FormatInt(lastSeen, 10))
		}
	}
}

func (client *FaunaClient) storeLastTxnTime(header http.Header) (err error) {
	if client.isTxnTimeEnabled {
		var newTxnTime int64

		if newTxnTime, err = parseTxnTimeHeader(header); err == nil {
			client.SyncLastTxnTime(newTxnTime)
		}
	}

	return
}

func parseTxnTimeHeader(header http.Header) (txnTime int64, err error) {
	if lastSeenHeader := header.Get(headerTxnTime); lastSeenHeader != "" {
		txnTime, err = strconv.ParseInt(lastSeenHeader, 10, 64)
	}
	return
}

func basicAuth(secret string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(secret + ":"))
	return fmt.Sprintf("Basic %s", encoded)
}
