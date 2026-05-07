package runner

import (
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

type RequestResult struct {
	FinalTime   time.Duration
	DumpHeaders []byte
	Body        []byte
	DNSDur      time.Duration
	ConnDur     time.Duration
	TLSDur      time.Duration
	FirstByte   time.Duration
	TLSVersion  uint16
	SawTLS      bool
	Headers     http.Header
	StatusCode  int
}

func DoRequest(req *http.Request, initTime time.Time, verbose bool) (*RequestResult, error) {
	var traceFunc func() (time.Duration, time.Duration, time.Duration, time.Duration, uint16, bool)
	if verbose {
		req, traceFunc = WithTrace(req, initTime)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	finalTime := time.Since(initTime)
	if err != nil {
		return &RequestResult{FinalTime: finalTime}, err
	}
	defer resp.Body.Close()

	var dnsDur, connDur, tlsDur, firstByteDur time.Duration
	var tlsVersion uint16
	var sawTLS bool
	if verbose && traceFunc != nil {
		dnsDur, connDur, tlsDur, firstByteDur, tlsVersion, sawTLS = traceFunc()
	}

	dumpHeaders, err := httputil.DumpResponse(resp, false)
	if err != nil {
		return &RequestResult{FinalTime: finalTime, DNSDur: dnsDur, ConnDur: connDur, TLSDur: tlsDur, FirstByte: firstByteDur, TLSVersion: tlsVersion, SawTLS: sawTLS}, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RequestResult{FinalTime: finalTime, DumpHeaders: dumpHeaders, DNSDur: dnsDur, ConnDur: connDur, TLSDur: tlsDur, FirstByte: firstByteDur, TLSVersion: tlsVersion, SawTLS: sawTLS, Headers: resp.Header, StatusCode: resp.StatusCode}, err
	}

	return &RequestResult{
		FinalTime:   finalTime,
		DumpHeaders: dumpHeaders,
		Body:        bodyBytes,
		DNSDur:      dnsDur,
		ConnDur:     connDur,
		TLSDur:      tlsDur,
		FirstByte:   firstByteDur,
		TLSVersion:  tlsVersion,
		SawTLS:      sawTLS,
		Headers:     resp.Header,
		StatusCode:  resp.StatusCode,
	}, nil
}
