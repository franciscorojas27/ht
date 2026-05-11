package runner

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/stretchr/testify/require"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestNewClient(t *testing.T) {
	client := NewClient(false)
	require.NotNil(t, client)
	require.Nil(t, client.Transport)

	client = NewClient(true)
	require.NotNil(t, client)
	transport, ok := client.Transport.(*http.Transport)
	require.True(t, ok)
	require.NotNil(t, transport.TLSClientConfig)
	require.True(t, transport.TLSClientConfig.InsecureSkipVerify)
}

func TestDoRequestSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("{\"ok\":true}"))
	}))
	defer srv.Close()

	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	res, err := DoRequest(req, time.Now(), false, nil)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, http.StatusCreated, res.StatusCode)
	require.Contains(t, string(res.DumpHeaders), "201")
	require.Equal(t, "{\"ok\":true}", string(res.Body))
	require.Equal(t, "application/json", res.Headers.Get("Content-Type"))
}

func TestDoRequestError(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.invalid", nil)
	require.NoError(t, err)

	client := &http.Client{Transport: roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	})}

	res, err := DoRequest(req, time.Now(), false, client)
	require.Error(t, err)
	require.NotNil(t, res)
}

func TestWithTraceManualHooks(t *testing.T) {
	color.NoColor = true
	initTime := time.Now()
	req, err := http.NewRequest("GET", "https://example.com", nil)
	require.NoError(t, err)

	req, traceFn := WithTrace(req, initTime)
	require.NotNil(t, traceFn)

	trace := httptrace.ContextClientTrace(req.Context())
	require.NotNil(t, trace)

	trace.DNSStart(httptrace.DNSStartInfo{Host: "localhost"})
	trace.DNSDone(httptrace.DNSDoneInfo{Err: errors.New("dns")})
	trace.DNSStart(httptrace.DNSStartInfo{Host: "localhost"})
	trace.DNSDone(httptrace.DNSDoneInfo{})

	trace.ConnectStart("tcp", "localhost:443")
	trace.ConnectDone("tcp", "localhost:443", errors.New("conn"))
	trace.ConnectStart("tcp", "localhost:443")
	trace.ConnectDone("tcp", "localhost:443", nil)

	trace.TLSHandshakeStart()
	trace.TLSHandshakeDone(tls.ConnectionState{Version: tls.VersionTLS13}, errors.New("tls"))
	trace.TLSHandshakeStart()
	trace.TLSHandshakeDone(tls.ConnectionState{Version: tls.VersionTLS13}, nil)

	trace.GotConn(httptrace.GotConnInfo{Reused: false})
	trace.GotConn(httptrace.GotConnInfo{Reused: true})

	trace.WroteRequest(httptrace.WroteRequestInfo{Err: errors.New("write")})
	trace.WroteRequest(httptrace.WroteRequestInfo{Err: nil})

	trace.GotFirstResponseByte()

	_, _, _, _, tlsVersion, sawTLS := traceFn()
	require.True(t, sawTLS)
	require.NotZero(t, tlsVersion)
}
