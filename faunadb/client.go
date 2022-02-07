package faunadb

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"golang.org/x/net/http2"
)

const (
	apiVersion        = "4"
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

// Headers configures the http headers for a FaunaClient.
func Headers(headers map[string]string) ClientConfig { return func(cli *FaunaClient) { cli.headers = headers } }

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

type faunaRequest struct {
	headers map[string]string
}

type faunaResponse struct {
	response *http.Response
	ctx      context.Context
	cncl     context.CancelFunc
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
	streamEndpoint   string
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
	Event      StreamEvent
	Streaming  bool
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
	if client.queryTimeoutMs <= 0 {
		client.queryTimeoutMs = uint64(requestTimeout / time.Millisecond)
	}
	streamURI := "stream"
	if client.endpoint[len(client.endpoint)-1] != '/' {
		streamURI = "/" + streamURI
	}
	client.streamEndpoint = client.endpoint + streamURI

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
		}
	}

	if client.observer == nil {
		client.observer = func(queryResult *QueryResult) {}
	}

	infoHeaders := map[string]string{
		"Content-Type":             "application/json; charset=utf-8",
		"X-FaunaDB-API-Version":    apiVersion,
		"X-Fauna-Driver":           headerFaunaDriver,
		"X-Runtime-Environment-OS": getRuntimeEnvironmentOs(),
		"X-Runtime-Environment":    getRuntimeEnvironment(),
		"X-GO-Version":             runtime.Version(),
		"X-Query-Timeout":          strconv.FormatUint(client.queryTimeoutMs, 10),
	}
	if len(client.headers) == 0 {
		client.headers = infoHeaders
	} else {
		for k, v := range infoHeaders {
			client.headers[k] = v
		}
	}

	if client.queryTimeoutMs > 0 {
		client.headers["X-Query-Timeout"] = strconv.FormatUint(client.queryTimeoutMs, 10)
	}

	return client
}

