package runner

import (
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"

	"github.com/fatih/color"
)

func WithTrace(req *http.Request, initTime time.Time) (*http.Request, func() (time.Duration, time.Duration, time.Duration, time.Duration, uint16, bool)) {
	var (
		dnsStart     time.Time
		connStart    time.Time
		tlsStart     time.Time
		dnsDur       time.Duration
		connDur      time.Duration
		tlsDur       time.Duration
		firstByteDur time.Duration
		tlsVersion   uint16
		sawTLS       bool
	)

	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now()
			color.Cyan("[TRACE] DNS start  : host=%s", info.Host)
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			dur := time.Since(dnsStart)
			dnsDur = dur
			if info.Err != nil {
				color.Red("[TRACE] DNS done   : error=%v (elapsed=%s)", info.Err, dur)
			} else {
				color.Green("[TRACE] DNS done   : addrs=%v (elapsed=%s)", info.Addrs, dur)
			}
		},
		ConnectStart: func(network, addr string) {
			connStart = time.Now()
			color.Cyan("[TRACE] Connect start: network=%s addr=%s", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			dur := time.Since(connStart)
			connDur = dur
			if err != nil {
				color.Red("[TRACE] Connect done : error=%v (elapsed=%s)", err, dur)
			} else {
				color.Green("[TRACE] Connect done : addr=%s (elapsed=%s)", addr, dur)
			}
		},
		TLSHandshakeStart: func() {
			tlsStart = time.Now()
			color.Cyan("[TRACE] TLS handshake start")
		},
		TLSHandshakeDone: func(state tls.ConnectionState, err error) {
			dur := time.Since(tlsStart)
			tlsDur = dur
			tlsVersion = state.Version
			sawTLS = true
			if err != nil {
				color.Red("[TRACE] TLS handshake error: %v (elapsed=%s)", err, dur)
			} else {
				color.Green("[TRACE] TLS handshake done  : version=%#x (elapsed=%s)", state.Version, dur)
			}
		},
		GotConn: func(info httptrace.GotConnInfo) {
			reuse := "new"
			if info.Reused {
				reuse = "reused"
			}
			color.Yellow("[TRACE] Got connection: %s (wasIdle=%v, idleTime=%s)", reuse, info.WasIdle, info.IdleTime)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err != nil {
				color.Red("[TRACE] Write request : error=%v", info.Err)
			} else {
				color.Green("[TRACE] Write request : completed")
			}
		},
		GotFirstResponseByte: func() {
			if !initTime.IsZero() {
				firstByteDur = time.Since(initTime)
			} else {
				firstByteDur = 0
			}
			color.Blue("[TRACE] First response byte received")
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	return req, func() (time.Duration, time.Duration, time.Duration, time.Duration, uint16, bool) {
		return dnsDur, connDur, tlsDur, firstByteDur, tlsVersion, sawTLS
	}
}