func getRuntimeEnvironmentOs() string {
	envOS := runtime.GOOS
	switch envOS {
	case "windows", "darwin", "linux":
		return envOS
	default:
		return "unknown"
	}
}
func getRuntimeEnvironment() string {
	var env = map[string]string{
		"NETLIFY_IMAGES_CDN_DOMAIN":                 "Netlify",
		"VERCEL":                                    "Vercel",
		"AWS_LAMBDA_FUNCTION_VERSION":               "AWS Lambda",
		"GOOGLE_CLOUD_PROJECT":                      "GCP Compute Instances",
		"WEBSITE_FUNCTIONS_AZUREMONITOR_CATEGORIES": "Azure Cloud Functions",
	}
	for k := range env {
		if _, ok := os.LookupEnv(k); ok {
			return env[k]
		}

		if _, ok := os.LookupEnv("PATH"); ok && strings.Contains(os.Getenv("PATH"), ".heroku") {
			return "Heroku"
		}

		if _, ok := os.LookupEnv("_"); ok && strings.Contains(os.Getenv("_"), "google") {
			return "GCP Cloud Functions"
		}

		if _, ok := os.LookupEnv("WEBSITE_INSTANCE_ID"); ok {
			if _, ok = os.LookupEnv("ORYX_ENV_TYPE"); ok &&
				strings.Contains(os.Getenv("ORYX_ENV_TYPE"), "AppService") {

				return "Azure Compute"
			}
		}
	}
	return "Unknown"
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
	var response faunaResponse
	var payload []byte

	startTime := time.Now()

	if payload, err = client.prepareRequestBody(expr); err == nil {
		body := bytes.NewReader(payload)

		response, err = client.performRequest(body, client.endpoint, false, configs)

		httpResponse := response.response

		if httpResponse != nil {
			defer func() {
				_, _ = io.Copy(ioutil.Discard, httpResponse.Body) // Discard remaining bytes so the connection can be reused
				_ = httpResponse.Body.Close()
				response.cncl()
			}()
		}

		if err == nil {
			if err = checkForResponseErrors(httpResponse); err == nil {
				value, err = client.parseResponse(httpResponse, expr, false, startTime)
			}
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

// Stream creates a stream subscription to the result of the given expression.
//
// The subscription returned by this method does not issue any requests until
// the subscription's Start method is called. Make sure to
// subscribe to the events of interest, otherwise the received events are simply
// ignored.
func (client *FaunaClient) Stream(query Expr, config ...StreamConfig) StreamSubscription {
	return newSubscription(client, query, config...)
}

func (client *FaunaClient) startStream(subscription *StreamSubscription) (err error) {
	var payload []byte
	var response faunaResponse

	startTime := time.Now()
	payload, err = client.prepareRequestBody(subscription.query)
	if err != nil {
		return
	}
	body := ioutil.NopCloser(bytes.NewReader(payload))

	var endpoint strings.Builder
	endpoint.WriteString(client.streamEndpoint)

	if len(subscription.config.Fields) > 0 {
		endpoint.WriteString("?fields=")
		count := len(subscription.config.Fields)
		for i := 0; i < count; i++ {
			endpoint.WriteString(string(subscription.config.Fields[i]))
			if i <= count-1 {
				endpoint.WriteString(",")
			}
		}
	}

	response, err = client.performRequest(body, endpoint.String(), true, nil)
	if err != nil {
		return
	}

	httpResponse := response.response
	_ = client.storeLastTxnTime(httpResponse.Header)
	if err = checkForResponseErrors(httpResponse); err != nil {
		httpResponse.Body.Close()
		response.cncl()
		return
	}

	go func() {
		<-subscription.closed
		httpResponse.Body.Close()
		response.cncl()
	}()

	go func() {

		for {
			var obj Obj

			if val, err := client.parseResponse(httpResponse, subscription.query, true, startTime); err != nil {
				if err == io.EOF || err.Error() == "http2: response body closed" {
					subscription.Close()
					break
				}
				subscription.events <- ErrorEvent{
					err: err,
				}
			} else {
				if err = val.Get(&obj); err == nil {
					var event StreamEvent
					if event, err = unMarshalStreamEvent(obj); err == nil {
						client.SyncLastTxnTime(event.Txn())
						subscription.events <- event
					} else {
						subscription.events <- ErrorEvent{
							err: err,
						}
					}
				} else {
					subscription.events <- ErrorEvent{
						err: err,
					}
				}
			}
		}
	}()

	return
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
		streamEndpoint:   client.streamEndpoint,
		headers:          client.headers,
		http:             client.http,
		isTxnTimeEnabled: client.isTxnTimeEnabled,
		queryTimeoutMs:   client.queryTimeoutMs,
		lastTxnTime:      client.lastTxnTime,
		observer:         observer,
	}
}

func (client *FaunaClient) performRequest(body io.Reader, endpoint string, streaming bool, configs []QueryConfig) (response faunaResponse, err error) {
	var request *http.Request
	var timeout = time.Duration(client.queryTimeoutMs) * time.Millisecond
	if streaming {
		response.ctx, response.cncl = context.WithCancel(context.Background())

	} else {
		response.ctx, response.cncl = context.WithTimeout(context.Background(), timeout)
	}
	if request, err = client.prepareRequest(response.ctx, body, endpoint, configs); err == nil {
		response.response, err = client.http.Do(request)
	}

	return
}

func (client *FaunaClient) prepareRequestBody(expr Expr) (payload []byte, err error) {
	payload, err = json.Marshal(expr)
	return
}

func (client *FaunaClient) prepareRequest(ctx context.Context, body io.Reader, endpoint string, configs []QueryConfig) (request *http.Request, err error) {
	if request, err = http.NewRequestWithContext(ctx, "POST", endpoint, body); err == nil {
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

	return
}

func (client *FaunaClient) parseResponse(response *http.Response, expr Expr, streaming bool, startTime time.Time) (value Value, err error) {
	var parsedResponse Value

	if !streaming {
		if err = client.storeLastTxnTime(response.Header); err != nil {
			return
		}
	}

	if parsedResponse, err = parseJSON(response.Body); err == nil {
		if streaming {
			value = parsedResponse
		} else {
			value, err = parsedResponse.At(resource).GetValue()
		}
		client.callObserver(response, expr, streaming, value, startTime)
	} else {
		return nil, err
	}

	return
}

func (client *FaunaClient) callObserver(response *http.Response, expr Expr, streaming bool, value Value, startTime time.Time) {
	var event StreamEvent
	if streaming {
		var obj Obj
		if err := value.Get(&obj); err != nil {
			event = ErrorEvent{
				err: err,
			}
		} else {
			event, _ = unMarshalStreamEvent(obj)
		}

		value = nil
	}
	queryResult := &QueryResult{
		client,
		expr,
		value,
		event,
		streaming,
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
	return fmt.Sprintf("Bearer %s", secret)
}
